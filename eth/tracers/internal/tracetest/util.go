package tracetest

import (

	// Force-load native and js packages, to trigger registration
	_ "github.com/ethereum/go-ethereum/eth/tracers/js"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
)

const (
	memoryPadLimit = 1024 * 1024
)

// To generate a new callTracer test, copy paste the makeTest method below into
// a Geth console and call it with a transaction hash you which to export.

/*
// makeTest generates a callTracer test by running a prestate reassembled and a
// call trace run, assembling all the gathered information into a test case.
var makeTest = function(tx, rewind) {
  // Generate the genesis block from the block, transaction and prestate data
  var block   = eth.getBlock(eth.getTransaction(tx).blockHash);
  var genesis = eth.getBlock(block.parentHash);

  delete genesis.gasUsed;
  delete genesis.logsBloom;
  delete genesis.parentHash;
  delete genesis.receiptsRoot;
  delete genesis.sha3Uncles;
  delete genesis.size;
  delete genesis.transactions;
  delete genesis.transactionsRoot;
  delete genesis.uncles;

  genesis.gasLimit  = genesis.gasLimit.toString();
  genesis.number    = genesis.number.toString();
  genesis.timestamp = genesis.timestamp.toString();

  genesis.alloc = debug.traceTransaction(tx, {tracer: "prestateTracer", rewind: rewind});
  for (var key in genesis.alloc) {
    var nonce = genesis.alloc[key].nonce;
    if (nonce) {
      genesis.alloc[key].nonce = nonce.toString();
    }
  }
  genesis.config = admin.nodeInfo.protocols.eth.config;

  // Generate the call trace and produce the test input
  var result = debug.traceTransaction(tx, {tracer: "callTracer", rewind: rewind});
  delete result.time;

  console.log(JSON.stringify({
    genesis: genesis,
    context: {
      number:     block.number.toString(),
      difficulty: block.difficulty,
      timestamp:  block.timestamp.toString(),
      gasLimit:   block.gasLimit.toString(),
      miner:      block.miner,
    },
    input:  eth.getRawTransaction(tx),
    result: result,
  }, null, 2));
}
*/

// func GetMemoryCopyPadded(m []byte, offset, size int64) ([]byte, error) {
// 	if offset < 0 || size < 0 {
// 		return nil, errors.New("offset or size must not be negative")
// 	}
// 	length := int64(len(m))
// 	if offset+size < length { // slice fully inside memory
// 		return memoryCopy(m, offset, size), nil
// 	}
// 	paddingNeeded := offset + size - length
// 	if paddingNeeded > memoryPadLimit {
// 		return nil, fmt.Errorf("reached limit for padding memory slice: %d", paddingNeeded)
// 	}
// 	cpy := make([]byte, size)
// 	if overlap := length - offset; overlap > 0 {
// 		copy(cpy, MemoryPtr(m, offset, overlap))
// 	}
// 	return cpy, nil
// }

// func memoryCopy(m []byte, offset, size int64) (cpy []byte) {
// 	if size == 0 {
// 		return nil
// 	}

// 	if len(m) > int(offset) {
// 		cpy = make([]byte, size)
// 		copy(cpy, m[offset:offset+size])

// 		return
// 	}

// 	return
// }

// // camel converts a snake cased input string into a camel cased output.
// func camel(str string) string {
// 	pieces := strings.Split(str, "_")
// 	for i := 1; i < len(pieces); i++ {
// 		pieces[i] = string(unicode.ToUpper(rune(pieces[i][0]))) + pieces[i][1:]
// 	}
// 	return strings.Join(pieces, "")
// }

// // MemoryPtr returns a pointer to a slice of memory.
// func MemoryPtr(m []byte, offset, size int64) []byte {
// 	if size == 0 {
// 		return nil
// 	}

// 	if len(m) > int(offset) {
// 		return m[offset : offset+size]
// 	}

// 	return nil
// }
