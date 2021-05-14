package testutil

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
)

// TestServer is an eth1x testutil server using go-ethereum
type TestServer struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	t        *testing.T
}

// NewTestServer creates a new eth1 server with go-ethereum
func NewTestServer(t *testing.T) *TestServer {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	cmd := []string{
		"--dev",
		"--http", "--http.addr", "0.0.0.0",
		"--ws", "--ws.addr", "0.0.0.0",
	}
	opts := &dockertest.RunOptions{
		Repository: "ethereum/client-go",
		Tag:        "latest",
		Cmd:        cmd,
	}
	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		t.Fatalf("Could not start go-ethereum: %s", err)
	}

	server := &TestServer{
		pool:     pool,
		resource: resource,
		t:        t,
	}
	if err := pool.Retry(func() error {
		return testHTTPEndpoint(server.Http())
	}); err != nil {
		server.Stop()
	}
	return server
}

// Fund funds a web3 address
func (t *TestServer) Fund(addr web3.Address, value *big.Int) error {
	txn := web3.Transaction{
		From:  t.Owner(),
		To:    &addr,
		Value: value,
	}
	if _, err := t.sendTxn(txn); err != nil {
		return err
	}
	return nil
}

// ProcessBlock advances one block on the chain
func (t *TestServer) ProcessBlock() error {
	zeroAddr := web3.Address{}
	txn := web3.Transaction{
		From:  t.Owner(),
		To:    &zeroAddr,
		Value: big.NewInt(1),
	}
	if _, err := t.sendTxn(txn); err != nil {
		return err
	}
	return nil
}

func (t *TestServer) sendTxn(txn web3.Transaction) (*web3.Receipt, error) {
	// gas price
	gasPrice, err := t.Provider().Eth().GasPrice()
	if err != nil {
		return nil, err
	}

	msg := &web3.CallMsg{
		From:     txn.From,
		To:       txn.To,
		GasPrice: gasPrice,
		Value:    txn.Value,
		Data:     txn.Input,
	}
	gasLimit, err := t.Provider().Eth().EstimateGas(msg)
	if err != nil {
		return nil, err
	}

	txn.GasPrice = gasPrice
	txn.Gas = gasLimit

	hash, err := t.Provider().Eth().SendTransaction(&txn)
	if err != nil {
		return nil, err
	}
	receipt, err := t.WaitForReceipt(hash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (t *TestServer) WaitForReceipt(hash web3.Hash) (*web3.Receipt, error) {
	var receipt *web3.Receipt
	var err error

	for {
		receipt, err = t.Provider().Eth().GetTransactionReceipt(hash)
		if err != nil {
			if err.Error() != "not found" {
				return nil, err
			}
		}
		if receipt != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return receipt, nil
}

// Stop closes the eth1 test server
func (t *TestServer) Stop() {
	if err := t.pool.Purge(t.resource); err != nil {
		t.t.Fatalf("Could not purge go-ethereum: %s", err)
	}
}

func (t *TestServer) Http() string {
	return fmt.Sprintf("http://%s:8545", t.resource.Container.NetworkSettings.IPAddress)
}

// Provider returns the jsonrpc provider
func (t *TestServer) Provider() *jsonrpc.Client {
	provider, _ := jsonrpc.NewClient(t.Http())
	return provider
}

// Owner returns the account with balance on go-ethereum
func (t *TestServer) Owner() web3.Address {
	owner, _ := t.Provider().Eth().Accounts()
	return owner[0]
}

func testHTTPEndpoint(endpoint string) error {
	resp, err := http.Post(endpoint, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
