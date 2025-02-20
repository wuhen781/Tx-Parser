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

func NewMemory() *memory {
	return &memory{
		transactions:           make(map[string][]Transaction),
		subscribers:            make(map[string]*subscriberInfo),
		lastUpdatedBlcokNumber: -1,
		transOffetsInLastBlcok: 0,
	}
}

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

func (this *memory) GetSubscribeFromBlockNumber(blockNumber int) []string {
	anses := make([]string, 0)
	for addr, info := range this.subscribers {
		if info.fromBlockNumber <= blockNumber { //From this point on, observations were made
			anses = append(anses, addr)
		}
	}
	return anses
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
