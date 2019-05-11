package bits

import (
	"math"
	"math/bits"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ToBits(n int64) string {
	return strconv.FormatInt(n, 2)
}

func Ones(n int64) int {
	return bits.OnesCount64(uint64(n))
}

func LeadingZeros(n int64) int {
	return bits.LeadingZeros64(uint64(n))
}

func TrailingZeros(n int64) int {
	return bits.TrailingZeros64(uint64(n))
}

func TestSetLeastSignificantBits(t *testing.T) {
	for i := 1; i < 63; i++ {
		ans := SetLeastSignificantBits(uint8(i), 0)
		bits := ToBits(ans)
		assert.Equal(t, i, Ones(ans), "Expected to have %d set bits, found %s", i, bits)
		assert.Equal(t, 64-i, LeadingZeros(ans), "Expected to have %d leading unset bits, found %s", 62-i, bits)
		assert.Equal(t, 0, TrailingZeros(ans), "Expected to have 0 trailing unset bits, found %s", bits)
	}
	for i := 1; i < 63; i++ {
		n := int64(1) << uint8(i)
		bits := ToBits(n)
		assert.Equal(t, 1, Ones(n), "Expected to have one set bit, found %s", bits)
		assert.Equal(t, 63-i, LeadingZeros(n), "Expected to have %d leading unset bits, found %s", 63-i, bits)
		assert.Equal(t, i, TrailingZeros(n), "Expected to have %d trailing unset bits, found %s", i, bits)

		for j := 1; j < i; j++ {
			ans := SetLeastSignificantBits(uint8(j), n)
			bits = ToBits(ans)
			assert.Equal(t, 1+j, Ones(ans))
			assert.Equal(t, 0, TrailingZeros(ans), "Expected to have 0 trailing unset bits, found %s", bits)
			assert.Equal(t, 63-i, LeadingZeros(ans), "Expected to have %d leading unset bits, found %s", 63-i, bits)
		}
	}
}

func TestUnsetLeastSignificantBits(t *testing.T) {
	for i := 1; i < 63; i++ {
		max := int64(math.MaxInt64)
		ans := UnsetLeastSignificantBits(uint8(i), max)

		// Note postive numbers has unset sign bit.
		// t.Logf("%s %d", ToBits(ans), i)
		bits := ToBits(ans)
		assert.Equal(t, 63-i, Ones(ans), "Expected %d set bits, found %s", 63-i, bits)
		assert.Equal(t, i, TrailingZeros(ans), "Expected %d trailing unset bits, found %s", i, bits)
		assert.Equal(t, 1, LeadingZeros(ans), "Expected 1 leading postive sign bit, found %s", bits)
	}
}

func TestUnsetBitsFromRight(t *testing.T) {
	for i := 1; i < 61; i++ {
		max := int64(math.MaxInt64)
		ans := UnsetBitsFromRight(2, max, uint8(i))

		// Note postive numbers has unset sign bit - 2 for mask.
		// t.Logf("%s %d", ToBits(ans), i)
		bits := ToBits(ans)
		assert.Equal(t, 61, Ones(ans), "Expected 63 set bits, found %s", bits)
		assert.Equal(t, 1, LeadingZeros(ans), "Expected 1 leading postive sign bits, found %s", bits)
		assert.Equal(t, 0, TrailingZeros(ans), "Expected 0 trailing bits, found %s", bits)

		arr := strings.Split(bits, "00")
		assert.Len(t, arr, 2)
		assert.Len(t, arr[0], 61-i)
		assert.Len(t, arr[1], i)
	}
}

func TestIsUnsetBitsFromRight(t *testing.T) {
	i, err := strconv.ParseInt("101000011", 2, 64)
	assert.NoError(t, err)

	assert.False(t, IsUnsetBitsFromRight(2, i, 1))
	assert.True(t, IsUnsetBitsFromRight(4, i, 2))
	assert.False(t, IsUnsetBitsFromRight(5, i, 2))
	assert.True(t, IsUnsetBitsFromRight(2, i, 2))

	assert.True(t, IsUnsetBitsFromRight(3, i, 3))
	assert.False(t, IsUnsetBitsFromRight(4, i, 3))

	assert.True(t, IsUnsetBitsFromRight(1, i, 7))
	assert.False(t, IsUnsetBitsFromRight(1, i, 6))
	assert.False(t, IsUnsetBitsFromRight(1, i, 8))

	assert.True(t, IsUnsetBitsFromRight(0, i, 0))
}
