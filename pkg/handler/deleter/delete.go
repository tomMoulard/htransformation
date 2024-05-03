package deleter

import (
	"net/http"
	"strings"

	"github.com/tomMoulard/htransformation/pkg/types"
)

func Validate(types.Rule) error {
	return nil
}

func Handle(rw http.ResponseWriter, req *http.Request, rule types.Rule) {
	if rule.SetOnResponse {
		rw.Header().Del(rule.Name)

		return
	}

	if strings.EqualFold(rule.Header, "Host") {
		req.Host = ""
	} else {
		req.Header.Del(rule.Header)
	}
}
