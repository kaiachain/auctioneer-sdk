This document mainly describes FAQ on Kaia MEV ecosystem.

### 1. Subscription
- Pending transaction subscription is not allowed two more connections per searcher address.
- Once connected, the connection is automatically closed after 24 hours.

### 2. API Latency

The Auctioneer utilizes an L7 Load Balancer with the HTTPS protocol.
The initial handshake consumes time depending on the network state.
To bypass this initial delay when sending subsequent bid APIs, it is strongly recommended to establish a keep-alive connection.
Additionally, to prevent being blocked by the Auctioneer API server, do not send the `ping` API too many times within a short period.
Finally, since the Auctioneer server is running in the GCP KR (Korea) region, you are recommended to host your infrastructure in a geographically close region to minimize latency and reduce geographic delay.

### 3. Biding timing

The timing of your bid submission is highly sensitive to the CN (Consensus Node) mining time.
If the auction starts late (close to the mining time), the bid transaction would be inserted into the block after the next one (block number +2 instead of +1).
This means you should set your target block number to +2.
However, the target block number is inherently sensitive to the CN mining schedule: if you target block +2 but the transaction is inserted at block +1 due to earlier processing, the bid will fail.
Therefore, it is recommended to maximize inclusion probability by sending your bid transaction two times: once with a target block number of +1 and once with a target block number of +2.
