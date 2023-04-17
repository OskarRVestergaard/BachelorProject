package networkservice

import (
	"github.com/OskarRVestergaard/BachelorProject/production/Peer"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func SetupPeers(noOfPeers int, noOfNames int) ([]*Peer.Peer, []string) {

	listOfPeers := make([]*Peer.Peer, noOfPeers)

	var connectedPeers []string

	pkList := make([]string, noOfNames)

	for i := 0; i < noOfPeers; i++ {
		var p Peer.Peer
		freePort, _ := GetFreePort()
		port := strconv.Itoa(freePort)
		listOfPeers[i] = &p
		p.RunPeer("127.0.0.1:" + port)
		// TODO maybe go p.RunPeer
	}
	listOfPeers[0].Connect("Piplup is best water pokemon", 18079)
	connectedPeers = append(connectedPeers, listOfPeers[0].IpPort)
	time.Sleep(250 * time.Millisecond)
	for i := 1; i < noOfPeers; i++ {

		ipPort := connectedPeers[rand.Intn(len(connectedPeers))]
		ip := ipPort[0:(len(ipPort) - 6)]
		port := ipPort[len(ipPort)-5:]

		port2, _ := strconv.Atoi(port)
		listOfPeers[i].Connect(ip, port2)

		time.Sleep(250 * time.Millisecond)
	}
	println("finished setting up connections")
	println("Starting simulation")

	for i := 0; i < noOfNames; i++ {
		pkList[i] = listOfPeers[i%noOfPeers].CreateAccount()
	}
	time.Sleep(250 * time.Millisecond)
	return listOfPeers, pkList
}

// TODO move to util
func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
