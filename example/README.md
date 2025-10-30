This document introduces the example usage of auctioneer SDK code for submitting a bid and subscribing to pending transactions.

# Contract

Before interacting with the auctioneer, it's necessary to deposit the minimum required amount to the AuctionDepositVault contract. Also it's strongly recommended to understand the overall auction process through the [KIP-249](https://kips.kaia.io/KIPs/kip-249). For better understanding, it introduces the detailed guide for searcher to interact with the AuctionEntryPoint and AuctionDepositVault contracts.

### Contract: AuctionEntryPoint

Please refer to the [ENTRYPOINT.md](./ENTRYPOINT.md) for the details of the AuctionEntryPoint contract.

### Contract: AuctionDepositVault

Please refer to the [DEPOSIT.md](./DEPOSIT.md) for the details of the AuctionDepositVault contract.

# SDK Usage

The auctioneer SDK provides a couple of functions to interact with auctioneer endpoint.
The provided functions helps for:

- Bid submission
- Subscribe pending transactions

### Code: Bid submission

[SDK Example: Submitting a Bid](./submitbid.go)

This example demonstrates how to submit a bid using the Auctioneer RESTful API.
It begins by generating a target transaction (in this case, a mock arbitrage transaction) and then constructs a bid that targets a specific slot for that transaction.
In real-world scenarios, such target transactions are typically created by the ecosystem, but here they are simulated for demonstration purposes.
Refer to the `genBid()` function to understand the structure and contents of a bid.
Signing the bid using the searcher's EIP-712 signature is a critical step in the process. Once the bid is fully prepared and signed, it can be submitted to the Auctioneer's RESTful API endpoint.

### Code: Subscribe pending transaction

[SDK Example: Subscribing to Pending Transactions](./subscribe_pendingtx.go)

Searchers who have deposited the minimum required amount can subscribe to live pending transactions via the Auctioneer RESTful API.
This API requires a valid EIP-712 signature from the deposit address to authenticate the subscription request.
Additionally, nonce value is required to prevent replay request from third party, in just case of where the signed message is exposed accidentally.
Refer to the relevant part of the code for instructions on how to generate this signature.
