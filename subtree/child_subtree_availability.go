package subtree

// ChildSubtreeAvailability represents child subtree availability information
type ChildSubtreeAvailability struct {
	Bitstream      *int `json:"bitstream,omitempty"`
	AvailableCount int  `json:"availableCount"`
	Constant       *int `json:"constant,omitempty"`
}