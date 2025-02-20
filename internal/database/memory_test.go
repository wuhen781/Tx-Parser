package database

import (
	"testing"
)

func TestMemoryDb(t *testing.T) {
	db := NewMemoryDb()

	// Test AddSubscribe and GetSubscribeFromBlockNumber
	db.AddSubscribe("0x123", 5)
	db.AddSubscribe("0x789", 10)
	subs := db.GetSubscribeFromBlockNumber(6)
	if len(subs) != 1 || subs[0] != "0x123" {
		t.Errorf("Expected [0x123], got %v", subs)
	}

	subs = db.GetSubscribeFromBlockNumber(10)
	if len(subs) != 2 {
		t.Errorf("Expected 2 subscribers, got %d", len(subs))
	}

	// Test SetTransactions and GetTransactions
	tx := Transaction{From: "0x123", To: "0x789", Value: "100", BlockNumber: 5}
	db.SetTransactions([]Transaction{tx})

	txs := db.GetTransactions("0x123")
	if len(txs) != 1 || txs[0].Value != "100" {
		t.Errorf("Expected transaction with value 100, got %v", txs)
	}

	txs = db.GetTransactions("0x123")
	if len(txs) > 0 {
		t.Errorf("Expected empty transaction , got %v", txs)
	}

	tx = Transaction{From: "0x123", To: "0x789", Value: "200", BlockNumber: 5}
	db.SetTransactions([]Transaction{tx})

	txs = db.GetTransactions("0x123")
	if len(txs) != 1 || txs[0].Value != "200" {
		t.Errorf("Expected transaction with value 200, got %v", txs)
	}

	txs = db.GetTransactions("0x789")
	if len(txs) != 2 {
		t.Errorf("Expected transaction 2 , got %v", txs)
	}

	if txs[0].Value != "100" {
		t.Errorf("Expected transaction with value 100, got %v", txs)
	}

	if txs[1].Value != "200" {
		t.Errorf("Expected transaction with value 200, got %v", txs)
	}

	// Test LastUpdatedBlockNumber
	db.SetLastUpdatedBlockNumber(20)
	if db.GetLastUpdatedBlockNumber() != 20 {
		t.Errorf("Expected last updated block number to be 20")
	}

	// Test Transaction Offsets
	db.SetTransOffetsInLastBlock(3)
	if db.GetTransOffetsInLastBlock() != 3 {
		t.Errorf("Expected transaction offsets to be 3")
	}
}
