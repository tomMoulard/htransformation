package deleter

import (
	"net/http"

	"github.com/tomMoulard/htransformation/pkg/types"
)

func Validate(rule types.Rule) error {
	return nil
}

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) {
	req.Header.Del(rule.Header)
}
