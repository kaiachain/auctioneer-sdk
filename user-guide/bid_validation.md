| Items | Auctioneer Check | Proposer Check | Contract Check |
|-------|------------------|----------------|----------------|
| Verification of EIP-712 domain fields                                                                                                                         | V | V | V |
| Verification of auctioneer sigantrue                                                                                                                          | - | V | V |
| Verification of searcher signature                                                                                                                            | V | V | V |
| The value of block number of auction bid should be either currentBlockNumber == auctionTx.blockNumber + 1 or currentBlockNumber == auctionTx.blockNumber + 2  | V | V | V |
| callGasLimit <= MaxGasLimit(`10_000_000`)                                                                                                                     | V | - | V |
| pool contains targetTxHash                                                                                                                                    | - | V | - |
| Sender can submit one bidding per auction race                                                                                                                | V | - | - |
| A winner of auction race cannot be a winner for another race in the same block                                                                                | V | - | - |
| Check of deposit balance to pay the bid                                                                                                                       | V | - | V |
| Check bid data size                                                                                                                                           | V | V | - |
| Check of searcher's nonce                                                                                                                                     | V | - | V |
| When multiple bids are submitted for a single target tx                                                                                                       | The auction race is unique and cannot be reopened once it was closed| Replace by highest bid | - |
| When submitting consecutive bids for different target txs                                                                                                     | Allow only one submission for conseuctive two block numbers for different target txs | - | - |


An additional clarification for the last item:
if Searcher A is the winner in block N, they cannot be selected as the winner if winning bids with different targets are found in blocks N-1 or N+1.

For example,
```
    Assume that
    - auctioneer's current block number = N
    - searcher A's bidding is winning bid in N+1
    -> searcher A cannot be a winner in N+2 for all auction race execpt for the same target.
```
