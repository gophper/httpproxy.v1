package main

import (
	"httpproxy.v1/server"
	"httpproxy.v1/utils"
)

func main() {
	utils.Panel(server.Server{})
}
