package htransformation

import (
	"context"
	"net/http"
	"regexp"
)

// Rule struct so that we get traefik config
type Rule struct {
	Name   string   `yaml:"Name"`
	Header string   `yaml:"Header"`
	Value  string   `yaml:"Value"`
	Values []string `yaml:"Values"`
	Sep    string   `yaml:"Sep"`
	Type   string   `yaml:"Type"`
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
			if rule.Header == "" || rule.Value == "" {
				rw.Write([]byte("not done Rename: " + rule.Name + "\n"))
				continue
				// TODO: Add logs
			}
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
		case "Set":
			if rule.Header == "" || rule.Value == "" {
				rw.Write([]byte("not done Set : " + rule.Name + "\n"))
				continue
				// TODO: Add logs
			}
			req.Header.Set(rule.Header, rule.Value)
		case "Del":
			if rule.Header == "" {
				rw.Write([]byte("not done Del : " + rule.Name + "\n"))
				continue
				// TODO: Add logs
			}
			req.Header.Del(rule.Header)
		//JOIN application
		// If header found, then joining the value
		// If no header found, then skiping
		case "Join":
			if rule.Header == "" || rule.Value == "" || rule.Sep == "" {
				rw.Write([]byte("not done Join : " + rule.Name + "\n"))
				continue
				// TODO: Add logs
			}
			if val, ok := req.Header[rule.Header]; ok {
				req.Header.Del(rule.Header)
				req.Header.Add(rule.Header, val[0]+rule.Sep+rule.Value)
			}
		default:
		}
	}
	u.next.ServeHTTP(rw, req)
}
