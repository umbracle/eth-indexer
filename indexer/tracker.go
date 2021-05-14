package indexer

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/umbracle/eth-indexer/indexer/proto"
	"github.com/umbracle/eth-indexer/sdk"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"github.com/umbracle/go-web3/tracker"
)

type trackerSrv struct {
	logger   hclog.Logger
	srv      *Server
	tracker  *tracker.Tracker
	provider *jsonrpc.Client
}

var (
	zeroAddr = web3.Address{}
)

func (t *trackerSrv) setupTracker(indexer *sdk.Provider) error {
	provider, err := jsonrpc.NewClient(t.srv.config.JSONRPCEndpoint)
	if err != nil {
		return err
	}
	t.provider = provider

	t.logger.Info("start tracker", "batch", t.srv.config.BatchSize)

	tConfig := tracker.DefaultConfig()
	tConfig.BatchSize = t.srv.config.BatchSize
	tConfig.EtherscanFastTrack = true

	// use an sql store for the tracker
	store, err := New("./test.db")
	if err != nil {
		return err
	}

	t.tracker = tracker.NewTracker(provider.Eth(), tConfig)
	t.tracker.SetStore(store)

	go func() {
		if err := t.tracker.Start(context.Background()); err != nil {
			t.logger.Error("failed to track", "err", err)
		}
	}()

	// wait for the tracker to be ready
	<-t.tracker.ReadyCh

	filter := indexer.GetFilter()
	track := &proto.Track{
		Name:       "-",
		StartBlock: filter.StartBlock,
	}
	if filter.FromAddr != zeroAddr {
		track.FromAddr = filter.FromAddr.String()
	}
	if filter.ToAddr != zeroAddr {
		track.ToAddr = filter.ToAddr.String()
	}

	if err := t.srv.state.UpsertTrack(track); err != nil {
		return err
	}
	if err := t.startTrack(track, indexer); err != nil {
		return err
	}
	return nil
}

func filterConfigFromTracker(t *proto.Track) (*tracker.FilterConfig, error) {
	config := &tracker.FilterConfig{
		Async:   false,
		To:      []web3.Address{},
		Address: []web3.Address{},
		Hash:    t.Name,
		Start:   t.StartBlock,
		Topics:  []*web3.Hash{},
	}

	if t.ToAddr != "" {
		var addr web3.Address
		if err := addr.UnmarshalText([]byte(t.ToAddr)); err != nil {
			return nil, err
		}
		config.To = []web3.Address{addr}
	}
	if t.FromAddr != "" {
		var addr web3.Address
		if err := addr.UnmarshalText([]byte(t.FromAddr)); err != nil {
			return nil, err
		}
		config.Address = []web3.Address{addr}
	}
	if t.Topic != "" {
		var hash web3.Hash
		if err := hash.UnmarshalText([]byte(t.Topic)); err != nil {
			return nil, err
		}
		config.Topics = append(config.Topics, &hash)
	}
	return config, nil
}

func processEvents(logs []*web3.Log) []*sdk.Action {
	res := []*sdk.Action{}

	act := &sdk.Action{
		BlockNum: logs[0].BlockNumber,
	}

	for _, log := range logs {
		if log.BlockNumber != act.BlockNum {
			sort.Sort(act.Events)
			res = append(res, act)
			act = &sdk.Action{
				BlockNum: log.BlockNumber,
			}
		}
		act.Events = append(act.Events, *proto.DecodeEvent(log))
	}

	sort.Sort(act.Events)
	res = append(res, act)
	return res
}

func (t *trackerSrv) startTrack(track *proto.Track, indexer *sdk.Provider) error {
	fConfig, err := filterConfigFromTracker(track)
	if err != nil {
		return err
	}

	filter, err := t.tracker.NewFilter(fConfig)
	if err != nil {
		return err
	}

	lastBlock, err := filter.GetLastBlock()
	if err != nil {
		return err
	}
	if lastBlock != nil {
		t.logger.Debug("last block", "block", lastBlock.Number)
	}

	go func() {
		for {
			select {
			case num := <-filter.SyncCh:
				fmt.Printf("--- %s %s num %d\n", track.Name, time.Now(), num)

			case evnt := <-filter.EventCh:
				if len(evnt.Added) == 0 {
					continue
				}

				actions := processEvents(evnt.Added)
				for _, act := range actions {
					diffs, err := indexer.Process(act)
					if err != nil {
						fmt.Printf("Failed to process: %v", err)
						return
					} else {
						if err := t.srv.state.ApplyDiff(diffs, true); err != nil {
							fmt.Printf("Failed to apply diff: %v", err)
							return
						}
					}
				}

			case <-filter.DoneCh:
				fmt.Printf("-- evnt done %s\n", track.Name)
			}
		}
	}()

	go func() {
		if err := filter.Sync(context.Background()); err != nil {
			t.logger.Error("failed to sync", "err", err)
		}
	}()

	t.logger.Debug("start track")
	return nil
}
