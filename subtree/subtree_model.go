package subtree

// Subtree represents a 3D Tiles subtree
type Subtree struct {
	SubtreeHeader               *SubtreeHeader
	SubtreeJson                 string
	SubtreeBinary               []byte
	ChildSubtreeAvailability    []bool
	TileAvailability            []bool
	TileAvailabilityConstant    int
	ContentAvailability         []bool
	ContentAvailabilityConstant int
}

// NewSubtree creates a new Subtree instance
func NewSubtree() *Subtree {
	return &Subtree{
		SubtreeHeader:               NewSubtreeHeader(),
		TileAvailabilityConstant:    0,
		ContentAvailabilityConstant: 0,
	}
}