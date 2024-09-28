package handle

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"syscall"
	"uaProxy/bootstrap"

	"github.com/sirupsen/logrus"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tcp"
)

const SO_ORIGINAL_DST = 80

func HandleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	logrus.Debugf("clientConn. remoteAdd: %s", clientConn.RemoteAddr().String())
	// logrus.Debugf("clientConn. LocalAddr: %s", clientConn.LocalAddr().String())

	d, err := tcp.GetOriginalDestination(clientConn)
	if err != nil {
		logrus.Errorf("failed to get original destination")
	}
	logrus.Debugf("%s, ip: %s, port: %s", d.Network, d.Address.IP().String(), d.Port)

	dial := net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// 如果是linux系统，使用 syscall 设置 SO_MARK
				if runtime.GOOS == "linux" {
					// 使用 syscall 设置 SO_MARK
					syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, 0xff)
				}

			})
		},
	}
	serverConn, err := dial.Dial(strings.ToLower(d.Network.String()), net.JoinHostPort(d.Address.IP().String(), d.Port.String()))
	// serverConn, err := dial.Dial(strings.ToLower(d.Network.String()), net.JoinHostPort("192.168.6.157", "1234"))
	if err != nil {
		logrus.Errorf("failed to dial destination: %v", err)
		return
	}
	// defer serverConn.Close()

	buffSize := 10
	buff := make([]byte, buffSize)
	realBuffSize, err := clientConn.Read(buff)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Debugf("size:%d, content:%s", realBuffSize, string(buff))
	serverConn.Write(buff[:realBuffSize])

	if realBuffSize > 0 && isEnglishLetter(buff[0]) && isHTTP(buff[:realBuffSize]) {
		fmt.Println("This is HTTP traffic.")
		scanner := bufio.NewScanner(clientConn)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" { // 空行表示请求头结束
				serverConn.Write([]byte("\r\n"))
				break
			}
			templine := strings.ToLower(line)
			if strings.Contains(templine, "user-agent") {
				line = "User-Agent: " + bootstrap.C.UA
			}
			logrus.Println("Received header:", line)
			serverConn.Write([]byte(line + "\r\n"))
		}
	} else {
		fmt.Println("This is non-HTTP traffic.")
	}

	relay(serverConn, clientConn)
}

func relay(l, r net.Conn) {
	go io.Copy(l, r)
	io.Copy(r, l)
}
