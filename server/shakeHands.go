package server

import (
	"errors"
	"fmt"
	"httpproxy.v1/utils"
	"net"
	"strconv"
	"strings"
)

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

	sd := strings.Split(string(data), "#")
	if len(sd) != 2 {
		return cp, errors.New(fmt.Sprintf("shake msg content error:%v-%v", string(data), sd))
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
