package internal

import (
	"container/list"
	"errors"
	"fmt"
)

// 顾名思义：LRU cache
type LRU struct {
	cap int
	q   *list.List // doubly link q
	kv  map[string]*list.Element
}

func NewLRU(cap int) LRU {
	return LRU{
		cap: cap,
		q:   list.New(),
		kv:  make(map[string]*list.Element, cap),
	}
}

// 使用的长度，即：里面有多少个元素
func (lru *LRU) Len() int {
	return lru.q.Len()
}

// 最大容量
func (lru *LRU) Capacity() int {
	return lru.cap
}

// 往LRU里面添加内容
// 永远添加到head
//
// @param: key 要添加的内容的key
// @param: value 要添加的内容
func (lru *LRU) Add(key string, value interface{}) *LRU {
	pNewNode := NewNode(key, value)
	if pElement, ok := lru.kv[key]; ok { // 存在，则找到queue的位置，挪动到最前面
		pOldNode := (pElement.Value).(*Node)
		if !pOldNode.Equals(pNewNode) {
			pOldNode.UpdateWith(pNewNode)
		}
		lru.q.MoveToFront(pElement)
	} else {
		pElement := lru.q.PushFront(pNewNode)
		lru.kv[key] = pElement
	}
	return lru
}

// 通过key查找内容
// 如果找到返回内容，如果没有找到则有error
func (lru *LRU) Get(key string) (interface{}, error) {
	if pElement, ok := lru.kv[key]; ok { // 存在，则找到queue的位置，挪动到最前面
		pNode := (pElement.Value).(*Node)
		return pNode.Value(), nil
	}
	return nil, errors.New(fmt.Sprintf("『%s』does not exist", key))
}
