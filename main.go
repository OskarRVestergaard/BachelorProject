package main

import "example.com/packages/ledger"

/*
	type Ledger struct {
		Accounts map[string]int
		mutex    sync.Mutex
		TA       int
	}

	func (l *Ledger) updateLedger(from string, to string, value int) {
		l.mutex.Lock()

		l.Accounts[from] -= value
		l.Accounts[to] += value

		l.TA = l.TA + 1
		l.mutex.Unlock()
	}
*/
func main() {
	ledger.EstablishNetwork()
	//For debugging true Ledger
	//establishNetwork()
}

/*
func establishNetwork() {
	peersQt := 5
	tau := 10
	names := 5
	listOfPeers := make([]*peer.Peer, peersQt)

	var connectedPeers []string

	pkList := make([]string, names)

	for i := 0; i < peersQt; i++ {
		var p peer.Peer
		port := strconv.Itoa(18080 + i)
		listOfPeers[i] = &p
		p.RunPeer("127.0.0.1:" + port)

	}
	listOfPeers[0].Connect("Piplup is best water pokemon", 18079)
	connectedPeers = append(connectedPeers, listOfPeers[0].IpPort)
	time.Sleep(250 * time.Millisecond)
	for i := 1; i < peersQt; i++ {

		ipPort := connectedPeers[rand.Intn(len(connectedPeers))]
		ip := ipPort[0:(len(ipPort) - 6)]
		port := ipPort[len(ipPort)-5:]

		port2, _ := strconv.Atoi(port)
		listOfPeers[i].Connect(ip, port2)

		time.Sleep(250 * time.Millisecond)
	}
	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < names; i++ {
		pkList[i] = listOfPeers[i%peersQt].CreateAccount()
	}

	Led := new(Ledger)
	Led.TA = 0
	Led.Accounts = make(map[string]int)
	rand.Seed(time.Now().Unix())
	p1 := ""
	p2 := ""
	value := 0

	for j := 0; j < tau; j++ {
		for i := 0; i < peersQt; i++ {
			p1 = pkList[rand.Intn(names/peersQt)*peersQt+i]
			p2 = pkList[rand.Intn(names)]
			value = rand.Intn(100) + 1
			Led.updateLedger(p1, p2, value)
			go listOfPeers[i].FloodSignedTransaction(p1, p2, value)
		}
	}
	time.Sleep(3000 * time.Millisecond)
	for i := 0; i < peersQt; i++ {
		listOfPeers[i].PrintLedger()
	}

	l := Led.Accounts
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	println()
	println("----------------Should be--------------")
	print("                            ")
	for _, value := range keys {
		print("[" + strconv.Itoa(l[value]) + "]")
	}
	print(" with the amount of SignedTransactions: " + strconv.Itoa(Led.TA) + " ")
	println()

}

*/
