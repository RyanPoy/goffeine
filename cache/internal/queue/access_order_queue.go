package queue

import (
	"container/list"
	"errors"
	"sync"
)

var (
	EmptyError = errors.New("Queue is empty")
)

// 顾名思义：AccessOrderQueue。里面封装了一个map和一个list
// 注意：这个LinkedHashMap是能存放的数据容量取决于cap
type AccessOrderQueue struct {
	cap     int
	queue   *list.List // doubly link queue
	hashmap sync.Map   // map[interface{}]*list.Element
}

func New(cap int) *AccessOrderQueue {
	return &AccessOrderQueue{
		cap:     cap,
		queue:   list.New(),
		hashmap: sync.Map{},
	}
}

func (q *AccessOrderQueue) IsFull() bool {
	return q.Len() >= q.cap
}

func (q *AccessOrderQueue) IsEmpty() bool {
	return q.Len() <= 0
}

// 使用的长度，即：里面有多少个元素
func (q *AccessOrderQueue) Len() int {
	return q.queue.Len()
}

// 最大容量
func (q *AccessOrderQueue) Capacity() int {
	return q.cap
}

// 重新分配最大容量
func (q *AccessOrderQueue) ReCapacity(cap int) {
	if q.cap < cap { // 当前容量比目标容量小，才需要重新分配
		q.cap = cap
	}
}

func (q *AccessOrderQueue) Contains(value interface{}) bool {
	_, ok := q.hashmap.Load(value)
	return ok
}

func (q *AccessOrderQueue) getQueueElementOfValue(value interface{}) *list.Element {
	r, _ := q.hashmap.Load(value)
	return r.(*list.Element)
}

// 添加内容
// 永远添加到tail
//
// @param: value 要添加的内容
func (q *AccessOrderQueue) Push(value interface{}) {
	if q.Contains(value) { // 存在，则找到queue的位置，并且挪动到tail
		pElement := q.getQueueElementOfValue(value)
		q.queue.MoveToBack(pElement)
	} else if !q.IsFull() { // 不存在，且空间没有满
		pElement := q.queue.PushBack(value)
		q.hashmap.Store(value, pElement)
	} else { // 不存在，空间也满了
		q.Pop()
		pElement := q.queue.PushBack(value)
		q.hashmap.Store(value, pElement)
	}
}

// 删除内容
// 永远删除head
//
// @param: value 要添加的内容
func (q *AccessOrderQueue) Pop() (interface{}, error) {
	if q.IsEmpty() {
		return nil, EmptyError
	}
	pElement := q.queue.Front()
	v := q.queue.Remove(pElement)
	q.hashmap.Delete(pElement.Value)
	return v, nil
}

func (q *AccessOrderQueue) First() (interface{}, error) {
	if q.IsEmpty() {
		return nil, EmptyError
	}
	return q.queue.Front().Value, nil
}

func (q *AccessOrderQueue) Last() (interface{}, error) {
	if q.IsEmpty() {
		return nil, EmptyError
	}
	return q.queue.Back().Value, nil
}

func (q *AccessOrderQueue) RemoveFirst() {
	q.Pop()
}

func (q *AccessOrderQueue) RemoveLast() {
	if q.IsEmpty() {
		return
	}
	pElement := q.queue.Back()
	q.queue.Remove(pElement)
	q.hashmap.Delete(pElement.Value)
}
