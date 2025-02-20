package database

type Db interface {
	AddSubscribe(address string, blockNumber int) bool
	GetSubscribeFromBlockNumber(blockNumber int) []string
	GetTransactions(address string) []Transaction
}
