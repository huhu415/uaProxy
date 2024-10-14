package handle

import (
	"bufio"
	"io"
	"net"
	"net/http"

	"uaProxy/bootstrap"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tcp"
)

func HandleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	logrus.Debugf("clientConn. remoteAdd: %s", clientConn.RemoteAddr().String())
	// logrus.Debugf("clientConn. LocalAddr: %s", clientConn.LocalAddr().String())

	d, err := tcp.GetOriginalDestination(clientConn)
	if err != nil {
		logrus.Errorf("failed to get original destination")
	}
	logrus.Debugf("%s, ip: %s, port: %s", d.Network, d.Address.IP().String(), d.Port)

	serverConn, err := dialDestination(d)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer serverConn.Close()

	bufioReader := bufio.NewReader(clientConn)
	peekBuff, _ := bufioReader.Peek(10)
	logrus.Debug(string(peekBuff))
	go io.Copy(clientConn, serverConn)
	if len(peekBuff) > 0 && isEnglishLetter(peekBuff[0]) && isHTTP(peekBuff) {
		logrus.Debug("this is http request")
		for {
			req, err := http.ReadRequest(bufioReader)
			if err != nil {
				return
			}

			if ua := req.Header.Get("User-Agent"); ua != "" {
				logrus.Debug(ua)
				bootstrap.GiveParserRecord().ParserAndRecord(ua)
			}
			req.Header.Set("User-Agent", viper.GetString("User-Agent"))

			if er := req.Write(serverConn); er != nil {
				logrus.Error(er)
				return
			}
		}
	} else {
		io.Copy(serverConn, bufioReader)
	}
}
