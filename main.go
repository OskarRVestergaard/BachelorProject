package main

import (
	"github.com/OskarRVestergaard/BachelorProject/production/utils"
	"math/rand"
)

func main() {
	var t = createGraph(1, 4, 4)
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

func createGraph(seed int, n int, k int) [][]bool {
	if !utils.PowerOfTwo(n) {
		panic("n must be a power of two")
	}
	if !utils.PowerOfTwo(k) {
		panic("k must be a power of two")
	}

	edges := make([][]bool, n*k, n*k)
	for i := range edges {
		edges[i] = make([]bool, n*k, n*k)
	}

	var d = CalculateD()

	source := rand.NewSource(5)
	rando := rand.New(source)

	preds := make([][]int, n, n)
	for i := range preds {
		preds[i] = make([]int, d, d)
		for k := range preds[i] {
			preds[i][k] = -1
		}
		for j := 0; j < d; j++ {
			newNumber := false
			for !newNumber {
				random := rando.Intn(n)
				if !numberAlreadyChosen(random, preds[i]) {
					preds[i][j] = random
					newNumber = true
				}
			}
		}
	}

	for i := range preds {
		for j := range preds[i] {
			edges[preds[i][j]][n+i] = true
		}
	}

	for i := 0; i < len(edges)-n; i++ {
		for j := 0; j < len(edges)-n; j++ {
			if edges[i][j] { // == true
				edges[i+n][j+n] = true
			}
		}
	}

	return [][]bool{}
}
func CalculateD() int {
	return 2
}

func numberAlreadyChosen(n int, lst []int) bool {
	for _, b := range lst {
		if b == n {
			return true
		}
	}
	return false
}
