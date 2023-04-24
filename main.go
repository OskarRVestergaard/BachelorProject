package main

import (
	networkservice2 "github.com/OskarRVestergaard/BachelorProject/test/test_utils"
	"time"
)

func main() {

	noOfPeers := 2
	noOfMsgs := 2
	noOfNames := 2
	listOfPeers, pkList := networkservice2.SetupPeers(noOfPeers, noOfNames) //setup peer
	networkservice2.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)      //send msg

	time.Sleep(3000 * time.Millisecond)

	for i := 0; i < noOfPeers; i++ {
		listOfPeers[i].PrintLedger()
	}

	//pkFromByteArray, _ := x509.ParseECPrivateKey(pkByteArray)
	//print(&privateKey.PublicKey.X)
	//
	//noOfPeers := 2
	//noOfNames := 2
	////listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames) //setup peer
	//listOfPeers, _ := test_utils.SetupPeers(noOfPeers, noOfNames) //setup peer
	////println(pkList)
	//println(listOfPeers[1].IpPort)
	////println(listOfPeers[1].PublicToSecret[])
	//for k, v := range listOfPeers {
	//	fmt.Println(k, "value is", v)
	//}
	//println("finished s

}
