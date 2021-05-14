package proto

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/umbracle/go-web3"
)

func decodeHex(str string) ([]byte, error) {
	if !strings.HasPrefix(str, "0x") {
		return nil, fmt.Errorf("0x prefix not found in data")
	}
	buf, err := hex.DecodeString(str[2:])
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (e *Event) ToLog() (*web3.Log, error) {
	log := &web3.Log{}
	log.TransactionIndex = e.TxIndex
	if err := log.TransactionHash.UnmarshalText([]byte(e.TxHash)); err != nil {
		return nil, err
	}
	log.BlockNumber = e.BlockNum
	if err := log.BlockHash.UnmarshalText([]byte(e.BlockHash)); err != nil {
		return nil, err
	}
	if err := log.Address.UnmarshalText([]byte(e.Address)); err != nil {
		return nil, err
	}

	if e.Topics != "" {
		log.Topics = []web3.Hash{}
		for _, item := range strings.Split(e.Topics, ",") {
			var topic web3.Hash
			if err := topic.UnmarshalText([]byte(item)); err != nil {
				return nil, err
			}
			log.Topics = append(log.Topics, topic)
		}
	} else {
		log.Topics = nil
	}

	if e.Data != "" {
		buf, err := decodeHex(e.Data)
		if err != nil {
			return nil, err
		}
		log.Data = buf
	} else {
		log.Data = nil
	}
	return log, nil
}

func DecodeEvent(log *web3.Log) *Event {
	topics := []string{}
	for _, topic := range log.Topics {
		topics = append(topics, topic.String())
	}
	evnt := &Event{
		LogIndex:  log.LogIndex,
		TxIndex:   log.TransactionIndex,
		TxHash:    log.TransactionHash.String(),
		BlockNum:  log.BlockNumber,
		BlockHash: log.BlockHash.String(),
		Address:   log.Address.String(),
		Topics:    strings.Join(topics, ","),
		Removed:   false,
	}
	if len(topics) != 0 {
		evnt.TopicID = topics[0]
	}
	if log.Data != nil {
		evnt.Data = "0x" + hex.EncodeToString(log.Data)
	}
	return evnt
}

/*
func (t *Transaction) ToTxn() (*web3.Transaction, error) {
	txn := &web3.Transaction{}
	if err := txn.Hash.UnmarshalText([]byte(t.Hash)); err != nil {
		return nil, err
	}
	if err := txn.From.UnmarshalText([]byte(t.FromAddr)); err != nil {
		return nil, err
	}
	if t.ToAddr != "" {
		if err := txn.To.UnmarshalText([]byte(t.ToAddr)); err != nil {
			return nil, err
		}
	}
	if t.Input != "" {
		buf, err := decodeHex(t.Input)
		if err != nil {
			return nil, err
		}
		txn.Input = buf
	} else {
		txn.Input = nil
	}
	txn.GasPrice = t.GasPrice
	txn.Gas = t.Gas
	value, ok := new(big.Int).SetString(t.Value, 10)
	if !ok {
		return nil, fmt.Errorf("cannot decode difficulty")
	}
	txn.Value = value
	txn.Nonce = t.Nonce
	if err := txn.BlockHash.UnmarshalText([]byte(t.BlockHash)); err != nil {
		return nil, err
	}
	txn.BlockNumber = t.BlockNum
	txn.TxnIndex = t.TxIndex
	return txn, nil
}

func DecodeTransaction(txn *web3.Transaction) *Transaction {
	obj := &Transaction{
		Hash:      txn.Hash.String(),
		FromAddr:  txn.From.String(),
		Input:     "0x" + hex.EncodeToString(txn.Input),
		GasPrice:  txn.GasPrice,
		Gas:       txn.Gas,
		Value:     txn.Value.String(),
		Nonce:     txn.Nonce,
		TxIndex:   txn.TxnIndex,
		BlockHash: txn.BlockHash.String(),
		BlockNum:  txn.BlockNumber,
	}
	if txn.To != nil {
		obj.ToAddr = txn.To.String()
	}
	return obj
}
*/
