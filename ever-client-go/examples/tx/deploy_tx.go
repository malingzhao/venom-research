package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/markgenuine/ever-client-go/domain"

	govenom "github.com/markgenuine/ever-client-go"
)

func main() {

	// init connection
	venom, err := govenom.NewEver("", []string{"https://gql.venom.foundation/graphql"}, "")
	if err != nil {
		log.Fatal(err)
	}

	defer venom.Client.Destroy()

	// DON't store seed in code, it for demo only
	HDPATH := "m/44'/396'/0'/0/0"
	params := &domain.ParamsOfMnemonicDeriveSignKeys{
		Phrase: "",
		Path:   HDPATH,
	}
	keyPair, err := venom.Crypto.MnemonicDeriveSignKeys(params)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("PublicKey is: ", keyPair.Public)

	// get state init
	walletCode := "te6cckEBBgEA/AABFP8A9KQT9LzyyAsBAgEgAgMABNIwAubycdcBAcAA8nqDCNcY7UTQgwfXAdcLP8j4KM8WI88WyfkAA3HXAQHDAJqDB9cBURO68uBk3oBA1wGAINcBgCDXAVQWdfkQ8qj4I7vyeWa++COBBwiggQPoqFIgvLHydAIgghBM7mRsuuMPAcjL/8s/ye1UBAUAmDAC10zQ+kCDBtcBcdcBeNcB10z4AHCAEASqAhSxyMsFUAXPFlAD+gLLaSLQIc8xIddJoIQJuZgzcAHLAFjPFpcwcQHLABLM4skB+wAAPoIQFp4+EbqOEfgAApMg10qXeNcB1AL7AOjRkzLyPOI+zYS/"
	pubKey := "0x" + keyPair.Public
	builder := []*domain.BuilderOp{
		domain.NewBuilderOp(domain.BuilderOpInteger{Size: 256, Value: pubKey}),
		domain.NewBuilderOp(domain.BuilderOpInteger{Size: 64, Value: 0})}
	dataParams := &domain.ParamsOfEncodeBoc{Builder: builder}
	walletData, err := venom.Boc.EncodeBoc(dataParams)

	psi := &domain.ParamsOfEncodeStateInit{Code: walletCode, Data: walletData.Boc}
	stateInit, err := venom.Boc.EncodeStateInit(psi)

	tvcHash, err := venom.Boc.GetBocHash(&domain.ParamsOfGetBocHash{Boc: stateInit.StateInit})
	address := "0:" + tvcHash.Hash
	log.Print("address is ", address)

	// Get ABI
	fileAbi, err := os.Open("./contracts/Wallet.abi.json")

	if err != nil {
		log.Fatalf("Can't open file %s, error: %s", "..contracts/Wallets.abi.json", err)
	}

	byteAbi, err := io.ReadAll(fileAbi)
	nn := &domain.AbiContract{}
	err = json.Unmarshal(byteAbi, &nn)
	walletAbi := domain.NewAbiContract(nn)

	dst, _ := json.Marshal(address)
	signID, err := venom.Net.GetSignatureID()
	log.Print("signID = ", &signID.SignatureID)

	// gen msg body
	body, err := venom.Abi.EncodeMessageBody(&domain.ParamsOfEncodeMessageBody{
		Abi:         walletAbi,
		IsInternal:  false,
		Signer:      domain.NewSigner(domain.SignerKeys{Keys: keyPair}),
		SignatureID: signID.SignatureID,
		Address:     address,
		CallSet: &domain.CallSet{
			FunctionName: "sendTransaction",
			Input:        json.RawMessage(fmt.Sprintf(`{"dest": %s,"value": 100000000, "bounce": false, "flags":3, "payload":""}`, dst)),
		},
	})

	log.Print("err", err)
	log.Print("body: ", body.Body)

	extMessage, err := venom.Boc.EncodeExternalInMessage(&domain.ParamsOfEncodeExternalInMessage{
		Dst:  address,
		Init: stateInit.StateInit,
		Body: body.Body,
	})
	log.Print("err", err)
	log.Print("Message: ", extMessage.Message)
	log.Print("Message hash: ", extMessage.MessageID)

	// https://github.com/markgenuine/ever-client-go/blob/master/usecase/processing/processing_test.go#L245
	shardBlockID, err := venom.Processing.SendMessage(&domain.ParamsOfSendMessage{Message: extMessage.Message, SendEvents: false}, nil)
	log.Print("err", err)
	log.Print("shardBlockID ", shardBlockID)

	result, err := venom.Processing.WaitForTransaction(&domain.ParamsOfWaitForTransaction{Message: extMessage.Message, ShardBlockID: shardBlockID.ShardBlockID, SendEvents: false, Abi: walletAbi}, nil)
	log.Print(result.Fees)
}
