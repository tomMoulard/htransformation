package utils

import (
	"net/http"
	"strings"
)

func AddHeader(req *http.Request, header string, value string) {
	if strings.EqualFold(header, "Host") {
		req.Host += value
	} else {
		req.Header.Add(header, value)
	}
}

func SetHeader(req *http.Request, header string, value string) {
	if strings.EqualFold(header, "Host") {
		req.Host = value
	} else {
		req.Header.Set(header, value)
	}
}

func DeleteHeader(req *http.Request, header string) {
	if strings.EqualFold(header, "Host") {
		req.Host = ""
	} else {
		req.Header.Del(header)
	}
}
