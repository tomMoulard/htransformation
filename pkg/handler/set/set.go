package set

import (
	"net/http"
	"strings"

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

	if strings.EqualFold(rule.Header, "Host") {
		req.Host = rule.Value
	} else {
		req.Header.Set(rule.Header, rule.Value)
	}
}
