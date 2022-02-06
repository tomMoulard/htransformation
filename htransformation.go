package htransformation

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// RuleType define the possible types of rules
type RuleType string

const (
	// Set will set the value of a header
	Set RuleType = "Set"
	// Join will concatenate the values of headers
	Join RuleType = "Join"
	// Delete will delete the value of a header
	Delete RuleType = "Del"
	// Rename will rename a header
	Rename RuleType = "Rename"
	// RewriteValueRule will replace the value of a header with the provided value
	RewriteValueRule RuleType = "RewriteValueRule"
	// EmptyType defines an empty type rule
	EmptyType RuleType = ""
)

// Rule struct so that we get traefik config
type Rule struct {
	Name         string   `yaml:"Name"`
	Header       string   `yaml:"Header"`
	Value        string   `yaml:"Value"`
	ValueReplace string   `yaml:"ValueReplace"`
	Values       []string `yaml:"Values"`
	Sep          string   `yaml:"Sep"`
	Type         RuleType `yaml:"Type"`
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
	name  string
	next  http.Handler
	rules []Rule
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	for _, rule := range config.Rules {
		if rule.Header == "" || rule.Type == EmptyType {
			return nil, fmt.Errorf("can't use '%s', some required fields are empty",
				rule.Name)
		}

		if rule.Type == Join && (len(rule.Values) == 0 || rule.Sep == "") {
			return nil, fmt.Errorf("can't use '%s', some required fields are empty",
				rule.Name)
		}

		if rule.Type == RewriteValueRule && rule.ValueReplace == "" {
			return nil, fmt.Errorf("can't use %s, some required fields are empty",
				rule.Name)
		}
	}

	return &HeadersTransformation{
		name:  name,
		next:  next,
		rules: config.Rules,
	}, nil
}

// Iterate over every headers to match the ones specified in the config and
// return nothing if regexp failed.
func (u *HeadersTransformation) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, rule := range u.rules {
		switch rule.Type {
		case Rename:
			for headerName, headerValues := range req.Header {
				matched, err := regexp.Match(rule.Header, []byte(headerName))
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
				if matched {
					req.Header.Del(headerName)
					for _, val := range headerValues {
						req.Header.Set(rule.Value, val)
					}
				}
			}
		case RewriteValueRule:
			for headerName, headerValues := range req.Header {
				matched, err := regexp.Match(rule.Header, []byte(headerName))
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				if !matched {
					continue
				}

				req.Header.Del(headerName)
				for _, headerValue := range headerValues {
					replacedHeaderValue := rule.ValueReplace
					r := regexp.MustCompile(rule.Value)
					captures := r.FindStringSubmatch(headerValue)
					if len(captures) == 0 || captures[0] == "" {
						req.Header.Add(headerName, headerValue)
					}

					for j, capture := range captures[1:] {
						placeholder := fmt.Sprintf("$%d", j+1)
						replacedHeaderValue = strings.ReplaceAll(replacedHeaderValue, placeholder, capture)
					}

					req.Header.Add(headerName, replacedHeaderValue)
				}
			}
		case Set:
			req.Header.Set(rule.Header, rule.Value)
		case Delete:
			req.Header.Del(rule.Header)
		case Join:
			if val, ok := req.Header[rule.Header]; ok {
				req.Header.Del(rule.Header)
				tmpVal := val[0]
				for _, value := range rule.Values {
					tmpVal += rule.Sep + value
				}
				req.Header.Add(rule.Header, tmpVal)
			}
		default:
		}
	}

	u.next.ServeHTTP(rw, req)
}
