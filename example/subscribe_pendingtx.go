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
	searcherKey, _ := crypto.HexToECDSA("a38f5bbf491d6e175050cde649f012ceeb766d8b0de976492d15ef0b0e2de1ec")
	// TODO: Replace the url with correct auctioneer endpoint
	url, err := auction_sdk.GetDialUrl("https://kaia-auctioneer-qa.in.kaia.io", searcherKey)
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
	var tx types.Transaction
	if err := json.Unmarshal(msg, &tx); err != nil {
		log.Fatal(err)
	}
	return &tx
}
