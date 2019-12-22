package main

import (
	"bytes"
	"errors"
	"fmt"
	"httpproxy.v1/config"
	"httpproxy.v1/utils"
	"log"
	"net"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	established = []byte("HTTP/1.1 200 Connection established\r\n\r\n")
	incrId      int64
)

const curVersion = "1.0.1"
const curMode = "client"

func main() {
	utils.Panel(Client{})
}

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

func handleClientRequest(socket net.Conn) {

	var (
		cp     utils.ClientPair
		client *utils.ClientConn
		server *utils.ClientConn
		err    error
		closed = false
	)
	defer func() {
		if !closed {
			socket.Close()
		}
	}()

	cp, err = shakeHands(socket)
	if err != nil {
		log.Println("Shake with server error:", err)
		return
	}

	client = cp.InputClient
	server = cp.OutputClient
	defer func() {
		if !closed {
			server.Close()
		}
	}()

	go func() {
		defer server.Close()
		if _, er, ew, err := Copy(server, client); err != nil {
			cp.TraceV1("client send to server error", er, ew, err)
		}
		cp.TraceV1("sended", client.Bufr, server.Bufw, "client logout")
	}()

	if _, er, ew, err := Copy2(client, server); err != nil {
		cp.TraceV1("server return to client error:", er, ew, err)
	}
	cp.TraceV1("received", server.Bufr, client.Bufw, "server logout,both close")
	client.Close()
	closed = true
}

//cryptFlag=1:encode,cryptFlag=2:decode
func Copy(server, client *utils.ClientConn) (waitWritten []byte, er, ew, err error) {
	if _, er, ew, err := Copy(server, client); err != nil {
		if err == utils.ErrRebuild {

		}
	}

}

func Copy2(server, client *utils.ClientConn, needWritten []byte) (waitWritten []byte, er, ew, err error) {
	var (
		size       = 32 * 1024
		encodeData []byte
		writeData  []byte
		nr, nw     int
	)
	buf := make([]byte, size)

	for {
		if nr = len(needWritten); nr > 0 {
			writeData = needWritten
		} else {
			nr, er = client.Read(buf, true)
			writeData = buf[0:nr]
		}

		if nr > 0 {
			encodeData, err = utils.EncryptAES(writeData)
			if err != nil {
				break
			}
			//重写nr
			nr = len(encodeData)
			nw, ew = server.Write(encodeData, true)
			//写服务端非正常断开
			if ew != nil && ew != utils.EOF {
				err = utils.ErrRebuild
				break
			}
			if nr != nw {
				err = utils.ErrShortWrite
				waitWritten = encodeData[nw:nr]
				break
			}
		}
		if er != nil {
			break
		}
	}
	return
}

func Copy3(dstClient, srcClient *utils.ClientConn, cryptFlag int) (written int64, er, ew, err error) {
	var (
		size   = 100
		data   []byte
		nr, nw int
	)
	buf := make([]byte, size)

	for {
		nr, er = srcClient.Read(buf, true)
		if nr > 0 {
			if cryptFlag == 1 {
				data, err = utils.EncryptAES(buf[0:nr])
				if err != nil {
					break
				}
			} else {
				data, err = utils.DecryptAES(buf[0:nr])
				if err != nil {
					break
				}
			}
			nr = len(data)
			nw, ew = dstClient.Write(data, true)
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				break
			}
			if nr != nw {
				err = utils.ErrShortWrite
				break
			}
		}
		if er != nil {
			break
		}
	}
	return
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

func allocConnectId() int64 {
	return atomic.AddInt64(&incrId, 1)
}

//Initiate a handshake request
func shakeHands(clientConn net.Conn) (cp utils.ClientPair, err error) {
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
