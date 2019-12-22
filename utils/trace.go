package utils

import (
	"fmt"
	"httpproxy.v1/config"
	"os"
	"strings"
	"time"
)

var filePath string

func init() {
	filePath = config.GetConfig("sys", "logDir")
}

func (c ClientPair) TraceV1(args ...interface{}) {
	if !config.Trace {
		return
	}
	fileName := strings.Replace(fmt.Sprintf("%s-%s-client", c.OutputClient.Conn.LocalAddr(), c.InputClient.Conn.RemoteAddr()), ":", "", -1)
	c.traceFile(fileName, args)
}

func (c ClientPair) TraceV2(args ...interface{}) {
	if !config.Trace {
		return
	}
	fileName := strings.Replace(fmt.Sprintf("%s-%s-server", c.InputClient.Conn.RemoteAddr(), c.OutputClient.Conn.LocalAddr()), ":", "", -1)
	c.traceFile(fileName, args)
}

func (c ClientPair) traceFile(fileName string, args ...interface{}) {
	var (
		err error
	)

	var logArgs []interface{}
	for _, value := range args {
		logArgs = append(logArgs, value)
		logArgs = append(logArgs, "\r\n")
	}
	content := fmt.Sprintln(c.InputClient.Conn.RemoteAddr(), c.OutputClient.Conn.LocalAddr(), ":", args)

	fileObj, err := os.OpenFile(fmt.Sprintf(filePath+"%s.log", fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fd_time := time.Now().Format("2006-01-02 15:04:05==>")
	_, err = fmt.Fprintf(fileObj, fd_time+content)
	if err != nil {
		fmt.Println("error:", err)
	}
}
