package payment

import "testing"

func TestMintTokensTo(t *testing.T) {
	invokeSmartContract("test_id", "AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn", 1, "apisit@o3.network")

	// invokeSmartContract("test_id", "AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn", 1, "apisit@o3.network")
	// txid = b487a1d1ca01d774e8f2470e36fa5d40af9c3d3a55cb411876f7706db1e5676e
}

func TestTransfer(t *testing.T) {}
