package deleter

import (
	"net/http"

	"github.com/tomMoulard/htransformation/pkg/types"
)

func Validate(rule types.Rule) error {
	return nil
}

func Handle(rw http.ResponseWriter, req *http.Request, rule types.Rule) {
	if rule.SetOnResponse {
		rw.Header().Del(rule.Name)

		return
	}

	req.Header.Del(rule.Header)
}
