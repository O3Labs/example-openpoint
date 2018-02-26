package payment

import (
	"testing"

	"github.com/o3labs/openpoint/platform/config"
)

func TestMintTokensTo(t *testing.T) {
	config.Load("local")
	invokeSmartContract("test_id", "AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn", 1000, "apisit@o3.network")

	// invokeSmartContract("test_id", "AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn", 1, "apisit@o3.network")
	// txid = b487a1d1ca01d774e8f2470e36fa5d40af9c3d3a55cb411876f7706db1e5676e
}

func TestTransfer(t *testing.T) {

	config.Load("local")

	invokeSmartContractTransfer("test_id", "AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn", "AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y", 1, "apisit@o3.network")

}
