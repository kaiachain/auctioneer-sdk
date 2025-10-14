This document mainly describes the auctioneer API usage and its MEV ecosystem.

### Auctiioneer RESTful APIs
The auctioneer server exposes four APIs.

- `api/v1/auction/config`
    - Searchers can check out configuration of auctineer
    ```shell
    > curl -X GET http://34.64.220.227:8080/api/v1/auction/config
    > {"version":"v0.0.1","auctioneer":"0x83d26f1534084bc37c476878cC0C3604165BCF4e","entrypoint":"0x74567C431d0D72A0d3Be077F884cCB9A4CF059c6","vaultaddr":"0x976281659A1794Db2FFfb3DF4c400aa24E5f7f95","auctionwindow":"200ms","minBid":"1e-11","searcherSubscription"
    :true,"searchers":1001}
    ```

- `api/v1/auction/send`
    - Searchers can send a bid via this API. The simple usage is as follows:
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
    - Reference to input structure: https://github.com/kaiachain/auctioneer/blob/main/api/rest_interface.go#L52

- `api/v1/subscribe/pendingtxs`
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

- `api/v1/subscribe/nonce`
     ```shell
     > curl "http://localhost:8080/api/v1/subscribe/nonce?address=<your address>"
     > <number>
     ```

### Explorer and Dashboard links
- TBU

### KIP-249
- https://kips.kaia.io/KIPs/kip-249
