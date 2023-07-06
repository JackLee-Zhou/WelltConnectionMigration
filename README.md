
### WalletConnect v1.0 protocol (version=1) the parameters are:
    key - symmetric key used for encryption
    bridge - url of the bridge server for relaying messages

### For WalletConnect v2.0 protocol (version=2) the parameters are:
    symKey - symmetric key used for encrypting messages over relay
    methods - jsonrpc methods supported for pairing topic
    relay-protocol - transport protocol for relaying messages
    relay-data - (optional) transport data for relaying messages

### Example of a WalletConnect URI for v1.0 And v2.0 protocol:
### 1.0

wc:8a5e5bdc-a0e4-4702-ba63-8f1a5655744f@1?bridge=https%3A%2F%2Fbridge.walletconnect.org&key=41791102999c339c844880b23950704cc43aa840f3739e365323cda4dfa89e7a

### 2.0
wc:7f6e504bfad60b485450578e05678ed3e8e8c4751d3c6160be17160d63ec90f9@2?relay-protocol=irn&symKey=587d5484ce2a2a6ee3ba1962fdd7e8588e06200c46823bd18fbd67def96ad303&methods=[wc_sessionPropose],[wc_authRequest,wc_authBatchRequest]"


wc:259c4fff2200e0c29185c5b252f1244dfbb8ec516e9d1ad5f98f0e12f7cb1765@2?relay-protocol=irn&symKey=a8c1a8fcd514c3c007e6b551ab4c078d9286f6e1e8e9fe721aa9e740ffe84eb7

wss://relay.walletconnect.com/?**auth**=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJkaWQ6a2V5Ono2TWtmNlBwaVdrRjFWeWlSRFdIb0dSYlhub0s0ZEFBb1VVUGtGVE5SNFl4alVuZSIsInN1YiI6ImIzNGI1MWE4OGQ0YWE3NTY4Nzg2OGI4Y2JmNzViMGM4MGQ4YzYyZTRmMzMyY2VmNmMyY2IxYTg1YWZjYzY1MjgiLCJhdWQiOiJ3c3M6Ly9yZWxheS53YWxsZXRjb25uZWN0LmNvbSIsImlhdCI6MTY4ODQ1MzU1OCwiZXhwIjoxNjg4NTM5OTU4fQ.-E1ZQSecKZwNwB97YnRRvk9joXhBJUOg8PjJkN52UKR8AL0n1_C75wuFakGGO5dvq_bTZqV_tonMh5qCuXN6Ag&**projectId**=f19f5c8e2b1ea7fbd382583761c167b3&**ua=wc-2**%2Fjs-2.8.6%2Fwindows10-chrome-114.0.0%2Fbrowser%3Areact-wallet.walletconnect.com&useOnCloseEvent=true

> 这些数据是 需要 JWT 验证的，看前端发送的原始数据结构是什么，publicKey 怎么发过来的
> eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJkaWQ6a2V5Ono2TWtmNlBwaVdrRjFWeWlSRFdIb0dSYlhub0s0ZEFBb1VVUGtGVE5SNFl4alVuZSIsInN1YiI6ImIzNGI1MWE4OGQ0YWE3NTY4Nzg2OGI4Y2JmNzViMGM4MGQ4YzYyZTRmMzMyY2VmNmMyY2IxYTg1YWZjYzY1MjgiLCJhdWQiOiJ3c3M6Ly9yZWxheS53YWxsZXRjb25uZWN0LmNvbSIsImlhdCI6MTY4ODQ1MzU1OCwiZXhwIjoxNjg4NTM5OTU4fQ.-E1ZQSecKZwNwB97YnRRvk9joXhBJUOg8PjJkN52UKR8AL0n1_C75wuFakGGO5dvq_bTZqV_tonMh5qCuXN6Ag
> 
> **header.ployload.signature**


ueyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9
eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9
