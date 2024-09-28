package handle

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
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
