package model

import "github.com/wuhen781/Tx-Parser/internal/database"
import "errors"

var ErrBlockNumberNotInitialed = errors.New("lastBlockNumber is not initialed")
var ErrBLockNumberIsUpdated = errors.New("lastBlockNumber is updated")

type ModelParser struct {
	db database.Db
}

func NewModelParser() *ModelParser {
	db := database.NewMemoryDb()
	return &ModelParser{
		db,
	}
}

func (this *ModelParser) GetTransactions(address string) []database.Transaction {
	return this.db.GetTransactions(address)
}

func (this *ModelParser) AddSubscribe(address string, blockNumber int) bool {
	return this.db.AddSubscribe(address, blockNumber)
}

func (this *ModelParser) updateTransactionsByLastBlockNumber(currentBlockNumber int, transactions []database.Transaction) error {
	lastBlockNumber := this.db.GetLastUpdatedBlockNumber()
	if lastBlockNumber < 0 {
		return ErrBlockNumberNotInitialed
	}
	transLen := len(transactions)
	transOffset := this.db.GetTransOffetsInLastBlock()
	transactions = transactions[transOffset:transLen] //skip the part that has already been saved
	subscribers := this.db.GetSubscribeFromBlockNumber(lastBlockNumber)
	subMaps := make(map[string]bool)
	for _, subscriber := range subscribers {
		subMaps[subscriber] = true
	}
	subscrbTrans := make([]database.Transaction, 0)
	for _, transaction := range transactions {
		from, to := transaction.From, transaction.To
		if subMaps[from] == true || subMaps[to] == true {
			subscrbTran := database.Transaction{
				From:        transaction.From,
				To:          transaction.To,
				Value:       transaction.Value,
				BlockNumber: transaction.BlockNumber,
				Gas:         transaction.Gas,
				GasPrice:    transaction.GasPrice,
				Hash:        transaction.Hash,
				Nonce:       transaction.Nonce,
				Timestamp:   transaction.Timestamp,
				Type:        "",
			}
			subscrbTrans = append(subscrbTrans, subscrbTran)
		}
	}
	if this.db.SetTransactions(subscrbTrans) {
		if currentBlockNumber == lastBlockNumber {
			this.db.SetTransOffetsInLastBlock(transLen)
		} else {
			this.db.SetTransOffetsInLastBlock(0)
			this.db.SetLastUpdatedBlockNumber(lastBlockNumber + 1)
		}
	}
	return nil
}
