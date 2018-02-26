package payment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/coz"
	"github.com/o3labs/neo-utils/neoutils/neorpc"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
	"github.com/o3labs/openpoint/platform/config"
	"github.com/o3labs/openpoint/platform/errors"
	"github.com/o3labs/openpoint/platform/models"
	"github.com/o3labs/openpoint/platform/services"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type ControllerInterface interface {
	ProcessPaymentToken(w http.ResponseWriter, r *http.Request)
}

type Controller struct{}

func NewController() ControllerInterface {
	return &Controller{}
}

type StripeTokenData struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Card   struct {
		ID                string      `json:"id"`
		Object            string      `json:"object"`
		AddressCity       interface{} `json:"address_city"`
		AddressCountry    interface{} `json:"address_country"`
		AddressLine1      interface{} `json:"address_line1"`
		AddressLine1Check interface{} `json:"address_line1_check"`
		AddressLine2      interface{} `json:"address_line2"`
		AddressState      interface{} `json:"address_state"`
		AddressZip        string      `json:"address_zip"`
		AddressZipCheck   string      `json:"address_zip_check"`
		Brand             string      `json:"brand"`
		Country           string      `json:"country"`
		CvcCheck          string      `json:"cvc_check"`
		DynamicLast4      interface{} `json:"dynamic_last4"`
		ExpMonth          int         `json:"exp_month"`
		ExpYear           int         `json:"exp_year"`
		Funding           string      `json:"funding"`
		Last4             string      `json:"last4"`
		Metadata          struct {
		} `json:"metadata"`
		Name               interface{} `json:"name"`
		TokenizationMethod interface{} `json:"tokenization_method"`
	} `json:"card"`
	ClientIP string `json:"client_ip"`
	Created  int    `json:"created"`
	Livemode bool   `json:"livemode"`
	Type     string `json:"type"`
	Used     bool   `json:"used"`
}

//PRIVATE NET ADDRESS
//address AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn
//WIF L1AUSgSr36KzXynHGVveTzpcd6mfhCoXX7Njku7kBawTvpvE87JS
//contract owner is AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y

func invokeSmartContract(stripeChargeID string, toAddress string, numberOfTokens int, toEmail string) {
	var openPointSmartContract = neoutils.UseSmartContract(config.Env.NEO.SmartContractScriptHash)

	//this is a private net private key
	//private key must be moved out from the code here.
	//it could be set as env variable on the server
	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	to := smartcontract.ParseNEOAddress(toAddress)
	if to == nil {
		//invalid neo address
		return
	}
	args := []interface{}{to, numberOfTokens}
	log.Printf("neo config %+v", config.Env.NEO)
	client := coz.NewClient(config.Env.NEO.CoZEndpoint)
	//get unspent of the contract owner
	unspentCoz, err := client.GetUnspentByAddress(privateNetwallet.Address)
	if err != nil {
		log.Printf("error gen tx %v", err)
		return
	}
	gasBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	for _, v := range unspentCoz.GAS.Unspent {
		gasTX1 := smartcontract.UTXO{
			Index: v.Index,
			TXID:  v.Txid,
			Value: v.Value,
		}
		log.Printf("utxo value = %.8f", v.Value)
		gasBalance.UTXOs = append(gasBalance.UTXOs, gasTX1)
	}

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}

	unspent.Assets[smartcontract.GAS] = &gasBalance

	transactionID := stripeChargeID
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Description] = []byte(transactionID)
	tx, err := openPointSmartContract.GenerateInvokeFunctionRawTransaction(*privateNetwallet, unspent, attributes, "mintTokensTo", args)
	if err != nil {
		log.Printf("error gen tx %v", err)
		return
	}
	log.Printf("%x", tx)

	neoclient := neorpc.NewClient(config.Env.NEO.FullNodeURL)
	result := neoclient.SendRawTransaction(fmt.Sprintf("%x", tx))
	log.Printf("%+v", result)
	if result.Result == true {
		//send email here
		body := fmt.Sprintf("We sent %v tokens to the NEO Address %v", numberOfTokens, toAddress)
		email := models.Email{
			From:    "apisit@o3.network",
			To:      toEmail,
			Subject: fmt.Sprintf("TEST: Purchased of %v tokens from O3 Labs", numberOfTokens),
			HTML:    body,
		}
		errEmail := services.NewEmailService().SendEmail(email)
		if errEmail != nil {
			log.Printf("err email = %v", errEmail)
		}
	}
}

