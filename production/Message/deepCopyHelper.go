package Message

import (
	"github.com/OskarRVestergaard/BachelorProject/Task1/Models"
	"github.com/OskarRVestergaard/BachelorProject/production/models"
	"github.com/OskarRVestergaard/BachelorProject/production/models/PoWblockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/models/SpaceMintBlockchain"
	"github.com/OskarRVestergaard/BachelorProject/production/sha256"
	"github.com/OskarRVestergaard/BachelorProject/production/strategies/lottery_strategy"
)

func deepCopyByteSlice(slice []byte) []byte {
	oldSlice := slice
	sliceCopy := make([]byte, len(oldSlice))
	copy(sliceCopy, oldSlice)
	return sliceCopy
}

func MakeDeepCopyOfTransaction(transaction models.SignedTransaction) (copyOfTransaction models.SignedTransaction) {
	deepCopyTransaction := models.SignedTransaction{
		Id:        transaction.Id,
		From:      transaction.From,
		To:        transaction.To,
		Amount:    transaction.Amount,
		Signature: deepCopyByteSlice(transaction.Signature),
	}
	return deepCopyTransaction
}

func deepCopyTransactions(transactions []models.SignedTransaction) (copyOfTransactions []models.SignedTransaction) {
	transactionsCopy := make([]models.SignedTransaction, len(transactions))
	for i, transaction := range transactions {
		transactionsCopy[i] = MakeDeepCopyOfTransaction(transaction)
	}
	return transactionsCopy
}

func MakeDeepCopyOfPoWBlock(block PoWblockchain.Block) (copyOfBlock PoWblockchain.Block) {
	deepCopyBlock := PoWblockchain.Block{
		IsGenesis: block.IsGenesis,
		Vk:        block.Vk,
		Slot:      block.Slot,
		Draw:      MakeDeepCopyOfWinningParams(block.Draw),
		BlockData: PoWblockchain.BlockData{
			Hardness:     block.BlockData.Hardness,
			Transactions: deepCopyTransactions(block.BlockData.Transactions),
		},
		ParentHash: block.ParentHash,
		Signature:  deepCopyByteSlice(block.Signature),
	}
	return deepCopyBlock
}

func MakeDeepCopyOfPoSBlock(block SpaceMintBlockchain.Block) (copyOfBlock SpaceMintBlockchain.Block) {
	deepCopyBlock := SpaceMintBlockchain.Block{
		IsGenesis:           block.IsGenesis,
		ParentHash:          block.ParentHash,
		HashSubBlock:        deepCopyHashSubBlock(block.HashSubBlock),
		TransactionSubBlock: deepCopyTransactionSubBlock(block.TransactionSubBlock),
		SignatureSubBlock:   deepCopySignatureSubBlock(block.SignatureSubBlock),
	}
	return deepCopyBlock
}

func deepCopyOpeningTriple(triple Models.OpeningTriple) Models.OpeningTriple {
	copyOfOpeningValues := make([]sha256.HashValue, len(triple.OpenValues))
	for i, value := range triple.OpenValues {
		copyOfOpeningValues[i] = value
	}
	copyOfTriple := Models.OpeningTriple{
		Index:      triple.Index,
		Value:      triple.Value,
		OpenValues: copyOfOpeningValues,
	}
	return copyOfTriple
}

func deepCopyOpeningTriples(triples []Models.OpeningTriple) []Models.OpeningTriple {
	triplesCopy := make([]Models.OpeningTriple, len(triples))
	for i, triple := range triples {
		triplesCopy[i] = deepCopyOpeningTriple(triple)
	}
	return triplesCopy
}

func deepCopyPoSpaceLotteryDraw(draw lottery_strategy.PoSpaceLotteryDraw) lottery_strategy.PoSpaceLotteryDraw {
	copyOfDraw := lottery_strategy.PoSpaceLotteryDraw{
		Vk:                        draw.Vk,
		ParentHash:                draw.ParentHash,
		ProofOfSpaceA:             deepCopyOpeningTriples(draw.ProofOfSpaceA),
		ProofOfCorrectCommitmentB: deepCopyOpeningTriples(draw.ProofOfCorrectCommitmentB),
	}
	return copyOfDraw
}

