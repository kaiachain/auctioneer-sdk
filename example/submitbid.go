package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	auction_sdk "github.com/kaiachain/auctioneer-sdk"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
)

// Target contract example code
// // SPDX-License-Identifier: LGPL-3.0-only
// pragma solidity ^0.8.18;

// contract C {
//     uint256 counter;

//     function inc() public {
//         counter += 1;
//     }

//     function get() public view returns (uint256) {
//         return counter;
//     }
// }

const TARGET_CONTRACT_ABI_STR = `
[{
  "inputs": [],
  "name": "inc",
  "outputs": [],
  "stateMutability": "nonpayable",
  "type": "function"
}]
`

type SendBid struct {
	Sender       common.Address `json:"sender"`
	ToAddr       common.Address `json:"toAddr"`
	TargetTxHash common.Hash    `json:"targetTxHash"`
	TargetTxRaw  hexutil.Bytes  `json:"targetTxRaw"`
	TargetBlkNum uint64         `json:"targetBlkNum"`
	Bid          hexutil.Big    `json:"bid"`
	Nonce        uint64         `json:"nonce"`
	CallGasLimit uint64         `json:"callGasLimit"`
	CallData     hexutil.Bytes  `json:"callData"`
	SearcherSig  hexutil.Bytes  `json:"searcherSig"`
}

var (
	TARGET_CONTRACT_ABI, _ = abi.JSON(bytes.NewReader([]byte(TARGET_CONTRACT_ABI_STR)))

	chainId   = big.NewInt(8127)
	GAS_LIMIT = uint64(10000000)
	signer    = types.LatestSignerForChainID(chainId)
	// TODO: Replace this with your actual key
	searcherKey, _ = crypto.HexToECDSA("e059d5ced4fe8b0420d1c9761842c4806c2cfa555448b238c9b5a3c8ff546730")
	// TODO: Not required. See the description of `genTx()`
	userKey, _ = crypto.HexToECDSA("199b8876d8091e0cbc251c18301dc85691f669ed0d4963df3018a3b2e6c3b461")
	searcher   = crypto.PubkeyToAddress(searcherKey.PublicKey)
	user       = crypto.PubkeyToAddress(userKey.PublicKey)

	entrypoint = common.HexToAddress("0xC259f758eD00Dcf743F28dB8193154Aa9B3862de")
	// TODO: Replace this with your target contract
	targetContract = common.HexToAddress("0x75B5608722ca06eE159Cc9850CEF470fe100105B")
)

func main() {
	var (
		headCh = make(chan *types.Header)
		// TODO: Replace this with your entrypoint url
		c, _ = client.Dial("ws://35.216.106.245:8552")
	)
	c.SubscribeNewHead(context.Background(), headCh)

	for {
		select {
		case header := <-headCh:
			var (
				targetBlkNum = big.NewInt(header.Number.Int64() + 1)
				bid          = genBid(c, targetBlkNum)
				jsonBid, _   = json.Marshal(bid)
			)
			fmt.Println("target block number:", targetBlkNum)
			// TODO: Replace the auctioneer URL
			resp, err := http.Post("http://localhost:8080/api/v1/auction/send", "application/json", bytes.NewBuffer(jsonBid))
			if err != nil {
				panic(err)
			}
			if resp.StatusCode != http.StatusOK {
				panic(fmt.Sprintf("status code = %d error msg = %s", resp.StatusCode, mustParseErrorMsg(resp.Body)))
			}
			return
		}
	}
}

func genBid(c *client.Client, headerNum *big.Int) auction_sdk.SendBid {
	tx := genTx(c)
	nonce, err := auction_sdk.GetEntrypointNonce(c, entrypoint, searcher)
	if err != nil {
		panic(err)
	}
	fmt.Println("txHash:", tx.Hash().String())
	fmt.Printf("entrypoint nonce(%s): %d\n", searcher.String(), nonce)
	auctionBid := auction_sdk.NewAuctionBid(
		searcher,
		targetContract,
		tx.Hash(),
		headerNum,
		mustDecodeStrToKaia("0.02"),
		nonce.Uint64(),
		GAS_LIMIT,
		genContractCall(),
		nil,
		toRlp(tx),
	)
	sendBid, err := auction_sdk.AuctionBidToSendBid(chainId, entrypoint, auctionBid, searcherKey)
	if err != nil {
		panic(err)
	}
	return *sendBid
}

// TODO: This is example calldata. Replace it with your desired one
func genContractCall() []byte {
	data, _ := TARGET_CONTRACT_ABI.Pack("inc")
	return data
}

func toRlp(tx *types.Transaction) []byte {
	raw, err := tx.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return raw
}

func getGasPriceAndNonce(c *client.Client, addr common.Address) (*big.Int, uint64) {
	var (
		gasPrice, _ = c.SuggestGasPrice(context.Background())
		nonce, _    = c.PendingNonceAt(context.Background(), addr)
	)
	fmt.Println("account nonce:", nonce)
	return gasPrice, nonce
}

// TODO: This is unnecessary transaction in real-word scenario because it will be replaced with actual arbitrage transaction
func genTx(c *client.Client) *types.Transaction {
	var (
		gasPrice, nonce = getGasPriceAndNonce(c, user)
		tx, _           = types.SignTx(types.NewTransaction(nonce, user, common.Big1, GAS_LIMIT, gasPrice, nil), signer, userKey)
	)
	return tx
}

func mustDecodeStrToKaia(v string) *big.Int {
	kaia := new(big.Int).SetUint64(params.KAIA)
	factor, ok := new(big.Rat).SetString(v)
	if !ok {
		panic("String value is incorrect")
	}
	result := new(big.Int).Mul(kaia, factor.Num())
	result.Div(result, factor.Denom())
	return result
}

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
