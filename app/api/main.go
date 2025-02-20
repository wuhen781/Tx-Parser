package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/wuhen781/Tx-Parser/internal/service"
)

var parser = service.NewEthParser()

func main() {
	initLogging()
	defer closeLogging()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	//The default is the maximum Block Number when accepting the first subscribe request
	//parser.SetLastUpdatedBlockNumber(21889934)

	server := &http.Server{Addr: ":8081"}

	var wg sync.WaitGroup
	wg.Add(2) // Add two goroutines to wait for

	// Start background worker
	go func() {
		defer wg.Done()
		parser.UpdateTransactionsInBackGroundRegularly(ctx, 3)
	}()

	// Start HTTP server
	go func() {
		defer wg.Done()
		http.HandleFunc("/currentBlock", getCurrentBlockHandler)
		http.HandleFunc("/subscribe", subscribeHandler)
		http.HandleFunc("/transactions", getTransactionsHandler)

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("Received shutdown signal...")

	cancel() // Stop background tasks

	// Gracefully shut down HTTP server
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(ctxShutdown); err != nil {
		fmt.Printf("HTTP server Shutdown: %v\n", err)
	} else {
		fmt.Println("HTTP server stopped")
	}

	wg.Wait() // Wait for both goroutines to finish
	fmt.Println("All background tasks stopped, exiting.")
}

var logFile *os.File

func initLogging() {
	var err error
	logFile, err = os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func closeLogging() {
	if logFile != nil {
		logFile.Close()
	}
}

func getCurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	block, err := parser.GetCurrentBlock()
	if err != nil {
		log.Printf("GetCurrentBlock Unavailable: %v", err)
		http.Error(w, "GetCurrentBlock Unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]int{"currentBlock": block}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm() // Parse form data
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	address := r.FormValue("address")
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}
	subscribed := parser.Subscribe(address)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"subscribed": subscribed}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}
	transactions := parser.GetTransactions(address)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
