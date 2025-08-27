package subtree

// Read reads a bit array from a byte array with the specified offset and length
func Read(subtreeBinary []byte, offset int, length int) []bool {
	// Ensure we don't go out of bounds
	if offset < 0 || offset >= len(subtreeBinary) {
		return []bool{}
	}
	
	end := offset + length
	if end > len(subtreeBinary) {
		end = len(subtreeBinary)
	}
	
	// Extract the relevant bytes
	slicedBytes := subtreeBinary[offset:end]
	
	// Convert bytes to bit array
	bitCount := len(slicedBytes) * 8
	bitArray := make([]bool, bitCount)
	
	for i, b := range slicedBytes {
		for j := 0; j < 8; j++ {
			if (b & (1 << j)) != 0 {
				bitArray[i*8+j] = true
			}
		}
	}
	
	return bitArray
}