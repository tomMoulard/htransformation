package header

import (
	"net/http"
	"strings"
)

func Delete(req *http.Request, header string) {
	if strings.EqualFold(header, "Host") {
		req.Host = ""
	} else {
		req.Header.Del(header)
	}
}
