package sdk

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	auction_module "github.com/kaiachain/kaia/kaiax/auction"
)

func getAuctionModuleBidWithoutAuctioneerSig(auction *AuctionBid) auction_module.Bid {
	bidData := auction_module.BidData{
		TargetTxHash: auction.TargetTxHash,
		BlockNumber:  auction.BlockNumber.Uint64(),
		Sender:       auction.Sender,
		To:           auction.To,
		Nonce:        auction.Nonce,
		Bid:          auction.Bid,
		CallGasLimit: auction.CallGasLimit,
		Data:         auction.Data,
	}
	return auction_module.Bid{BidData: bidData}
}

func getAuctionDigest(auction *AuctionBid, chainId *big.Int, verifyingContract common.Address) []byte {
	bid := getAuctionModuleBidWithoutAuctioneerSig(auction)
	return bid.GetHashTypedData(chainId, verifyingContract)
}

// GetEntrypointNonce retrieves searcher's entrypoint nonce
func GetEntrypointNonce(c *client.Client, entrypoint, searcher common.Address) (*big.Int, error) {
	var (
		data, _ = ENTRYPOINT_CONTRACT_ABI.Pack("nonces", searcher)
		callMsg = kaia.CallMsg{
			To:   &entrypoint,
			Data: data,
		}
	)
	output, err := c.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	var nonce *big.Int
	err = ENTRYPOINT_CONTRACT_ABI.UnpackIntoInterface(&nonce, "nonces", output)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

// AuctionBidToSendBid calculates searcher signature on given auction bid struct and returns it as form of `SendBid` struct
func AuctionBidToSendBid(
	chainId *big.Int,
	entrypoint common.Address,
	auctionBid *AuctionBid,
	searcherKey *ecdsa.PrivateKey,
) (*SendBid, error) {
	auctionTxHash := getAuctionDigest(auctionBid, chainId, entrypoint)
	searcherSig, err := crypto.Sign(auctionTxHash, searcherKey)
	if err != nil {
		return nil, err
	}
	searcherSig[crypto.RecoveryIDOffset] += 27
	auctionBid.SearcherSig = searcherSig
	return &SendBid{
		Sender:       auctionBid.Sender,
		ToAddr:       auctionBid.To,
		TargetTxHash: auctionBid.TargetTxHash,
		TargetTxRaw:  auctionBid.TargetTxRaw,
		TargetBlkNum: auctionBid.BlockNumber.Uint64(),
		Bid:          hexutil.Big(*auctionBid.Bid),
		Nonce:        auctionBid.Nonce,
		CallGasLimit: auctionBid.CallGasLimit,
		CallData:     auctionBid.Data,
		SearcherSig:  auctionBid.SearcherSig,
	}, nil
}

// SubmitBid submits a bid to auctioneer
func SubmitBid(host string, bid SendBid) (*http.Response, error) {
	jsonBid, err := json.Marshal(bid)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/api/v1/auction/send", host)
	return http.Post(url, "application/json", bytes.NewBuffer(jsonBid))
}
