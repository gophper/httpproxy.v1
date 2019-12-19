package utils

import (
	"fmt"
	"gostudy/github.com/name5566/leaf/log"
	"os"
	"strings"
	"time"
)

type ClientPair struct {
	InputClient  *ClientConn
	OutputClient *ClientConn
}

var filePath = "E:\\code\\go\\src\\gostudy\\httpproxy\\logs\\"

func (c ClientPair) Info(arg ...interface{}) {

	fileName := strings.Replace(fmt.Sprintf("%s-%s-client", c.OutputClient.Conn.LocalAddr(), c.InputClient.Conn.RemoteAddr()), ":", "", -1)
	fpath := fmt.Sprintf(filePath+"%s.txt", fileName)
	str := fmt.Sprintln(c.InputClient.Conn.RemoteAddr(), c.OutputClient.Conn.LocalAddr(), ":", arg)
	tracefile(fpath, str)
	fmt.Println(str)
}

func (c ClientPair) InfoV2(arg ...interface{}) {

	fileName := strings.Replace(fmt.Sprintf("%s-%s-server", c.InputClient.Conn.RemoteAddr(), c.OutputClient.Conn.LocalAddr()), ":", "", -1)
	fpath := fmt.Sprintf(filePath+"%s.txt", fileName)
	str := fmt.Sprintln(c.InputClient.Conn.RemoteAddr(), c.OutputClient.Conn.LocalAddr(), ":", arg)
	tracefile(fpath, str)
	fmt.Println(str)

}
func (c ClientPair) Warn(arg ...interface{}) {

}
func (c ClientPair) Error(arg ...interface{}) {

}

func tracefile(fileName, content string) {
	var (
		err error
	)
	fileObj, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fd_time := time.Now().Format("2006-01-02 15:04:05======");
	_, err = fmt.Fprintf(fileObj, fd_time+content)

	if err != nil {
		log.Error("error", err)
	}

}