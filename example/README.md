This document mainly describes how to deposit and witdraw on Auction contract

<!-- TODO: add a contract link on Kairos / Mainnet -->
# Contract: AuctionDepositVault

AuctionDepositVault contract is a contract where searchers can deposit amount of bid to be used for bidding.

### Deposit
There are two moethods to deposit, `deposit()` and `depositFor(address searcher)`.

The first method deposits using the sender's balance and the deposit amount is setteld with the sender.
While, the second method deposits with an account whose have enough balance and settled with the parameter address.


The command below shows how to deposit via `cast` tool.
<!-- TODO: Change AuctionDepositVault address -->
<!-- TODO: Change entrypoint url -->
<!-- cast send --private-key 0x03b7971e86b6cb750875fa8846d2ae0deb58e80f1801da97a6d79f1c2463f7c3 0x7C9bf955F7534eB2BAfD4d213DD62FA067922C58 "deposit()" --rpc-url "http://35.216.106.245:8551" --confirmations 0 --value 10000000000000000000 -->
```
cast send --private-key <your private key> <AuctionDepositVault contract address> "deposit()" --rpc-url "http://35.216.106.245:8551" --confirmations 0 --value 10000000000000000000
```
e.g., `cast send --private-key 0x... 0x... "deposit()" --rpc-url "http://35.216.106.245:8551" --confirmations 0 --value 10000000000000000000`

