package main

import (
	"fmt"
	"net"

	"uaProxy/bootstrap"
	"uaProxy/handle"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	bootstrap.LoadConfig()

	// 监听代理端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("redir-port")))
	if err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Proxy server listening on port: \033[1;34m%d\033[0m, UA is set to \033[1;34m%s\033[0m\n",
		viper.GetInt("redir-port"), viper.GetString("User-Agent"))

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			logrus.Error("Error accepting connection:", err)
			continue
		}
		go handle.HandleConnection(clientConn)
	}
}
