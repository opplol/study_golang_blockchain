package blockchain

import (
	"errors"
	"go_crypo_coin/utils"
	"time"
)

const (
	minerReward int = 50
)

type mempoll struct {
	Txs []*Tx
}

var Mempoll *mempoll = &mempoll{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxId   string
	Index  int
	Amount int
}

func isOnMemPool(UTxOut *UTxOut) bool {
	for _, tx := range Mempoll.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == UTxOut.TxId && input.Index == UTxOut.Index {
				return true
			}
		}
	}
	return false

	// Use label for break for loop
	// 	exsits := false
	// Outer:
	// 	for _, tx := range Mempoll.Txs {
	// 		for _, input := range tx.TxIns {
	// 			if input.TxID == UTxOut.TxId && input.Index == UTxOut.Index {
	// 				exsits = true
	// 				break Outer
	// 			}
	// 		}
	// 	}
	// 	exsits
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("not enoguh money")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UnspandTxOutsByAddress(from, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxId, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		changeTxout := &TxOut{from, change}
		txOuts = append(txOuts, changeTxout)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil

}

func (m *mempoll) AddTx(to string, amount int) error {
	tx, err := makeTx("nico", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempoll) txToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("nico")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
