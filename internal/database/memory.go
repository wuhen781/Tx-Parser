package database

type memoryDb struct {
	transactions           map[string][]Transaction
	subscribers            map[string]*subscriberInfo
	lastUpdatedBlockNumber int
	transOffetsInLastBlock int
}

type subscriberInfo struct {
	fromBlockNumber int
	lastGetOffset   int
}

func NewMemoryDb() *memoryDb {
	return &memoryDb{
		transactions:           make(map[string][]Transaction),
		subscribers:            make(map[string]*subscriberInfo),
		lastUpdatedBlockNumber: -1,
		transOffetsInLastBlock: 0,
	}
}

func (this *memoryDb) AddSubscribe(address string, blockNumber int) bool {
	this.subscribers[address] = &subscriberInfo{
		fromBlockNumber: blockNumber,
		lastGetOffset:   0,
	}
	return true
}

func (this *memoryDb) GetSubscribeFromBlockNumber(blockNumber int) []string {
	anses := make([]string, 0)
	for addr, info := range this.subscribers {
		if info.fromBlockNumber <= blockNumber { //From this point on, observations were made
			anses = append(anses, addr)
		}
	}
	return anses
}

func (this *memoryDb) SetTransactions(transactions []Transaction) bool {
	for _, tx := range transactions {
		from, to := tx.From, tx.To
		if _, ok := this.transactions[from]; !ok {
			this.transactions[from] = make([]Transaction, 0)
		}
		if _, ok := this.transactions[to]; !ok {
			this.transactions[to] = make([]Transaction, 0)
		}
		tx.Type = "outbound"
		this.transactions[from] = append(this.transactions[from], tx)
		tx.Type = "inbound"
		this.transactions[to] = append(this.transactions[to], tx)
	}
	return true
}

func (this *memoryDb) GetTransactions(address string) []Transaction {
	subscriberinfo, ok := this.subscribers[address]
	if !ok {
		return make([]Transaction, 0)
	}
	transactions := this.transactions[address]
	tranLen := len(transactions)
	ans := transactions[subscriberinfo.lastGetOffset:tranLen]
	this.subscribers[address].lastGetOffset = tranLen
	return ans
}

func (this *memoryDb) GetLastUpdatedBlockNumber() int {
	return this.lastUpdatedBlockNumber
}

func (this *memoryDb) SetLastUpdatedBlockNumber(blockNumber int) bool {
	this.lastUpdatedBlockNumber = blockNumber
	return true
}

func (this *memoryDb) GetTransOffetsInLastBlock() int {
	return this.transOffetsInLastBlock
}

func (this *memoryDb) SetTransOffetsInLastBlock(offset int) bool {
	this.transOffetsInLastBlock = offset
	return true
}
