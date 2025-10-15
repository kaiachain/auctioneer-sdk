This document mainly describes the AuctionEntryPoint contract and its usage.

# Contract: AuctionEntryPoint

AuctionEntryPoint contract is a contract where the bid transaction is executed. The validator will execute the auction transaction requested by searchers. Please find the execution sequence of BidTx through the AuctionEntryPoint contract in [KIP-249](https://kips.kaia.io/KIPs/kip-249).

### Nonce

The AuctionEntryPoint contract manages the monotonic-increasing nonce for each searcher to prevent replay attack. You can check the current nonce of a searcher by calling the `nonces(address searcher)` function. Please note that the nonce will be increased even if searcher's requested call is failed.

```
cast call <AuctionEntryPoint contract address> "nonces(address)(uint256)" <searcher address> --rpc-url <rpc url>
```

### Access Control

When executing the auction transaction, the AuctionEntryPoint contract will call to `bid.to` with `bid.data` as a calldata. It means that the `msg.sender` from the perspective of `bid.to` is the AuctionEntryPoint, not the searcher. It will make a security issue if you simply grants the AuctionEntryPoint contract to call the `bid.to` since any other searchers can put the `bid.to` as their own contract without any permission check.

To prevent this issue, it's highly recommended to introduce the appropriate access control to the `bid.to` contract. For example, you can include the signature of searcher in `bid.data` as a permission check.
