package service

type Transaction struct {
	From        string
	To          string
	Value       string
	BlockNumber int
	Gas         int
	GasPrice    string
	Hash        string
	Nonce       int
	Timestamp   int64
	Type        string
}
