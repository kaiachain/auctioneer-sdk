package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"time"

	auction_sdk "github.com/kaiachain/auctioneer-sdk"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
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

var (
	TARGET_CONTRACT_ABI, _ = abi.JSON(bytes.NewReader([]byte(TARGET_CONTRACT_ABI_STR)))

	chainId   = big.NewInt(8127)
	GAS_LIMIT = uint64(10000000)
	signer    = types.LatestSignerForChainID(chainId)
	// TODO: Replace this with your actual key
	searcherKey, _ = crypto.HexToECDSA("e059d5ced4fe8b0420d1c9761842c4806c2cfa555448b238c9b5a3c8ff546730")
	// TODO: Not required in real world scenario. See the description of `genTx()`
	userKey, _ = crypto.HexToECDSA("199b8876d8091e0cbc251c18301dc85691f669ed0d4963df3018a3b2e6c3b461")
	searcher   = crypto.PubkeyToAddress(searcherKey.PublicKey)
	user       = crypto.PubkeyToAddress(userKey.PublicKey)

	entrypoint = common.HexToAddress("0x0ac2872Ed033e55897c2595f33cd5C7FA8D24878")
	// TODO: Replace this with Kairos target contract
	targetContract = common.HexToAddress("0x75B5608722ca06eE159Cc9850CEF470fe100105B")

	AUCTIONEER_HOST = "kaia-auctioneer-qa.in.kaia.io"
)

func main() {
	ips, err := net.LookupIP(AUCTIONEER_HOST)
	if err != nil {
		panic(err)
	}
	var (
		headCh = make(chan *types.Header)
		// TODO: Replace this with Kairos entrypoint url
		c, _   = client.Dial("ws://35.216.106.245:8552")
		client = &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					if addr == fmt.Sprintf("%s:443", AUCTIONEER_HOST) {
						addr = fmt.Sprintf("%s:443", ips[0].String())
					}
					dialer := &net.Dialer{}
					return dialer.DialContext(ctx, network, addr)
				},
			},
		}
	)
	go pingLoop(client)
	// wait for a first TLs handshake
	time.Sleep(time.Second * 3)

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

			// TODO: Replace with Kairos auctioneer URL
			req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/v1/auction/send", AUCTIONEER_HOST), bytes.NewBuffer(jsonBid))

			{
				fmt.Println(string(jsonBid))
			}
			if err != nil {
				panic(err)
			}
			req.Host = AUCTIONEER_HOST
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				panic(fmt.Sprintf("status code = %d error msg = %s", resp.StatusCode, mustParseErrorMsg(resp.Body)))
			}
			return
		}
	}
}

// this ping loop sends a ping request to keep the HTTPS connection to avoid TLS handshake overhead when a bid is submitted
func pingLoop(c *http.Client) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/v1/ping", AUCTIONEER_HOST), nil)
	if err != nil {
		panic(err)
	}
	for {
		resp, err := c.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK {
			panic("resp status code is not 200")
		}
		time.Sleep(time.Second * 100)
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
		mustDecodeStrToKaia("1.0001"),
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

// NOTE: This is example calldata. Replace it with your desired one
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

// NOTE: This is unnecessary transaction in real-word scenario because it will be replaced with actual arbitrage transaction
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
