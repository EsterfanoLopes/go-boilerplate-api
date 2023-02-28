//go:build tools
// +build tools

package tools

import (
	_ "github.com/cespare/reflex"
	_ "github.com/kyoh86/richgo"
	_ "github.com/pressly/goose/cmd/goose"
	_ "github.com/swaggo/swag/cmd/swag"
	_ "github.com/vektra/mockery/v2"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
