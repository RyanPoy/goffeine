package node

type Position int

const (
	WindowPosition Position = iota
	ProbationPosition
	ProtectedPosition
)

type GoffeineNode struct {
	Key      string
	Value    any
	Position Position
}

func New(key string, value any, position Position) *GoffeineNode {
	return &GoffeineNode{Key: key, Value: value, Position: position}
}
