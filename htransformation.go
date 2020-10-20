package htransformation

import (
	"context"
	"net/http"
	"regexp"
)

type Transform struct {
	Name   string `yaml:"Name"`
	Rename string `yaml:"Rename"`
	With   string `yaml:"With"`
	Type   string `yaml:"Type"`
}
type Set struct {
	Name   string `yaml:"Name"`
	Header string `yaml:"Header"`
	Value  string `yaml:"Value"`
}

// Config holds configuration to be passed to the plugin
type Config struct {
	Transformations []Transform
    Setters []Set
}

// CreateConfig populates the Config data object
func CreateConfig() *Config {
	return &Config{
		Transformations: []Transform{},
		Setters: []Set{},
	}
}

// HeadersTransformation holds the necessary components of a Traefik plugin
type HeadersTransformation struct {
	next			http.Handler
	transformations []Transform
	setters         []Set
	name			string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &HeadersTransformation{
		transformations: config.Transformations,
        setters:         config.Setters,
		next:            next,
		name:            name,
	}, nil
}

// Iterate over every headers to match the ones specified in the config and
// return nothing if regexp failed.
func (u *HeadersTransformation) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for headerName, headerValues := range req.Header {
		for _, trans := range u.transformations {
			matched, err := regexp.Match(trans.Rename, []byte(headerName))
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			if matched {
				req.Header.Del(headerName)
				for _, val := range headerValues {
					req.Header.Set(trans.With, val)
				}
			}
		}
	}
    for _, set := range u.setters {
        req.Header.Set(set.Header, set.Value)
    }
	u.next.ServeHTTP(rw, req)
}
