package server

import (
	"httpproxy.v1/utils"
	"log"
	"net"
)

func handleClientRequest(socket net.Conn) (err error) {

	var (
		cp      utils.ClientPair
		remote  *utils.ClientConn
		server  *utils.ClientConn
		closed  = false
		delayed = false
	)

	defer func() {
		if !closed {
			_ = socket.Close()
		}
	}()

	cp, err = shakeHands(socket)
	if err != nil {
		log.Println("Shake with server error", err)
		return
	}
	remote = cp.OutputClient
	server = cp.InputClient
	defer func() {
		if !closed && !delayed {
			_ = remote.Close()
		}
	}()

	go func() {
		if _, er, ew, err := utils.Copy(remote, server, 2); er != nil || ew != nil || err != nil {
			log.Println("server to remote error:", er, ew, err)
			if er != nil && er != utils.EOF {
				utils.DelayUnsetPair(cp) //如果客户端和服务端非EOF断开，则延迟关闭remote，等待新server和remote关联
				delayed = true
				log.Println("server exception:", er)
				return
			}
		}
		remote.Close()
		cp.TraceV2("sended", server.Bufr, remote.Bufw, "sender logout")

	}()

	if _, er, ew, err := utils.Copy(server, remote, 1); er != nil || ew != nil || err != nil {
		cp.TraceV2("server to client error:", er, ew, err)
		if ew != nil && ew != utils.EOF {
			utils.DelayUnsetPair(cp)
			delayed = true
			log.Println("server exception:", er)
		}
	}
	cp.TraceV2("received", remote.Bufr, server.Bufw, "receiver logout,both close")
	server.Close()

	return
}
