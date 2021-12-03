package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/umbracle/eth-indexer/indexer/proto"
	"google.golang.org/grpc"
)

// Commands returns the cli commands
func Commands() map[string]cli.CommandFactory {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	meta := &Meta{
		UI: ui,
	}

	return map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				UI: ui,
			}, nil
		},
		"apply": func() (cli.Command, error) {
			return &ApplyCommand{
				meta,
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

func (m *Meta) IndexerClient() (proto.IndexerServiceClient, error) {
	conn, err := m.Conn()
	if err != nil {
		return nil, err
	}
	return proto.NewIndexerServiceClient(conn), nil
}
