package set

import (
	"net/http"

	"github.com/tommoulard/htransformation/pkg/types"
)

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) error {
	req.Header.Set(rule.Header, rule.Value)

	return nil
}
