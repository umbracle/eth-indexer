package command

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/gops/agent"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"github.com/umbracle/eth-indexer/indexer"
)

// ServerCommand is the ServerCommand to run the agent
type ServerCommand struct {
	UI cli.Ui
}

// Help implements the cli.ServerCommand interface
func (c *ServerCommand) Help() string {
	return ""
}

// Synopsis implements the cli.ServerCommand interface
func (c *ServerCommand) Synopsis() string {
	return ""
}

// Run implements the cli.ServerCommand interface
func (c *ServerCommand) Run(args []string) int {
	if err := agent.Listen(agent.Options{}); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	flags := flag.NewFlagSet("server", flag.ContinueOnError)
	flags.Usage = func() {}

	var endpoint string
	var database string
	var batchSize uint64
	var provider string

	flags.StringVar(&endpoint, "endpoint", "", "")
	flags.StringVar(&database, "database", "postgres://postgres@localhost:5432/postgres?sslmode=disable", "")
	flags.Uint64Var(&batchSize, "batch-size", 5000, "")
	flags.StringVar(&provider, "provider", "pancake", "")

	if err := flags.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("failed to parse args: %v", err))
		return 1
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "indexer",
		Level: hclog.LevelFromString("debug"),
	})

	config := &indexer.Config{
		GRPCAddr:        &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 6001},
		JSONRPCEndpoint: endpoint,
		Database:        database,
		BatchSize:       batchSize,
		Provider:        provider,
	}
	srv, err := indexer.NewServer(config, logger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to start the server: %v", err))
		return 1
	}
	return c.handleSignals(srv.Stop)
}

func (c *ServerCommand) handleSignals(closeFn func()) int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	var sig os.Signal
	select {
	case sig = <-signalCh:
	}

	c.UI.Output(fmt.Sprintf("Caught signal: %v", sig))
	c.UI.Output("Gracefully shutting down agent...")

	gracefulCh := make(chan struct{})
	go func() {
		if closeFn != nil {
			closeFn()
		}
		close(gracefulCh)
	}()

	select {
	case <-signalCh:
		return 1
	case <-time.After(5 * time.Second):
		return 1
	case <-gracefulCh:
		return 0
	}
}
