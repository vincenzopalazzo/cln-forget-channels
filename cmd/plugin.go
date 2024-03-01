package main

import (
	core "github.com/cln-reckless/cln4go.plugin/pkg/plugin"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

func main() {
	state := core.State{}
	plugin := plugin.New(&state, true, plugin.DummyOnInit[*core.State])
	plugin.RegisterOption("foo", "string", "Hello Go", "An example of option", false)
	plugin.RegisterRPCMethod("hello", "", "an example of rpc method", core.Hello)
	plugin.RegisterNotification("shutdown", core.OnShutdown)
	plugin.Start()
}
