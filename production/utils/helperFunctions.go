package utils

import (
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/blockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/signature_strategy"
)

func GetSomeKey[t comparable](m map[t]t) t {
	for k := range m {
		return k
	}
	panic("Cant get key from an empty map!")
}

func TransactionHasCorrectSignature(signatureStrategy signature_strategy.SignatureInterface, signedTrans blockchain.SignedTransaction) bool {
	transByteArray := signedTrans.ToByteArrayWithoutSign()
	hashedMessage := sha256.HashByteArrayToByteArray(transByteArray)
	publicKey := signedTrans.From
	signature := signedTrans.Signature
	result := signatureStrategy.Verify(publicKey, hashedMessage, signature)
	return result
}

func MakeDeepCopyOfTransaction(transaction blockchain.SignedTransaction) (copyOfTransaction blockchain.SignedTransaction) {
	oldSign := transaction.Signature
	signatureCopy := make([]byte, len(oldSign))
	copy(signatureCopy, oldSign)
	deepCopyTransaction := blockchain.SignedTransaction{
		Id:        transaction.Id,
		From:      transaction.From,
		To:        transaction.To,
		Amount:    transaction.Amount,
		Signature: signatureCopy,
	}
	return deepCopyTransaction
}

func MakeDeepCopyOfBlock(block blockchain.Block) (copyOfBlock blockchain.Block) {

	oldSign := block.Signature
	signatureCopy := make([]byte, len(oldSign))
	copy(signatureCopy, oldSign)

	oldHash := block.ParentHash
	hashCopy := make([]byte, len(oldHash))
	copy(hashCopy, oldHash)

	oldTransactions := block.BlockData.Transactions
	transactionsCopy := make([]blockchain.SignedTransaction, len(oldTransactions))
	for i, transaction := range oldTransactions {
		transactionsCopy[i] = MakeDeepCopyOfTransaction(transaction)
	}

	deepCopyBlock := blockchain.Block{
		IsGenesis: block.IsGenesis,
		Vk:        block.Vk,
		Slot:      block.Slot,
		Draw:      MakeDeepCopyOfWinningParams(block.Draw),
		BlockData: blockchain.BlockData{
			Hardness:     block.BlockData.Hardness,
			Transactions: transactionsCopy,
		},
		ParentHash: hashCopy,
		Signature:  signatureCopy,
	}
	return deepCopyBlock
}

func MakeDeepCopyOfWinningParams(params lottery_strategy.WinningLotteryParams) (copyOfParams lottery_strategy.WinningLotteryParams) {

	oldHash := params.ParentHash
	hashCopy := make([]byte, len(oldHash))
	copy(hashCopy, oldHash)

	deepCopyParams := lottery_strategy.WinningLotteryParams{
		Vk:         params.Vk,
		ParentHash: hashCopy,
		Counter:    params.Counter,
	}
	return deepCopyParams
}

func MakeDeepCopyOfMessage(msg blockchain.Message) (copyOfMessage blockchain.Message) {

	oldBlocks := msg.MessageBlocks
	blocksCopy := make([]blockchain.Block, len(oldBlocks))
	for i, block := range oldBlocks {
		blocksCopy[i] = MakeDeepCopyOfBlock(block)
	}

	oldPeers := msg.PeerMap
	peersCopy := make(map[string]models.Void, len(oldPeers))
	for k, v := range oldPeers {
		peersCopy[k] = v
	}

	deepCopyMessage := blockchain.Message{
		MessageType:       msg.MessageType,
		MessageSender:     msg.MessageSender,
		SignedTransaction: MakeDeepCopyOfTransaction(msg.SignedTransaction),
		MessageBlocks:     blocksCopy,
		PeerMap:           peersCopy,
	}
	return deepCopyMessage
}
