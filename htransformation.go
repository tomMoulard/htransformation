package htransformation

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type RuleType string

const (
	Set              RuleType = "Set"
	Join             RuleType = "Join"
	Delete           RuleType = "Del"
	Rename           RuleType = "Rename"
	RewriteValueRule RuleType = "RewriteValueRule"
	EmptyType        RuleType = ""
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
	next  http.Handler
	rules []Rule
	name  string
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
		if rule.Type == RewriteValueRule && (len(rule.ValueReplace) == 0) {
			return nil, fmt.Errorf("can't use %s, some required fields are empty",
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
				if matched {
					req.Header.Del(headerName)
					for _, headerValue := range headerValues {
						replacedHeaderValue := rule.ValueReplace
						r := regexp.MustCompile(rule.Value)
						captures := r.FindAllStringSubmatch(headerValue, -1)
						if len(captures) == 0 || len(captures[0]) == 0 {
							continue
						}
						for j, capture := range captures[0][1:] {
							placeholder := fmt.Sprintf("$%d", j+1)
							replacedHeaderValue = strings.ReplaceAll(replacedHeaderValue, placeholder, capture)
						}
						req.Header.Add(headerName, replacedHeaderValue)
					}
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
