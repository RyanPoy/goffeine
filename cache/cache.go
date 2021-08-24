package cache

import (
	"fmt"
	"goffeine/cache/internal/node"
	"goffeine/cache/internal/queue"
	"goffeine/cache/internal/sketch"
	"math/rand"
	"sync"
)

func percentMainOf(c *Cache) float64 {
	if c != nil {
		return float64(c.wsize * 1.0 / c.cap)
	}
	return 0.99
}

func percentMainProtectedOf(c *Cache) float64 {
	if c != nil {
		return float64(c.ptsize * 1.0 / c.cap)
	}
	return 0.8
}

/*
1.队列内部设置为无界，不要在插入元素时对队列的长度做判断
2.hashmap在cache中也保留一份
3.frequencysketch 的大小设置为元素总大小n
*/
type Cache struct {
	cap        int
	sketch     *sketch.FrequencySketch
	windowQ    *queue.AccessOrderQueue
	probationQ *queue.AccessOrderQueue
	protectedQ *queue.AccessOrderQueue
	hashmap    sync.Map
	weight     int //集合当前权重，容量
	wsize      int //window大小
	ptsize     int // protectedQ size
}

func New(cap int) Cache {
	//没有异常如何判断是否成功？如果cap<0？？
	percentMain := percentMainOf(nil)
	percentMainProtected := percentMainProtectedOf(nil)

	wsize := cap - int(float64(cap)*percentMain)
	ptsize := int(percentMainProtected * float64(cap-wsize))
	pbsize := cap - wsize - ptsize

	fmt.Printf("======初始化大小：window：%d，probatio：%d，protected：%d\n", wsize, pbsize, ptsize)
	var cache = Cache{
		cap:        cap,
		sketch:     sketch.New(cap),
		windowQ:    queue.New(),
		probationQ: queue.New(),
		protectedQ: queue.New(),
		hashmap:    sync.Map{},
		weight:     0,
		wsize:      wsize,
		ptsize:     ptsize,
	}
	cache.windowQ.MaxWeight = wsize
	cache.probationQ.MaxWeight = pbsize
	cache.protectedQ.MaxWeight = ptsize
	return cache
}

func NewWith(cap, wsize, ptsize int) Cache {
	fmt.Printf("======初始化大小：window：%d，protected：%d\n", wsize, ptsize)
	var cache = Cache{
		cap:        cap,
		sketch:     sketch.New(cap),
		windowQ:    queue.New(),
		probationQ: queue.New(),
		protectedQ: queue.New(),
		hashmap:    sync.Map{},
		weight:     0,
		wsize:      wsize,
		ptsize:     ptsize,
	}
	cache.windowQ.MaxWeight = wsize
	cache.probationQ.MaxWeight = cap - wsize - ptsize
	cache.protectedQ.MaxWeight = ptsize
	return cache
}

func (c *Cache) Capacity() int {
	return c.cap
}

func (c *Cache) Weight() int {
	return c.weight
}

func (c *Cache) Contains(pNode *node.Node) bool {
	return c.windowQ.Contains(pNode) || c.probationQ.Contains(pNode) || c.protectedQ.Contains(pNode)
}

// 往cache里面添加内容
func (c *Cache) Put(key string, value interface{}) {
	fmt.Println("========执行put方法=======")
	//增加频率
	//c.sketch.Increment( []byte(key))
	//获取map的key，看key是否已经存在
	n, ok := c.hashmap.Load(key)
	if !ok {
		//不存在则新建node放入map
		newnode := node.New(key, value)
		c.weight += newnode.Weight()
		c.hashmap.Store(key, newnode)
		//fmt.Println("添加node",newnode)
		//写操作后的维护，调整队列或驱逐，后期可改为异步
		c.afterWrite(newnode)
	} else {
		oldnode := n.(*node.Node)
		//存在则覆盖nodevalue值
		//fmt.Println("old值：",oldnode.Value())
		oldnode.SetValue(value)
		//fmt.Println("new值：",oldnode.Value())
		//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
		c.afterRead(oldnode)
	} //判断权重是否上限
}

