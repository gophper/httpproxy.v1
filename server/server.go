package server

import (
	"log"
	"net"
	"runtime"
)

const curVersion = "1.0.1"
const curMode = "server"

type Server struct{}

func (c Server) GetCurVersion() string {
	return curVersion
}

func (c Server) GetCurMode() string {
	return curMode
}

func (c Server) Setup(address string) {

	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
	}

	for {
		socket, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					var buf [2 << 10]byte
					stack := string(buf[:runtime.Stack(buf[:], true)])
					log.Println("panic error:", err, "stack:", stack)
				}
			}()

			if err := handleClientRequest(socket); err != nil {
				log.Println(err)
			}
		}()
	}
}
