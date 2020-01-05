cd ..
cd ..
set GOPATH=%cd%
cd src
cd httpproxy.v1

go build -o ./deploy/bin/client.exe  client.go
go build -o ./deploy/bin/server.exe  server.go