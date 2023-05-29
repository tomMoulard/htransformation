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

func Handle(rw http.ResponseWriter, req *http.Request, rule types.Rule) {
	if rule.SetOnResponse {
		rw.Header().Set(rule.Header, rule.Value)

		return
	}

	req.Header.Set(rule.Header, rule.Value)
}
