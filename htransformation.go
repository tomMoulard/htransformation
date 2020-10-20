package htransformation

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type transformations struct {
	Name   string
	Rename string
	With   string
	Type   string
}

// Config holds configuration to be passed to the plugin
type Config struct {
	Transformations []transformations
}

// CreateConfig populates the Config data object
func CreateConfig() *Config {
	return &Config{}
}

// UIDDemo holds the necessary components of a Traefik plugin
type UIDDemo struct {
	next            http.Handler
	transformations []transformations
	name            string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &UIDDemo{
		transformations: config.Transformations,
		next:            next,
		name:            name,
	}, nil
}

func (u *UIDDemo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for headerName, headerValues := range r.Header {
		for _, trans := range u.transformations {
			matched, err := regexp.Match(trans.Rename, []byte(headerName))
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			if matched {
				rw.Header().Del(headerName)
				for _, val := range headerValues {
					rw.Header().Add(t.With, val)
				}
			}
		}
	}
	u.next.ServeHTTP(rw, req)
}
