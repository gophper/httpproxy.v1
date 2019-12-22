package main

import (
	"errors"
	"fmt"
	"httpproxy.v1/utils"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"
)

const curVersion = "1.0.1"
const curMode = "server"

func main() {
	utils.Panel(Server{})
}

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

//Initiate a handshake request
func shakeHands(client net.Conn) (cp utils.ClientPair, err error) {
	var (
		b         [1024]byte
		connectId int64
	)
	n, err := client.Read(b[:])
	if err != nil || n == 0 {
		return cp, errors.New(fmt.Sprintf("read shake head error:%v,readed:%v", err, n))
	}

	data, err := utils.DecryptAES(b[:n])
	if err != nil {
		return cp, errors.New(fmt.Sprintf("decode data error:%v,raw data:%v", err, data))
	}

	sd := strings.Split("#", string(data))
	if len(sd) != 2 {
		return cp, errors.New(fmt.Sprintf("shake msg content error:%v", string(data)))
	}

	remote, err := net.Dial("tcp", sd[0])
	if err != nil {
		return cp, errors.New(fmt.Sprint("connect server error:", err))
	}
	fmt.Println("server dial to:", string(data))

	n, err = client.Write([]byte("ok"))
	if err != nil || n == 0 {
		return cp, errors.New(fmt.Sprintf("send ok to client error:%v,sended:%v", err, n))
	}

	if connectId, err = strconv.ParseInt(sd[1], 10, 0); err != nil {
		return cp, errors.New(fmt.Sprintf("send ok to client error:%v", err))
	}

	if ocp := utils.GetPair(connectId); ocp != nil {
		return utils.RebuildPair(connectId, client), nil
	}

	return utils.NewClientPair(client, remote, connectId), nil
}
