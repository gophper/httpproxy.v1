package main

import (
	"errors"
	"fmt"
	"gostudy/httpproxy/utils"
	"log"
	"net"
	"runtime"
)

func main() {
	setup(l)
}

func setup(l net.Listener) {
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Panic(err)
	}

	for {
		socketC, err := l.Accept()
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

			if err := handleClientRequest(&utils.ClientConn{socketC, make([]byte, 0), make([]byte, 0)}); err != nil {
				log.Println(err)
			}
		}()
	}
}

func handleClientRequest(client *utils.ClientConn) (err error) {

	var (
		cp     *utils.ClientPair
		server *utils.ClientConn
		closed = false
	)

	defer func() {
		if !closed {
			client.Close()
		}
	}()

	cp, err = shakeHands(client)
	if err != nil {
		log.Println("Shake with server error", err)
		return
	}
	server = cp.OutputClient
	defer func() {
		if !closed {
			server.Close()
		}
	}()

	go func() {
		defer server.Close()
		if _, er, ew, err := utils.Copy(server, client, 2); er != nil || ew != nil || err != nil {
			log.Println("server to remote error:", er, ew, err)
		}
		cp.TraceV2("sended", client.Bufr, server.Bufw, "sender logout")

	}()

	if _, er, ew, err := utils.Copy(client, server, 1); er != nil || ew != nil || err != nil {
		cp.TraceV2("server to client error:", er, ew, err)
	}
	cp.TraceV2("received", server.Bufr, client.Bufw, "receiver logout,both close")
	client.Close()

	return
}

//Initiate a handshake request
func shakeHands(client *utils.ClientConn) (cp *utils.ClientPair, err error) {
	var b [1024]byte
	n, err := client.Read(b[:], false)
	if err != nil || n == 0 {
		return nil, errors.New(fmt.Sprintf("Read shake head error:%v,readed:%v", err, n))
	}

	data, err := utils.DecryptAES(b[:n])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Decode data error:%v,raw data:%v", err, data))
	}

	socketS, err := net.Dial("tcp", string(data))
	if err != nil {
		return nil, errors.New(fmt.Sprint("Connect server error:", err))
	}
	fmt.Println("server dial to:", string(data))

	n, err = client.Write([]byte("ok"), false)
	if err != nil || n == 0 {
		return nil, errors.New(fmt.Sprintf("Send ok to client error:%v,sended:%v", err, n))
	}

	return &utils.ClientPair{
		client,
		&utils.ClientConn{socketS, make([]byte, 0), make([]byte, 0)},},
		nil
}
