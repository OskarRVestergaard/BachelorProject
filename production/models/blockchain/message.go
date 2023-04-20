package blockchain

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
)

type Message struct {
	MessageType       string
	MessageSender     string
	SignedTransaction SignedTransaction
	MessageBlocks     []Block
	PeerMap           map[string]models.Void
}
