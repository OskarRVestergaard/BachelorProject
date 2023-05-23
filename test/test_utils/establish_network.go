package test_utils

import (
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/peer_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/peer_strategy/PowPeer"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/peer_strategy/SpacemintPeer"
	"net"
	"strconv"
	"time"
)

func SetupPeers(noOfPeers int, noOfNames int, useProofOfSpace bool, blockChainConstants peer_strategy.PeerConstants) ([]peer_strategy.PeerInterface, []string) {

	startTime := time.Now()
	listOfPeers := make([]peer_strategy.PeerInterface, noOfPeers)

	pkList := make([]string, noOfNames)

	for i := 0; i < noOfPeers; i++ {
		var p peer_strategy.PeerInterface
		if useProofOfSpace {
			p = &SpacemintPeer.PoSpacePeer{}
			p.ActivatePeer(startTime, blockChainConstants.SlotLength) //TODO PLACE AFTER START MINING
		} else {
			p = &PowPeer.PoWPeer{}
			p.ActivatePeer(startTime, blockChainConstants.SlotLength) //TODO PLACE AFTER START MINING
		}
		freePort, _ := GetFreePort()
		port := strconv.Itoa(freePort)
		listOfPeers[i] = p
		p.RunPeer("127.0.0.1:"+port, blockChainConstants)
	}
	time.Sleep(150 * time.Millisecond)
	for i := 0; i < noOfPeers; i++ {
		for j := 0; j < noOfPeers; j++ {
			if j == i {
				continue
			}
			addr := listOfPeers[j].GetAddress()
			listOfPeers[i].Connect(addr.Ip, addr.Port)
		}
	}

	for i := 0; i < noOfNames; i++ {
		pkList[i] = listOfPeers[i%noOfPeers].CreateAccount()
	}
	time.Sleep(500 * time.Millisecond)
	println("finished setting up connections")
	println("Starting simulation")

	return listOfPeers, pkList
}

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
