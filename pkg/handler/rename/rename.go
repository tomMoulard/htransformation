package rename

import (
	"net/http"

	"github.com/tommoulard/htransformation/pkg/types"
)

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) {
	for headerName, headerValues := range req.Header {
		if matched := rule.Regexp.Match([]byte(headerName)); !matched {
			continue
		}

		req.Header.Del(headerName)

		for _, val := range headerValues {
			req.Header.Set(rule.Value, val)
		}
	}
}
