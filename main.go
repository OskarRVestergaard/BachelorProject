package main

import (
	"math/rand"
)

func main() {
	var t = createGraph()
	print(t)
	//noOfPeers := 2
	//noOfMsgs := 2
	//noOfNames := 2
	//listOfPeers, pkList := networkservice2.SetupPeers(noOfPeers, noOfNames) //setup peer
	//networkservice2.SendMsgs(noOfMsgs, noOfPeers, listOfPeers, pkList)      //send msg
	//
	//time.Sleep(3000 * time.Millisecond)
	//
	//for i := 0; i < noOfPeers; i++ {
	//	listOfPeers[i].PrintLedger()
	//}

	//pkFromByteArray, _ := x509.ParseECPrivateKey(pkByteArray)
	//print(&privateKey.PublicKey.X)
	//
	//noOfPeers := 2
	//noOfNames := 2
	////listOfPeers, pkList := test_utils.SetupPeers(noOfPeers, noOfNames) //setup peer
	//listOfPeers, _ := test_utils.SetupPeers(noOfPeers, noOfNames) //setup peer
	////println(pkList)
	//println(listOfPeers[1].IpPort)
	////println(listOfPeers[1].publicToSecret[])
	//for k, v := range listOfPeers {
	//	fmt.Println(k, "value is", v)
	//}
	//println("finished s

}

func createGraph() [][]bool {
	var bipartiteExpanders = 5
	var sourceNodes = 5
	edges := make([][]bool, sourceNodes, sourceNodes)
	for i := range edges {
		edges[i] = make([]bool, (bipartiteExpanders * sourceNodes), (bipartiteExpanders * sourceNodes))
	}
	//edges[0][1] = true
	//edges[0][2] = true
	//edges[0][3] = true
	//edges[1][3] = true
	//edges[1][4] = true
	//edges[2][4] = true
	//edges[2][5] = true
	for i1 := 0; i1 < sourceNodes; i1++ {
		for i2 := 0; i2 < bipartiteExpanders; i2++ {
			edges[i1][1] = true
		}
	}

	var d = CalculateD()
	for i := 0; i < bipartiteExpanders; i++ {
		for j := 0; j <= sourceNodes; j++ {
			for t := 0; t < d; t++ {
				var random = rand.Intn(sourceNodes - 1)
				edges[i][random] = true
			}
			//var random = rand.Intn(sourceNodes - 1)
			//edges[i][random] = true
		}
		//rand.Intn(100)
		//edges[]
	}
	return edges
}
func CalculateD() int {
	return 1

}
