package sdk

import (
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/umbracle/eth-indexer/indexer/proto"
	protosdk "github.com/umbracle/eth-indexer/sdk/proto"
)

type Backend interface {
	// Returns the available schemas for the backend
	GetSchemas() GetSchemasResponse

	// Process an action for the backend
	Process(ac *Action) ([]*protosdk.Diff, *ErrorEvent)

	// Returns filters that encode how to parse stuff
	GetFilter() *FilterByAddr

	// TODO: Return access to the state to query info
}

type GetSchemasResponse struct {
	Schemas []*Table
}

type Action struct {
	BlockNum uint64
	Events   Events
}

type Events []proto.Event

func (e Events) Len() int {
	return len(e)
}

func (e Events) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Events) Less(i, j int) bool {
	return e[i].LogIndex < e[j].LogIndex
}

func UUID() string {
	return uuid.New().String()
}

func DecodeEvent(evnt interface{}, vals map[string]interface{}) error {
	if err := mapstructure.Decode(vals, &evnt); err != nil {
		return err
	}
	return nil
}
