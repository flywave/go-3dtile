package subtree

// ContentAvailability represents content availability information
type ContentAvailability struct {
	AvailableCount int  `json:"availableCount"`
	Bitstream      *int `json:"bitstream,omitempty"`
	Constant       *int `json:"constant,omitempty"`
}