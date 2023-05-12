package Message

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/models/SpaceMintBlockchain"
)

type Message struct {
	MessageType       string
	MessageSender     string
	SignedTransaction models.SignedTransaction
	PoWMessageBlocks  []PoWblockchain.Block
	SpaceMintBlocks   []SpaceMintblockchain.Block
	PeerMap           map[string]models.Void
}
