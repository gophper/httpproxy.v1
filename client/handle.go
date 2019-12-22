package client

import (
	"httpproxy.v1/utils"
	"log"
	"net"
)

func handleClientRequest(socket net.Conn) {

	var (
		cp     utils.ClientPair
		cc     *utils.ClientConn
		sc     *utils.ClientConn
		err    error
		closed = false
	)
	defer func() {
		if !closed {
			socket.Close()
		}
	}()

	cp, err = ShakeHands(socket)
	if err != nil {
		log.Println("Shake with server error:", err)
		return
	}

	cc = cp.InputClient
	sc = cp.OutputClient
	defer func() {
		if !closed {
			sc.Close()
		}
	}()

	go func() {
		defer sc.Close()
		if er, ew, err := Client2Server(sc, cc); err != nil {
			cp.TraceV1("client send to server error", er, ew, err)
		}
		cp.TraceV1("sended", cc.Bufr, sc.Bufw, "client logout")
	}()

	if er, ew, err := Server2Client(cc, sc); err != nil {
		cp.TraceV1("server return to client error:", er, ew, err)
	}
	cp.TraceV1("received", sc.Bufr, cc.Bufw, "server logout,both close")
	cc.Close()
	closed = true
}
