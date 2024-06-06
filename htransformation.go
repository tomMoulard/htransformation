package htransformation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tomMoulard/htransformation/pkg/handler/deleter"
	"github.com/tomMoulard/htransformation/pkg/handler/join"
	"github.com/tomMoulard/htransformation/pkg/handler/rename"
	"github.com/tomMoulard/htransformation/pkg/handler/rewrite"
	"github.com/tomMoulard/htransformation/pkg/handler/set"
	"github.com/tomMoulard/htransformation/pkg/types"
)

// HeadersTransformation holds the necessary components of a Traefik plugin.
type HeadersTransformation struct {
	name     string
	next     http.Handler
	handlers []types.Handler
}

// Config holds configuration to be passed to the plugin.
type Config struct {
	Rules []types.Rule
}

// CreateConfig populates the Config data object.
func CreateConfig() *Config {
	return &Config{
		Rules: []types.Rule{},
	}
}

// New instantiates and returns the required components used to handle an HTTP request.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	handlerBuilder := map[types.RuleType]func(types.Rule) (types.Handler, error){
		types.Delete:           deleter.New,
		types.Join:             join.New,
		types.Rename:           rename.New,
		types.RewriteValueRule: rewrite.New,
		types.Set:              set.New,
	}

	handlers := make([]types.Handler, 0, len(config.Rules))

	for _, rule := range config.Rules {
		newHandler, ok := handlerBuilder[rule.Type]
		if !ok {
			return nil, fmt.Errorf("%w: %s", types.ErrInvalidRuleType, rule.Name)
		}

		h, err := newHandler(rule)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", err, rule.Name)
		}

		if err := h.Validate(); err != nil {
			return nil, fmt.Errorf("%w: %s", err, rule.Name)
		}

		handlers = append(handlers, h)
	}

	return &HeadersTransformation{
		name:     name,
		next:     next,
		handlers: handlers,
	}, nil
}

// Iterate over every header to match the ones specified in the config and
// return nothing if regexp failed.
func (u *HeadersTransformation) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	for _, handler := range u.handlers {
		handler.Handle(responseWriter, request)
	}

	u.next.ServeHTTP(responseWriter, request)
}
