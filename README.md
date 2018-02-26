### Draft
##### for debugging
$GOPATH/bin/go-bindata -debug -pkg=helper -o platform/helper/templates.go templates/...

$GOPATH/bin/go-bindata -pkg=helper -o platform/helper/templates.go templates/...
$GOPATH/bin/go-bindata -pkg=config -o=platform/config/json.go config/...


##### build for linux
env GOOS=linux GOARCH=amd64 go build -o o3-loyaltysystem platform/main.go