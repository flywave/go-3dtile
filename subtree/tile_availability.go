package subtree

// TileAvailability represents tile availability information
type TileAvailability struct {
	Bitstream      *int `json:"bitstream,omitempty"`
	AvailableCount *int `json:"availableCount,omitempty"`
	Constant       *int `json:"constant,omitempty"`
}