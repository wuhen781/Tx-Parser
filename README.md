# Ethereum Blockchain Parser

This project implements an Ethereum blockchain parser in Go. The parser interacts with the Ethereum blockchain using the JSON-RPC API and provides functionality to:
- Get the current block number.
- Subscribe to Ethereum addresses.
- Get transactions for a subscribed address.

### Prerequisites

- Go 1.23.4 or later
- Access to the Ethereum JSON-RPC endpoint (e.g., https://ethereum-rpc.publicnode.com)

### Installation

1. Clone the repository:
```sh
git clone https://github.com/wuhen781/Tx-Parser
cd Tx-Parser

2. Build main.go:
go build -o ethparser_api app/api/main.go

3. Run ethparser_api
# start the api server
./ethparser_api

# To get the current block number
curl localhost:8081/currentBlock

# To subscribe to an address
curl -X POST -d "address=0xdac17f958d2ee523a2206206994597c13d831ec7" http://localhost:8081/subscribe

# To get transactions for an address
curl localhost:8081/transactions?address=0xyouraddress
