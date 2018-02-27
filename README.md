### Blockchain Loyalty Program
This backend demonstates the ability to buy loyalty points of business. Points are stored on the NEO blockchain so users are able to freely send to other users. Points can be redeemed at the business. Business can launch its own loyalty program easily. 

There are several components.

#### [Smart Contract](https://github.com/O3Labs/Blockchain-Loyalty-Program)
It's a smart contract written in C# that conforms NEP5 protocol with additional methods that only owner of the smart contract can invoke such as `mintTokensTo`, `burnTokens`, `useTokens`

#### [neo-utils](https://github.com/O3Labs/neo-utils) 
neo-utils is used to invok a smart contract.

#### [NEO-transaction-watcher](https://github.com/O3Labs/neo-transaction-watcher) 
This is used to watch an incoming transaction of a particular smart contract. In this case, It parses data of the transaction and verify with external source such as payment processor. If it's valid then send user an email about the transaction or even push notification to the mobile device.

#### [NEP9](https://github.com/O3Labs/NEP9-go)
This will be used to parse the NEP9 QR code shows on the terminal at the business to let user scans it to invoke the smart contract easily.

---
### Note
This application is just a demo so it skips some of the components needed to make it more resilient such as Messaging Queue and other things. A complete solution will be built soon.

This [article](https://www2.deloitte.com/us/en/pages/financial-services/articles/making-blockchain-real-customer-loyalty-rewards-programs.html) contains good information about why putting customer loyalty rewards programs on the blockchain makes sense. *We believe NEO is a great fit for it.*

----

### NEO Smart contract script hash
Testnet: `75b7a2153b7f4bea152a5ecf19f8855c11e74517`

---

### Run locally

##### for debugging
```
$GOPATH/bin/go-bindata -debug -pkg=helper -o platform/helper/templates.go templates/...
```

##### prerequisite for staging/production
```
$GOPATH/bin/go-bindata -pkg=helper -o platform/helper/templates.go templates/...
$GOPATH/bin/go-bindata -pkg=config -o=platform/config/json.go config/...
```

##### build for linux
```
env GOOS=linux GOARCH=amd64 go build -o o3-loyaltysystem platform/main.go
```

```bash
go run platform/main.go -mode=[local|staging|production]
```

http://localhost:8080/web/


--- 

#### Online demo
http://54.168.104.133:8080/web/

Test card is : `4242 4242 4242 4242` with any expiration date and zip code.  
Try purchase some points with a valid NEO address on TestNet. You will be notified by email about the purchase.