For the RPC URL, you can use any EN node including [public entrypoint node](https://docs.kaia.io/references/public-en/).

To check your deposit balance, you can call the function `depositBalances(address)(uint256)`.
<!-- TODO: Change AuctionDepositVault address -->
<!-- TODO: Change entrypoint url -->
```
cast call <AuctionDepositVault contract address> "depositBalances(address)(uint256)" <your account address> --rpc-url "http://35.216.106.245:8551"
```
e.g., `cast call 0x... "depositBalances(address)(uint256)" 0x... --rpc-url "http://35.216.106.245:8551"`

### DepositFor

<!-- TODO: Change AuctionDepositVault address -->
<!-- TODO: Change entrypoint url -->
<!-- cast send --private-key 0x03b7971e86b6cb750875fa8846d2ae0deb58e80f1801da97a6d79f1c2463f7c3 0x7C9bf955F7534eB2BAfD4d213DD62FA067922C58 "depositFor(address)" 0x583254e7c638c0f685da164bbe043316a7693077 --rpc-url "http://35.216.106.245:8551" --confirmations 0 --value 10000000000000000000 -->
```
cast send --private-key <your private key> <AuctionDepositVault contract address> "depositFor(address)" <account address to settle deposit> --rpc-url "http://35.216.106.245:8551" --confirmations 0 --value 10000000000000000000
```
e.g., `cast send --private-key 0x... 0x... "depositFor(address)" 0x... --rpc-url "http://35.216.106.245:8551" --confirmations 0 --value 10000000000000000000`

### Withdraw
There are two steps to withdraw your deposit amount.
1. Call `reserveWithdraw()`.
<!-- cast send --private-key 0x03b7971e86b6cb750875fa8846d2ae0deb58e80f1801da97a6d79f1c2463f7c3 0x7C9bf955F7534eB2BAfD4d213DD62FA067922C58 "reserveWithdraw()" --rpc-url "http://35.216.106.245:8551" --confirmations 0 -->
```
cast send --private-key <your private key> <AuctionDepositVault contract address> "reserveWithdraw()" --rpc-url "http://35.216.106.245:8551" --confirmations 0
```
e.g., `cast send --private-key 0x... 0x... "reserveWithdraw()" --rpc-url "http://35.216.106.245:8551" --confirmations 0`

This step reserves withdraw and makes lock time for 60 second. After 60 seconds, the reserved withdraw amount can be transferred through `withdraw()` function.


2. Call `withdraw()`.

You can call `withdraw()` with the following command example.

```
cast send --private-key <your private key> <AuctionDepositVault contract address> "withdraw()" --rpc-url "http://35.216.106.245:8551" --confirmations 0
```
e.g., `cast send --private-key 0x... 0x... "withdraw()" --rpc-url "http://35.216.106.245:8551" --confirmations 0`

<!-- cast send --private-key 0x03b7971e86b6cb750875fa8846d2ae0deb58e80f1801da97a6d79f1c2463f7c3 0x7C9bf955F7534eB2BAfD4d213DD62FA067922C58 "withdraw()" --rpc-url "http://35.216.106.245:8551" --confirmations 0 -->

# Auctioneer

This section mainly describes the auctioneer API usage.

This document mainly describes the auctioneer API usage and its MEV ecosystem.

### Auctiioneer RESTful APIs

The auctioneer server exposes four APIs.

`api/v1/auction/config`
- Searchers can check out configuration of auctineer

```shell
> curl -X GET http://34.64.220.227:8080/api/v1/auction/config
> {"version":"v0.0.1","auctioneer":"0x83d26f1534084bc37c476878cC0C3604165BCF4e","entrypoint":"0x74567C431d0D72A0d3Be077F884cCB9A4CF059c6","vaultaddr":"0x976281659A1794Db2FFfb3DF4c400aa24E5f7f95","auctionwindow":"200ms"}
```


`api/v1/auction/send`
- Searchers can send their bid via this API. The simple usage is as follows:

```shell
> curl -X POST http://<auctioneer server url>:8080/api/v1/auction/send \
    -H "Content-Type: application/json" \
    -d '{
        "targetTxRaw": "0xf8674785066720b30083015f909496bd8e216c0d894c0486341288bf486d5686c5b601808207f4a0a97fa83b989a6d66acc942d1cbd70f548c21e24eefea12e72f8c27ba4369a434a01900811315ba3c64055e9778470f438128b54a46712cc032f25a1487
        "sender": "0x96Bd8E216c0D894C0486341288Bf486d5686C5b6",
        "toAddr": "0x96Bd8E216c0D894C0486341288Bf486d5686C5b6",
        "targetTxHash": "0xacb81e7c775471be3e286a461701436f74b7bf7b951096f979b8557d870f246e",
        "targetBlkNum": 1,
        "bid": "0x1",
        "nonce": 1,
        "callGasLimit": 1,
        "callData": "0x1",
        "searcherSig": "0x9f92d0f25af58f402d39ecebea9d06d24d665d5ccea0bc61188596106e8944386e0204ed795701a30d9552dca5ba4ba054ac3aeeebc4570712be89aee08798be01"
    }'

If validation fails, return the corresponding error message:
> {"code":412,"message":"Invalid signature"}

If validation succeed, return a true
> {"success":true}
```
- Reference to input structure: https://github.com/kaiachain/auctioneer/blob/main/api/rest_interface.go#L51

`api/v1/subscribe/pendingtxs`

```shell
> websocat "ws://<auctioneer endpoint url>:8080/api/v1/subscribe/pendingtxs?sig=0x7ba88835c8124cc7799574a96cdfd1043664848d778906f6743c31876f6b4b9f60d2ea5f85b0ebd81914a36fc50f823b2254c5522dbb81d0e4e1082d635007c401&nonce=<searhcer's nonce>"

If the signature is invalid, corresponding error message will return
> {"code":412,"message":"invalid signature length"}

If already connected,
> {"code":412,"message":"Already connected"}

If nonce differs,
> {"code: 412, "message":"searcher's nonce differ. expected nonce = 3, given nonce = 2"}

If successful, pending transaction is continuously relayed.
```

`api/v1/subscribe/nonce`
```shell
> curl "http://localhost:8080/api/v1/subscribe/nonce?address=<your address>"
> <number>
```

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
Additionally, nonce value is reqruied to prevent replay request from thrid party, in just case of where the signed message is exposed accidentally.
Refer to the relevant part of the code for instructions on how to generate this signature.
