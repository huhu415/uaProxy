package handle

import (
	"bufio"
	"io"
	"net"
	"net/http"

	"uaProxy/bootstrap"

	"github.com/sirupsen/logrus"
)

const PEEKSIZE = 8 // 不能太小, 如果是1的话, 容易全部都判定为http方法

func HandleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	// logrus.Debugf("clientConn. remoteAdd: %s", clientConn.RemoteAddr().String())

	serverConn, err := GetDestConn(clientConn)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer serverConn.Close()

	bufioReader := bufio.NewReader(clientConn)
	peekBuff, _ := bufioReader.Peek(PEEKSIZE)
	logrus.Debugf("locationIp: %s, remoteIp: %s, peekBuff: %v", clientConn.RemoteAddr().String(), serverConn.RemoteAddr().String(), peekBuff)
	go io.Copy(clientConn, serverConn)

	if len(peekBuff) > 0 && isEnglishBigLetter(peekBuff[0]) && isHTTP(peekBuff) {
		handleHTTPConnection(bufioReader, serverConn)
	} else {
		handleNonHTTPConnection(bufioReader, serverConn)
	}
}

func handleHTTPConnection(bufioReader *bufio.Reader, serverConn net.Conn) {
	logrus.Debug("this is http request")
	for {
		req, err := http.ReadRequest(bufioReader)
		if err != nil {
			return
		}

		if ua := req.Header.Get("User-Agent"); ua != "" {
			logrus.Debug(ua)
			if bootstrap.GiveParserRecord().ParserAndRecord(ua) {
				req.Header.Set("User-Agent", bootstrap.C.UserAgent)
			}
		}

		if er := req.Write(serverConn); er != nil {
			logrus.Error(er)
			return
		}

		if req.Header.Get("Upgrade") == "websocket" && req.Header.Get("Connection") == "Upgrade" {
			logrus.Debug("this is websocket request")
			break
		}
	}

	handleNonHTTPConnection(bufioReader, serverConn)
}

func handleNonHTTPConnection(bufioReader *bufio.Reader, serverConn net.Conn) {
	io.Copy(serverConn, bufioReader)
}
