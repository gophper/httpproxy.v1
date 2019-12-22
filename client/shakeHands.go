package client

import (
	"bytes"
	"errors"
	"fmt"
	"httpproxy.v1/config"
	"httpproxy.v1/utils"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	established = []byte("HTTP/1.1 200 Connection established\r\n\r\n")
	incrId      int64
)

//The client handshake with the server to communicate the connection address and port
func ShakeHands(clientConn net.Conn) (cp utils.ClientPair, err error) {
	var (
		head      [1024]byte
		ackBytes  [1024]byte
		connectId = allocConnectId()
	)

	n, err := clientConn.Read(head[:])
	if err != nil || n == 0 {
		return cp, errors.New(fmt.Sprintf("read http head error:%v,readed:%v", err, n))
	}

	address, method, err := parseHost(head[:n])
	sc, err := net.Dial("tcp", config.ServerHost)
	if err != nil {
		return cp, errors.New(fmt.Sprint("connect server error:", err))
	}
	fmt.Println("connected with server:", config.ServerHost)

	encryptData, err := utils.EncryptAES([]byte(address + "#" + strconv.FormatInt(connectId, 10)))
	if err != nil {
		return cp, errors.New(fmt.Sprint("encode data error:", err))
	}

	server := utils.NewClientConn(sc)
	_, err = server.Write(encryptData, false)
	if err != nil {
		return cp, errors.New(fmt.Sprint("send head to server error:", err))
	}

	n, err = server.Read(ackBytes[:], false)
	if err != nil || n == 0 {
		return cp, errors.New(fmt.Sprintf("read ack error:%v,readed:%v", err, n))
	}

	if method == "CONNECT" {
		n, err = clientConn.Write([]byte(established))
		if err != nil || n == 0 {
			return cp, errors.New(fmt.Sprintf("send established to client erro:%v,readed:%v", err, n))
		}
	}

	if string(ackBytes[0:n]) == "" {
		return cp, errors.New(fmt.Sprint("receive ack message error:", string(ackBytes[0:n])))
	}

	return utils.NewClientPair(clientConn, sc, connectId), nil

}

//Parse the connection address and port from the http header
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

func allocConnectId() int64 {
	return atomic.AddInt64(&incrId, 1)
}
