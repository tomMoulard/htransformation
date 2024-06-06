package set

import (
	"net/http"

	"github.com/tomMoulard/htransformation/pkg/types"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

type Set struct {
	rule *types.Rule
}

func New(rule types.Rule) (types.Handler, error) {
	return &Set{rule: &rule}, nil
}

func (s *Set) Validate() error {
	if s.rule.Header == "" {
		return types.ErrMissingRequiredFields
	}

	return nil
}

func (s *Set) Handle(rw http.ResponseWriter, req *http.Request) {
	if s.rule.SetOnResponse {
		rw.Header().Set(s.rule.Header, s.rule.Value)

		return
	}

	header.Set(req, s.rule.Header, s.rule.Value)
}
