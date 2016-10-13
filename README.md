netproxy
--------
netproxy 是一个工作在4层的透明代理，你可以用它来代理某一个服务，并且控制网络的通断，丢包率，发包速率等，可以帮助你模拟一些网络场景，去更好的测试你的程序。

## 使用指南
### netproxy
netproxy 是代理主程序，启动后会加载配置，监听多个代理地址和一个api地址，其中api地址是给cli使用的。

```
./netproxy --help

Usage of ./netproxy:
  -config string
    	config file
  -cpus int
    	use cpu nums (default 1)
  -log-file string
    	which file to record log, if not set stdout to use.
  -log-level string
    	log level. (default "info")

```

#### 配置文件：
```
{
    "apiAddr": ":8080",
    "units": [
        {
            "src": ":6379",
            "target": "192.168.159.130:6379",
            "desc": "zk server proxy",
            "ctl": {
                "in": {
                    "lossRate": 0,
                    "delayMs": 0
                },
                "out": {
                    "lossRate": 0,
                    "delayMs": 0
                }
            },
            "timeoutConnect": 5,
            "timeoutWrite": 30
        }
    ]
}
```
#### 配置参数说明：
* apiAddr
restful的http接口，提供给cli程序使用，用来修改某一个客户端链接的丢包率，延迟设置，并且实时生效。
* units
代理设置，可以设置多个代理，那么netproxy就会监听多个端口。其中src表示netproxy监听的端口，target代表代理的实际服务的地址，ctl是网络控制。

### cli
cli是一个客户端命令行工具，通过api的restfulhttp接口和netproxy通信，动态的修改某一个客户端的延迟和丢包设置。主要包括list和update两个子命令
```
./cli --help
A simple command line client for netproxy.

Usage:
  cli [command]

Available Commands:
  list        List the proxies
  update      Update the proxy ctl
  pause       Pause proxy on addr
  resume      Resume proxy on addr

Flags:
      --endpoints string   netproxt api address (default "127.0.0.1:8080")

Use "cli [command] --help" for more information about a command.
```


#### list命令
list 命令列出所有proxy的信息

```
./cli list --help
List the proxies

Usage:
  cli list [flags]

Global Flags:
      --endpoints string   netproxt api address (default "127.0.0.1:8080")```

#### update命令
update 命令用于修改某一个客户端的丢包和超时配置

```
./cli update --help
Update the proxy ctl

Usage:
  cli update <proxyAddr> [flags]

Flags:
      --in-delayMs int     set the client receive packet delay.
      --in-lossRate int    set the client receive packet loss rate.
      --out-delayMs int    set the client sent packet delay.
      --out-lossRate int   set the client sent packet loss rate.

Global Flags:
      --endpoints string   netproxy api address (default "127.0.0.1:8080")
```

#####  参数说明
* in-delayMs
控制客户端程序收包的延迟
* in-lossRate
控制客户端程序收包的丢包率
* out-delayMs
控制客户端程序发包的延迟
* out-lossRate
控制客户端程序发包的丢包率

#### pause命令
pause 命令用于暂停一个proxy的监听，模拟网络中断不通

```
./cli pause --help
Pause proxy on addr

Usage:
  cli pause <addr> [flags]

Global Flags:
      --endpoints string   netproxy api address (default "127.0.0.1:8080")

```


#### resume命令
resume 命令用于恢复一个proxy的监听，模拟网络中断恢复

```
./cli pause --help
Resume proxy on addr

Usage:
  cli resume <addr> [flags]

Global Flags:
      --endpoints string   netproxy api address (default "127.0.0.1:8080")
```




