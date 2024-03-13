package plugin

import (
	"fmt"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

func WithdrawCompletedTx(cln *plugin.Plugin[*State], request map[string]any) (map[string]any, error) {
	listFunds, err := cln.State.Rpc("listfunds", map[string]any{})
	if err != nil {
		return nil, err
	}

	outputs, found := listFunds["outputs"].([]any)
	if !found {
		return nil, fmt.Errorf("Outputs insid the object is not found: `%s`", stringify(listFunds))
	}

	cln.Log("info", stringify(outputs))
	utxo := make([]string, 0)
	for _, output := range outputs {
		cln.Log("info", fmt.Sprintf("output %s", stringify(output)))
		output := output.(map[string]any)
		switch output["status"].(string) {
		case "confirmed":
			cln.Log("info", fmt.Sprintf("found a eligible tx `%s`", stringify(output)))
			utxo = append(utxo, fmt.Sprintf("%s:%d", output["txid"], uint32(output["output"].(float32))))
		default:
			cln.Log("info", fmt.Sprintf("found a tx that it is not eligible `%s`", stringify(output)))
			continue
		}
	}

	if len(utxo) == 0 {
		return nil, fmt.Errorf("none to witdraw from `%s`", stringify(outputs))
	}

	address, found := request["destination"]
	if !found {
		return nil, fmt.Errorf("please specify an destination address where to withdraw")
	}

	amount, found := request["amount"]
	if !found {
		return nil, fmt.Errorf("please specify the amount of sats to sent to the address `%s`", address)
	}

	req := map[string]any{
		"destination": address,
		"satoshi":     amount,
		"utxos":       utxo,
	}
	cln.Log("info", fmt.Sprintf("calling withdraw with body `%s`", stringify(req)))
	return cln.State.Rpc("withdraw", req)
}
