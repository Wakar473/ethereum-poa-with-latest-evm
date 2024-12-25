package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie/utils"
	"github.com/holiman/uint256"
)

var zeroTreeIndex uint256.Int

type mode byte

const (
	AccessWitnessReadFlag  = mode(1)
	AccessWitnessWriteFlag = mode(2)
)

// AccessEvents lists the locations of the state that are being accessed
// during the production of a block.
type AccessEvents struct {
	branches map[branchAccessKey]mode
	chunks   map[chunkAccessKey]mode

	pointCache *utils.PointCache
}

type branchAccessKey struct {
	addr      common.Address
	treeIndex uint256.Int
}

type chunkAccessKey struct {
	branchAccessKey
	leafKey byte
}

func NewAccessEvents(pointCache *utils.PointCache) *AccessEvents {
	return &AccessEvents{
		branches:   make(map[branchAccessKey]mode),
		chunks:     make(map[chunkAccessKey]mode),
		pointCache: pointCache,
	}
}

// AddAccount returns the gas to be charged for each of the currently cold
// member fields of an account.
func (ae *AccessEvents) AddAccount(addr common.Address, isWrite bool) uint64 {
	var gas uint64
	gas += ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.BasicDataLeafKey, isWrite)
	gas += ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.CodeHashLeafKey, isWrite)
	return gas
}

// touchAddressAndChargeGas adds any missing access event to the access event list, and returns the cold
// access cost to be charged, if need be.
func (ae *AccessEvents) touchAddressAndChargeGas(addr common.Address, treeIndex uint256.Int, subIndex byte, isWrite bool) uint64 {
	stemRead, selectorRead, stemWrite, selectorWrite, selectorFill := ae.touchAddress(addr, treeIndex, subIndex, isWrite)

	var gas uint64
	if stemRead {
		gas += params.WitnessBranchReadCost
	}
	if selectorRead {
		gas += params.WitnessChunkReadCost
	}
	if stemWrite {
		gas += params.WitnessBranchWriteCost
	}
	if selectorWrite {
		gas += params.WitnessChunkWriteCost
	}
	if selectorFill {
		gas += params.WitnessChunkFillCost
	}
	return gas
}

// ContractCreateInitGas returns the access gas costs for the initialization of
// a contract creation.
func (ae *AccessEvents) ContractCreateInitGas(addr common.Address) uint64 {
	var gas uint64
	gas += ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.BasicDataLeafKey, true)
	gas += ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.CodeHashLeafKey, true)
	return gas
}

func (ae *AccessEvents) ContractCreatePreCheckGas(addr common.Address) uint64 {
	var gas uint64
	gas += ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.BasicDataLeafKey, false)
	gas += ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.CodeHashLeafKey, false)
	return gas
}

func newBranchAccessKey(addr common.Address, treeIndex uint256.Int) branchAccessKey {
	var sk branchAccessKey
	sk.addr = addr
	sk.treeIndex = treeIndex
	return sk
}

func newChunkAccessKey(branchKey branchAccessKey, leafKey byte) chunkAccessKey {
	var lk chunkAccessKey
	lk.branchAccessKey = branchKey
	lk.leafKey = leafKey
	return lk
}

// touchAddress adds any missing access event to the access event list.
func (ae *AccessEvents) touchAddress(addr common.Address, treeIndex uint256.Int, subIndex byte, isWrite bool) (bool, bool, bool, bool, bool) {
	branchKey := newBranchAccessKey(addr, treeIndex)
	chunkKey := newChunkAccessKey(branchKey, subIndex)

	// Read access.
	var branchRead, chunkRead bool
	if _, hasStem := ae.branches[branchKey]; !hasStem {
		branchRead = true
		ae.branches[branchKey] = AccessWitnessReadFlag
	}
	if _, hasSelector := ae.chunks[chunkKey]; !hasSelector {
		chunkRead = true
		ae.chunks[chunkKey] = AccessWitnessReadFlag
	}

	// Write access.
	var branchWrite, chunkWrite, chunkFill bool
	if isWrite {
		if (ae.branches[branchKey] & AccessWitnessWriteFlag) == 0 {
			branchWrite = true
			ae.branches[branchKey] |= AccessWitnessWriteFlag
		}

		chunkValue := ae.chunks[chunkKey]
		if (chunkValue & AccessWitnessWriteFlag) == 0 {
			chunkWrite = true
			ae.chunks[chunkKey] |= AccessWitnessWriteFlag
		}
		// TODO: charge chunk filling costs if the leaf was previously empty in the state
	}
	return branchRead, chunkRead, branchWrite, chunkWrite, chunkFill
}

