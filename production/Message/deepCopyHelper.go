package Message

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
)

func MakeDeepCopyOfTransaction(transaction models.SignedTransaction) (copyOfTransaction models.SignedTransaction) {
	oldSign := transaction.Signature
	signatureCopy := make([]byte, len(oldSign))
	copy(signatureCopy, oldSign)
	deepCopyTransaction := models.SignedTransaction{
		Id:        transaction.Id,
		From:      transaction.From,
		To:        transaction.To,
		Amount:    transaction.Amount,
		Signature: signatureCopy,
	}
	return deepCopyTransaction
}

func MakeDeepCopyOfBlock(block PoWblockchain.Block) (copyOfBlock PoWblockchain.Block) {

	oldSign := block.Signature
	signatureCopy := make([]byte, len(oldSign))
	copy(signatureCopy, oldSign)

	hashCopy := block.ParentHash //Array is by default copied by value

	oldTransactions := block.BlockData.Transactions
	transactionsCopy := make([]models.SignedTransaction, len(oldTransactions))
	for i, transaction := range oldTransactions {
		transactionsCopy[i] = MakeDeepCopyOfTransaction(transaction)
	}

	deepCopyBlock := PoWblockchain.Block{
		IsGenesis: block.IsGenesis,
		Vk:        block.Vk,
		Slot:      block.Slot,
		Draw:      MakeDeepCopyOfWinningParams(block.Draw),
		BlockData: PoWblockchain.BlockData{
			Hardness:     block.BlockData.Hardness,
			Transactions: transactionsCopy,
		},
		ParentHash: hashCopy,
		Signature:  signatureCopy,
	}
	return deepCopyBlock
}

func MakeDeepCopyOfWinningParams(params lottery_strategy.WinningLotteryParams) (copyOfParams lottery_strategy.WinningLotteryParams) {

	hashCopy := params.ParentHash //Array is by default copied by value

	deepCopyParams := lottery_strategy.WinningLotteryParams{
		Vk:         params.Vk,
		ParentHash: hashCopy,
		Counter:    params.Counter,
	}
	return deepCopyParams
}

func MakeDeepCopyOfMessage(msg Message) (copyOfMessage Message) {

	oldBlocks := msg.MessageBlocks
	blocksCopy := make([]PoWblockchain.Block, len(oldBlocks))
	for i, block := range oldBlocks {
		blocksCopy[i] = MakeDeepCopyOfBlock(block)
	}

	oldPeers := msg.PeerMap
	peersCopy := make(map[string]models.Void, len(oldPeers))
	for k, v := range oldPeers {
		peersCopy[k] = v
	}

	deepCopyMessage := Message{
		MessageType:       msg.MessageType,
		MessageSender:     msg.MessageSender,
		SignedTransaction: MakeDeepCopyOfTransaction(msg.SignedTransaction),
		MessageBlocks:     blocksCopy,
		PeerMap:           peersCopy,
	}
	return deepCopyMessage
}
