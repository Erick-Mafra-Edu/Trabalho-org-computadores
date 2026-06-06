package cache

import "fmt"

func isPowerOfTwo(value uint64) bool {
	return value > 0 && (value&(value-1)) == 0
}

func log2PowerOfTwo(value uint64) uint {
	var bits uint
	for value > 1 {
		value >>= 1
		bits++
	}
	return bits
}

func maxAddressValue(bits uint) uint64 {
	if bits >= 64 {
		return ^uint64(0)
	}
	return (uint64(1) << bits) - 1
}

func addressBinary(value uint64, bits uint) string {
	return fmt.Sprintf("%0*b", int(bits), value)
}
