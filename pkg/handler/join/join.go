package join

import (
	"net/http"

	"github.com/tommoulard/htransformation/pkg/types"
)

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) {
	if val, ok := req.Header[rule.Header]; ok {
		req.Header.Del(rule.Header)

		tmpVal := val[0]

		for _, value := range rule.Values {
			tmpVal += rule.Sep + value
		}

		req.Header.Add(rule.Header, tmpVal)
	}
}
