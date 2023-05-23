package peer_strategy

import (
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"time"
)

type PeerInterface interface {
	RunPeer(IpPort string, constants PeerConstants)
	CreateAccount() string
	GetBlockTree() interface{}
	StartMining(n int) error
	StopMining() error
	GetAddress() network.Address
	Connect(ip string, port int)
	FloodSignedTransaction(from string, to string, amount int)
	ActivatePeer(startTime time.Time, slotLength time.Duration)
}
