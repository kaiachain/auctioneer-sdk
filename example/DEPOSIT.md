This document mainly describes how to deposit and withdraw on Auction contract.

<!-- TODO: add a contract link on Kairos / Mainnet -->

# Contract: AuctionDepositVault

AuctionDepositVault contract is a contract where searchers can deposit amount of bid to be used for bidding.

### Deposit

There are two methods to deposit, `deposit()` and `depositFor(address searcher)`.

The first method deposits using the sender's balance and the deposit amount is settled with the sender.
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
