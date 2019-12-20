httpproxy
========

##描述

使用Golang开发的http(s)代理
>本代理特点：
- 简单对称加密
- 流量追踪

>待实现或者优化：
- 支持更多的协议
- 流量负载均衡

##安装方法
- window: ./build.bat
- linux: ./build_linux.sh

##使用方法

>ini配置文件格式样列

	[sys]
    # 监听端口
    port = 8083
    
    # 监听地址
    host = 
    
    # 读套接字超时秒数
    timeout = 120
    
    [log]
    # 是否开启流量追踪
    trace = true
    # 日志目录
    logDir = /data/logs/httpproxy/logs/
    
    [crypt]
    # 256位对称加密密钥
    key = EEvjHa5qx5MHycO0o46AMz2yldO5RKBUpy37Y1lzHqqoWwERgp/XGTJ3UmHe7ZvICMJXHyyhPLd4bTtVVlDMPxK+iSGEvbt2bw+XI5IWjTG2hZlR54xdg7CerPFFnKWWJUrw2hNMAOTZApBxi2drdAbbCXXFzVgD9n70ME5uOs/r9da1zjhy/o8KpjVfHMaUf333eipg/Wz4FXvqU/xA+qJ835FIq8Ca5iYpXNEr4SezgS7VpPMFiOi8KNDvPg0U6UbUmDRineIEacTYZXm/R+z5Zg7y3MEvF1r/GzfuQkFouIe6SV7Sra8MIBjlT+AL3cs2ikOGIqk5TRrKcCRksQ==

