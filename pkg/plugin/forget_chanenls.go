package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

// / # How to forget about a channel?
// /
// / Channels may end up stuck during funding and never
// / confirm on-chain. There is a variety of causes, the
// / most common ones being that the funds have been double-spent, or
// / the funding fee was too low to be confirmed.
// /
// / This is unlikely to happen in normal operation, as CLN tries to use
// / sane defaults and prevents double-spends whenever possible, but using custom feerates
// / or when the bitcoin backend has no good fee estimates it is still possible.
func forgetChannels(cln *plugin.Plugin[*State], request map[string]any) (map[string]any, error) {
	rescanOutputs, found := request["rescan"].(bool)
	if found && rescanOutputs {
		// Before forgetting about a channel it is important to ensure that
		// the funding transaction will never be confirmable by double-spending the funds.
		//
		// To do so you have to rescan the UTXOs using dev-rescan-outputs to reset
		// any funds that may have been used in the funding transaction, then move all the funds
		// to a new address.
		rescan, err := cln.State.Rpc("dev-rescan-outputs", map[string]any{})
		if err != nil {
			return nil, err
		}

		cln.Log("debug", fmt.Sprintf("dev-rescan-outpus returned %s", rescan))
		newAddr, err := cln.State.Rpc("newaddr", map[string]any{})
		if err != nil {
			return nil, err
		}
		addr := newAddr["p2tr"]
		cln.Log("debug", fmt.Sprintf("Sending the output to a new addr %s", addr))
		_, err = cln.State.Rpc("withdraw", map[string]any{
			"destination": addr,
			"satoshi":     "all",
		})

		if err != nil {
			return nil, err
		}
	}

	// FIXME: give the possibility to force close the unconfirmed chanenls,
	// BTW this is not suggested by the developer of this plugins

	// Now there are two reason that a channel is in the state of
	// CHANNELD_AWAITING_LOCKIN
	//
	// 1. Funding transaction isn't confirmed yet. In this case we have to wait longer, or,
	// 	in the case of a transaction that'll never confirm, forget the channel safely.
	//
	// 2. The peer hasn't sent a lockin message. This message acknowledges that
	// 	the node has seen sufficiently many confirmations to consider the channel funded.
	return checkChannelsToForget(cln, request)
}

type TxStatus = int

const (
	UNCONFIRMED_TX = iota
	MEM_POOL_DISCARDED
	CONFIRMED_TX
)

func stringify(obj any) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func checkChannelsToForget(cln *plugin.Plugin[*State], request map[string]any) (map[string]any, error) {
	listFunds, err := cln.State.Rpc("listfunds", map[string]any{})
	if err != nil {
		return nil, err
	}
	channels_modified := make([]map[string]any, 0)
	channels := listFunds["channels"].([]any)
	cln.Log("debug", fmt.Sprintf("channels to look inside size %d", len(channels)))
	for _, channel := range channels {
		cln.Log("debug", fmt.Sprintf("channel: `%s`", stringify(channel)))
		channel := channel.(map[string]any)
		peer_id := channel["peer_id"].(string)
		funding_txid := channel["funding_txid"].(string)
		funding_output := channel["funding_output"].(float64)
		channel_state := channel["state"].(string)
		cln.Log("debug", fmt.Sprintf("checking channel with state %s", channel_state))
		// Check if the status of the channel is correct with what we want
		if channel_state != "CHANNELD_AWAITING_LOCKIN" {
			continue
		}

		status, err := checkFundingTransaction(cln, funding_txid, uint32(funding_output))
		if err != nil {
			cln.Log("debug", fmt.Sprintf("error while checking the funding tx %d", status))
			return nil, err
		}
		short_channel_id, found := channel["channel_id"]
		if found {
			short_channel_id = short_channel_id.(string)
			cln.Log("debug", fmt.Sprintf("looking channel %s with status %s for confirmation", short_channel_id, channel_state))
		} else {
			short_channel_id = nil
			cln.Log("debug", fmt.Sprintf("looking channel without short channel id with status %s for confirmation", channel_state))
		}
		var state, action string
		switch status {
		case MEM_POOL_DISCARDED:
			action = "forget the channel"
			state = "MEM_POOL_DISCARDED"
			if err := forgetChannel(cln, channel); err != nil {
				return nil, err
			}
		case CONFIRMED_TX:
			action = "reconnect with the peer"
			state = "CONFIRMED_TX"
			_, err := cln.State.Rpc("reconnect", map[string]any{
				"peer_id": peer_id,
			})
			if err != nil {
				return nil, err
			}
		case UNCONFIRMED_TX:
			action = "logged"
			state = "UNCONFIRMED_TX"
			cln.Log("info", fmt.Sprintf("channel `%s` with peer `%s` still waiting on confirmed utxo `%s`", peer_id, short_channel_id, funding_txid))
		}
		cln.Log("info", fmt.Sprintf("channel funding tx state %s", state))
		channels_modified = append(channels_modified, map[string]any{
			"peer_id":          peer_id,
			"state":            state,
			"action":           action,
			"funding_txid":     funding_txid,
			"chanenl_state":    channel_state,
			"short_channel_id": short_channel_id,
		})
	}
	return map[string]any{
		"channels": channels_modified,
	}, nil
}

func checkFundingTransaction(cln *plugin.Plugin[*State], funding_tx string, funding_output uint32) (TxStatus, error) {
	// 1. check if there is an utxo
	cln.Log("debug", fmt.Sprintf("checking for funding tx %s", funding_tx))
	requestLink := fmt.Sprintf("https://mempool.space/api/tx/%s", funding_tx)
	res, err := http.Get(requestLink)
	if err != nil {
		return UNCONFIRMED_TX, nil
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		cln.Log("debug", fmt.Sprintf("return MEM_POOL_DISCARDED: %s - %s", err, res.Status))
		return MEM_POOL_DISCARDED, nil
	}

	str, err := io.ReadAll(res.Body)
	if err != nil {
		cln.Log("debug", fmt.Sprintf("return UNCONFIRMED_TX %s", err))
		return UNCONFIRMED_TX, err
	}
	cln.Log("debug", fmt.Sprintf("funding tx status `%s`", str))
	var response map[string]any
	if err := json.Unmarshal(str, &response); err != nil {
		cln.Log("debug", fmt.Sprintf("return UNCONFIRMED_TX %s", err))
		return UNCONFIRMED_TX, nil
	}
	statusJson, found := response["status"].(map[string]any)
	if !found {
		cln.Log("debug", fmt.Sprintf("return UNCONFIRMED_TX %s", err))
		return UNCONFIRMED_TX, nil
	}
	status, found := statusJson["confirmed"].(bool)
	if !found || status {
		cln.Log("debug", fmt.Sprintf("return CONFIRMED_TX %s", err))
		return CONFIRMED_TX, nil
	}
	cln.Log("info", fmt.Sprintf("Funding transaction %s looks discarded", funding_tx))
	return MEM_POOL_DISCARDED, nil
}

func forgetChannel(cln *plugin.Plugin[*State], channel map[string]any) error {
	// 1. Call the forget channel.
	request := map[string]any{
		"id": channel["peer_id"],
	}
	if val, ok := channel["short_channel_id"]; ok {
		request["short_channel_id"] = val
	}
	cln.Log("debug", fmt.Sprintf("call forget the channel with payload `%s`", stringify(request)))
	forgetChannel, err := cln.State.Rpc("dev-forget-channel", request)
	if err != nil {
		return err
	}
	cln.Log("info", fmt.Sprintf("forgetting channel `%s`", stringify(forgetChannel)))
	return nil
}
