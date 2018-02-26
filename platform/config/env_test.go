package config_test

import (
	"log"
	"testing"

	"github.com/o3labs/openpoint/platform/config"
)

func TestEnvConfig(t *testing.T) {

	config.SetEnv(`
	{
	"baseUrl": "https://staging.getchannel.co",
	"database":{
		"name":"channel",
		"host":"",
		"port":"",
		"user":"",
		"password":""
	},
	"email": {
		"from": "Channel <channel@mogohichi.com>"
	},
	"aws": {
		"region": "us-west-2",
		"s3": {
			"bucket": "co.getchannel.files",
			"region": "us-west-2"
		},
		"cloudfront": {
			"endpoint": "https://files.getchannel.co",
			"keypairID": "APKAIWZCPMHPFRFFHV5A",
			"keyPath": "./keys/pk-APKAIWZCPMHPFRFFHV5A.pem"
		}
	}
}
			`)
	config.LoadEnv()
	log.Printf("%+v", config.Env)
}
