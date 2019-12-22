package client

import (
	"errors"
	"fmt"
	"httpproxy.v1/config"
	"httpproxy.v1/utils"
	"net"
)

//Read data from the client socket cyclically and write to the server socket
//Support disconnection then reconnection
func Server2Client(client, server *utils.ClientConn) (er, ew, err error) {
start:
	if er, ew, err = scCopy(client, server); err != nil {
		if err == utils.ErrRebuild {
			if reconect > reconectLimist {
				return nil, nil, errors.New(fmt.Sprint("reconnect times exceeded :", err))
			}
			reconect++
			sc, err := net.Dial("tcp", config.ServerHost)
			if err != nil {
				return nil, nil, errors.New(fmt.Sprint("reconnect server error:", err))
			}
			fmt.Println("reconnected with server:", config.ServerHost)
			server.Conn = sc
			goto start
		}
	}
	return er, ew, err
}

//Perform io copy
func scCopy(client, server *utils.ClientConn) (er, ew, err error) {
	var (
		size       = 32 * 1024
		decodeData []byte
		readedData []byte
		nr, nw     int
	)
	buf := make([]byte, size)

	for {

		nr, er = server.Read(buf, true)
		readedData = buf[0:nr]

		if nr > 0 {
			decodeData, err = utils.DecryptAES(readedData)
			if err != nil {
				break
			}
			nr = len(decodeData)
			nw, ew = client.Write(decodeData, true)
			if ew != nil {
				break
			}
			if nr != nw {
				err = utils.ErrShortWrite
				break
			}
		}
		//Read server disconnected abnormally
		if er != nil && ew != utils.EOF {
			err = utils.ErrRebuild
			break
		}
	}
	return
}
