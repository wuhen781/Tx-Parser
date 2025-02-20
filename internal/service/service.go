package service

import "github.com/wuhen781/Tx-Parser/pkg/ethclient"
import "github.com/wuhen781/Tx-Parser/internal/database"
import "github.com/wuhen781/Tx-Parser/internal/model"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type EthParser struct {
}

func NewEthParser() *EthParser {
	return &EthParser{}
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
	modelParser := model.NewModelParser()
	return modelParser.AddSubscribe(address, blockNumber)
}

func (this *EthParser) GetTransactions(address string) []database.Transaction {
	modelParser := model.NewModelParser()
	return modelParser.GetTransactions(address)
}
