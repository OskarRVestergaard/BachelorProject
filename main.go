package main

import (
	"example.com/packages/service"
	"fmt"
)

func main() {
	//ledger.EstablishNetwork()
	//privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//print(&privateKey.PublicKey.X)

	noOfPeers := 2
	noOfNames := 2
	//listOfPeers, pkList := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	listOfPeers, _ := service.SetupPeers(noOfPeers, noOfNames) //setup peer
	//println(pkList)
	println(listOfPeers[1].IpPort)
	//println(listOfPeers[1].PublicToSecret[])
	for k, v := range listOfPeers {
		fmt.Println(k, "value is", v)
	}
	//println("finished s

}
