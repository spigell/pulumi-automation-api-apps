package snapshots

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/spigell/pulumi-automation-api-apps/hetzner-snapshots-manager/manager"
	"github.com/spigell/pulumi-automation-api-apps/common/pulumi"
)

const (
)

var (
	ErrSnapshotNotFound = errors.New("snapshot not found")
	ErrValidationServerMissing = errors.New(manager.ErrServerMissing)
)

type SnapshotResponce struct {
	Body   Body `json:"body"`
	Error  string
	Status string
	ID     string
}

type Body struct {
	ID int
}


// A client function for getting last snapshot info.
// It can be used in tests or pulumi program.
func GetLastSnapshot(client *http.Client, serverName string) (*SnapshotResponce, error) {
	var responce *SnapshotResponce

	if serverName == "" {
		return nil, fmt.Errorf("validate: %w", ErrValidationServerMissing)
	}

	url := url.URL{
		Scheme:   "http",
		Host:     os.Getenv(pulumi.EnvAutomaionAPIAddr),
		Path:     manager.SnapshotsAPIPath,
		RawQuery: fmt.Sprintf("server=%s", serverName),
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &responce)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return responce, ErrSnapshotNotFound
	}

	if responce.Status != "OK" {
		return responce, fmt.Errorf("status: %s", responce.Status)
	}

	if responce.Error != "" {
		return responce, fmt.Errorf(responce.Error)
	}

	return responce, nil
}
