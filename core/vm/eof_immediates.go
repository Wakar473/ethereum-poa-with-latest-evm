package vm

import(
	
)

var terminals [256]bool
var immediates [256]uint8

func Immediates(op OpCode) int {
	return int(immediates[op])
}