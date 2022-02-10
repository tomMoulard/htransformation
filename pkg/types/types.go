package types

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
	// EmptyType defines an empty type rule.
	EmptyType RuleType = ""
)

// Rule struct so that we get traefik config.
type Rule struct {
	Name         string   `yaml:"Name"`
	Header       string   `yaml:"Header"`
	Value        string   `yaml:"Value"`
	ValueReplace string   `yaml:"ValueReplace"`
	Values       []string `yaml:"Values"`
	Sep          string   `yaml:"Sep"`
	Type         RuleType `yaml:"Type"`
}
