package handle

import (
	"net"
	"net/http"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	vnet "github.com/v2fly/v2ray-core/v5/common/net"
)

var anyMethods = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete,
	http.MethodTrace, "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK",
	http.MethodConnect,
}

func isHTTP(peek []byte) bool {
	tempPeekString := strings.ToUpper(string(peek))
	logrus.Debug(tempPeekString)
	for _, m := range anyMethods {
		if strings.HasPrefix(tempPeekString, m) {
			return true
		}
	}
	return false
}

func isEnglishLetter(b byte) bool {
	// 检查是否为大写字母或小写字母
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func dialDestination(d vnet.Destination) (net.Conn, error) {
	dial := net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// 只有linux才可以使用 syscall 设置 SO_MARK, 目的是防止设置了本机透明代理时, 代理循环
				syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, 0xff)
			})
		},
	}
	return dial.Dial(strings.ToLower(d.Network.String()), net.JoinHostPort(d.Address.IP().String(), d.Port.String()))
	// return dial.Dial(strings.ToLower(d.Network.String()), net.JoinHostPort("192.168.6.157", "1234"))
}
