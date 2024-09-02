package wallet

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const walletFile = "./tmp/wallets.json"

type Wallets struct {
	Wallets map[string]*Wallet
}

type SerializableWallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	if err != nil {
		return nil, err
	}

	return &wallets, nil
}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.Open(walletFile)
	if err != nil {
		return err
	}

	var serializedWallet map[string]*SerializableWallet

	decoder := json.NewDecoder(fileContent)
	err = decoder.Decode(&serializedWallet)
	if err != nil {
		return err
	}

	for address, serializedWallet := range serializedWallet {
		privKey, err := x509.ParseECPrivateKey(serializedWallet.PrivateKey)
		if err != nil {
			return err
		}

		ws.Wallets[address] = &Wallet{
			PrivateKey: *privKey,
			PublicKey:  serializedWallet.PublicKey,
		}
	}

	return nil
}

func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	serializedWallets := make(map[string]*SerializableWallet)

	for _, wallet := range ws.Wallets {
		privKeyBytes, err := x509.MarshalECPrivateKey(&wallet.PrivateKey)
		if err != nil {
			log.Panic(err)
		}

		serializedWallet := &SerializableWallet{
			PrivateKey: privKeyBytes,
			PublicKey:  wallet.PublicKey,
		}

		address := fmt.Sprintf("%s", wallet.Address())

		serializedWallets[address] = serializedWallet
	}

	jsonEncoder := json.NewEncoder(&content)
	err := jsonEncoder.Encode(serializedWallets)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0777)
	if err != nil {
		log.Panic(err)
	}
}
