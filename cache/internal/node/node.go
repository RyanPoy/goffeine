package node

import (
	"errors"
)

type Position int

const (
	WINDOW Position = iota
	PROBATION
	PROTECTED
)

type Node struct {
	Key      string
	KeyHash  []byte
	Value    interface{}
	Location Position
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
		Location: WINDOW,
		Weight:   weight,
	}
}

func (n *Node) InWindow() {
	n.Location = WINDOW
}

func (n *Node) InProbation() {
	n.Location = PROBATION
}

func (n *Node) InProtected() {
	n.Location = PROTECTED
}

func (n *Node) IsInWindow() bool {
	return n.Location == WINDOW
}

func (n *Node) IsInProbation() bool {
	return n.Location == PROBATION
}

func (n *Node) IsInProtected() bool {
	return n.Location == PROTECTED
}

func (n *Node) Position() Position {
	return n.Location
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
