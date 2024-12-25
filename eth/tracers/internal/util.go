package internal

import "github.com/holiman/uint256"

const (
	memoryPadLimit = 1024 * 1024
)

func MemoryPtr(m []byte, offset, size int64) []byte {
	if size == 0 {
		return nil
	}

	if len(m) > int(offset) {
		return m[offset : offset+size]
	}

	return nil
}

func StackBack(st []uint256.Int, n int) *uint256.Int {
	return &st[len(st)-n-1]
}
