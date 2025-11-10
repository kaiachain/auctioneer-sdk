package sdk

import (
	"bytes"

	"github.com/kaiachain/kaia/accounts/abi"
)

const ENTRYPOINT_ABI_STR = `
[{
  "inputs": [
    {
      "internalType": "address",
      "name": "owner",
      "type": "address"
    }
  ],
  "name": "nonces",
  "outputs": [
    {
      "internalType": "uint256",
      "name": "",
      "type": "uint256"
    }
  ],
  "stateMutability": "view",
  "type": "function"
}]
`

var ENTRYPOINT_CONTRACT_ABI, _ = abi.JSON(bytes.NewReader([]byte(ENTRYPOINT_ABI_STR)))
