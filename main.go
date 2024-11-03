package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	"uaProxy/bootstrap"
	"uaProxy/handle"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bootstrap.LoadConfig()
	if bootstrap.C.Stats {
		p := bootstrap.C.StatsConfig
		fmt.Printf("path: \033[1;34m%s\033[0m, start recording...\n", p)
		bootstrap.NewParserRecord(ctx, p)
	}

	go server(ctx)

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Infof("Shutting down server...")
}

func server(ctx context.Context) {
	// 监听代理端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", bootstrap.C.RedirPort))
	if err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	headersList := []string{}
	for key, value := range bootstrap.C.Headers {
		headersList = append(headersList, (fmt.Sprintf("\033[1;34m%s\033[0m: \033[1;34m%s\033[0m", key, value)))
	}

	fmt.Printf("uaProxy server listening on port: \033[1;34m%d\033[0m, stats mode is \033[1;34m%v\033[0m, Headers: {%s}\n",
		bootstrap.C.RedirPort, bootstrap.C.Stats, strings.Join(headersList, ", "))

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Shutting down server...")
			return
		default:
			clientConn, err := listener.Accept()
			if err != nil {
				logrus.Error("Error accepting connection:", err)
				continue
			}
			go handle.HandleConnection(clientConn)
		}
	}
}
