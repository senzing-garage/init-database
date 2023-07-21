//go:build darwin

package cmd

import "github.com/senzing/go-common/option"

var ContextVariablesForOsArch = []option.ContextVariable{
	option.SenzingDirectory,
	option.ConfigPath,
	option.ResourcePath,
	option.SupportPath,
}
