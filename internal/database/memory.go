package database

type memory struct {
	transactions           map[string][]Transaction
	subscribers            map[string]*subscriberInfo
	lastUpdatedBlcokNumber int
	transOffetsInLastBlcok int
}

type subscriberInfo struct {
	fromBlockNumber int
}

type Db interface {
	AddSubscribe(address string, blockNumber int) bool
	GetSubscribeFromBlockNumber(blockNumber int) []string
	GetTransactions(address string) []Transaction
	GetLastUpdatedBlcokNumber() int
	SetLastUpdatedBlcokNumber(blockNumber int) bool
	GetTransOffetsInLastBlock() int
	SetTransOffetsInLastBlock(offset int) bool
}

func (this *memory) GetLastUpdatedBlockNumber() int {
	return this.lastUpdatedBlcokNumber
}

func (this *memory) SetLastUpdatedBlockNumber(blockNumber int) bool {
	this.lastUpdatedBlcokNumber = blockNumber
	return this.lastUpdatedBlcokNumber
}

func (this *memory) GetTransOffetsInLastBlock() int {
	return this.transOffetsInLastBlcok
}

func (this *memory) SetTransOffetsInLastBlock(offset int) bool {
	this.transOffetsInLastBlcok = offset
	return this.transOffetsInLastBlcok
}
