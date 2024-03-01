//! Core Lightning Integration testing
//!
//! Author: Vincenzo Palazzo <vincenzopalazzo@member.fsf.org>
use std::sync::Once;

use clightning_testing::cln;
use serde_json::json;

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
