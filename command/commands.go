package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"google.golang.org/grpc"
)

// Commands returns the cli commands
func Commands() map[string]cli.CommandFactory {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	_ = Meta{
		UI: ui,
	}

	return map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				UI: ui,
			}, nil
		},
	}
}

// Meta is a helper utility for the commands
type Meta struct {
	UI   cli.Ui
	addr string
}

// FlagSet adds some default commands to handle grpc connections with the server
func (m *Meta) FlagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)
	f.StringVar(&m.addr, "address", "127.0.0.1:6001", "Address of the http api")
	return f
}

// Conn returns a grpc connection
func (m *Meta) Conn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(m.addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}
	return conn, nil
}
