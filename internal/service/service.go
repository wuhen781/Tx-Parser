package service

import "github.com/wuhen781/Tx-Parser/pkg/ethclient"
import "github.com/wuhen781/Tx-Parser/internal/database"
import "github.com/wuhen781/Tx-Parser/internal/model"
import "log"
import "time"
import "context"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []database.Transaction
}

type EthParser struct {
	modelParser *model.ModelParser
}

func NewEthParser() *EthParser {
	return &EthParser{
		modelParser: model.NewModelParser(),
	}
}

func (this *EthParser) GetCurrentBlock() (int, error) {
	client := ethclient.NewEthclient("")
	blockNumber, err := client.GetCurrentBlock()
	if err != nil {
		log.Printf("Error getting current block: %v", err)
	}
	return blockNumber, err
}

func (this *EthParser) Subscribe(address string) bool {
	client := ethclient.NewEthclient("")
	blockNumber, err := client.GetCurrentBlock()
	if err != nil {
		log.Printf("Error getting current block: %v", err)
		return false
	}
	return this.modelParser.AddSubscribe(address, blockNumber)
}

func (this *EthParser) GetTransactions(address string) []database.Transaction {
	return this.modelParser.GetTransactions(address)
}

//Update transactions regularly in the background
func (this *EthParser) UpdateTransactionsInBackGroundRegularly(ctx context.Context, interval int) {
	timer := time.NewTimer(0)
	defer timer.Stop()

	client := ethclient.NewEthclient("")

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}

		lastBlockNumber := this.modelParser.GetDb().GetLastUpdatedBlockNumber()
		if lastBlockNumber < 0 {
			timer.Reset(time.Duration(0) * time.Second)
			continue
		}
		currentBlockNumber, err := client.GetCurrentBlock()
		transactions, err2 := client.GetBlockByNumber(lastBlockNumber)
		if err != nil {
			log.Printf("Error getting current block: %v", err)
		} else if err2 != nil {
			log.Printf("Error getting block by number %v", err)
		} else {
			dbTransactions := make([]database.Transaction, len(transactions))
			for idx, transaction := range transactions {
				dbTx := database.Transaction{
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
				dbTransactions[idx] = dbTx
			}
			err3 := this.modelParser.UpdateTransactionsByLastBlockNumber(currentBlockNumber, dbTransactions)
			if err3 != nil {
				log.Printf("Error updateTransactionsByLastBlockNumber %v", err)
			} else {
				log.Printf("Debug updateTransactionsByLastBlockNumber lastBlockNumber = %d , currentBlockNumber = %d, len(transactions) = %d", lastBlockNumber, currentBlockNumber, len(dbTransactions))
			}
		}
		timer.Reset(time.Duration(interval) * time.Second)
	}
}

func (this *EthParser) SetLastUpdatedBlockNumber(blockNumber int) bool {
	return this.modelParser.SetLastUpdatedBlockNumber(blockNumber)
}
