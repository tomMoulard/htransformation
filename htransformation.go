package htransformation

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Rule struct so that we get traefik config
type Rule struct {
	Name                string   `yaml:"Name"`
	Header              string   `yaml:"Header"`
	Value               string   `yaml:"Value"`
	Values              []string `yaml:"Values"`
	ValueIsHeaderPrefix string   `yaml:"vauleIsHeaderPrefix"`
	Sep                 string   `yaml:"Sep"`
	Type                string   `yaml:"Type"`
}

// Config holds configuration to be passed to the plugin
type Config struct {
	Rules []Rule
}

// CreateConfig populates the Config data object
func CreateConfig() *Config {
	return &Config{
		Rules: []Rule{},
	}
}

// HeadersTransformation holds the necessary components of a Traefik plugin
type HeadersTransformation struct {
	next  http.Handler
	rules []Rule
	name  string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	for _, rule := range config.Rules {
		if rule.Header == "" || rule.Type == "" {
			return nil, fmt.Errorf("can't use '%s', some required fields are empty",
				rule.Name)
		}
		if rule.Type == "Join" && (len(rule.Values) == 0 || rule.Sep == "") {
			return nil, fmt.Errorf("can't use '%s', some required fields are empty",
				rule.Name)
		}
		if rule.ValueIsHeaderPrefix != "" && rule.Value == "" && len(rule.Values) == 0 {
			return nil, fmt.Errorf("can't use '%s', cannot set ValueIsHeaderPrefix without passing in Value/Values",
				rule.Name)
		}
	}
	return &HeadersTransformation{
		rules: config.Rules,
		next:  next,
		name:  name,
	}, nil
}

// Iterate over every headers to match the ones specified in the config and
// return nothing if regexp failed.
func (u *HeadersTransformation) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, rule := range u.rules {
		switch rule.Type {
		case "Rename":
			for headerName, headerValues := range req.Header {
				matched, err := regexp.Match(rule.Header, []byte(headerName))
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
				if matched {
					req.Header.Del(headerName)
					for _, val := range headerValues {
						req.Header.Set(getValue(rule.Value, rule.ValueIsHeaderPrefix, req), val)
					}
				}
			}
		case "Set":
			req.Header.Set(rule.Header, getValue(rule.Value, rule.ValueIsHeaderPrefix, req))
		case "Del":
			req.Header.Del(rule.Header)
		case "Join":
			if val, ok := req.Header[rule.Header]; ok {
				tmp_val := val[0]
				for _, value := range rule.Values {
					tmp_val += rule.Sep + getValue(value, rule.ValueIsHeaderPrefix, req)
				}
				// Delete after creating the tmp_val, so that if the values refer itself, it won't be empty.
				req.Header.Del(rule.Header)
				req.Header.Add(rule.Header, tmp_val)
			}
		default:
		}
	}
	u.next.ServeHTTP(rw, req)
}

// getValue checks if prefix exists, the given prefix is present, and then proceeds to read the existing header (after stripping the prefix) to return as value
func getValue(ruleValue string, vauleIsHeaderPrefix string, req *http.Request) string {
	actualValue := ruleValue
	if vauleIsHeaderPrefix != "" && strings.HasPrefix(ruleValue, vauleIsHeaderPrefix) {
		header := strings.TrimPrefix(ruleValue, vauleIsHeaderPrefix)
		// If the resulting value after removing the prefix is empty (value was only prefix), we return the actual value, which is the prefix itself.
		// This is because doing a req.Header.Get("") would not fly well.
		if header != "" {
			actualValue = req.Header.Get(header)
		}
	}
	return actualValue
}
