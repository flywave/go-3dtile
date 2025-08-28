package subtree

// FromString creates a boolean array from a string representation
func FromString(bits string) []bool {
	bitArray := make([]bool, len(bits))
	for i, c := range bits {
		if c == '1' {
			bitArray[i] = true
		}
	}
	return bitArray
}

// ToByteArray converts a boolean array to a byte array
func ToByteArray(bits []bool) []byte {
	if len(bits) == 0 {
		return []byte{}
	}

	// Calculate the number of bytes needed
	byteLength := (len(bits)-1)/8 + 1
	result := make([]byte, byteLength)

	// Convert bits to bytes
	for i, bit := range bits {
		if bit {
			byteIndex := i / 8
			bitIndex := i % 8
			result[byteIndex] |= 1 << bitIndex
		}
	}

	return result
}

// AsString converts a boolean array to a string representation
func AsString(bitArray []bool) string {
	result := make([]byte, len(bitArray))
	for i, bit := range bitArray {
		if bit {
			result[i] = '1'
		} else {
			result[i] = '0'
		}
	}
	return string(result)
}

// Count counts the number of bits with the specified value
func Count(bitArray []bool, whereClause bool) int {
	count := 0
	for _, bit := range bitArray {
		if bit == whereClause {
			count++
		}
	}
	return count
}
