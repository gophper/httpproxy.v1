package utils

import (
	"errors"
	"fmt"
	"net"
)

var ErrShortWrite = errors.New("short write")
var EOF = errors.New("EOF")

type ClientConn struct {
	Conn net.Conn
	Bufw []byte
	Bufr []byte
}

//cryptFlag=1:encode,cryptFlag=2:decode
func Copy(dstClient, srcClient *ClientConn, cryptFlag int) (written int64, er, ew, err error) {
	var (
		size   = 32 * 1024
		data   []byte
		nr, nw int
	)
	buf := make([]byte, size)

	for {
		nr, er = srcClient.Read(buf, true)
		if nr > 0 {
			if cryptFlag == 1 {
				data, err = EncryptAES(buf[0:nr])
				if err != nil {
					break
				}
			} else {
				data, err = DecryptAES(buf[0:nr])
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
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			break
		}
	}
	return
}

func (c *ClientConn) Write(b []byte, trace bool) (n int, err error) {
	n, err = c.Conn.Write(b)
	if trace {
		c.Bufw = append(c.Bufw, b[0:n]...)
	}
	return
}

func (c *ClientConn) Read(b []byte, trace bool) (n int, err error) {
	n, err = c.Conn.Read(b)
	if trace {
		c.Bufr = append(c.Bufr, b[0:n]...)
	}
	return
}
func (c *ClientConn) Close() error {
	return c.Conn.Close()
}

func (c *ClientConn) RemoteAddr() string {
	return fmt.Sprintf("remote:%v,local:%v", c.Conn.RemoteAddr(), c.Conn.LocalAddr())
}
