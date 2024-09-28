package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"uaProxy/bootstrap"
	"uaProxy/handle"

	"github.com/sirupsen/logrus"
)

var anyMethods = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete,
	http.MethodTrace, "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK",
	http.MethodConnect, // 这个webdav没有
}

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

	logrus.Infoln("Proxy server listening on port", bootstrap.C.RedirPort)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handle.HandleConnection(clientConn)
	}
}
