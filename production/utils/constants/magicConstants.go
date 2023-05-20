package constants

import "time"

var BlockPaymentAmountLimit = 20
var BlockSpaceCommitAmountLimit = 32
var BlockPenaltyAmountLimit = 32

var Hardness = 23
var SlotLength = 4000 * time.Millisecond

const GraphK = 4

const Alpha = 0.25
const Beta = 0.5
