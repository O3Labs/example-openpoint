package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	LocalEnv      = "local"
	StagingEnv    = "staging"
	ProductionEnv = "production"
)

type config struct {
	Name          string `json:"name"`
	Root          string `json:"root"`
	BaseURL       string `json:"baseUrl"`
	PublicBaseURL string `json:"publicBaseUrl"`
	Port          int    `json:"port"`
	Resource      struct {
		File struct {
			Path string `json:"path"`
		} `json:"file"`
	} `json:"resource"`
	Email struct {
		From string `json:"from"`
	} `json:"email"`
	Aws struct {
		Region string `json:"region"`
		S3     struct {
			CertBucket string `json:"certBucket"`
			Bucket     string `json:"bucket"`
			Region     string `json:"region"`
		} `json:"s3"`
		Cloudfront struct {
			Endpoint  string `json:"endpoint"`
			KeypairID string `json:"keypairID"`
			KeyPath   string `json:"keyPath"`
		} `json:"cloudfront"`
	} `json:"aws"`
	NEO struct {
		CoZEndpoint             string `json:"cozEndpoint"`
		SmartContractScriptHash string `json:"smartContractScriptHash"`
		FullNodeURL             string `json:"fullNodeURL"`
	} `json:"neo"`
}

var Env config

func Load(file string) {
	b, err := Asset("config/" + file + ".json")
	//b, err := ioutil.ReadFile("config/" + file + ".json")
	if err != nil {
		panic(fmt.Sprintf("Cannot load config file\n%+v\n", err))
		return
	}
	SetEnv(string(b))
	LoadEnv()
}

func LoadEnv() {
	osEnv := os.Getenv("env")
	err := json.Unmarshal([]byte(osEnv), &Env)
	if err != nil {
		panic("error parsing env")
	}
}

func SetEnv(env string) {
	os.Setenv("env", env)
}