func deepCopyHashSubBlock(subBlock SpaceMintBlockchain.HashSubBlock) SpaceMintBlockchain.HashSubBlock {
	deepCopyOfHashSubBlock := SpaceMintBlockchain.HashSubBlock{
		Slot:                      subBlock.Slot,
		SignatureOnParentSubBlock: deepCopyByteSlice(subBlock.SignatureOnParentSubBlock),
		Draw:                      deepCopyPoSpaceLotteryDraw(subBlock.Draw),
	}

	return deepCopyOfHashSubBlock
}

func deepCopySpaceCommit(spaceCommitment SpaceMintBlockchain.SpaceCommitment) SpaceMintBlockchain.SpaceCommitment {
	copyOfSpaceCommitment := SpaceMintBlockchain.SpaceCommitment{
		Id:         spaceCommitment.Id,
		N:          spaceCommitment.N,
		PublicKey:  spaceCommitment.PublicKey,
		Commitment: spaceCommitment.Commitment,
	}
	return copyOfSpaceCommitment
}

func deepCopySpaceCommitments(spaceCommitments []SpaceMintBlockchain.SpaceCommitment) []SpaceMintBlockchain.SpaceCommitment {
	spaceCommitmentsCopy := make([]SpaceMintBlockchain.SpaceCommitment, len(spaceCommitments))
	for i, spaceCommit := range spaceCommitments {
		spaceCommitmentsCopy[i] = deepCopySpaceCommit(spaceCommit)
	}
	return spaceCommitmentsCopy
}

func deepCopyTransactionSubBlock(subBlock SpaceMintBlockchain.TransactionSubBlock) (copyOfSubBlock SpaceMintBlockchain.TransactionSubBlock) {

	spaceMintTransactionCopy := SpaceMintBlockchain.SpacemintTransactions{
		Payments:         deepCopyTransactions(subBlock.Transactions.Payments),
		SpaceCommitments: deepCopySpaceCommitments(subBlock.Transactions.SpaceCommitments),
		Penalties:        nil, //TODO Implement
	}

	deepCopyOfTransactionSubBlock := SpaceMintBlockchain.TransactionSubBlock{
		Slot:         subBlock.Slot,
		Transactions: spaceMintTransactionCopy,
	}
	return deepCopyOfTransactionSubBlock
}

func deepCopySignatureSubBlock(subBlock SpaceMintBlockchain.SignatureSubBlock) SpaceMintBlockchain.SignatureSubBlock {

	deepCopyOfSignatureSubBlock := SpaceMintBlockchain.SignatureSubBlock{
		Slot:                                  subBlock.Slot,
		SignatureOnCurrentTransactionSubBlock: deepCopyByteSlice(subBlock.SignatureOnCurrentTransactionSubBlock),
		SignatureOnParentSubBlock:             deepCopyByteSlice(subBlock.SignatureOnParentSubBlock),
	}
	return deepCopyOfSignatureSubBlock
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

	oldPoWBlocks := msg.PoWMessageBlocks
	blocksCopy := make([]PoWblockchain.Block, len(oldPoWBlocks))
	for i, block := range oldPoWBlocks {
		blocksCopy[i] = MakeDeepCopyOfPoWBlock(block)
	}

	oldPoSBlocks := msg.SpaceMintBlocks
	PoSBlocksCopy := make([]SpaceMintBlockchain.Block, len(oldPoSBlocks))
	for i, block := range oldPoSBlocks {
		PoSBlocksCopy[i] = MakeDeepCopyOfPoSBlock(block)
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
		PoWMessageBlocks:  blocksCopy,
		SpaceMintBlocks:   PoSBlocksCopy,
		PeerMap:           peersCopy,
	}
	return deepCopyMessage
}
