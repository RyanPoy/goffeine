package lru

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

// 顾名思义：LRU。里面封装了一个map和一个list
// 注意：这个LinkedHashMap是能存放的数据容量取决于cap
type LRU struct {
	cap     int
	queue   *list.List // doubly link queue
	hashmap sync.Map   // map[string]*list.Element
}

func New(cap int) *LRU {
	return &LRU{
		cap:     cap,
		queue:   list.New(),
		hashmap: sync.Map{},
	}
}

func (lru *LRU) IsFull() bool {
	return lru.Len() >= lru.cap
}

// 使用的长度，即：里面有多少个元素
func (lru *LRU) Len() int {
	return lru.queue.Len()
}

// 最大容量
func (lru *LRU) Capacity() int {
	return lru.cap
}

// 重新分配最大容量
func (lru *LRU) ReCapacity(cap int) *LRU {
	if lru.cap < cap { // 当前容量比目标容量小，才需要重新分配
		lru.cap = cap
	}
	return lru
}

func (lru *LRU) Contains(key string) bool {
	_, ok := lru.hashmap.Load(key)
	return ok
}

// 往LRU里面添加内容
// 永远添加到tail
//
// @param: key 要添加的内容的key
// @param: value 要添加的内容
// @return 如果有自动淘汰的数据，则返回它。否则，返回nil
func (lru *LRU) Add(key string, value interface{}) interface{} {
	var elementEliminated interface{} = nil
	if lru.Contains(key) {
		// 存在，则找到queue的位置，并且挪动到tail
		hashMapValue, _ := lru.hashmap.Load(key)
		pElement := hashMapValue.(*list.Element)
		if pElement.Value != value {
			// 不相等，表示要进行更新内容的操作
			pElement.Value = value
		}
		lru.queue.MoveToBack(pElement)
	} else if !lru.IsFull() {
		// 不存在，且空间没有满
		pElement := lru.queue.PushBack(value)
		lru.hashmap.Store(key, pElement)
	} else { // 不存在，空间也满了
		elementEliminated = lru.Eliminate()
		pElement := lru.queue.PushBack(value)
		lru.hashmap.Store(key, pElement)
	}
	return elementEliminated
}

// 通过key查找内容
// 如果找到返回内容，如果没有找到则有error
func (lru *LRU) Get(key string) (interface{}, error) {
	if hashMapValue, ok := lru.hashmap.Load(key); ok { // 存在，则找到queue的位置，挪动到最前面
		pElement := hashMapValue.(*list.Element)
		return pElement.Value, nil
	}
	return nil, errors.New(fmt.Sprintf("『%s』does not exist", key))
}

// 删除一个key
// 如果key存在，删除，并且返回，如果key不存在，返回nil
func (lru *LRU) Remove(key string) interface{} {
	if hashMapValue, ok := lru.hashmap.Load(key); ok { // 存在，则找到queue的位置，挪动到最前面
		pElement := hashMapValue.(*list.Element)
		lru.queue.Remove(pElement)
		lru.hashmap.Delete(key)
		return pElement.Value
	}
	return nil
}

// 自动淘汰最近最少使用的，从尾部淘汰
func (lru *LRU) Eliminate() interface{} {
	if pElement := lru.queue.Front(); pElement != nil {
		return lru.queue.Remove(pElement)
	}
	return nil
}
