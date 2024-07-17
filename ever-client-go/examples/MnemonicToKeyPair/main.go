package main

import (
	"github.com/markgenuine/ever-client-go/domain"
	clientgw "github.com/markgenuine/ever-client-go/gateway/client"
	"github.com/markgenuine/ever-client-go/usecase/crypto"
	"log"
)

func main() {
	//ever, err := goever.NewEver("", []string{"https://gql-testnet.venom.foundation/"}, "")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer ever.Client.Destroy()

	HDPATH := "m/44'/396'/0'/0/0"
	params := &domain.ParamsOfMnemonicDeriveSignKeys{
		Phrase: "action inject penalty envelope rabbit element slim tornado dinner pizza off blood",
		Path:   HDPATH,
	}
	//keyPair, err := ever.Crypto.MnemonicDeriveSignKeys(params)
	//if err != nil {
	//	log.Fatal(err)
	//}

	gateway, _ := clientgw.NewClientGateway(domain.ClientConfig{})

	newCrypto := crypto.NewCrypto(domain.ClientConfig{}, gateway)
	keyPair, _ := newCrypto.MnemonicDeriveSignKeys(params)

	log.Print("PublicKey is: ", keyPair.Public)

	log.Print("SecretKey is: ", keyPair.Secret)
	//2024/07/10 15:11:03 PublicKey is: e85f61aaef0ea43afc14e08e6bd46c3b996974c495a881baccc58760f6349300
	//2024/07/10 1
	//5:11:03 SecretKey is: bb2903d025a330681e78f3bcb248d7d89b861f3e8a480eb74438ec0299319f7a

}
