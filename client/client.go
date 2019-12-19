package main

import (
	"bytes"
	"errors"
	"fmt"
	"gostudy/httpproxy/utils"
	"log"
	"net"
	"net/url"
	"runtime"
	"strings"
)

var established = []byte("HTTP/1.1 200 Connection established\r\n\r\n")

func main() {
	l, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		log.Println("Received:", client.RemoteAddr())
		go handleClientRequest(&utils.ClientConn{client, make([]byte, 0), make([]byte, 0)})
	}
}

func handleClientRequest(client *utils.ClientConn) {
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
		err    error
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
			cp.Info("---------->", "\n", client.Bufr, "\n", server.Bufw, )
		}()
		if _, er, ew, err := utils.Copy(server, client, 1); err != nil {
			//浏览器客户端已经关闭，如果有往服务端写的数据，必须等写完再关闭所有，但是目前是同步读写，所以不存在这种情况
			if er != nil {
				//nothing to do
			}
			//服务端异常通道关闭,尝试进行重连
			if ew != nil {
				//todo 客户端和服务端通道重连
			}

			_ = server.Close()
			_ = client.Close()

			log.Println("client logout", client.RemoteAddr(), "->", server.RemoteAddr())
			log.Println("client send to server error", er, ew, err)
		}
	}()
	go func() {
		defer func() {
			ch <- struct{}{}
			cp.Info("<----------", "\n", server.Bufr, "\n", client.Bufw)
		}()

		if _, er, ew, err := utils.Copy(client, server, 2); err != nil {

			//服务端通道异常关闭，如果客户端正在往浏览器写数据，则等待写完成再关闭客户端对浏览器的通道。
			// 如果客户端正在读取浏览器的数据，则进行重连（todo），暂时改为关闭所有通道，并抛出醒目错误
			if er != nil {
				server.Close()
				client.Close()
			}
			//客户端通道异常关闭，如果有正在往服务端写数据，则等待写完再关闭；
			// 如果有从服务端读则关闭私有云通道，并抛出醒目错误
			if ew != nil {
				//todo 客户端和服务端通道重连
			}

			log.Println("server logout", client.RemoteAddr(), "->", server.RemoteAddr())
			log.Println("server return to client error", err)
		}
	}()

	<-ch
	cp.Info("both close")
	defer server.Close()
	defer client.Close()
}

func parseHost(data []byte) (address, method string, err error) {

	var (
		host string
	)
	_, err = fmt.Sscanf(string(data[:bytes.IndexByte(data, '\n')]), "%s%s", &method, &host)
	if err != nil {
		log.Println("domain parse error", err, string(data))
		return
	}
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println("url parse error", err)
		return
	}

	if hostPortURL.Opaque == "443" {
		address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}
	return
}

//Initiate a handshake request
func shakeHands(client *utils.ClientConn) (cp *utils.ClientPair, err error) {
	var head [1024]byte
	n, err := client.Read(head[:], false)
	if err != nil || n == 0 {
		return nil, errors.New(fmt.Sprintf("Read http head error:%v,readed:%v", err, n))
	}

	address, method, err := parseHost(head[:n])
	sc, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		return nil, errors.New(fmt.Sprint("Connect server error:", err))
	}

	encryptData, err := utils.EncryptAES([]byte(address), utils.Key)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Encode data error:", err))
	}
	server := &utils.ClientConn{sc, make([]byte, 0), make([]byte, 0)}
	_, err = server.Write(encryptData, false)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Send head to server error:", err))
	}

	var ackBytes [1024]byte
	n, err = server.Read(ackBytes[:], false)
	if err != nil || n == 0 {
		return nil, errors.New(fmt.Sprintf("Read ack error:%v,readed:%v", err, n))
	}

	if method == "CONNECT" {
		n, err = client.Write([]byte(established), false)
		if err != nil || n == 0 {
			return nil, errors.New(fmt.Sprintf("Send established to client erro:%v,readed:%v", err, n))
		}
	}

	if string(ackBytes[0:n]) == "" {
		return nil, errors.New(fmt.Sprint("Receive ack message error:", string(ackBytes[0:n])))
	}

	return &utils.ClientPair{client, server,}, nil

}
