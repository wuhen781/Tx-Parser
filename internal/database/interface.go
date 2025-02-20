package database

type Db interface {
	AddSubscribe(address string, blockNumber int) bool
	GetSubscribeFromBlockNumber(blockNumber int) []string
	GetTransactions(address string) []Transaction
	SetTransactions(transactions []Transaction) bool
	GetLastUpdatedBlcokNumber() int
	SetLastUpdatedBlcokNumber(blockNumber int) bool
	GetTransOffetsInLastBlock() int
	SetTransOffetsInLastBlock(offset int) bool
}