func invokeSmartContractTransfer(stripeChargeID string, fromAddress string, toAddress string, numberOfTokens int, toEmail string) {
	var openPointSmartContract = neoutils.UseSmartContract(config.Env.NEO.SmartContractScriptHash)

	//this is a private net private key
	//private key must be moved out from the code here.
	//it could be set as env variable on the server
	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	to := smartcontract.ParseNEOAddress(toAddress)
	if to == nil {
		//invalid neo address
		return
	}

	from := smartcontract.ParseNEOAddress(fromAddress)
	if from == nil {
		//invalid neo address
		return
	}
	args := []interface{}{from, to, numberOfTokens}
	log.Printf("neo config %+v", config.Env.NEO)
	client := coz.NewClient(config.Env.NEO.CoZEndpoint)
	//get unspent of the contract owner
	unspentCoz, err := client.GetUnspentByAddress(privateNetwallet.Address)
	if err != nil {
		log.Printf("error gen tx %v", err)
		return
	}
	gasBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	for _, v := range unspentCoz.GAS.Unspent {
		gasTX1 := smartcontract.UTXO{
			Index: v.Index,
			TXID:  v.Txid,
			Value: v.Value,
		}
		log.Printf("utxo value = %.8f", v.Value)
		gasBalance.UTXOs = append(gasBalance.UTXOs, gasTX1)
	}

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}

	unspent.Assets[smartcontract.GAS] = &gasBalance

	transactionID := stripeChargeID
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Description] = []byte(transactionID)
	tx, err := openPointSmartContract.GenerateInvokeFunctionRawTransaction(*privateNetwallet, unspent, attributes, "transfer", args)
	if err != nil {
		log.Printf("error gen tx %v", err)
		return
	}
	log.Printf("%x", tx)
	neoclient := neorpc.NewClient(config.Env.NEO.FullNodeURL)
	result := neoclient.SendRawTransaction(fmt.Sprintf("%x", tx))
	log.Printf("%+v", result)
	if result.Result == true {
		//send email here
		body := fmt.Sprintf("We sent %v tokens to the NEO Address %v", numberOfTokens, toAddress)
		email := models.Email{
			From:    "apisit@o3.network",
			To:      toEmail,
			Subject: fmt.Sprintf("TEST: Purchased of %v tokens from O3 Labs", numberOfTokens),
			HTML:    body,
		}
		errEmail := services.NewEmailService().SendEmail(email)
		if errEmail != nil {
			log.Printf("err email = %v", errEmail)
		}
	}
}

func (c *Controller) ProcessPaymentToken(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Email          string          `json:"email"`
		NEOAddress     string          `json:"neoAddress"`
		NumberOfTokens string          `json:"numberOfTokens"`
		PaymentToken   StripeTokenData `json:"paymentToken"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		//handle error
		errors.BadRequest("%v", err).Write(w)
		return
	}
	log.Printf("data =%+v", data)
	tokenAmount, err := strconv.Atoi(data.NumberOfTokens)
	if err != nil {
		errors.BadRequest("%v", "Token amount must be more a number").Write(w)
		return
	}
	if tokenAmount <= 0 {
		errors.BadRequest("%v", "Token amount must be more than zero").Write(w)
		return
	}
	//let make it simple that 1 point = 1 cent
	amountToCharge := uint64(tokenAmount)
	stripe.Key = "sk_test_t6OGPA3a3ZhScXoZNeueqN72"

	chargeParams := &stripe.ChargeParams{
		Amount:   amountToCharge,
		Currency: "usd",
		Desc:     fmt.Sprintf("Charge of %v tokens to %v", tokenAmount, data.NEOAddress),
	}
	chargeParams.SetSource(data.PaymentToken.ID) // obtained with Stripe.js
	ch, errStripe := charge.New(chargeParams)
	if errStripe != nil {
		if e, ok := errStripe.(*stripe.Error); ok {
			log.Printf("stripe charge response error %v", e.Msg)
			errors.BadRequest("%v", e.Msg).Write(w)
			return
		}
		errors.BadRequest("%v", errStripe).Write(w)
		return
	}
	log.Printf("%v", ch)

	//BEST is to put this in the messaging queue and process later with retry and backoff mechanism
	invokeSmartContract(ch.ID, data.NEOAddress, tokenAmount, data.Email)

	json.NewEncoder(w).Encode(models.EmptySuccessResponse())
}
