package Message

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/models/SpaceMintBlockchain"
)

type Message struct {
	MessageType       string
	MessageSender     string
	SignedTransaction models.SignedPaymentTransaction
	SpaceCommitment   SpaceMintBlockchain.SpaceCommitment
	PoWMessageBlocks  []PoWblockchain.Block
	SpaceMintBlocks   []SpaceMintBlockchain.Block
	PeerMap           map[string]models.Void
}
