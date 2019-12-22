set GOOS=linux
set GOPACH=amd64
go build -tags="jsoniter" -o deploy/server server/server.go
go build -tags="jsoniter" -o deploy/client client/client.go
