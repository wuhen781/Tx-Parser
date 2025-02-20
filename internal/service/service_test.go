package service

import (
	"context"
	"testing"
	"time"
)

func TestEthParser(t *testing.T) {
	parser := NewEthParser()

	// Test GetCurrentBlock
	blockNumber, err := parser.GetCurrentBlock()
	if err != nil {
		t.Errorf("Error getting current block: %v", err)
	}
	if blockNumber < 0 {
		t.Errorf("Expected valid block number, got %d", blockNumber)
	}

	// Test Subscribe
	if !parser.Subscribe("0x123") {
		t.Errorf("Failed to subscribe to address")
	}

	// Test GetTransactions
	txs := parser.GetTransactions("0x123")
	if len(txs) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(txs))
	}
}

func TestUpdateTransactionsInBackGroundRegularly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parser := NewEthParser()
	parser.Subscribe("0x123")

	// Run the update function in a separate goroutine
	go func() {
		parser.UpdateTransactionsInBackGroundRegularly(ctx, 1)
	}()

	// Allow some time for execution
	time.Sleep(10 * time.Second)

	// Stop the background process
	cancel()
}
