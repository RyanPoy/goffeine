package internal

import (
	"container/list"
	"errors"
	"fmt"
)

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

func (lru *LRU) Len() int {
	return lru.q.Len()
}

func (lru *LRU) Capacity() int {
	return lru.cap
}

// 添加一个node到LRU
// 永远添加到head
//
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

func (lru *LRU) Get(key string) (interface{}, error) {
	if pElement, ok := lru.kv[key]; ok { // 存在，则找到queue的位置，挪动到最前面
		pNode := (pElement.Value).(*Node)
		return pNode.Value(), nil
	}
	return nil, errors.New(fmt.Sprintf("『%s』does not exist", key))
}
