package htransformation

import (
	"fmt"
	"context"
	"net/http"
	"log"
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
	next		http.Handler
	transformations []transformations
	name		string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &UIDDemo{
		transformations:	config.Transformations,
		next:		next,
		name:		name,
	}, nil
}

func (u *UIDDemo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	uid := "toto"
	http.Error(rw, "toto", http.StatusInternalServerError)
	// header injection to backend service
	// req.Header.Set(u.headerName, uid)
	// header injection to client response
	// rw.Header().Add(u.headerName, uid)
	for i, t := range u.transformations {
		rw.Header().Add(t.With, uid)
	}

	u.next.ServeHTTP(rw, req)
}
