package rename

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/tomMoulard/htransformation/pkg/types"
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

		switch {
		case rule.SetOnResponse:
			rw.Header().Del(headerName)
		case strings.EqualFold(headerName, "Host"):
			req.Host = ""
		default:
			req.Header.Del(headerName)
		}

		for _, val := range headerValues {
			switch {
			case rule.SetOnResponse:
				rw.Header().Set(rule.Value, val)
			case strings.EqualFold(rule.Value, "Host"):
				req.Host = val
			default:
				req.Header.Set(rule.Value, val)
			}
		}
	}

	req.Header.Set("Host", originalHost)
}
