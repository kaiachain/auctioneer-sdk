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

const (
	// PermissionlessHFBlock is the target block at which v3.0 activates.
	// TODO: replace with params.MainnetChainConfig.PermissionlessCompatibleBlock.Uint64()
	// once the kaia module version includes the permissionless HF block config.
	PermissionlessHFBlock = ^uint64(0)

	auctionTxTypeV3 = "AuctionTx(bytes32 targetTxHash,uint256 blockNumber,address sender,address to,uint256 nonce,uint256 bid,uint256 maxGasPrice,uint256 callGasLimit,bytes data)"
	domainName      = "KAIA_AUCTION"
	domainVersion   = "0.0.1"
)

var (
	auctionTxTypeHashV3 = crypto.Keccak256Hash([]byte(auctionTxTypeV3))
	domainTypeHash      = crypto.Keccak256Hash([]byte(auction_module.EIP712DomainType))
	domainNameHash      = crypto.Keccak256Hash([]byte(domainName))
	domainVersionHash   = crypto.Keccak256Hash([]byte(domainVersion))
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

// getV3AuctionDigest computes the v3.0 EIP-712 digest locally (includes maxGasPrice).
// TODO: replace with auction_module.Bid.GetHashTypedData once kaia core adds v3 typehash support.
func getV3AuctionDigest(auction *AuctionBid, chainId *big.Int, verifyingContract common.Address) []byte {
	maxGasPrice := auction.MaxGasPrice
	if maxGasPrice == nil {
		maxGasPrice = new(big.Int)
	}

	structEncoded := make([]byte, 0, 11*32)
	structEncoded = append(structEncoded, auctionTxTypeHashV3.Bytes()...)
	structEncoded = append(structEncoded, auction.TargetTxHash.Bytes()...)
	structEncoded = append(structEncoded, common.LeftPadBytes(auction.BlockNumber.Bytes(), 32)...)
	structEncoded = append(structEncoded, common.LeftPadBytes(auction.Sender.Bytes(), 32)...)
	structEncoded = append(structEncoded, common.LeftPadBytes(auction.To.Bytes(), 32)...)
	structEncoded = append(structEncoded, common.LeftPadBytes(new(big.Int).SetUint64(auction.Nonce).Bytes(), 32)...)
	structEncoded = append(structEncoded, common.LeftPadBytes(auction.Bid.Bytes(), 32)...)
	structEncoded = append(structEncoded, common.LeftPadBytes(maxGasPrice.Bytes(), 32)...)
	structEncoded = append(structEncoded, common.LeftPadBytes(new(big.Int).SetUint64(auction.CallGasLimit).Bytes(), 32)...)
	structEncoded = append(structEncoded, crypto.Keccak256Hash(auction.Data).Bytes()...)
	structHash := crypto.Keccak256Hash(structEncoded)

	domainEncoded := make([]byte, 0, 5*32)
	domainEncoded = append(domainEncoded, domainTypeHash.Bytes()...)
	domainEncoded = append(domainEncoded, domainNameHash.Bytes()...)
	domainEncoded = append(domainEncoded, domainVersionHash.Bytes()...)
	domainEncoded = append(domainEncoded, common.LeftPadBytes(chainId.Bytes(), 32)...)
	domainEncoded = append(domainEncoded, common.LeftPadBytes(verifyingContract.Bytes(), 32)...)
	domainSeparator := crypto.Keccak256Hash(domainEncoded)

	return crypto.Keccak256([]byte{0x19, 0x01}, domainSeparator.Bytes(), structHash.Bytes())
}

func getAuctionDigest(auction *AuctionBid, chainId *big.Int, verifyingContract common.Address) []byte {
	if auction.BlockNumber.Uint64() < PermissionlessHFBlock {
		bid := getAuctionModuleBidWithoutAuctioneerSig(auction)
		return bid.GetHashTypedData(chainId, verifyingContract)
	}
	return getV3AuctionDigest(auction, chainId, verifyingContract)
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
		MaxGasPrice:  (*hexutil.Big)(auctionBid.MaxGasPrice),
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
