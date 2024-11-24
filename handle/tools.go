package handle

import (
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	vnet "github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tcp"
)

var anyMethod = [][]byte{
	[]byte(http.MethodGet),
	[]byte(http.MethodPost),
	[]byte(http.MethodPut),
	[]byte(http.MethodPatch),
	[]byte(http.MethodHead),
	[]byte(http.MethodOptions),
	[]byte(http.MethodDelete),
	[]byte(http.MethodTrace),
	[]byte(http.MethodConnect),
	[]byte("PROPFIND"),
	[]byte("PROPPATCH"),
	[]byte("MKCOL"),
	[]byte("COPY"),
	[]byte("MOVE"),
	[]byte("LOCK"),
	[]byte("UNLOCK"),
	[]byte("LINK"),
	[]byte("UNLINK"),
	[]byte("PURGE"),
	[]byte("VIEW"),
	[]byte("REPORT"),
	[]byte("SEARCH"),
	[]byte("CHECKOUT"),
	[]byte("CHECKIN"),
	[]byte("MERGE"),
	[]byte("SUBSCRIBE"),
	[]byte("UNSUBSCRIBE"),
	[]byte("NOTIFY"),
}

func isHTTP(peek []byte) bool {
	if len(peek) == 0 {
		return false
	}

	for _, method := range anyMethod {
		if compareMethod(peek, method) {
			return true
		}
	}
	return false
}

func compareMethod(peek, method []byte) bool {
	for i := 0; i < len(peek) && i < len(method); i++ {
		if peek[i] != method[i] {
			return false
		}
	}
	return true
}

func isEnglishBigLetter(b byte) bool {
	// 检查是否为大写字母
	return (b >= 'A' && b <= 'Z')
}

func GetDestConn(clientConn net.Conn) (net.Conn, error) {
	var dest vnet.Destination
	var err error

	dest, err = tcp.GetOriginalDestination(clientConn)
	if err != nil {
		logrus.Errorf("failed to get original destination")
		return nil, err
	}

	logrus.Debugf("clientIp: %s, remoteIp: %s:%s", clientConn.RemoteAddr().String(), dest.Address.IP().String(), dest.Port)
	return dialDestination(dest)
}

func dialDestination(d vnet.Destination) (net.Conn, error) {
	dial := net.Dialer{
		Control: getDialerControl(),
	}
	return dial.Dial(strings.ToLower(d.Network.String()), net.JoinHostPort(d.Address.IP().String(), d.Port.String()))
	// return dial.Dial(strings.ToLower(d.Network.String()), net.JoinHostPort("192.168.6.157", "1234"))
}
