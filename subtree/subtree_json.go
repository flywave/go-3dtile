package subtree

// SubtreeJson represents the JSON structure of a subtree
type SubtreeJson struct {
	Buffers                  []Buffer                   `json:"buffers"`
	BufferViews              []BufferView               `json:"bufferViews"`
	TileAvailability         *TileAvailability          `json:"tileAvailability"`
	ContentAvailability      []ContentAvailability      `json:"contentAvailability"`
	ChildSubtreeAvailability *ChildSubtreeAvailability  `json:"childSubtreeAvailability,omitempty"`
}