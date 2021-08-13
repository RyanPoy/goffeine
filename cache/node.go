package cache

type Node struct {
	key string
	keyHash []byte
	value interface{}
}

func NewNode(key string, value interface{}) Node {
	return Node {
		key: key,
		keyHash: []byte(key),
		value: value,
	}
}

func (n *Node) Value() interface{} {
	return n.value
}
