package queue

import (
	"container/list"
	"goffeine/cache/internal/node"
	"sync"
)

// 顾名思义：AccessOrderQueue。里面封装了一个map和一个list
// 注意：这个LinkedHashMap是能存放的数据容量取决于cap
type AccessOrderQueue struct {
	queue     *list.List // doubly link queue
	data      sync.Map   // map[[]byte]*list.Element
	MaxWeight int
}

func New() *AccessOrderQueue {
	return &AccessOrderQueue{
		queue:     list.New(),
		data:      sync.Map{},
		MaxWeight: 0,
	}
}

func (q *AccessOrderQueue) Weight() int {
	w := 0
	for e := q.queue.Front(); e != nil; e = e.Next() {
		w += e.Value.(*node.Node).Weight()
	}
	return w
}

func (q *AccessOrderQueue) IsEmpty() bool {
	return q.Weight() <= 0
}

func (q *AccessOrderQueue) Contains(pNode *node.Node) bool {
	_, ok := q.data.Load(pNode.Key())
	return ok
}

func (q *AccessOrderQueue) GetQueueElementBy(pNode *node.Node) *list.Element {
	r, _ := q.data.Load(pNode.Key())
	return r.(*list.Element)
}


// 添加内容
// 永远添加到tail
// @param: value 要添加的内容
func (q *AccessOrderQueue) Push(pNode *node.Node) {
	if q.Contains(pNode) { // 存在，则找到queue的位置，并且挪动到tail
		pElement := q.GetQueueElementBy(pNode)
		q.queue.MoveToBack(pElement)
	} else { // 不存在，空间也满了
		pElement := q.queue.PushBack(pNode)
		q.data.Store(pNode.Key(), pElement)
	}
}

func (q *AccessOrderQueue) MoveToBack(pNode *node.Node) {
	if q.Contains(pNode) { // 存在，则找到queue的位置，并且挪动到tail
		pElement := q.GetQueueElementBy(pNode)
		q.queue.MoveToBack(pElement)
	} else {
		q.Push(pNode)
	}
}

func (q *AccessOrderQueue) Remove(pNode *node.Node) {
	//移除nod结点
	if q.Contains(pNode) {
		pElement := q.GetQueueElementBy(pNode)
		q.data.Delete(pNode.Key())
		q.queue.Remove(pElement)
	}
}

func (q *AccessOrderQueue) AddFirst(pNode *node.Node) {
	//添加到队头
	if q.Contains(pNode) { // 存在，则找到queue的位置，并且挪动到tail
		pElement := q.GetQueueElementBy(pNode)
		q.queue.MoveToFront(pElement)
	}
	pElement := q.queue.PushFront(pNode)
	q.data.Store(pNode.Key(), pElement)
}
func (q *AccessOrderQueue) MoveToFront(pNode *node.Node) {
	//将nod结点移到队首
	if q.Contains(pNode) { // 存在，则找到queue的位置，并且挪动到tail
		pElement := q.GetQueueElementBy(pNode)
		q.queue.MoveToFront(pElement)
	} else {
		q.AddFirst(pNode)
	}
}

// 删除内容
// 永远删除head
// @param: value 要添加的内容
func (q *AccessOrderQueue) Pop() (*node.Node, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	pElement := q.queue.Front()
	v := q.queue.Remove(pElement)
	q.data.Delete(pElement.Value)
	return v.(*node.Node), true
}

func (q *AccessOrderQueue) First() (*node.Node, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	pElement := q.queue.Front()
	return pElement.Value.(*node.Node), true
}

func (q *AccessOrderQueue) Last() (*node.Node, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	pElement := q.queue.Back()
	return pElement.Value.(*node.Node), true
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
	q.data.Delete(pElement.Value)
}
