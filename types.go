package sdk

import (
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
)

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

type AuctionBid struct {
	TargetTxRaw  []byte
	TargetTxHash common.Hash
	BlockNumber  *big.Int
	Sender       common.Address
	To           common.Address
	Nonce        uint64
	Bid          *big.Int
	CallGasLimit uint64
	Data         []byte
	SearcherSig  []byte
}

func NewAuctionBid(sender, toAddr common.Address, targetTxHash common.Hash, targetBlkNum,
	bid *big.Int, nonce, callGasLimit uint64, callData, sig, targetTxRaw []byte,
) *AuctionBid {
	return &AuctionBid{
		TargetTxRaw:  targetTxRaw,
		Sender:       sender,
		To:           toAddr,
		TargetTxHash: targetTxHash,
		BlockNumber:  targetBlkNum,
		Bid:          bid,
		Nonce:        nonce,
		CallGasLimit: callGasLimit,
		Data:         callData,
		SearcherSig:  sig,
	}
}
