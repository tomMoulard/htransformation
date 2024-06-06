package deleter

import (
	"net/http"

	"github.com/tomMoulard/htransformation/pkg/types"
	"github.com/tomMoulard/htransformation/pkg/utils/header"
)

type Deleter struct {
	rule *types.Rule
}

func New(rule types.Rule) (types.Handler, error) {
	return &Deleter{rule: &rule}, nil
}

func (d *Deleter) Validate() error {
	return nil
}

func (d *Deleter) Handle(rw http.ResponseWriter, req *http.Request) {
	if d.rule.SetOnResponse {
		rw.Header().Del(d.rule.Header)

		return
	}

	header.Delete(req, d.rule.Header)
}
