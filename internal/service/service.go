package service

import "github.com/wuhen781/Tx-Parser/pkg/ethclient"

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type EthParser struct {
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

}

func (this *EthParser) GetTransactions(address string) []Transaction {

}