func (ae *AccessEvents) AddTxOrigin(originAddr common.Address) {
	ae.touchAddressAndChargeGas(originAddr, zeroTreeIndex, utils.BasicDataLeafKey, true)
	ae.touchAddressAndChargeGas(originAddr, zeroTreeIndex, utils.CodeHashLeafKey, false)
}

func (ae *AccessEvents) AddTxDestination(addr common.Address, sendsValue bool) {
	ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.BasicDataLeafKey, sendsValue)
	ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.CodeHashLeafKey, false)
}


func (ae *AccessEvents) Merge(other *AccessEvents) {
	for k := range other.branches {
		ae.branches[k] |= other.branches[k]
	}
	for k, chunk := range other.chunks {
		ae.chunks[k] |= chunk
	}
}
// SlotGas returns the amount of gas to be charged for a cold storage access.
func (ae *AccessEvents) SlotGas(addr common.Address, slot common.Hash, isWrite bool) uint64 {
	treeIndex, subIndex := utils.StorageIndex(slot.Bytes())
	return ae.touchAddressAndChargeGas(addr, *treeIndex, subIndex, isWrite)
}

func (ae *AccessEvents) BasicDataGas(addr common.Address, isWrite bool) uint64 {
	return ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.BasicDataLeafKey, isWrite)
}

func (ae *AccessEvents) CodeHashGas(addr common.Address, isWrite bool) uint64 {
	return ae.touchAddressAndChargeGas(addr, zeroTreeIndex, utils.CodeHashLeafKey, isWrite)
}

func (ae *AccessEvents) MessageCallGas(destination common.Address) uint64 {
	var gas uint64
	gas += ae.touchAddressAndChargeGas(destination, zeroTreeIndex, utils.BasicDataLeafKey, false)
	return gas
}

// CodeChunksRangeGas is a helper function to touch every chunk in a code range and charge witness gas costs
func (ae *AccessEvents) CodeChunksRangeGas(contractAddr common.Address, startPC, size uint64, codeLen uint64, isWrite bool) uint64 {
	// note that in the case where the copied code is outside the range of the
	// contract code but touches the last leaf with contract code in it,
	// we don't include the last leaf of code in the AccessWitness.  The
	// reason that we do not need the last leaf is the account's code size
	// is already in the AccessWitness so a stateless verifier can see that
	// the code from the last leaf is not needed.
	if (codeLen == 0 && size == 0) || startPC > codeLen {
		return 0
	}

	endPC := startPC + size
	if endPC > codeLen {
		endPC = codeLen
	}
	if endPC > 0 {
		endPC -= 1 // endPC is the last bytecode that will be touched.
	}

	var statelessGasCharged uint64
	for chunkNumber := startPC / 31; chunkNumber <= endPC/31; chunkNumber++ {
		treeIndex := *uint256.NewInt((chunkNumber + 128) / 256)
		subIndex := byte((chunkNumber + 128) % 256)
		gas := ae.touchAddressAndChargeGas(contractAddr, treeIndex, subIndex, isWrite)
		var overflow bool
		statelessGasCharged, overflow = math.SafeAdd(statelessGasCharged, gas)
		if overflow {
			panic("overflow when adding gas")
		}
	}
	return statelessGasCharged
}
