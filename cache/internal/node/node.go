package node

import (
	"errors"
	"sync"
)

type Node struct {
	key string
	keyHash []byte
	value interface{}
	hashmap sync.Map
	position int//window：0，probation：1，protected：2
	weight int
}

func New(key string, value interface{}) *Node {
	return &Node{
		key: key,
		keyHash: []byte(key),
		value: value,
		position: 0,
		weight:1,
	}
}
func (n* Node) Weight() int {
	return n.weight
}
func (n* Node) SetWeight(a int)  {
	 n.weight=a
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
func (n *Node) SetPosition( position int) {
	n.position=position
}
func (n *Node) Keyhash() []byte{
	return n.keyHash
}
func (n *Node) Equals(n2 *Node) bool {
	return n.key == n2.key && n.value == n2.value
}
func (n *Node) SetValue(value interface{}){
	n.value=value
}
//func (n *Node) GetPreviousInAccessOrder() *Node{
//	//返回node结点前一个结点
//	return nil
//}
//func (n *Node) GetNextInAccessOrder() *Node{
//	//返回node结点后一个结点
//	return nil
//}
//func (n *Node) SetWeight(weight int){
//	n.weight=weight
//}
//func (n *Node)Weight() int{
//	return n.weight
//}
func (n *Node) UpdateWith(n2 *Node) error {
	if n.key != n2.key {
		return errors.New("The keys of two nodes are different")
	}
	n.value = n2.value
	return nil
}
