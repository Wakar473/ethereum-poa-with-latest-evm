package stateless

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type HeaderReader interface {
	// GetHeader retrieves a block header from the database by hash and number,
	GetHeader(hash common.Hash, number uint64) *types.Header
}
type Witness struct {
	context *types.Header // Header to which this witness belongs to, with rootHash and receiptHash zeroed out

	Headers []*types.Header     // Past headers in reverse order (0=parent, 1=parent's-parent, etc). First *must* be set.
	Codes   map[string]struct{} // Set of bytecodes ran or accessed
	State   map[string]struct{} // Set of MPT state trie nodes (account and storage together)

	chain HeaderReader // Chain reader to convert block hash ops to header proofs
	lock  sync.Mutex   // Lock to allow concurrent state insertions
}

func (w *Witness) AddBlockHash(number uint64) {
	// Keep pulling in headers until this hash is populated
	for int(w.context.Number.Uint64()-number) > len(w.Headers) {
		tail := w.Headers[len(w.Headers)-1]
		w.Headers = append(w.Headers, w.chain.GetHeader(tail.ParentHash, tail.Number.Uint64()-1))
	}
}
