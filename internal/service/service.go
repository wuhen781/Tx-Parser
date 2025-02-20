package service

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type EthParser struct {
}

func (this *EthParser) GetCurrentBlock() (int, error) {
	ethclient := NewEthclient("")
	blockNumber, err := ethclient.GetCurrentBlock()
	if err != nil {
		log.Printf("Error getting current block: %v", err)
	}
	return blockNumber, err
}

func (this *EthParser) Subscribe(address string) bool {

}

func (this *EthParser) GetTransactions(address string) []Transaction {

}
