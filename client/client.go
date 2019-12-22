package client

import (
	"log"
	"net"
	"runtime"
)

const curVersion = "1.0.1"
const curMode = "client"

type Client struct{}

func (c Client) GetCurVersion() string {
	return curVersion
}

func (c Client) GetCurMode() string {
	return curMode
}

func (c Client) Setup(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
	}

	for {
		socket, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		log.Println("Received:", socket.RemoteAddr())
		go func() {
			defer func() {
				if err := recover(); err != nil {
					var buf [2 << 10]byte
					stack := string(buf[:runtime.Stack(buf[:], true)])
					log.Println("panic error:", err, "stack:", stack)
				}
			}()

			handleClientRequest(socket)
		}()

	}
}
