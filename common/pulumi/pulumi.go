package pulumi

import (
	"context"
	"fmt"
	"sync"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"go.uber.org/zap"
)

const (
	// EnvAutomaionAPIAddr is an environment variable name for API server address
	// gosec G101
	EnvAutomaionAPIAddr = "AUTOMATION_API_HTTP_ADDR"
)

type Pulumi struct {
	mode string
	stack  auto.Stack
	logger *zap.Logger
	ctx    context.Context
	diff bool
}

func New(ctx context.Context, logger *zap.Logger, name, workDir, command string, diff bool) (*Pulumi, error) {
	selected, err := auto.SelectStackLocalSource(ctx, name, workDir)
	if err != nil {
		return nil, fmt.Errorf("get stack: %w", err)
	}

	logger.Info(fmt.Sprintf("selected stack is %q\n", name))

	return &Pulumi{
		mode:   command,
		stack:  selected,
		logger: logger,
		ctx:    ctx,
		diff: diff,
	}, nil
}

func (p *Pulumi) GetMode() string {
	return p.mode
}

func (p *Pulumi) IsPreview() bool {
	return p.mode == "preview"
}

func (p *Pulumi) GetConfig(name string) (string, error) {
	cfg, err := p.stack.GetConfig(p.ctx, name)
	if err != nil {
		return "", fmt.Errorf("get config value: %w", err)
	}

	return cfg.Value, nil
}

func (p *Pulumi) AttachToAPIServer(addr string) {
	p.stack.Workspace().SetEnvVar("AUTOMATION_API_HTTP_ADDR", addr)
	p.logger.Info(fmt.Sprintf("successfully configured for attaching to API server: addr `%s`", addr))
}

func collectEvents(eventChannel <-chan events.EngineEvent, events *[]events.EngineEvent) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go (func() {
		for event := range eventChannel {
			*events = append(*events, event)
		}
		wg.Done()
	})()
	return &wg
}

func GetEvent([]events.EngineEvent, string) events.EngineEvent {
	return events.EngineEvent{}
}
