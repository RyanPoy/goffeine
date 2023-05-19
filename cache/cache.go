package cache

import (
	"fmt"
	"goffeine/cache/internal/node"
	"goffeine/cache/internal/queue"
	"goffeine/cache/internal/sketch"
	"math/rand"
	"sync"
	"time"
)

/*
1.队列内部设置为无界，不要在插入元素时对队列的长度做判断
2.hashmap在cache中也保留一份
3.frequencysketch 的大小设置为元素总大小n
*/
type Cache struct {
	maxWeight   int
	sketch      *sketch.FrequencySketch
	windowQ     *queue.AccessOrderQueue
	probationQ  *queue.AccessOrderQueue
	protectedQ  *queue.AccessOrderQueue
	hashmap     sync.Map
	Weight      int //集合当前权重，容量
	wMaxWeight  int //window大小
	ptMaxWeight int // protectedQ size
}

func New(maxWeight int) Cache {
	//没有异常如何判断是否成功？如果cap<0？？
	percentMain, percentMainProtected := 0.99, 0.8
	wMaxWeight := maxWeight - int(float64(maxWeight)*percentMain)
	ptMaxWeight := int(percentMainProtected * float64(maxWeight-wMaxWeight))
	return NewWith(maxWeight, wMaxWeight, ptMaxWeight)
}

func NewWith(maxWeight, wMaxWeight, ptMaxWeight int) Cache {
	fmt.Printf("======初始化大小：window：%d，protected：%d\n", wMaxWeight, ptMaxWeight)
	var cache = Cache{
		maxWeight:   maxWeight,
		sketch:      sketch.New(maxWeight),
		windowQ:     queue.New(),
		probationQ:  queue.New(),
		protectedQ:  queue.New(),
		hashmap:     sync.Map{},
		Weight:      0,
		wMaxWeight:  wMaxWeight,
		ptMaxWeight: ptMaxWeight,
	}
	cache.windowQ.MaxWeight = wMaxWeight
	cache.probationQ.MaxWeight = maxWeight - wMaxWeight - ptMaxWeight
	cache.protectedQ.MaxWeight = ptMaxWeight
	return cache
}

func (c *Cache) percentMainOf() float64 {
	return float64(c.wMaxWeight * 1.0 / c.maxWeight)
}

func (c *Cache) percentMainProtectedOf() float64 {
	return float64(c.ptMaxWeight * 1.0 / c.maxWeight)
}

// 往cache里面添加内容
func (c *Cache) Put(key string, value interface{}) {
	c.PutWithWeight(key, value, 1) // node的默认权重为1
}

func (c *Cache) PutWithWeight(key string, value interface{}, weight int) {
	fmt.Println("========执行put方法=======")
	//获取map的key，看key是否已经存在
	n, ok := c.hashmap.Load(key)
	if !ok {
		//不存在则新建node放入map
		newNode := node.New(key, value)
		newNode.Weight = weight
		c.Weight += weight
		c.hashmap.Store(key, newNode)
		//fmt.Println("添加node",newNode)
		//写操作后的维护，调整队列或驱逐，后期可改为异步
		c.afterWrite(newNode)
	} else {
		oldNode := n.(*node.Node)
		c.Weight -= oldNode.Weight
		//存在则覆盖nodevalue值
		//fmt.Println("old值：",oldNode.Value())
		oldNode.Value = value
		oldNode.Weight = weight
		c.Weight += weight

		//fmt.Println("new值：",oldNode.Value())
		//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
		c.afterRead(oldNode)
	} //判断权重是否上限
}

func (c *Cache) Get(key string) interface{} {
	value, _ := c.GetWithWeight(key)
	return value
}

func (c *Cache) GetWithWeight(key string) (interface{}, int) {
	fmt.Println("========执行Get方法=======")
	n, ok := c.hashmap.Load(key)
	if !ok {
		return nil, 0
	}
	pNode := n.(*node.Node)
	//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
	c.afterRead(pNode)
	return pNode.Value, pNode.Weight
}

