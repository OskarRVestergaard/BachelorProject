package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func main() {
	//ledger.EstablishNetwork()
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	print(&privateKey.PublicKey.X)
}
