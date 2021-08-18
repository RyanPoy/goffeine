package node

import "errors"

type Node struct {
	key string
	keyHash []byte
	value interface{}
}

func New(key string, value interface{}) *Node {
	return &Node{
		key: key,
		keyHash: []byte(key),
		value: value,
	}
}

func (n *Node) Key() string {
	return n.key
}

func (n *Node) KeyHash() []byte {
	return n.keyHash
}

func (n *Node) Value() interface{} {
	return n.value
}

func (n *Node) Equals(n2 *Node) bool {
	return n.key == n2.key && n.value == n2.value
}

func (n *Node) UpdateWith(n2 *Node) error {
	if n.key != n2.key {
		return errors.New("The keys of two nodes are different")
	}
	n.value = n2.value
	return nil
}