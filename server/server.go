package main

import (
	"errors"
	"fmt"
	"gostudy/httpproxy/utils"
	"io"
	"log"
	"net"
	"runtime"
)

func main() {
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
			if err := handleClientRequest(&utils.ClientConn{socketC, make([]byte, 0), make([]byte, 0)}); err != nil {
				log.Println(err)
			}
		}()
	}
}

func handleClientRequest(client *utils.ClientConn) (err error) {
	defer func() {
		if err := recover(); err != nil {
			var buf [2 << 10]byte
			stack := string(buf[:runtime.Stack(buf[:], true)])
			log.Println("panic error:", err, "stack:", stack)
		}
	}()
	var (
		cp     *utils.ClientPair
		server *utils.ClientConn
	)

	cp, err = shakeHands(client)
	if err != nil {
		log.Println("Shake with server error", err)
		return
	}

	server = cp.OutputClient
	ch := make(chan struct{}, 2)

	go func() {
		defer func() {
			ch <- struct{}{}
			cp.InfoV2("---------->", "\n", client.Bufr, "\n", server.Bufw)
		}()
		if _, er, ew, err := utils.Copy(server, client, 2); er != nil || ew != nil || err != nil {
			//浏览器客户端已经关闭，关闭服务端通道
			if er != nil {
				server.Close()
				client.Close()
			}
			//服务端异常通道关闭
			if ew != nil {
				//todo 客户端和服务端通道重连
			}
			log.Println("server to remote error:", er, ew, err)
			log.Println("client logout", client.RemoteAddr(), "->", server.RemoteAddr())
			return

		}
	}()
	go func() {
		defer func() {
			ch <- struct{}{}
			cp.InfoV2("<----------", "\n", server.Bufr, "\n", client.Bufw)
		}()
		if _, er, ew, err := utils.Copy(client, server, 1); er != nil || ew != nil || err != nil {

			//服务端通道异常关闭，如果客户端正在往浏览器写数据，则等待写完成再关闭客户端对浏览器的通道。如果客户端正在读取浏览器的数据，则进行重连（todo），暂时改为关闭所有通道，并抛出醒目错误
			if er != nil {
				server.Close()
				client.Close()
			}
			//客户端通道异常关闭，如果有正在往服务端写数据，则等待写完再关闭，如果有从服务端读则关闭私有云通道，并抛出醒目错误
			if ew != nil {
				//todo 客户端和服务端通道重连
			}


			if err == io.EOF {
				log.Println("server logout", client.RemoteAddr(), "->", server.RemoteAddr())
			} else {
				log.Println("server to client error", err)
			}
		}
	}()

	<-ch

	cp.InfoV2("both close")
	defer server.Close()
	defer client.Close()
	return
}

//Initiate a handshake request
func shakeHands(client *utils.ClientConn) (cp *utils.ClientPair, err error) {
	var b [1024]byte
	n, err := client.Read(b[:], false)
	if err != nil || n == 0 {
		return nil, errors.New(fmt.Sprintf("Read shake head error:%v,readed:%v", err, n))
	}

	data, err := utils.DecryptAES(b[:n], utils.Key)
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

	server := &utils.ClientConn{socketS, make([]byte, 0), make([]byte, 0)}
	return &utils.ClientPair{client, server,}, nil
}
