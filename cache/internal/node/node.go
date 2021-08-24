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
	key      string
	keyHash  []byte
	value    interface{}
	position int //window：0，probation：1，protected：2
	weight   int
}

func New(key string, value interface{}) *Node {
	return &Node{
		key:      key,
		keyHash:  []byte(key),
		value:    value,
		position: WINDOW,
		weight:   1,
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

func (n *Node) Weight() int {
	return n.weight
}

func (n *Node) SetWeight(a int) {
	n.weight = a
}

func (n *Node) KeyHash() []byte {
	return n.keyHash
}

func (n *Node) Value() interface{} {
	return n.value
}
func (n *Node) Key() string {
	return n.key

}
func (n *Node) Position() int {
	return n.position
}

func (n *Node) Keyhash() []byte {
	return n.keyHash
}
func (n *Node) Equals(n2 *Node) bool {
	return n.key == n2.key && n.value == n2.value
}
func (n *Node) SetValue(value interface{}) {
	n.value = value
}

func (n *Node) UpdateWith(n2 *Node) error {
	if n.key != n2.key {
		return errors.New("The keys of two nodes are different")
	}
	n.value = n2.value
	return nil
}
