package main

import (
	core "github.com/cln-reckless/cln4go.plugin/pkg/plugin"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

func main() {
	state := core.State{}
	plugin := plugin.New(&state, true, core.OnInit)
	plugin.RegisterRPCMethod("forget-channels", "", "A dangerus command that will help to clean up broken core lightning with a list of channel that will never confirm", core.ForgetChannels)
	plugin.RegisterNotification("shutdown", core.OnShutdown)
	plugin.Start()
}
