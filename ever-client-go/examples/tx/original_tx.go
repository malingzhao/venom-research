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
	value, err := venom.Client.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Version bindings is: ", value.Version)

	//add own mnemonic
	HDPATH := "m/44'/396'/0'/0/0"
	params := &domain.ParamsOfMnemonicDeriveSignKeys{
		Phrase: "",
		Path:   HDPATH,
	}
	keyPair, err := venom.Crypto.MnemonicDeriveSignKeys(params)
	fmt.Println("the priv is ", keyPair.Secret)
	//keyPair := &domain.KeyPair{
	//	Public: "a9caf7dd7d2a79057a3b2885f6e22edacf2c4ba8e078381ab3ca4d882c0c43a5",
	//	Secret: "fcbe5566e95f446cbed70f7b17c868dcce3cc0a8e25ef63e451762e594a2f07f",
	//}
	if err != nil {
		log.Fatal(err)
	}

	// generate init data, it's use for calculate address by pubkey or for deploy new wallet
	pubKey := "0x" + keyPair.Public
	builder := []*domain.BuilderOp{
		domain.NewBuilderOp(domain.BuilderOpInteger{Size: 256, Value: pubKey}),
		domain.NewBuilderOp(domain.BuilderOpInteger{Size: 64, Value: 0})}
	dataParams := &domain.ParamsOfEncodeBoc{Builder: builder}
	walletData, err := venom.Boc.EncodeBoc(dataParams)

	walletCode := "te6cckEBBgEA/AABFP8A9KQT9LzyyAsBAgEgAgMABNIwAubycdcBAcAA8nqDCNcY7UTQgwfXAdcLP8j4KM8WI88WyfkAA3HXAQHDAJqDB9cBURO68uBk3oBA1wGAINcBgCDXAVQWdfkQ8qj4I7vyeWa++COBBwiggQPoqFIgvLHydAIgghBM7mRsuuMPAcjL/8s/ye1UBAUAmDAC10zQ+kCDBtcBcdcBeNcB10z4AHCAEASqAhSxyMsFUAXPFlAD+gLLaSLQIc8xIddJoIQJuZgzcAHLAFjPFpcwcQHLABLM4skB+wAAPoIQFp4+EbqOEfgAApMg10qXeNcB1AL7AOjRkzLyPOI+zYS/"

	psi := &domain.ParamsOfEncodeStateInit{Code: walletCode, Data: walletData.Boc}
	stateInit, err := venom.Boc.EncodeStateInit(psi)

	tvcHash, err := venom.Boc.GetBocHash(&domain.ParamsOfGetBocHash{Boc: stateInit.StateInit})
	address := "0:" + tvcHash.Hash
	log.Print("address is ", address)

	// Send external message
	fileAbi, err := os.Open("./contracts/Wallet.abi.json")

	if err != nil {
		log.Fatalf("Can't open file %s, error: %s", "..contracts/Wallets.abi.json", err)
	}

	byteAbi, err := io.ReadAll(fileAbi)
	nn := &domain.AbiContract{}
	err = json.Unmarshal(byteAbi, &nn)
	walletAbi := domain.NewAbiContract(nn)

	dst, _ := json.Marshal("0:7671cdd006d85a989395beba848848ae5f505fc37ff8f2c4402ce57f4d72fb49")

	// signID := 1
	signID, err := venom.Net.GetSignatureID()
	log.Print("signID = ", &signID.SignatureID)

	// use DeploySet only if need to deploy wallet
	unsignedMessage, err := venom.Abi.EncodeMessage(&domain.ParamsOfEncodeMessage{
		Abi:     walletAbi,
		Signer:  domain.NewSigner(domain.SignerExternal{PublicKey: keyPair.Public}),
		Address: address,
		CallSet: &domain.CallSet{
			FunctionName: "sendTransaction",
			Input:        json.RawMessage(fmt.Sprintf(`{"dest": %s,"value": 10, "bounce": false, "flags":3, "payload":""}`, dst)),
		},
		SignatureID: signID.SignatureID,
	})

	log.Print("err", err)
	log.Print("Message address", unsignedMessage.Address)

	signature, err := venom.Crypto.Sign(&domain.ParamsOfSign{Unsigned: unsignedMessage.DataToSign, Keys: keyPair})
	log.Print("signature ", signature)

	signed, err := venom.Abi.AttachSignature(&domain.ParamsOfAttachSignature{Abi: walletAbi, PublicKey: keyPair.Public, Message: unsignedMessage.Message, Signature: signature.Signature})
	log.Print("err ", err)
	log.Print("signed ", signed.Message)

	// https://github.com/markgenuine/ever-client-go/blob/master/usecase/processing/processing_test.go#L245
	shardBlockID, err := venom.Processing.SendMessage(&domain.ParamsOfSendMessage{Message: signed.Message, SendEvents: false}, nil)
	log.Print("shardBlockID ", shardBlockID)

	result, err := venom.Processing.WaitForTransaction(&domain.ParamsOfWaitForTransaction{Message: signed.Message, ShardBlockID: shardBlockID.ShardBlockID, SendEvents: false, Abi: walletAbi}, nil)
	log.Print(result.Fees)
}
