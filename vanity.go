package main 

import (
	"fmt"
	secp256k1 "github.com/haltingstate/secp256k1-go"
	"encoding/hex"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"github.com/gp2112/bitcoin-vanity-gen/base58"
	qrc "github.com/gp2112/bitcoin-vanity-gen/qrcode"
	"crypto/rand"
	"regexp"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"math"
)

func main() {
	fmt.Print("Name: ")
	var word string
	fmt.Scan(&word)
	vanity(word)
}

func vanity(word string) {
	unpattern := false
	priv, address := "", ""
	count := 0
	fmt.Println("Finding...")
	for unpattern != true {
		count++
		priv, address = getKeyAddress()
		pattern, _ :=  regexp.Match(word, []byte(address[0:10]))
		unpattern = pattern
	}
/////////////////////////////////////////////////////// Got Adress!! /////////////////////////////////////////////////////////////
	fmt.Printf("%d addresses runned!", count)																					//
	fmt.Printf("\nPrivateKey(Hex): %s\nPrivateKey(WIF): %s\nAddress: %s\n", priv, hex_wif(priv), address)	
	fmt.Printf("Address Balance: %f %s\n", getBalance(address), "BTC")							
	fmt.Println("*---------------------------------------------------------------------------------------------------------*")
	fmt.Println("  Warning: Always check if the private key matchs with the address created, before send it coins!")	
	fmt.Println("*---------------------------------------------------------------------------------------------------------*")
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	var answer string
	fmt.Printf("\nGenerate qr code for your address? (y/n): ")
	fmt.Scanf("%s", &answer)
	if answer == "y" {
		fmt.Println("")
		qrc.Makeqr(address)
	}
}

func getKeyAddress() (string, string) {
	seed, _ := random_seed()
	pub, priv := secp256k1.GenerateDeterministicKeyPair(seed)
	hash_pub := sha256.New()
	hash_pub.Write(pub)
	ripe_hash := ripemd160.New()
	ripe_hash.Write(hash_pub.Sum(nil))
	ext_ripe := "00"+hex.EncodeToString(ripe_hash.Sum(nil))
	ext_hash := sha256.New()
	extripe_decoded, _ := hex.DecodeString(ext_ripe)
	ext_hash.Write(extripe_decoded)
	hash2 := sha256.New()
	hash2.Write(ext_hash.Sum(nil))
	checksum := hex.EncodeToString(hash2.Sum(nil))[0:8]
	check_ripe := ext_ripe+checksum
	checkripe_decoded, _ := hex.DecodeString(check_ripe)
	address := base58.Encode(checkripe_decoded)
	return hex.EncodeToString(priv), address
}

func hex_wif(priv string) string {
	priv_mainnet, _ := hex.DecodeString("80"+priv+"01")
	hash_mainnet := sha256.New()
	hash_mainnet.Write(priv_mainnet)
	hash2 := sha256.New()
	hash2.Write(hash_mainnet.Sum(nil))
	checksum := hex.EncodeToString(hash2.Sum(nil))[0:8]
	check_mainnet, _ := hex.DecodeString(hex.EncodeToString(priv_mainnet)+checksum)
	return base58.Encode(check_mainnet)
}

func random_seed() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getBalance(address string) float64 {
	resp, _ := http.Get("https://blockchain.info/rawaddr/"+address)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var addrs Address
	err := json.Unmarshal(data, &addrs)
	if err != nil {
		fmt.Println(err)
	}
	total_sent := float64(addrs.Total_sent)*math.Pow10(-8)
	return total_sent
}

type Address struct {
	Hash160 string
	Address string
	N_tx int
	Total_received int
	Total_sent int
	Final_balance int
}
