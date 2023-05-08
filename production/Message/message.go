package Message

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
)

type Message struct {
	MessageType       string
	MessageSender     string
	SignedTransaction models.SignedTransaction
	MessageBlocks     []PoWblockchain.Block
	PeerMap           map[string]models.Void
}
