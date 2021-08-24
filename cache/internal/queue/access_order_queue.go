package queue

import (
	"container/list"
	"errors"
	"goffeine/cache/internal/node"
	"sync"
)

var (
	EmptyError = errors.New("Queue is empty")
)


// 顾名思义：AccessOrderQueue。里面封装了一个map和一个list
// 注意：这个LinkedHashMap是能存放的数据容量取决于cap
type AccessOrderQueue struct {
	queue   *list.List // doubly link queue
	hashmap sync.Map   // map[[]byte]*list.Element
	weight int
	MaxWeight int
}

func New() *AccessOrderQueue {
	return &AccessOrderQueue{
		queue:   list.New(),
		hashmap: sync.Map{},
		weight:	0,
		MaxWeight: 0,
	}
}
func (q *AccessOrderQueue) Weight() int {
	 w:=0
	for e := q.queue.Front(); e != nil; e = e.Next() {
		w+=e.Value.(*node.Node).Weight()
	}
	q.weight=w
	return q.weight
}
func (q *AccessOrderQueue) SetWeiht(w int){
	q.weight=w
}
//func (q *AccessOrderQueue) IsFull() bool {
//	return q.Weight() >= q.cap
//}

func (q *AccessOrderQueue) IsEmpty() bool {
	return q.Weight() <= 0
}

// 使用的长度，即：里面有多少个元素
//func (q *AccessOrderQueue) Len() int {
//	return q.queue.Len()
//}





func (q *AccessOrderQueue) Contains(pNode *node.Node) bool {
	_, ok := q.hashmap.Load(pNode.Key())
	return ok
}

func (q *AccessOrderQueue) GetQueueElementBy(pNode *node.Node) *list.Element {
	r, _ := q.hashmap.Load(pNode.Key())
	return r.(*list.Element)
}

func (q *AccessOrderQueue) GetNextNodeBy(pNode *node.Node) *node.Node {
	r, _ := q.hashmap.Load(pNode.Key())
	if r==nil {
		return nil
	}
	p:=r.(*list.Element).Next()
	if(p ==nil) {return nil}
	return p.Value.(*node.Node)
}
func (q *AccessOrderQueue) GetPrevNodeBy(pNode *node.Node) *node.Node {
	r, _ := q.hashmap.Load(pNode.Key())
	if r==nil {
		return nil
	}
	p:=r.(*list.Element).Prev()
	if(p ==nil) {return nil}
	return p.Value.(*node.Node)
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
		q.hashmap.Store(pNode.Key(), pElement)
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


func (q *AccessOrderQueue) Remove(pNode *node.Node)  {
	//移除nod结点
	if (q.Contains(pNode)) {
		pElement := q.GetQueueElementBy(pNode)
		q.hashmap.Delete(pNode.Key())
		q.queue.Remove(pElement)
	}
}
func (q*AccessOrderQueue) AddFirst (pNode *node.Node){
	//添加到队头
	if q.Contains(pNode) { // 存在，则找到queue的位置，并且挪动到tail
		pElement := q.GetQueueElementBy(pNode)
		q.queue.MoveToFront(pElement)
	}
		pElement := q.queue.PushFront(pNode)
		q.hashmap.Store(pNode.Key(), pElement)
}
func (q *AccessOrderQueue) MoveToFront(pNode *node.Node){
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
	q.hashmap.Delete(pElement.Value)
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
	q.hashmap.Delete(pElement.Value)
}
