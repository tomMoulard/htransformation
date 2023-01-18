package htransformation

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/tomMoulard/htransformation/pkg/handler/deleter"
	"github.com/tomMoulard/htransformation/pkg/handler/join"
	"github.com/tomMoulard/htransformation/pkg/handler/rename"
	"github.com/tomMoulard/htransformation/pkg/handler/rewrite"
	"github.com/tomMoulard/htransformation/pkg/handler/set"
	"github.com/tomMoulard/htransformation/pkg/types"
)

// HeadersTransformation holds the necessary components of a Traefik plugin.
type HeadersTransformation struct {
	name         string
	next         http.Handler
	rules        []types.Rule
	ruleHandlers map[types.RuleType]func(http.ResponseWriter, *http.Request, types.Rule)
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
	ruleHandlers := map[types.RuleType]func(http.ResponseWriter, *http.Request, types.Rule){
		types.Delete:           deleter.Handle,
		types.Join:             join.Handle,
		types.Rename:           rename.Handle,
		types.RewriteValueRule: rewrite.Handle,
		types.Set:              set.Handle,
	}

	validateRules := map[types.RuleType]func(types.Rule) error{
		types.Delete:           deleter.Validate,
		types.Join:             join.Validate,
		types.Rename:           rename.Validate,
		types.RewriteValueRule: rewrite.Validate,
		types.Set:              set.Validate,
	}

	for i, rule := range config.Rules {
		if _, ok := ruleHandlers[rule.Type]; !ok {
			return nil, fmt.Errorf("%w: %s", types.ErrInvalidRuleType, rule.Name)
		}

		validate, ok := validateRules[rule.Type]
		if !ok {
			continue
		}

		if err := validate(rule); err != nil {
			return nil, fmt.Errorf("%w: %s", err, rule.Name)
		}

		if rule.Type == types.Rename || rule.Type == types.RewriteValueRule {
			re, err := regexp.Compile(rule.Header)
			if err != nil { // must be validated before
				return nil, fmt.Errorf("%w: %s", types.ErrInvalidRegexp, rule.Name)
			}

			config.Rules[i].Regexp = re
		}
	}

	return &HeadersTransformation{
		name:         name,
		next:         next,
		rules:        config.Rules,
		ruleHandlers: ruleHandlers,
	}, nil
}

// Iterate over every header to match the ones specified in the config and
// return nothing if regexp failed.
func (u *HeadersTransformation) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	for _, rule := range u.rules {
		ruleHandler, ok := u.ruleHandlers[rule.Type]
		if !ok {
			continue
		}

		ruleHandler(responseWriter, request, rule)
	}

	u.next.ServeHTTP(responseWriter, request)
}
