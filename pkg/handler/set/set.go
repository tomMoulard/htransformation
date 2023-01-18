package set

import (
	"net/http"

	"github.com/tomMoulard/htransformation/pkg/types"
)

func Validate(rule types.Rule) error {
	if rule.Header == "" {
		return types.ErrMissingRequiredFields
	}

	return nil
}

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) {
	req.Header.Set(rule.Header, rule.Value)
}