func (c *Cache) afterWrite(newNode *node.Node) {
	//后期使用读写缓冲区可增加此函数功能，当前主要是执行移动操作
	c.addTask(newNode) //将node添加到三个队列中的一个
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
	if nod.IsInWindow() {
		if c.wMaxWeight >= nod.Weight { //当前node值小于weight最大值，后期1改为node权重，代表node不超过window最大大小
			c.onAccess(nod)
		} else { //说明超过最大值，移到队首等待被清除
			c.windowQ.MoveToFirst(nod)
		}
	} else if nod.IsInProbation() {
		if c.probationQ.MaxWeight >= nod.Weight { //后期1改为node权重，代表node不超过probation最大大小
			c.onAccess(nod)
		} else { //说明超过最大值，移除然后放入window队头等待驱逐比较
			c.probationQ.Remove(nod)
			c.windowQ.LinkFirst(nod)
			//修改window,probation权重，此时不需要
		}
	} else {
		if c.protectedQ.MaxWeight >= nod.Weight { //后期1改为node权重，代表node不超过probation最大大小
			c.onAccess(nod)
		} else { //说明超过最大值，移除然后放入window队头等待驱逐比较
			c.protectedQ.Remove(nod)
			c.windowQ.LinkFirst(nod)
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
	if c.wMaxWeight >= nod.Weight { //说明没有超过最大值，移到队尾部
		c.windowQ.LinkLast(nod)
		c.onAccess(nod)
	} else { //代表node超过window最大大小
		c.windowQ.LinkFirst(nod)
	}
}
func (c *Cache) evictEntries() {
	//容量满时驱逐元素
	candidates := c.evictFromWindow()
	c.evictFromMain(candidates)
}

func (c *Cache) evictFromWindow() int {
	var candidates = 0
	for c.windowQ.Weight() > c.wMaxWeight { //后期改为权重大于最大权重
		var nod, _ = c.windowQ.First()

		if nod == nil {
			break
		}
		nod.InProbation()
		c.windowQ.Remove(nod)
		c.probationQ.LinkLast(nod)
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

func (c *Cache) evictFromMain(candidates int) {
	// 首先默认选择probation的队头和队尾作为victim和candidate，参与淘汰；
	//若 FrequencyCandidate < 5，则淘汰c；
	//若 5 <= FrequencyCandidate < FrequencyVictim:
	// 随机淘汰 victim 或者  candidate
	//若 FrequencyCandidate > FrequencyVictim 则淘汰v
	for c.Weight > c.maxWeight {
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
		if victim == candidate || candidate.Weight > c.maxWeight {
			c.evictEntry(candidate)
			candidates--
			continue
		}

		if c.admit(candidate, victim) { // 淘汰victim
			c.evictEntry(victim)
		} else {
			c.evictEntry(candidate)
		}
		candidates--
	}
}

func (c *Cache) evictEntry(nod *node.Node) bool {
	fmt.Println("=========执行evictEntry========")
	//尝试根据给定的删除原因驱逐条目。由于当前只会根据容量驱逐，因此不设其他参数
	//分别判断三种情况导致的驱逐，并更具情况设置是否复活，synchronized执行
	//移除node监听器
	//移除hashmap里的元素
	fmt.Println("=========执行evictEntry:", nod.Key)
	c.hashmap.Delete(nod.Key)
	c.Weight -= nod.Weight
	if nod.IsInWindow() {
		c.windowQ.Remove(nod)
	} else if nod.IsInProbation() {
		c.probationQ.Remove(nod)
	} else {
		c.protectedQ.Remove(nod)
	}
	//c.makeDead(nod) 修改nod所在队列的权重当前不需要
	return true
}

func (c *Cache) makeDead(nod *node.Node) { //加锁完成，修改权重
}

func (c *Cache) rnd() bool {
	rand.Seed(time.Now().Unix())
	return rand.Int()&127 == 0
}

func (c *Cache) admit(candidate *node.Node, victim *node.Node) bool { //window到probation晋升
	victimFreq := c.sketch.Frequency(victim)
	candidateFreq := c.sketch.Frequency(candidate)
	fmt.Printf("========执行admit方法：victim：%s：%d次,candidate:%s：%d=======\n", victim.Key, victimFreq, candidate.Key, candidateFreq)
	if victimFreq < candidateFreq {
		return true
	}
	if candidateFreq <= 5 {
		//最大频率为 15，在重置历史记录后减半为 7。
		// 攻击利用热门候选人被拒绝而有利于热门受害者。 温暖候选者的阈值减少了随机接受的次数，以尽量减少对命中率的影响。
		return false
	}
	return c.rnd()
}

func (c *Cache) removalTask(nod *node.Node) {
	//执行删除元素任务
}

func (c *Cache) onAccess(nod *node.Node) {
	//更新结点位置
	c.sketch.Increment(nod) //增加访问频率
	if nod.IsInWindow() {
		c.windowQ.MoveToLast(nod)
	} else if nod.IsInProtected() {
		c.protectedQ.MoveToLast(nod)
	} else { // IsInProbation
		if nod.Weight > c.protectedQ.MaxWeight {
			//若大小超过protected大小，则放入pb的尾部
			c.probationQ.MoveToLast(nod)
		} else {
			//修改pt的权重大小，但现在不考虑权重故不需要
			c.probationQ.Remove(nod)
			c.protectedQ.LinkLast(nod)
			nod.InProtected()
		}
	}
}

func (c *Cache) demoteFromMainProtected() {
	for c.protectedQ.Weight() > c.ptMaxWeight {
		nod, ok := c.protectedQ.First()
		if ok {
			c.protectedQ.Remove(nod)
			c.probationQ.LinkLast(nod)
			nod.InProbation()
		}
	}
}
