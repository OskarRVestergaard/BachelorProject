package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
)

func main() {

	//ledger.EstablishNetwork()
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		panic(err)
	}
	pkByteArray, _ := x509.MarshalECPrivateKey(privateKey)
	pkString := string(pkByteArray)

	pkByteArray2 := []byte(pkString)
	pkFromByteArray, _ := x509.ParseECPrivateKey(pkByteArray2)

	fmt.Println(privateKey)
	fmt.Println(pkFromByteArray)

	//pkFromByteArray, _ := x509.ParseECPrivateKey(pkByteArray)
	//print(&privateKey.PublicKey.X)
	//
	//noOfPeers := 2
	//noOfNames := 2
	////listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	//listOfPeers, _ := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	////println(pkList)
	//println(listOfPeers[1].IpPort)
	////println(listOfPeers[1].PublicToSecret[])
	//for k, v := range listOfPeers {
	//	fmt.Println(k, "value is", v)
	//}
	//println("finished s

}
