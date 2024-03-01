package plugin

import (
	"os"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

type State struct{}

// / Hello - Hello a payment to a node from a BOLT 12 offer including the
// / information of the payer.
func Hello(plugin *plugin.Plugin[*State], request map[string]any) (map[string]any, error) {
	return map[string]any{"message": "hello from cln4go.template"}, nil
}

// / OnShutdown - Kill the plugin when cln is going to shutdown
func OnShutdown(plugin *plugin.Plugin[*State], request map[string]any) {
	os.Exit(0)
}