func (c *Cache) PutWithWeight(key string, value interface{}, weight int) {
	fmt.Println("========执行put方法=======")
	//增加频率
	//c.sketch.Increment( []byte(key))
	//获取map的key，看key是否已经存在
	n, ok := c.hashmap.Load(key)
	if !ok {
		//不存在则新建node放入map
		newnode := node.New(key, value)
		newnode.SetWeight(weight)
		c.weight += weight
		c.hashmap.Store(key, newnode)
		//fmt.Println("添加node",newnode)
		//写操作后的维护，调整队列或驱逐，后期可改为异步
		c.afterWrite(newnode)
	} else {
		oldnode := n.(*node.Node)
		c.weight -= oldnode.Weight()
		//存在则覆盖nodevalue值
		//fmt.Println("old值：",oldnode.Value())
		oldnode.SetValue(value)
		oldnode.SetWeight(weight)
		c.weight += weight

		//fmt.Println("new值：",oldnode.Value())
		//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
		c.afterRead(oldnode)
	} //判断权重是否上限
}
func (c *Cache) GetWithWeight(key string) (interface{}, int) {
	fmt.Println("========执行Get方法=======")
	n, ok := c.hashmap.Load(key)
	if ok { // 找到了
		oldnode := n.(*node.Node)
		//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
		c.afterRead(oldnode)
		return oldnode.Value(), oldnode.Weight()
	} else {
		return nil, 0
	}
}
func (c *Cache) Get(key string) interface{} {
	fmt.Println("========执行Get方法=======")
	n, ok := c.hashmap.Load(key)
	if !ok {
		return nil
	}

	oldnode := n.(*node.Node)
	//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
	c.afterRead(oldnode)
	return oldnode.Value()
}

func (c *Cache) afterWrite(newnode *node.Node) {
	//后期使用读写缓冲区可增加此函数功能，当前主要是执行移动操作
	c.addTask(newnode) //将node添加到三个队列中的一个
	c.maintenance()
}
func (c *Cache) afterRead(node *node.Node) {
	//后期增加命中次数统计等操作时可增加此函数功能，当前主要是执行移动操作
	c.updateTask(node) //访问后更新元素
	c.maintenance()
}

func (c *Cache) maintenance() {
	//执行维护工作，查看是否过期，队列是否满等操作。
	//当前主要是执行队列满的驱逐策略
	c.evictEntries()
	c.demoteFromMainProtected() //收缩protected大小
}
func (c *Cache) updateTask(nod *node.Node) {
	//执行更新元素位置工作
	if nod.IsBelongsToWindow() {
		if c.wsize >= nod.Weight() { //当前node值小于weight最大值//后期1改为node权重，代表node不超过window最大大小
			c.onAccess(nod)
		} else if c.windowQ.Contains(nod) { //说明超过最大值，移到队首等待被清除
			c.windowQ.MoveToFront(nod)
		}
	} else if nod.IsBelongsToProbation() {
		if c.probationQ.MaxWeight >= nod.Weight() { //后期1改为node权重，代表node不超过probation最大大小
			c.onAccess(nod)
		} else { //说明超过最大值，移除然后放入window队头等待驱逐比较
			c.probationQ.Remove(nod)
			c.windowQ.AddFirst(nod)
			//修改window,probation权重，此时不需要
		}
	} else {
		if c.protectedQ.MaxWeight >= nod.Weight() { //后期1改为node权重，代表node不超过probation最大大小
			c.onAccess(nod)
		} else { //说明超过最大值，移除然后放入window队头等待驱逐比较
			c.protectedQ.Remove(nod)
			c.windowQ.AddFirst(nod)
			//修改window权重，此时不需要
		}
	}
	//由于更新没有数量变化，只有权重变化，因此下面不用修改权重。
}
func (c *Cache) addTask(nod *node.Node) {
	//执行增添元素任务
	//判断权重大小，并对fsketch进行扩容
	//增加未命中次数
	//过期等策略
	//获取锁执行写操作
	//当前只用判断大小并插入即可
	if c.wsize >= nod.Weight() { //说明没有超过最大值，移到队尾部
		c.windowQ.Push(nod)
		c.onAccess(nod)
	} else { //代表node超过window最大大小
		c.windowQ.AddFirst(nod)
	}
}
func (c *Cache) evictEntries() {
	//容量满时驱逐元素
	candidates := c.evictFromWindow()
	c.evictFromMain(candidates)
}

func (c *Cache) evictFromWindow() int {
	var candidates = 0
	for c.windowQ.Weight() > c.wsize { //后期改为权重大于最大权重
		var nod, _ = c.windowQ.First()

		if nod == nil {
			break
		}
		nod.InProbation()
		c.windowQ.Remove(nod)
		c.probationQ.Push(nod)
		candidates++
		//修改权重
	}
	return candidates
}

func (c *Cache) getVictim() (*node.Node, bool) {
	nod, ok := c.probationQ.First()
	if !ok {
		nod, ok = c.protectedQ.First()
	}
	if !ok {
		nod, ok = c.windowQ.First()
	}
	return nod, ok
}

func (c *Cache) getCandidate(fromWindow bool) (*node.Node, bool) {
	if fromWindow {
		return c.windowQ.First()
	}
	return c.probationQ.Last()
}

