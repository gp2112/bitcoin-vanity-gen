package main 

import (
	"fmt"
	secp256k1 "github.com/haltingstate/secp256k1-go"
	"encoding/hex"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"personal/base58"
	"crypto/rand"
	"regexp"
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
	fmt.Printf("%d addresses runned!", count)
	fmt.Printf("\nPrivateKey(Hex): %s\nPrivateKey(WIF): %s\nAddress: %s\n", priv, hex_wif(priv), address)
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
	priv_mainnet, _ := hex.DecodeString("80"+priv)
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
