package join

import (
	"net/http"
	"strings"

	"github.com/tomMoulard/htransformation/pkg/types"
)

func Validate(rule types.Rule) error {
	if len(rule.Values) == 0 || rule.Sep == "" {
		return types.ErrMissingRequiredFields
	}

	return nil
}

func Handle(rw http.ResponseWriter, req *http.Request, rule types.Rule) {
	val, ok := req.Header[rule.Header]
	if !ok {
		return
	}

	newHeaderVal := val[0]
	for _, value := range rule.Values {
		newHeaderVal += rule.Sep + getValue(value, rule.HeaderPrefix, req)
	}

	if rule.SetOnResponse {
		rw.Header().Set(rule.Name, newHeaderVal)

		return
	}

	req.Header.Set(rule.Header, newHeaderVal)
}

// getValue checks if prefix exists, the given prefix is present,
// and then proceeds to read the existing header (after stripping the prefix)
// to return as value.
func getValue(ruleValue, valueIsHeaderPrefix string, req *http.Request) string {
	actualValue := ruleValue

	if valueIsHeaderPrefix != "" && strings.HasPrefix(ruleValue, valueIsHeaderPrefix) {
		header := strings.TrimPrefix(ruleValue, valueIsHeaderPrefix)
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
