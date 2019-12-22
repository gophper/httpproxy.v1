package client

import (
	"errors"
	"fmt"
	"httpproxy.v1/config"
	"httpproxy.v1/utils"
	"net"
)

const reconectLimist = 3 //Limit on the number of reconnections
var reconect = 0         //Reconnected times

//Read data from the client socket cyclically and write to the server socket
//Support disconnection then reconnection
func Client2Server(server, client *utils.ClientConn) (er, ew, err error) {
	var (
		waitWritten []byte
	)
start:
	if waitWritten, er, ew, err = csCopy(server, client, waitWritten); err != nil {
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
func csCopy(server, client *utils.ClientConn, needWritten []byte) (waitWritten []byte, er, ew, err error) {
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
			//Rewrite nr
			nr = len(encodeData)
			nw, ew = server.Write(encodeData, true)
			//Write server is abnormally disconnected
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
