package indexer

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/umbracle/eth-indexer/indexer/proto"
)

func (s *Server) Apply(ctx context.Context, req *proto.Component) (*empty.Empty, error) {
	fmt.Println("X")
	return &empty.Empty{}, nil
}
