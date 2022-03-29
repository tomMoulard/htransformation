package join

import (
	"net/http"
	"strings"

	"github.com/tomMoulard/htransformation/pkg/types"
)

func Handle(_ http.ResponseWriter, req *http.Request, rule types.Rule) {
	if val, ok := req.Header[rule.Header]; ok {
		tmpVal := val[0]

		for _, value := range rule.Values {
			tmpVal += rule.Sep + getValue(value, rule.HeaderPrefix, req)
		}

		req.Header.Del(rule.Header)
		req.Header.Add(rule.Header, tmpVal)
	}
}

// getValue checks if prefix exists, the given prefix is present,
// and then proceeds to read the existing header (after stripping the prefix)
// to return as value.
func getValue(ruleValue, vauleIsHeaderPrefix string, req *http.Request) string {
	actualValue := ruleValue

	if vauleIsHeaderPrefix != "" && strings.HasPrefix(ruleValue, vauleIsHeaderPrefix) {
		header := strings.TrimPrefix(ruleValue, vauleIsHeaderPrefix)
		// If the resulting value after removing the prefix is empty,
		// we return the actual value,
		// which is the prefix itself.
		// This is because doing a req.Header.Get("") would not fly well.
		if header != "" {
			actualValue = req.Header.Get(header)
		}
	}

	return actualValue
}
