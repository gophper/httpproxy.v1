cd ..
cd ..
set GOPATH=%cd%
set GOOS=linux
set GOPACH=amd64
cd src
cd httpproxy.v1
go build -tags="jsoniter" -o deploy/bin/server ./server
go build -tags="jsoniter" -o deploy/bin/client ./client