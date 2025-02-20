package database

import "sync"

type memoryDb struct {
	transactions           map[string][]Transaction
	subscribers            map[string]*subscriberInfo
	subsMu                 sync.RWMutex
	transMu                sync.RWMutex
	lastUpdatedBlockNumber int
	transOffetsInLastBlock int
}

type subscriberInfo struct {
	fromBlockNumber int
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
	this.subsMu.Lock()
	defer this.subsMu.Unlock()
	if _, ok := this.subscribers[address]; ok {
		return true
	}
	this.subscribers[address] = &subscriberInfo{
		fromBlockNumber: blockNumber,
	}
	return true
}

func (this *memoryDb) GetSubscribeFromBlockNumber(blockNumber int) []string {
	this.subsMu.Lock()
	defer this.subsMu.Unlock()
	anses := make([]string, 0)
	for addr, info := range this.subscribers {
		if info.fromBlockNumber <= blockNumber { //From this point on, observations were made
			anses = append(anses, addr)
		}
	}
	return anses
}

func (this *memoryDb) SetTransactions(transactions []Transaction) bool {
	this.transMu.Lock()
	defer this.transMu.Unlock()
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
	this.transMu.Lock()
	transactions, ok := this.transactions[address]
	if !ok {
		this.transMu.Unlock()
		return make([]Transaction, 0)
	}
	delete(this.transactions, address)
	this.transMu.Unlock()
	anses := make([]Transaction, len(transactions))
	copy(anses, transactions)
	return anses
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
