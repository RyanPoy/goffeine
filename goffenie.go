package goffeine

import (
	"container/list"
	"goffeine/internal/node"
	"sync"
)

// A Goffeine represents a cache
// It is implemented with Window-TinyLFU algorithm
type Goffeine struct {
	data                 *sync.Map
	maximumSize          int
	window               list.List
	windowMaximumSize    int
	probation            list.List
	probationMaximumSize int
	protected            list.List
	protectedMaximumSize int
	counter              map[string]int
	expireMilliseconds   int64
	refreshMilliseconds  int64
}

func (g *Goffeine) MaximumSize() int           { return g.maximumSize }
func (g *Goffeine) ExpireMilliseconds() int64  { return g.expireMilliseconds }
func (g *Goffeine) RefreshMilliseconds() int64 { return g.refreshMilliseconds }
func (g *Goffeine) WindowMaximumSize() int     { return g.windowMaximumSize }
func (g *Goffeine) ProbationMaximumSize() int  { return g.probationMaximumSize }
func (g *Goffeine) ProtectedMaximumSize() int  { return g.protectedMaximumSize }

func (g *Goffeine) windowIsFull() bool    { return g.window.Len() == g.windowMaximumSize }
func (g *Goffeine) probationIsFull() bool { return g.probation.Len() == g.probationMaximumSize }
func (g *Goffeine) protectedIsFull() bool { return g.protected.Len() == g.protectedMaximumSize }

func (g *Goffeine) Get(key string) (any, bool) {
	ele, ok := g.data.Load(key)
	if !ok {
		return nil, false
	}
	gnode := ele.(*list.Element).Value.(*node.GoffeineNode)
	defer func() {
		go g.move(gnode)
	}()
	return gnode.Value, true
}

func (g *Goffeine) Put(key string, value any) {
	g.PutWithDelay(key, value, g.expireMilliseconds)
}

func (g *Goffeine) PutWithDelay(key string, value any, delayMilliseconds int64) {
	g.put(key, value, delayMilliseconds)
}

func (g *Goffeine) put(key string, value any, expireMilliseconds int64) {
	gnode := node.New(key, value, node.WindowPosition)
	v, ok := g.data.Load(key)
	if ok {
		oldEle := v.(*list.Element)
		gnode.Position = oldEle.Value.(*node.GoffeineNode).Position
		oldEle.Value = gnode

		if gnode.Position == node.WindowPosition {
			g.window.MoveToFront(oldEle)
		} else if gnode.Position == node.ProbationPosition {
			// todo
		} else if gnode.Position == node.ProtectedPosition {
			// todo
		}
		return
	}

	g.putToWindow(key, gnode)
	return
}

func (g *Goffeine) putToWindow(key string, gnode *node.GoffeineNode) {
	// not exist
	var ele *list.Element
	if g.windowIsFull() {
		// remove the key data
		ele = g.window.Back()
		oldKey := ele.Value.(*node.GoffeineNode).Key
		g.data.Delete(oldKey)

		// put value to back element and move it to front
		ele.Value = gnode
		g.window.MoveToFront(ele)
	} else {
		ele = g.window.PushFront(gnode)
	}
	g.data.Store(key, ele)

	//if g.windowIsFull() && g.probationIsFull() { // window 满了和probation都满了
	//	ele1 := g.window.Back()
	//	ele2 := g.probation.Back()
	//	g.evict(key, ele1, ele2)
	//}
	////else if g.window.Len() > g.windowSize
	////if g.window.Len() > g.windowSize { // window 满了，则要考虑和probation一起进行处理
	////	ele1 := g.window.Back()
	////	ele2 := g.probation.Back()
	////	g.compete(ele1, ele2)
	////}
	////g.data.Store(key, ele)
	//return ele
	//ele := v.(*list.Element)
	//ele.Value = value
}

func (g *Goffeine) move(gnode *node.GoffeineNode) {
	if gnode.Position == node.WindowPosition {
		// todo move to probation
		return
	}
	if gnode.Position == node.ProbationPosition {
		// todo move to protected
		return
	}
	return
}

//
//func (g *Goffeine) evict(key string, wEle *list.Element, pbEle *list.Element) {
//	cnt = g.counter[key]
//}
