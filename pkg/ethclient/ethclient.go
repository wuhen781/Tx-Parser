package ethclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const endpoint = "https://ethereum-rpc.publicnode.com"

type Ethclient struct {
}

func NewEthclient() *Ethclient {
	return &Ethclient{}
}

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
}

func (this *Ethclient) GetCurrentBlock() (int, error) {
	response, err := this.callRPC("eth_blockNumber", []interface{}{})
	if err != nil {
		return -1, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		err = fmt.Errorf("fail to unmarshal current block response: %w", err)
		return -1, err
	}

	blockNumberHex, ok := result["result"].(string)
	if !ok {
		err = fmt.Errorf("Type assertion failed: result['result'] is not a string")
		return -1, err
	}

	if !strings.HasPrefix(blockNumberHex, "0x") {
		return -1, fmt.Errorf("result is not a valid quantitity")
	}

	// the first 2 characters are "0x" , we get the rest of number
	blockNumber, errParse := strconv.ParseInt(resBody.Result[2:], 16, 64)
	if errParse != nil {
		return 0, fmt.Errorf("failed to strconv.ParseInt block number: %v", errParse)
	}

	return int(blockNumber), nil
}

func (p *EthereumParser) GetBlockByNumber(blcokNumber int) ([]Transaction, error) {
	var transactions []Transaction

	response, err := p.callRPC("eth_getBlockByNumber", []interface{}{intToHex(blockNumber), true})
	if err != nil {
		err = fmt.Errorf("fail to get block by number: %w", err)
		return transactions, err
	}

	var blockResult map[string]interface{}
	if err := json.Unmarshal(response, &blockResult); err != nil {
		err = fmt.Errorf("fail to unmarshal block response: %w", err)
		return transactions, err
	}

	block := blockResult["result"].(map[string]interface{})
	txs := block["transactions"].([]interface{})
	for _, tx := range txs {
		txMap := tx.(map[string]interface{})
		transactions = append(transactions, Transaction{
			From:        txMap["from"].(string),
			To:          txMap["to"].(string),
			Value:       txMap["value"].(string),
			BlockNumber: HexToInt(block["number"].(string)),
			Gas:         HexToInt(txMap["gas"].(string)),
			GasPrice:    txMap["gasPrice"].(string),
			Hash:        txMap["hash"].(string),
			Nonce:       HexToInt(txMap["nonce"].(string)),
			Timestamp:   int64(HexToInt(block["timestamp"].(string))),
		})
	}

	return transactions, nil
}

func (p *EthereumParser) callRPC(method string, params []interface{}) ([]byte, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(endpoint, "application/json", strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func intToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}
