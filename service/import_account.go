package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func ImportAccount(eosAccountName string) {
	//log.Printf("ENCRYPTED: %s\n", encrypted)
	// file write
	err := ioutil.WriteFile("default.wallet", []byte(eosAccountName), 0)
	if err != nil {
		panic(err)
	}

	err = os.Chmod("default.wallet", 0777)
	if err != nil {
		fmt.Println(err)
	}
}

func GetAccountFromWallet() string {
	// file read
	tempBytes, err := ioutil.ReadFile("default.wallet")
	if err != nil {
		log.Println("DoesNotExist_Wallet")
		return "DoesNotExist_Wallet"
	}

	return string(tempBytes[:])
}