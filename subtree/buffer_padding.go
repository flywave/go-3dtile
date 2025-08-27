package subtree

const boundary = 8

// AddPadding adds padding to a byte array
func AddPadding(bytes []byte, offset int) []byte {
	remainder := (offset + len(bytes)) % boundary
	padding := 0
	if remainder != 0 {
		padding = boundary - remainder
	}

	// Create padding bytes (whitespace)
	paddingBytes := make([]byte, padding)
	for i := range paddingBytes {
		paddingBytes[i] = ' '
	}

	// Concatenate original bytes with padding
	result := make([]byte, len(bytes)+padding)
	copy(result, bytes)
	copy(result[len(bytes):], paddingBytes)

	return result
}

// AddPaddingString adds padding to a string
func AddPaddingString(input string, offset int) string {
	bytes := []byte(input)
	paddedBytes := AddPadding(bytes, offset)
	return string(paddedBytes)
}

// AddBinaryPadding adds binary padding to a byte array
func AddBinaryPadding(bytes []byte, offset int) []byte {
	remainder := (offset + len(bytes)) % boundary
	padding := 0
	if remainder != 0 {
		padding = boundary - remainder
	}

	// Create padding bytes (zeros)
	paddingBytes := make([]byte, padding)

	// Concatenate original bytes with padding
	result := make([]byte, len(bytes)+padding)
	copy(result, bytes)
	copy(result[len(bytes):], paddingBytes)

	return result
}

// GetByteArray creates a byte array of the specified length filled with zeros
func GetByteArray(length int) []byte {
	arr := make([]byte, length)
	// Already filled with zeros by default
	return arr
}