func (c *Cache) getQueueByNode(nod *node.Node) *queue.AccessOrderQueue {
	if nod.IsBelongsToWindow() {
		return c.windowQ
	} else if nod.IsBelongsToProtected() {
		return c.probationQ
	}
	return c.protectedQ
}
func (c *Cache) evictFromMain(candidates int) {
	// 首先默认选择probation的队头和队尾作为victim和candidate，参与淘汰；
	//若 FrequencyCandidate < 5，则淘汰c；
	//若 5 <= FrequencyCandidate < FrequencyVictim:
	// 随机淘汰 victim 或者  candidte
	//若 FrequencyCandidate > FrequencyVictim 则淘汰v
	for c.weight > c.cap {
		victim, ok1 := c.getVictim()
		candidate, ok2 := c.getCandidate(candidates <= 0)
		if !ok1 && !ok2 {
			return
		}
		if ok1 && !ok2 { // victim有，candidate没有
			c.evictEntry(victim)
			continue
		}
		if !ok1 && ok2 { // victim没有，candidate有
			c.evictEntry(candidate)
			candidates--
			continue
		}
		// 执行到这里，则：ok1 && ok2 ，表示 victim 和 candidate都有
		if victim == candidate {
			c.evictEntry(candidate)
			candidates--
			continue
		}
		if candidate.Weight() > c.cap {
			c.evictEntry(candidate)
			candidates--
			continue
		}

		candidates--
		if c.admit(candidate, victim) { // 淘汰victim
			c.evictEntry(victim)
		} else {
			c.evictEntry(candidate)
		}
	}
}

func (c *Cache) evictEntry(nod *node.Node) bool {
	fmt.Println("=========执行evictEntry========")
	//尝试根据给定的删除原因驱逐条目。由于当前只会根据容量驱逐，因此不设其他参数
	//分别判断三种情况导致的驱逐，并更具情况设置是否复活，synchronized执行
	//移除node监听器
	//移除hashmap里的元素
	fmt.Println("=========执行evictEntry:", nod.Key())
	c.hashmap.Delete(nod.Key())
	c.weight -= nod.Weight()
	if nod.IsBelongsToWindow() {
		c.windowQ.Remove(nod)
	} else if nod.IsBelongsToProbation() {
		c.probationQ.Remove(nod)
	} else {
		c.protectedQ.Remove(nod)
	}
	//c.makeDead(nod) 修改nod所在队列的权重当前不需要
	return true
}

func (c *Cache) makeDead(nod *node.Node) { //加锁完成，修改权重
}
func (c *Cache) admit(candidate *node.Node, victim *node.Node) bool { //window到probation晋升
	victimFreq := c.sketch.Frequency(victim)
	candidateFreq := c.sketch.Frequency(candidate)
	fmt.Printf("========执行admit方法：victim：%s：%d次,candidate:%s：%d=======\n", victim.Key(), victimFreq, candidate.Key(), candidateFreq)
	if victimFreq < candidateFreq {
		return true
	} else if candidateFreq <= 5 {
		//最大频率为 15，在重置历史记录后减半为 7。
		// 攻击利用热门候选人被拒绝而有利于热门受害者。 温暖候选者的阈值减少了随机接受的次数，以尽量减少对命中率的影响。
		return false
	}
	random := rand.Int()
	return (random & 127) == 0
}

func (c *Cache) removalTask(nod *node.Node) {
	//执行删除元素任务
}
func (c *Cache) onAccess(nod *node.Node) {
	//更新结点位置
	c.sketch.Increment(nod) //增加访问频率
	if nod.IsBelongsToWindow() {
		c.reorder(c.windowQ, nod)
	} else if nod.IsBelongsToProbation() {
		c.reorderProbation(nod)
	} else {
		c.reorder(c.protectedQ, nod)
	}
}

func (c *Cache) reorder(queue *queue.AccessOrderQueue, nod *node.Node) { //将节点移动至指定队列尾部
	if queue.Contains(nod) {
		queue.MoveToBack(nod)
	}
}
func (c *Cache) reorderProbation(nod *node.Node) { //从probation队列中移动元素
	if !c.probationQ.Contains(nod) {
		return
	}
	if nod.Weight() > c.protectedQ.MaxWeight {
		//若大小超过protected大小，则放入pb的尾部
		c.reorder(c.probationQ, nod)
	} else {
		//修改pt的权重大小，但现在不考虑权重故不需要
		c.probationQ.Remove(nod)
		c.protectedQ.Push(nod)
		nod.InProtected()
	}
}

func (c *Cache) demoteFromMainProtected() {
	for c.protectedQ.Weight() > c.ptsize {
		nod, ok := c.protectedQ.First()
		if ok {
			c.protectedQ.Remove(nod)
			c.probationQ.Push(nod)
			nod.InProbation()
		}
	}
}
