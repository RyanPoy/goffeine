package cache

import (
	"container/list"
	"errors"
	"fmt"
)


type LRU struct {
	cap int
	queue *list.List // doubly link queue
	cache map[string]*list.Element
}

func NewLRU(cap int) LRU {
	return LRU{
		cap: cap,
		queue: list.New(),
		cache: make(map[string]*list.Element, cap),
	}
}

func (lru *LRU) Len() int {
	return lru.queue.Len()
}


func (lru *LRU) Capacity() int {
	return lru.cap
}


// 添加一个node到LRU
// 永远添加到head
//
func (lru *LRU) Add(key string, value interface{}) *LRU {
	if pElement, ok := lru.cache[key]; ok { // 存在，则找到queue的位置，挪动到最前面
		lru.queue.MoveToFront(pElement)
	} else {
		pElement := lru.queue.PushFront(NewNode(key, value))
		lru.cache[key] = pElement
	}
	return lru
}

func (lru *LRU) Get(key string) (interface{}, error) {
	if pElement, ok := lru.cache[key]; ok { // 存在，则找到queue的位置，挪动到最前面
		node := (pElement.Value).(Node)
		return node.Value(), nil
	}
	return nil, errors.New(fmt.Sprintf("『%s』does not exist", key))
}


