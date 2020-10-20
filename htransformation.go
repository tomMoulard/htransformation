package htransformation

import (
	"fmt"
	"context"
	"net/http"
)

// Config holds configuration to be passed to the plugin
type Config struct {
	HeaderName string
}

// CreateConfig populates the Config data object
func CreateConfig() *Config {
	return &Config{}
}

// UIDDemo holds the necessary components of a Traefik plugin
type UIDDemo struct {
	next		http.Handler
	headerName	string
	name		string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.HeaderName) == 0 {
		return nil, fmt.Errorf("HeaderName cannot be empty")
	}

	return &UIDDemo{
		headerName:	config.HeaderName,
		next:		next,
		name:		name,
	}, nil
}

func (u *UIDDemo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	uid := "toto"

	// header injection to backend service
	req.Header.Set(u.headerName, uid)
	// header injection to client response
	rw.Header().Add(u.headerName, uid)

	u.next.ServeHTTP(rw, req)
}
