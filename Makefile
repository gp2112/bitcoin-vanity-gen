
requirements:
		go get github.com/qpliu/qrencode-go/qrencode
		go get github.com/haltingstate/secp256k1-go
		go get github.com/btcsuite/btcutil/base58

run:
	go run qrcode.go vanity.go