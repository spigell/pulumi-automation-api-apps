package manager

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

func (m *Manager) Run(cmd *cobra.Command) error {
	onlyAPIServer, err := cmd.Flags().GetBool("only-api-server")
	if err != nil {
		m.Logger.Fatal(err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = m.APIServer.Run()
		if err != nil {
			m.Logger.Fatal(fmt.Sprint("running server: %w", err))
		}
	}()

	if onlyAPIServer {
		m.Logger.Info("will not go further since running in only-api-server mode. Press enter CTRL-C to stop")

		wg.Wait()

		return nil
	}

	events, err := m.Runner.Preview(m.ctx, true)
	if err != nil {
		return fmt.Errorf("running preview: %w", err)
	}

	if err := m.ProcessEvents(events); err != nil {
		return fmt.Errorf("process events: %w", err)
	}

	switch m.Runner.GetMode() {
	case "up":
		if err := m.Runner.Up(m.ctx); err != nil {
			return fmt.Errorf("running up: %w", err)
		}
	case "destroy":
		if err := m.Runner.Destroy(m.ctx); err != nil {
			return fmt.Errorf("running up: %w", err)
		}
	}

	m.APIServer.Close()

	wg.Wait()

	return nil
}
