package htransformation

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type Transform struct {
	Name   string `yaml:"Name"`
	Rename string `yaml:"Rename"`
	With   string `yaml:"With"`
	Type   string `yaml:"Type"`
}

//type Transformations []struct {
//	Transform struct {
//		Name   string `yaml:"Name"`
//		Rename string `yaml:"Rename"`
//		With   string `yaml:"With"`
//		Type   string `yaml:"Type"`
//	} `yaml:"Transform,omitempty"`
//}

// Config holds configuration to be passed to the plugin
type Config struct {
	Transformations []Transform
}

// CreateConfig populates the Config data object
func CreateConfig() *Config {
	return &Config{
		Transformations: []Transform{},
	}
}

// HeadersTransformation holds the necessary components of a Traefik plugin
type HeadersTransformation struct {
	next			http.Handler
	transformations []Transform
	name			string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &HeadersTransformation{
		transformations: config.Transformations,
		next:            next,
		name:            name,
	}, nil
}

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
	u.next.ServeHTTP(rw, req)
}
