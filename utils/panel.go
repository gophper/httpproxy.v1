package utils

import (
	"flag"
	"fmt"
	"httpproxy.v1/config"
	"os"
)

type ProxyProcesser interface {
	GetCurVersion() string
	GetCurMode() string
	Setup(string)
}

func Panel(c ProxyProcesser) {
	var (
		host, port string
		printVer   bool
	)
	config.InitConfig(c.GetCurMode())
	flag.BoolVar(&printVer, "v", false, "current version")
	flag.StringVar(&host, "b", "", "iP address for local monitoring")
	flag.StringVar(&port, "p", "", "port address for local listening")
	flag.IntVar(&config.Timeout, "t", 0, "timeout seconds")
	flag.BoolVar((*bool)(&config.Trace), "d", false, "log input and output")
	flag.Parse()

	if printVer {
		fmt.Println("httpproxy version:", c.GetCurVersion())
		os.Exit(0)
	}

	if port == "" {
		if p := config.GetConfig("sys", "port"); p != "" {
			port = p
		} else {
			fmt.Println("port needed !")
			os.Exit(0)
		}
	}

	if host == "" {
		if h := config.GetConfig("sys", "host"); h != "" {
			host = h
		}
	}

	if config.Timeout == 0 {
		if to := config.GetConfigInt("sys", "timeout"); to != 0 {
			config.Timeout = to
		}
	}

	if !config.Trace {
		if to := config.GetConfigBool("sys", "trace"); to {
			config.Trace = to
		}
	}

	c.Setup(fmt.Sprintf("%s:%s", host, port))
}
