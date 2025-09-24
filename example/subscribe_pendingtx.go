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
	searcherKey, _ := crypto.HexToECDSA("e059d5ced4fe8b0420d1c9761842c4806c2cfa555448b238c9b5a3c8ff546730")
	// TODO: Replace the url with correct auctioneer endpoint
	url, err := auction_sdk.GetDialUrl("http://localhost:8080", searcherKey)
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
		log.Printf("Received pending transaction: %s, %s", tx.Hash().String(), tx.Time().String())
	}
}

func convertMsgToTx(msg []byte) *types.Transaction {
	var tx types.Transaction
	if err := json.Unmarshal(msg, &tx); err != nil {
		log.Fatal(err)
	}
	return &tx
}
