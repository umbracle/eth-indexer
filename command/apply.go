package command

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/umbracle/eth-indexer/indexer/proto"
)

// ApplyCommand runs a script
type ApplyCommand struct {
	*Meta
}

// Help implements the cli.ApplyCommand interface
func (c *ApplyCommand) Help() string {
	return ""
}

// Synopsis implements the cli.ApplyCommand interface
func (c *ApplyCommand) Synopsis() string {
	return ""
}

// Run implements the cli.ApplyCommand interface
func (c *ApplyCommand) Run(args []string) int {
	flags := c.FlagSet("apply")
	if err := flags.Parse(args); err != nil {
		panic(err)
	}

	filename := args[0]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	clt, err := c.IndexerClient()
	if err != nil {
		panic(err)
	}
	if _, err := clt.Apply(context.Background(), &proto.Component{}); err != nil {
		panic(err)
	}
	return 0
}
