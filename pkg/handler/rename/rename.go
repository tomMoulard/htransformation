package rename

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/tomMoulard/htransformation/pkg/types"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

func Validate(rule types.Rule) error {
	if _, err := regexp.Compile(rule.Header); err != nil {
		return fmt.Errorf("%s: %w", types.ErrInvalidRegexp.Error(), err)
	}

	if rule.Value == "" {
		return types.ErrMissingRequiredFields
	}

	return nil
}

func Handle(rw http.ResponseWriter, req *http.Request, rule types.Rule) {
	originalHost := req.Header.Get("Host") // Eventually X-Forwarded-Host
	req.Header.Set("Host", req.Host)

	for headerName, headerValues := range req.Header {
		if matched := rule.Regexp.Match([]byte(headerName)); !matched {
			continue
		}

		if rule.SetOnResponse {
			rw.Header().Del(headerName)
		} else {
			header.Delete(req, headerName)
		}

		for _, val := range headerValues {
			if rule.SetOnResponse {
				rw.Header().Set(rule.Value, val)
			} else {
				header.Set(req, rule.Value, val)
			}
		}
	}

	req.Header.Set("Host", originalHost)
}
