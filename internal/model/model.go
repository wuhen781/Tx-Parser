package model

import "github.com/wuhen781/Tx-Parser/internal/database"

var ErrBlockNumberNotInitialed = errors.New("lastBlcokNumber is not initialed")
var ErrBLockNumberIsUpdated = errors.New("lastBlcokNumber is updated")

type ModelParser struct {
	db Db
}

func NewModelParser() *ModelParser {
	db := NewMemoryDb()
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

func (this *ModelParser) updateTransactionsByLastBlockNumber() error {
	lastBlcokNumber := this.db.GetLastUpdatedBlcokNumber()
	if lastBlcokNumber < 0 {
		return ErrBlockNumberNotInitialed
	}
	client := ethclient.NewEthclient("")
	currentBlcokNumber, err := client.GetCurrentBlock()
	if err != ni {
		fmt.Errorf("GetCurrentBlock error : %w", err)
		return err
	}
	transactions, err2 := client.GetBlockByNumber(lastBlcokNumber)
	if err2 != nil {
		fmt.Errorf("GetBlockByNumber error : %w", err2)
		return err2
	}
	transLen := len(transactions)
	transOffset := this.db.GetTransOffetsInLastBlock()
	transactions = transactions[transOffset:transLen] //skip the part that has already been saved
	subscribers := this.db.GetSubscribeFromBlockNumber(blockNumber)
	subMaps := make(map[string]bool)
	for _, subscriber := range subscribers {
		subMaps[subscriber] = true
	}
	subscrbTrans := make([]database.Transaction, 0)
	for _, transaction := range transactions {
		from, to := transaction.from, transaction.to
		if subMaps[from] == true {
			subscrbTran := data.Transaction{
				From: transaction.From,
				To:   transaction, To,
				Value:       transaction.Value,
				BlockNumber: transaction.BlockNumber,
				Gas:         transaction.Gas,
				GasPrice:    transaction.GasPrice,
				Hash:        transaction.Hash,
				Nonce:       transaction.Nonce,
				Timestamp:   transaction.Timestamp,
				Type:        "outbound",
			}
			subscrbTrans = append(subscrbTrans, subscrbTran)
		}
		if subMaps[to] == true {
			subscrbTran := data.Transaction{
				From: transaction.From,
				To:   transaction, To,
				Value:       transaction.Value,
				BlockNumber: transaction.BlockNumber,
				Gas:         transaction.Gas,
				GasPrice:    transaction.GasPrice,
				Hash:        transaction.Hash,
				Nonce:       transaction.Nonce,
				Timestamp:   transaction.Timestamp,
				Type:        "inbound",
			}
			subscrbTrans = append(subscrbTrans, subscrbTran)
		}
	}
	if this.db.SetTransactions(subscrbTrans) {
		if currentBlcokNumber == lastBlcokNumber {
			this.db.SetTransOffetsInLastBlock(transLen)
		} else {
			this.db.SetTransOffetsInLastBlock(0)
			this.db.SetLastUpdatedBlcokNumber(lastBlcokNumber + 1)
		}
	}
	return nil
}
