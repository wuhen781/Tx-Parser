package database

type Db interface {
	AddSubscribe(address string, blockNumber int) bool
	GetTransactions(address string) []Transaction
}
