### Blockchain Loyalty Program
This backend demonstates the ability to buy loyalty points of business. Points are stored on the NEO blockchain so users are able to freely send to other users. Points can be redeemed at the business. Business can launch its own loyalty program easily. 


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

### NEO Smart contract script hash
Testnet: `75b7a2153b7f4bea152a5ecf19f8855c11e74517`

### Run locally
```bash
go run platform/main.go =mode=[local|staging|production]
```

http://localhost:8080/web/



#### Online demo
http://54.168.104.133:8080/web/
