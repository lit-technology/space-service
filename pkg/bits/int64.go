package bits

// RangeInt64 creates a int64 range from unsetting and setting n least significant bits.
func RangeInt64(bits uint8, n int64) (int64, int64) {
	mask := int64(1<<(bits)) - 1
	start := n &^ mask
	end := n | mask
	return start, end
}

func SetLeastSignificantBits(bits uint8, n int64) int64 {
	return n | ((1 << (bits)) - 1)
}

func UnsetLeastSignificantBits(bits uint8, n int64) int64 {
	return n &^ ((1 << (bits)) - 1)
}

func UnsetBitsFromRight(bits uint8, n int64, right uint8) int64 {
	mask := (1 << bits) - 1
	mask <<= right
	return n &^ int64(mask)
}

func IsUnsetBitsFromRight(bits uint8, n int64, right uint8) bool {
	return n == UnsetBitsFromRight(bits, n, right)
}
