package rename

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/tommoulard/htransformation/pkg/types"
)

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) error {
	for headerName, headerValues := range req.Header {
		matched, err := regexp.Match(rule.Header, []byte(headerName))
		if err != nil {
			return fmt.Errorf("RenameHandler error: %w", err)
		}

		if !matched {
			continue
		}

		req.Header.Del(headerName)

		for _, val := range headerValues {
			req.Header.Set(rule.Value, val)
		}
	}

	return nil
}
