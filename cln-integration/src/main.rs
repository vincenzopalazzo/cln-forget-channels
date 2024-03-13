//! Core Lightning Integration testing
//!
//! Author: Vincenzo Palazzo <vincenzopalazzo@member.fsf.org>
use std::sync::Once;

use clightning_testing::cln;
use serde::Deserialize;
use serde_json::{json, Value};

mod utils;

#[allow(unused_imports)]
use utils::*;

static INIT: Once = Once::new();

fn init() {
    // ignore error
    INIT.call_once(|| {
        env_logger::init();
    });
}

#[tokio::test(flavor = "multi_thread")]
async fn test_init_plugin() -> anyhow::Result<()> {
    init();
    let pwd = std::env!("PWD");
    let plugin_name = std::env!("PLUGIN_NAME");
    let cln1 = cln::Node::with_params(
        &format!("--developer --experimental-splicing --plugin={pwd}/../{plugin_name}"),
        "regtest",
    )
    .await?;
    let splice: Result<serde_json::Value, _> = cln1.rpc().call("getinfo", json!({}));
    assert!(splice.is_ok(), "{:?}", splice);
    Ok(())
}

#[tokio::test(flavor = "multi_thread")]
#[ntest::timeout(560000)]
async fn test_simple_devforgetchanenls() -> anyhow::Result<()> {
    init();

    let node_one = node!();
    let btc = node_one.btc();
    let node_two = node!(btc.clone());
    open_channel(&node_two, &node_one, false)?;

    #[derive(Deserialize, Debug)]
    struct ForgetChannels {
        channels: Vec<Value>,
    }
    let forget_channels: ForgetChannels = node_one.rpc().call("forget-channels", json!({}))?;
    assert_eq!(forget_channels.channels.len(), 0, "{:?}", forget_channels);
    Ok(())
}

#[tokio::test(flavor = "multi_thread")]
#[ntest::timeout(560000)]
async fn test_simple_withdraw_only_confirmed_one() -> anyhow::Result<()> {
    init();

    let node_one = node!();
    let btc = node_one.btc();
    let node_two = node!(btc.clone());
    open_channel(&node_two, &node_one, false)?;
    let addr = node_one.rpc().newaddr(None)?.bech32.unwrap();
    fund_wallet(node_one.btc(), &addr, 8)?;
    wait_for_funds(&node_one)?;

    let withdraw: Result<Value, _> = node_one.rpc().call("withdraw-only-confirmed", json!({}));
    assert!(withdraw.is_err());
    log::info!(target: "test_simple_withdraw_only_confirmed_one", "{:?}", withdraw);
    Ok(())
}

#[tokio::test(flavor = "multi_thread")]
#[ntest::timeout(560000)]
async fn test_simple_withdraw_only_confirmed_two() -> anyhow::Result<()> {
    init();

    let node_one = node!();
    let btc = node_one.btc();
    let node_two = node!(btc.clone());
    open_channel(&node_two, &node_one, false)?;
    let addr = node_one.rpc().newaddr(None)?.bech32.unwrap();
    fund_wallet(node_one.btc(), &addr, 8)?;
    wait_for_funds(&node_one)?;

    let addr = node_one.rpc().newaddr(None)?.bech32.unwrap();
    let withdraw: Result<Value, _> = node_one.rpc().call(
        "withdraw-only-confirmed",
        json!({
            "destination": addr
        }),
    );
    assert!(withdraw.is_err());
    log::info!(target: "test_simple_withdraw_only_confirmed_two", "{:?}", withdraw);
    Ok(())
}

#[tokio::test(flavor = "multi_thread")]
#[ntest::timeout(560000)]
async fn test_simple_withdraw_only_confirmed_3() -> anyhow::Result<()> {
    init();

    let node_one = node!();
    let btc = node_one.btc();
    let node_two = node!(btc.clone());
    open_channel(&node_two, &node_one, false)?;

    let addr = node_one.rpc().newaddr(None)?.bech32.unwrap();
    fund_wallet(node_one.btc(), &addr, 8)?;
    wait_for_funds(&node_one)?;

    let addr = node_one.rpc().newaddr(None)?.bech32.unwrap();
    let withdraw: Result<Value, _> = node_one.rpc().call(
        "withdraw-only-confirmed",
        json!({
            "destination": addr,
            "amount": "all",
        }),
    );
    assert!(withdraw.is_err());
    log::info!(target: "test_simple_withdraw_only_confirmed_3", "{:?}", withdraw);
    Ok(())
}

#[tokio::test(flavor = "multi_thread")]
#[ntest::timeout(560000)]
async fn test_simple_withdraw_only_confirmed_4() -> anyhow::Result<()> {
    init();

    let node_one = node!();
    let btc = node_one.btc();
    let node_two = node!(btc.clone());
    open_channel(&node_two, &node_one, false)?;

    let addr = node_one.rpc().newaddr(None)?.bech32.unwrap();
    fund_wallet(node_one.btc(), &addr, 8)?;
    wait_for_funds(&node_one)?;

    let addr = node_one.rpc().newaddr(None)?.bech32.unwrap();
    let withdraw: Result<Value, _> = node_one.rpc().call(
        "withdraw-only-confirmed",
        json!({
            "destination": addr,
            "amount": "all",
        }),
    );

    log::info!(target: "test_simple_withdraw_only_confirmed_4", "{:?}", withdraw);
    node_one.print_logs()?;
    assert!(withdraw.is_ok());
    Ok(())
}
