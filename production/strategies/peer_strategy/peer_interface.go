package peer_strategy

import (
	"github.com/OskarRVestergaard/BachelorProject/production/network"
	"time"
)

type PeerInterface interface {
	RunPeer(IpPort string, startTime time.Time)
	CreateAccount() string
	GetBlockTree() interface{}
	StartMining() error
	StopMining() error
	GetAddress() network.Address
	Connect(ip string, port int)
	FloodSignedTransaction(from string, to string, amount int)
}
