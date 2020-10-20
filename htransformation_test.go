package htransformation_test

import (
	"testing"

	plug "github.com/tommoulard/htransformation"
)

func TestDummy(t *testing.T) {
	cfg := plug.CreateConfig()
	t.Log(cfg)
}
