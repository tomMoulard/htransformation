package types

import (
	"errors"
	"net/http"
	"regexp"
)

// RuleType define the possible types of rules.
type RuleType string

const (
	// Set will set the value of a header.
	Set RuleType = "Set"
	// Join will concatenate the values of headers.
	Join RuleType = "Join"
	// Delete will delete the value of a header.
	Delete RuleType = "Del"
	// Rename will rename a header.
	Rename RuleType = "Rename"
	// RewriteValueRule will replace the value of a header with the provided value.
	RewriteValueRule RuleType = "RewriteValueRule"
)

// Rule struct so that we get traefik config.
type Rule struct {
	Header       string         `yaml:"Header"`       // header value
	HeaderPrefix string         `yaml:"HeaderPrefix"` // header prefix to find header
	Name         string         `yaml:"Name"`         // rule name
	Regexp       *regexp.Regexp `yaml:"-"`            // Used for rewrite, rename header matching
	Sep          string         `yaml:"Sep"`          // separator to use for join
	Type         RuleType       `yaml:"Type"`         // Differentiate rule types
	Value        string         `yaml:"Value"`
	ValueReplace string         `yaml:"ValueReplace"` // value used as replacement in rewrite
	Values       []string       `yaml:"Values"`       // values to join
	// if SetOnResponse is true, the header will be changed on the response. It will be on the request otherwise (default).
	SetOnResponse bool `yaml:"SetOnResponse"`
}

var ErrMissingRequiredFields = errors.New("missing required fields")

var ErrInvalidRuleType = errors.New("invalid rule type")

var ErrInvalidRegexp = errors.New("invalid regexp")

type Handler interface {
	Validate() error
	Handle(rw http.ResponseWriter, req *http.Request)
}
