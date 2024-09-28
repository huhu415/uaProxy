package main

import (
	"fmt"
	"log"
	"net"
	"uaProxy/bootstrap"
	"uaProxy/handle"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := bootstrap.LoadConfig(); err != nil {
		logrus.Fatalf("Error loading config: %v", err)
	}

	// 监听代理端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", bootstrap.C.RedirPort))
	if err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	logrus.Infof("Proxy server listening on port: %d", bootstrap.C.RedirPort)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handle.HandleConnection(clientConn)
	}
}
