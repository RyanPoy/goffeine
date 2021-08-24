package node

import (
	"errors"
)

const (
	WINDOW    int = 0
	PROBATION int = 1
	PROTECTED int = 2
)

type Node struct {
	Key      string
	KeyHash  []byte
	Value    interface{}
	position int //window：0，probation：1，protected：2
	Weight   int
}

func New(key string, value interface{}) *Node {
	return NewWithWeight(key, value, 1)
}

func NewWithWeight(key string, value interface{}, weight int) *Node {
	return &Node{
		Key:      key,
		KeyHash:  []byte(key),
		Value:    value,
		position: WINDOW,
		Weight:   weight,
	}
}

func (n *Node) InWindow() {
	n.position = WINDOW
}

func (n *Node) InProbation() {
	n.position = PROBATION
}

func (n *Node) InProtected() {
	n.position = PROTECTED
}

func (n *Node) IsBelongsToWindow() bool {
	return n.position == WINDOW
}

func (n *Node) IsBelongsToProbation() bool {
	return n.position == PROBATION
}

func (n *Node) IsBelongsToProtected() bool {
	return n.position == PROTECTED
}

func (n *Node) Position() int {
	return n.position
}

func (n *Node) Equals(n2 *Node) bool {
	return n.Key == n2.Key
}

func (n *Node) UpdateWith(n2 *Node) error {
	if !n.Equals(n2) {
		return errors.New("The keys of two nodes are different")
	}
	n.Value = n2.Value
	n.Weight = n2.Weight
	return nil
}
