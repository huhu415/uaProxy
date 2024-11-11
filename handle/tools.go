package handle

import (
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	vnet "github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tcp"
)

var anyMethodSet = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodPatch:   {},
	http.MethodHead:    {},
	http.MethodOptions: {},
	http.MethodDelete:  {},
	http.MethodTrace:   {},
	http.MethodConnect: {},
	"PROPFIND":         {},
	"PROPPATCH":        {},
	"MKCOL":            {},
	"COPY":             {},
	"MOVE":             {},
	"LOCK":             {},
	"UNLOCK":           {},
	"LINK":             {},
	"UNLINK":           {},
	"PURGE":            {},
	"VIEW":             {},
	"REPORT":           {},
	"SEARCH":           {},
	"CHECKOUT":         {},
	"CHECKIN":          {},
	"MERGE":            {},
	"SUBSCRIBE":        {},
	"UNSUBSCRIBE":      {},
	"NOTIFY":           {},
}

func isHTTP(peek []byte) bool {
	tempPeekString := strings.ToUpper(string(peek))
	logrus.Debug(tempPeekString)

	first, _, ok := strings.Cut(tempPeekString, " ")
	if !ok {
		return false
	}

	if _, ok := anyMethodSet[first]; !ok {
		return false
	}
	return true
}

func isEnglishLetter(b byte) bool {
	// 检查是否为大写字母或小写字母
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func GetDestConn(clientConn net.Conn) (net.Conn, error) {
	var dest vnet.Destination
	var err error

	dest, err = tcp.GetOriginalDestination(clientConn)
	if err != nil {
		logrus.Errorf("failed to get original destination")
	}

	logrus.Debugf("%s, ip: %s, port: %s", dest.Network, dest.Address.IP().String(), dest.Port)
	return dialDestination(dest)
}

func dialDestination(d vnet.Destination) (net.Conn, error) {
	dial := net.Dialer{
		Control: getDialerControl(),
	}
	return dial.Dial(strings.ToLower(d.Network.String()), net.JoinHostPort(d.Address.IP().String(), d.Port.String()))
	// return dial.Dial(strings.ToLower(d.Network.String()), net.JoinHostPort("192.168.6.157", "1234"))
}
