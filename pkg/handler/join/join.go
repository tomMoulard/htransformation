package join

import (
	"net/http"
	"strings"

	"github.com/tomMoulard/htransformation/pkg/types"
)

type Join struct {
	rule *types.Rule
}

func New(rule types.Rule) (types.Handler, error) {
	return &Join{rule: &rule}, nil
}

func (j *Join) Validate() error {
	if len(j.rule.Values) == 0 || j.rule.Sep == "" {
		return types.ErrMissingRequiredFields
	}

	return nil
}

func (j *Join) Handle(rw http.ResponseWriter, req *http.Request) {
	var val []string
	if strings.EqualFold(j.rule.Header, "Host") {
		val = []string{req.Host}
	} else {
		var ok bool
		val, ok = req.Header[j.rule.Header]

		if !ok {
			return
		}
	}

	newHeaderVal := val[0]
	for _, value := range j.rule.Values {
		newHeaderVal += j.rule.Sep + getValue(value, j.rule.HeaderPrefix, req)
	}

	if j.rule.SetOnResponse {
		rw.Header().Set(j.rule.Name, newHeaderVal)

		return
	}

	if strings.EqualFold(j.rule.Header, "Host") {
		req.Host = newHeaderVal
	} else {
		req.Header.Set(j.rule.Header, newHeaderVal)
	}
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
		if header == "" {
			return actualValue
		}

		if strings.EqualFold(header, "Host") {
			actualValue = req.Host
		} else {
			actualValue = req.Header.Get(header)
		}
	}

	return actualValue
}
