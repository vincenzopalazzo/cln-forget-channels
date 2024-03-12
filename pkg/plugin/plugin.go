package plugin

import (
	"fmt"
	"os"

	json "github.com/mitchellh/mapstructure"

	"github.com/vincenzopalazzo/cln4go/client"
	"github.com/vincenzopalazzo/cln4go/plugin"
)

type State struct {
	rpc *client.UnixRPC
}

func (self *State) Rpc(method string, args map[string]any) (map[string]any, error) {
	return client.Call[map[string]any, map[string]any](self.rpc, method, args)
}

// / ForgetChannels - ForgetChannels a payment to a node from a BOLT 12 offer including the
// / information of the payer.
func ForgetChannels(cln *plugin.Plugin[*State], request map[string]any) (map[string]any, error) {
	return forgetChannels(cln, request)
}

// / OnShutdown - Kill the plugin when cln is going to shutdown
func OnShutdown(_ *plugin.Plugin[*State], request map[string]any) {
	os.Exit(0)
}

// / OnInit - Callback on the init RPC call
func OnInit(cln *plugin.Plugin[*State], conf map[string]any) map[string]any {
	clnConf := struct {
		LightningDir string `mapstructure:"lightning-dir"`
		RpcFile      string `mapstructure:"rpc-file"`
	}{}
	if err := json.Decode(conf, &clnConf); err != nil {
		return map[string]any{
			"disable": err,
		}
	}
	rpc, err := client.NewUnix(fmt.Sprintf("%s/%s", clnConf.LightningDir, clnConf.RpcFile))
	if err != nil {
		return map[string]any{
			"disable": err,
		}
	}
	cln.State.rpc = rpc
	return checkIfThereDeveloperIsEnable(cln)
}

func checkIfThereDeveloperIsEnable(cln *plugin.Plugin[*State]) map[string]any {
	listConfig, err := cln.State.Rpc("listconfigs", map[string]any{})
	if err != nil {
		cln.Log("debug", fmt.Sprintf("error while calling listconfigs: %s", err))
		return map[string]any{
			"disable": err,
		}
	}
	config, found := listConfig["configs"].(map[string]any)
	if !found {
		return map[string]any{
			"disable": "error while looking inside the configuration, the `configs` keys is not present",
		}
	}
	developerEnable := config["developer"].(map[string]any)
	developerValue, found := developerEnable["set"].(bool)
	if !found || !developerValue {
		return map[string]any{
			"disable": fmt.Sprintf("disabling because we need a `developer` mode enabled (%t)", developerValue),
		}
	}
	return map[string]any{}
}
