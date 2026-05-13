package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
	auction_sdk "github.com/kaiachain/auctioneer-sdk"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/crypto"
)

func mustParseErrorMsg(respBody io.ReadCloser) string {
	body, err := io.ReadAll(respBody)
	if err != nil {
		panic(err)
	}
	var errResp struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &errResp); err != nil {
		panic(err)
	}
	return errResp.Message
}

func main() {
	// TODO: Replace with your private key
	searcherKey, err := crypto.HexToECDSA("00000000000000000000000000000")
	if err != nil {
		fmt.Println("Check searcher's private key")
		panic(err)
	}
	url, err := auction_sdk.GetDialUrl("https://auction-mainnet.kaia.io", searcherKey)
	if err != nil {
		panic(err)
	}

	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		if resp != nil {
			fmt.Println("status:", resp.Status)
			fmt.Println("error message:", mustParseErrorMsg(resp.Body))
		}
		log.Fatal("Connection failed:", err)
	}
	defer conn.Close()
	log.Println("Connected to")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		tx := convertMsgToTx(msg)
		fmt.Printf("Received pending transaction: %s\n", tx.Hash().String())
	}
}

func convertMsgToTx(msg []byte) *types.Transaction {
	// NOTE: Uncomment this code block if you want to check out when a tx has been settled in CN
	// {
	// 	var tx map[string]any
	// 	if err := json.Unmarshal(msg, &tx); err != nil {
	// 		panic(err)
	// 	}
	// 	for k, v := range tx {
	// 		if k == "time" {
	// 			fmt.Println("tx time:", v)
	// 			break
	// 		}
	// 	}
	// }

	var tx types.Transaction
	if err := json.Unmarshal(msg, &tx); err != nil {
		log.Fatal(err)
	}
	return &tx
}
