package main

import (
	"httpproxy.v1/client"
	"httpproxy.v1/utils"
)

func main() {
	utils.Panel(client.Client{})
}
