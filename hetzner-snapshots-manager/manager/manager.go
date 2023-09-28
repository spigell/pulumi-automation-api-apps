package manager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spigell/pulumi-automation-api-apps/common/apiserver"
	"github.com/spigell/pulumi-automation-api-apps/common/log"
	"github.com/spigell/pulumi-automation-api-apps/common/pulumi"
	"github.com/spigell/pulumi-automation-api-apps/hetzner-snapshots-manager/hetzner"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"go.uber.org/zap"
)

const (
	serverType = "hcloud:index/server:Server"
)

type Manager struct {
	ctx       context.Context
	Runner    *pulumi.Pulumi
	APIServer *apiserver.Server
	Hetzner   *hetzner.API
	Logger    *zap.Logger
	Snapshots *hetzner.Snapshots
	Cleaner   *Cleaner
}

type Cleaner struct {
	maxKeep int
}

func New(ctx context.Context, command string) (*Manager, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, fmt.Errorf("create config: %w", err)
	}

	logger, _ := log.New(config.Verbose)

	runner, err := pulumi.New(ctx, logger, config.Stack.Name, config.Stack.Path, command, config.Diff)
	if err != nil {
		return nil, fmt.Errorf("create pulumi runner: %w", err)
	}
	hcloudToken := config.Token

	hetzner := hetzner.New(ctx, logger, hcloudToken)

	snapshots, err := hetzner.GatherSnapshotInfo()
	if err != nil {
		return nil, fmt.Errorf("retrieve info about snapshots: %w", err)
	}

	httpAddr := fmt.Sprintf("localhost:%d", config.APIServerPort)

	apiServer, err := apiserver.New(httpAddr, logger, getAllRoutes(snapshots))
	if err != nil {
		return nil, fmt.Errorf("create api server: %w", err)
	}

	// pass needed info to pulumi cli. Bad naming, but I like to `attach` :)
	runner.AttachToAPIServer(apiServer.Addr())

	// set token for pulumi program
	os.Setenv("HCLOUD_TOKEN", hcloudToken)

	return &Manager{
		ctx:       ctx,
		Runner:    runner,
		APIServer: apiServer,
		Logger:    logger,
		Hetzner:   hetzner,
		Snapshots: snapshots,
		Cleaner: &Cleaner{
			maxKeep: config.MaxKeep,
		},
	}, nil
}

func (m *Manager) ProcessEvents(events []events.EngineEvent) error {
	var wg sync.WaitGroup
	var err error
	for _, p := range events {
		if p.ResourcePreEvent != nil { //nolint:nestif
			metadata := p.ResourcePreEvent.Metadata
			if metadata.Type == serverType {
				switch {
				case metadata.Op == "delete":
					// Use goroutines since making snapshots is a time consuming operation.
					// We can do it in parallel.
					wg.Add(1)
					go func() error {
						err = m.processServerDeletion(metadata.Old.URN)
						if err != nil {
							return fmt.Errorf("process server deletion: %w", err)
						}
						wg.Done()
						return nil
					}()

				case m.Runner.GetMode() == "destroy":
					wg.Add(1)
					go func() error {
						err = m.processServerDeletion(metadata.Old.URN)
						if err != nil {
							return fmt.Errorf("process server deletion: %w", err)
						}
						wg.Done()
						return nil
					}()

				default:
					continue
				}
			}
		}
	}

	wg.Wait()

	return nil
}

func (m *Manager) processServerDeletion(urn string) error {
	m.Logger.Info(fmt.Sprintf("we are deleting a server resource %s. a snapshot creation needed", urn))

	// name of resource is the last :: separated part
	name := strings.Split(urn, "::")[len(strings.Split(urn, "::"))-1]

	// Decrease maxKeep by 1 since we are creating a new snapshot
	oldSnaphots, err := m.Snapshots.GetStalledSnapshots(name, m.Cleaner.maxKeep-1)
	if err != nil {
		switch {
		case errors.Is(err, hetzner.ErrSnapshotNotFound):
			oldSnaphots = []hetzner.SnapshotInfo{}
		default:
			return fmt.Errorf("get stalled snapshots: %w", err)
		}
	}

	m.Logger.Debug(fmt.Sprintf("stalled snapshots: %+v", oldSnaphots))

	for _, s := range oldSnaphots {
		if m.Runner.IsPreview() {
			m.Logger.Info(fmt.Sprintf("dry run, skipping a deletion for stalled snapshot with id: %d", s.ID))
			continue
		}
	}

	if m.Runner.IsPreview() {
		m.Logger.Info("dry run, skipping a snapshot creation")
		return nil
	}

	err = m.makeSnapshot(name)
	if err != nil {
		return fmt.Errorf("make snapshot: %w", err)
	}
	for _, s := range oldSnaphots {
		timeout := 1 * time.Minute
		ctx, cancel := context.WithTimeout(m.ctx, timeout)
		if err := m.Hetzner.DeleteSnapshot(ctx, s.ID); err != nil {
			cancel()
			return fmt.Errorf("delete the old snapshot: %w", err)
		}
		cancel()
	}

	return nil
}

func (m *Manager) makeSnapshot(id string) error {
	timeout := 20 * time.Minute
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	m.Logger.Info(fmt.Sprintf("the snapshot creation for %s may take some time. Please be patient. Max allowed time is %s",
		id, timeout.String(),
	))

	start := time.Now()
	err := m.Hetzner.CreateSnapshot(ctx, id)

	m.Logger.Debug(fmt.Sprintf("the snapshot creation took %s", time.Since(start).String()))

	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	return nil
}
