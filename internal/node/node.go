package node

type NodePosition int

const (
	WindowPosition NodePosition = iota
	ProbationPosition
	ProtectedPosition
)

type GoffeineNode struct {
	Key      string
	Value    any
	Position NodePosition
}

func New(key string, value any, position NodePosition) *GoffeineNode {
	return &GoffeineNode{Key: key, Value: value, Position: position}
}
