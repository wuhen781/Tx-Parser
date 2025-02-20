package model

import (
	"errors"
	"testing"

	"github.com/wuhen781/Tx-Parser/internal/database"
)

func TestModelParser(t *testing.T) {
	parser := NewModelParser()

	// Test AddSubscribe
	if !parser.AddSubscribe("0x123", 5) {
		t.Errorf("Failed to add subscription")
	}
	if !parser.AddSubscribe("0x456", 5) {
		t.Errorf("Failed to add subscription")
	}

	// Test GetTransactions (Initially Empty)
	txs := parser.GetTransactions("0x123")
	if len(txs) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(txs))
	}

	// Test updateTransactionsByLastBlockNumber with uninitialized block number
	err := parser.updateTransactionsByLastBlockNumber(10, []database.Transaction{})
	if !errors.Is(err, ErrBlockNumberNotInitialed) {
		t.Errorf("Expected ErrBlockNumberNotInitialed, got %v", err)
	}

	// Set last block number and add transactions
	parser.db.SetLastUpdatedBlockNumber(5)
	tx := database.Transaction{From: "0x123", To: "0x456", Value: "50", BlockNumber: 6}
	tx2 := database.Transaction{From: "0x789", To: "0x456", Value: "100", BlockNumber: 6}
	parser.updateTransactionsByLastBlockNumber(6, []database.Transaction{tx, tx2})
	lastBlockNumber := parser.db.GetLastUpdatedBlockNumber()
	if lastBlockNumber != 6 {
		t.Errorf("Expected lastBlockNumber with value 6, got %v", lastBlockNumber)
	}
	transOffset := parser.db.GetTransOffetsInLastBlock()
	if transOffset != 0 {
		t.Errorf("Expected transOffset with value 0, got %v", transOffset)
	}

	// Validate Transactions
	txs = parser.GetTransactions("0x123")
	if len(txs) != 1 || txs[0].Value != "50" {
		t.Errorf("Expected transaction with value 50, got %v", txs)
	}
	txs = parser.GetTransactions("0x123")
	if len(txs) != 0 {
		t.Errorf("Expected transaction 0 , got %v", txs)
	}
	txs = parser.GetTransactions("0x456")
	if len(txs) != 2 {
		t.Errorf("Expected transaction 2 , got %v", txs)
	}
	if txs[0].Type != "inbound" || txs[1].Type != "inbound" {
		t.Errorf("Expected transaction 2 inbound, got %v", txs)
	}

	tx = database.Transaction{From: "0x123", To: "0x456", Value: "50", BlockNumber: 6}
	tx2 = database.Transaction{From: "0x789", To: "0x456", Value: "100", BlockNumber: 6}
	parser.updateTransactionsByLastBlockNumber(6, []database.Transaction{tx, tx2})
	lastBlockNumber = parser.db.GetLastUpdatedBlockNumber()
	if lastBlockNumber != 6 {
		t.Errorf("Expected lastBlockNumber with value 6, got %v", lastBlockNumber)
	}
	transOffset = parser.db.GetTransOffetsInLastBlock()
	if transOffset != 2 {
		t.Errorf("Expected transOffset with value 0, got %v", transOffset)
	}

}
