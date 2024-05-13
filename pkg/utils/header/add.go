package header

import (
	"net/http"
	"strings"
)

func Add(req *http.Request, header string, value string) {
	if strings.EqualFold(header, "Host") {
		req.Host += value
	} else {
		req.Header.Add(header, value)
	}
}
