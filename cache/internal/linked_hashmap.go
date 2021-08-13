package internal

import (
	"container/list"
	"errors"
	"fmt"
)

// 顾名思义：LinkedHashMap。里面封装了一个map和一个list
// 注意：这个LinkedHashMap是能存放的数据容量取决于cap
type LinkedHashMap struct {
	cap     int
	queue   *list.List // doubly link queue
	hashmap map[string]*list.Element
}

func NewLRU(cap int) LinkedHashMap {
	return LinkedHashMap{
		cap:     cap,
		queue:   list.New(),
		hashmap: make(map[string]*list.Element, cap),
	}
}

func (m *LinkedHashMap) IsFull() bool {
	return m.Len() >= m.cap
}

// 使用的长度，即：里面有多少个元素
func (m *LinkedHashMap) Len() int {
	return m.queue.Len()
}

// 最大容量
func (m *LinkedHashMap) Capacity() int {
	return m.cap
}

// 重新分配最大容量
func (m *LinkedHashMap) ReCapacity(cap int) *LinkedHashMap {
	if m.cap < cap { // 当前容量比目标容量小，才需要重新分配
		m.cap = cap
	}
	return m
}

// 往LRU里面添加内容
// 永远添加到head
//
// @param: key 要添加的内容的key
// @param: value 要添加的内容
func (m *LinkedHashMap) Add(key string, value interface{}) *LinkedHashMap {
	pNewNode := NewNode(key, value)
	if pElement, ok := m.hashmap[key]; ok { // 存在，则找到queue的位置，挪动到最前面
		pOldNode := (pElement.Value).(*Node)
		if !pOldNode.Equals(pNewNode) {
			pOldNode.UpdateWith(pNewNode)
		}
		m.queue.MoveToFront(pElement)
	} else if !m.IsFull() { // 不存在，且空间没有满
		pElement := m.queue.PushFront(pNewNode)
		m.hashmap[key] = pElement
	} else { // 不存在，空间也满了
		m.Eliminate()
		pElement := m.queue.PushFront(pNewNode)
		m.hashmap[key] = pElement
	}
	return m
}

// 通过key查找内容
// 如果找到返回内容，如果没有找到则有error
func (m *LinkedHashMap) Get(key string) (interface{}, error) {
	if pElement, ok := m.hashmap[key]; ok { // 存在，则找到queue的位置，挪动到最前面
		pNode := (pElement.Value).(*Node)
		return pNode.Value(), nil
	}
	return nil, errors.New(fmt.Sprintf("『%s』does not exist", key))
}

// 删除一个key
// 如果key存在，删除，并且返回，如果key不存在，返回nil
func (m *LinkedHashMap) Remove(key string) interface{} {
	if pElement, ok := m.hashmap[key]; ok { // 存在，则找到queue的位置，挪动到最前面
		m.queue.Remove(pElement)
		delete(m.hashmap, key)

		pNode := (pElement.Value).(*Node)
		return pNode.Value()
	}
	return nil
}

// 自动淘汰，本质上就是淘汰最近最少使用的
func (m *LinkedHashMap) Eliminate() interface{} {
	if element := m.queue.Back(); element != nil {
		m.queue.Remove(element)
	}
	return nil
}
