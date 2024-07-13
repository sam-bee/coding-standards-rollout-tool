package csfixing

import (
	"testing"
)

func TestFix(t *testing.T) {
	conf := getConfig()
	Fix(conf)
}

func getConfig() ApplicationConfig {
	return BuildConfig(
		map[string]interface{}{
			"git": map[string]interface{}{
				"mainline-branch-name": "main",
				"remote-name":          "origin",
			},
			"codingstandards": map[string]interface{}{
				"command-to-run": "echo 'hello i am a command to fix coding standards in php or something'",
			},
		},
	)
}
