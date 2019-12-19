package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"gostudy/httpproxy/utils"
	"httpproxy/config"
	"log"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"
)

var established = []byte("HTTP/1.1 200 Connection established\r\n\r\n")

const curVersion = "1.0.1"

func main() {
	var (
		host, port, curMode string
		timeout             int
		trace, printVer     bool
	)

	flag.BoolVar(&printVer, "v", false, "current version")
	flag.StringVar(&host, "b", "", "iP address for local monitoring")
	flag.StringVar(&port, "p", "", "port address for local listening")
	flag.IntVar(&timeout, "t", 1200, "timeout seconds")
	flag.BoolVar((*bool)(&trace), "d", false, "log input and output")
	flag.Parse()

	if printVer {
		fmt.Println("httpproxy version:", curVersion)
		os.Exit(0)
	}
	if port == "" {
		fmt.Println("port needed !")
		os.Exit(0)
	}

	config.InitConfig(curMode)
	setup(fmt.Sprintf("%s:%s", host, port))
}

func setup(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		log.Println("Received:", client.RemoteAddr())
		go func() {
			defer func() {
				if err := recover(); err != nil {
					var buf [2 << 10]byte
					stack := string(buf[:runtime.Stack(buf[:], true)])
					log.Println("panic error:", err, "stack:", stack)
				}
			}()

			handleClientRequest(&utils.ClientConn{client, make([]byte, 0), make([]byte, 0)})
		}()

	}
}

func handleClientRequest(client *utils.ClientConn) {

	var (
		cp     *utils.ClientPair
		server *utils.ClientConn
		err    error
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
		if _, er, ew, err := utils.Copy(server, client, 1); err != nil {
			cp.TraceV1("client send to server error", er, ew, err)
		}
		cp.TraceV1("sended", client.Bufr, server.Bufw, "client logout")
	}()

	if _, er, ew, err := utils.Copy(client, server, 2); err != nil {
		cp.TraceV1("server return to client error:", er, ew, err)
	}
	cp.TraceV1("received", server.Bufr, client.Bufw, "server logout,both close")
	client.Close()
	closed = true

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
	var (
		head     [1024]byte
		ackBytes [1024]byte
	)
	n, err := client.Read(head[:], false)
	if err != nil || n == 0 {
		return nil, errors.New(fmt.Sprintf("Read http head error:%v,readed:%v", err, n))
	}

	address, method, err := parseHost(head[:n])
	sc, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		return nil, errors.New(fmt.Sprint("Connect server error:", err))
	}

	encryptData, err := utils.EncryptAES([]byte(address))
	if err != nil {
		return nil, errors.New(fmt.Sprint("Encode data error:", err))
	}
	server := &utils.ClientConn{sc, make([]byte, 0), make([]byte, 0)}
	_, err = server.Write(encryptData, false)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Send head to server error:", err))
	}

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
