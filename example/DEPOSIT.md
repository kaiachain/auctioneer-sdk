This document mainly describes how to deposit and withdraw on Auction contract.

# Contract: AuctionDepositVault

AuctionDepositVault contract is a contract where searchers can deposit amount of bid to be used for bidding.

### Deposit

There are two methods to deposit, `deposit()` and `depositFor(address searcher)`.

The first method deposits using the sender's balance and the deposit amount is settled with the sender.
While, the second method deposits with an account whose have enough balance and settled with the parameter address.

The command below shows how to deposit via `cast` tool.

```
cast send --private-key <your private key> <AuctionDepositVault contract address> "deposit()" --rpc-url <EN endpoint> --confirmations 0 --value 100000000000000000000
```

e.g., `cast send --private-key 0x... 0x... "deposit()" --rpc-url <EN endpoint> --confirmations 0 --value 10000000000000000000`

For the RPC URL, you can use any EN node including [public entrypoint node](https://docs.kaia.io/references/public-en/).

To check your deposit balance, you can call the function `depositBalances(address)(uint256)`.

```
cast call <AuctionDepositVault contract address> "depositBalances(address)(uint256)" <your account address> --rpc-url <EN endpoint>
```

e.g., `cast call 0x... "depositBalances(address)(uint256)" 0x... --rpc-url <EN endpoint>

### DepositFor

```
cast send --private-key <your private key> <AuctionDepositVault contract address> "depositFor(address)" <account address to settle deposit> --rpc-url <EN endpoint> --confirmations 0 --value 10000000000000000000
```

e.g., `cast send --private-key 0x... 0x... "depositFor(address)" 0x... --rpc-url <EN endpoint> --confirmations 0 --value 10000000000000000000`

### Withdraw

There are two steps to withdraw your deposit amount.

1. Call `reserveWithdraw()`.
```
cast send --private-key <your private key> <AuctionDepositVault contract address> "reserveWithdraw()" --rpc-url <EN endpoint> --confirmations 0
```

e.g., `cast send --private-key 0x... 0x... "reserveWithdraw()" --rpc-url <EN endpoint> --confirmations 0`

This step reserves withdraw and makes lock time for 60 second. After 60 seconds, the reserved withdraw amount can be transferred through `withdraw()` function.

2. Call `withdraw()`.

You can call `withdraw()` with the following command example.

```
cast send --private-key <your private key> <AuctionDepositVault contract address> "withdraw()" --rpc-url <EN endpoint> --confirmations 0
```
e.g., `cast send --private-key 0x... 0x... "withdraw()" --rpc-url <EN endpoint> --confirmations 0`
