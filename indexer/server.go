package indexer

import (
	"context"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/umbracle/eth-indexer/indexer/proto"
	"google.golang.org/grpc"
)

type Config struct {
	GRPCAddr        *net.TCPAddr
	JSONRPCEndpoint string
	Database        string
	BatchSize       uint64
	Provider        string
}

type Server struct {
	proto.UnimplementedIndexerServiceServer

	config     *Config
	logger     hclog.Logger
	grpcServer *grpc.Server

	//tracker *trackerSrv
	//state   *State
	//schemas map[string]*sdk.Table
}

func NewServer(config *Config, logger hclog.Logger) (*Server, error) {
	srv := &Server{
		config:     config,
		logger:     logger,
		grpcServer: grpc.NewServer(),
	}

	if err := srv.setupGRPCServer(config.GRPCAddr.String()); err != nil {
		return nil, err
	}

	/*	state, err := newState(config.Database)
		if err != nil {
			return nil, err
		}
		state.i = srv
		srv.state = state

		// srv.addSchema()

		srv.tracker = &trackerSrv{
			logger: logger.Named("tracker"),
			srv:    srv,
		}

		// srv.addIndexers()

		indexer, err := srv.setupIndexer()
		if err != nil {
			return nil, err
		}
		if err := srv.tracker.setupTracker(indexer); err != nil {
			return nil, err
		}
	*/

	return srv, nil
}

func (s *Server) setupGRPCServer(addr string) error {
	s.grpcServer = grpc.NewServer(s.withLoggingUnaryInterceptor())

	proto.RegisterIndexerServiceServer(s.grpcServer, s)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Error("failed to serve grpc server", "err", err)
		}
	}()

	s.logger.Info("Server started", "addr", addr)
	return nil
}

func (s *Server) withLoggingUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(s.loggingServerInterceptor)
}

func (s *Server) loggingServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	s.logger.Trace("Request", "method", info.FullMethod, "duration", time.Since(start), "error", err)
	return h, err
}

func (s *Server) Stop() {

}

/*
func (s *Server) setupIndexer() (*sdk.Provider, error) {
	provider, err := jsonrpc.NewClient(s.config.JSONRPCEndpoint)
	if err != nil {
		return nil, err
	}

	indexer := providers.BuiltinProviders[s.config.Provider]()
	if err := indexer.Init(); err != nil {
		return nil, err
	}
	indexer.SetClient(provider)
	indexer.SetStateResolver(s)

	sss := indexer.GetSchemas()
	for _, sch := range sss.Schemas {
		s.schemas[sch.Name] = sch
	}

	// write the tables
	for _, sch := range indexer.GetSchemas().Schemas {
		if err := s.state.UpsertTable(sch); err != nil {
			return nil, err
		}
	}
	return indexer, nil
}

func (s *Server) Stop() {
	// TODO
}

func (s *Server) GetObj2(table string, keys map[string]string) (*sdk.Obj, error) {
	raw, err := s.state.GetObj2(table, keys)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	return &sdk.Obj{Data: raw.Data}, nil
}

func (s *Server) GetObjs2(q *sdk.Query) ([]*sdk.Obj, error) {
	return nil, nil
}
*/
