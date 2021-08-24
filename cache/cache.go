package cache

import (
	"fmt"
	"goffeine/cache/internal/node"
	"goffeine/cache/internal/queue"
	"goffeine/cache/internal/sketch"
	"math/rand"
	"sync"
)
/*
1.队列内部设置为无界，不要在插入元素时对队列的长度做判断
2.hashmap在cache中也保留一份
3.frequencysketch 的大小设置为元素总大小n
*/
type Cache struct {
	cap        int
	sketch    *sketch.FrequencySketch
	windowQ    *queue.AccessOrderQueue
	probationQ *queue.AccessOrderQueue
	protectedQ *queue.AccessOrderQueue
	hashmap    sync.Map
	weight     int //集合当前权重，容量
	wsize      int //window大小
	percentMain           float64 //main空间比列
	percentMainProtected float64 //protected比例
}

func New(cap int) Cache {
	//没有异常如何判断是否成功？如果cap<0？？
	percentMain := 0.99
	percentMainProtected := 0.8
	wsize := cap - int(float64(cap)*percentMain)
	ptsize := int(percentMainProtected * float64(cap-wsize))
	pbsize := cap - wsize - ptsize

	fmt.Printf("======初始化大小：window：%d，probatio：%d，protected：%d\n",wsize,pbsize,ptsize)
	var cache = Cache{
		cap:                    cap,
		sketch:                 sketch.New(cap),
		windowQ:				queue.New(),
		probationQ:				queue.New(),
		protectedQ:				queue.New(),
		hashmap:                sync.Map{},
		weight:                 0,
		percentMain:           percentMain,
		percentMainProtected: percentMainProtected,
		wsize:                  wsize,
	}
	cache.windowQ.MaxWeight = wsize
	cache.probationQ.MaxWeight = pbsize
	cache.protectedQ.MaxWeight = ptsize
	return cache
}

func (c *Cache) Capacity() int {
	return c.cap
}

//func (c *Cache) Len() int {
//
//}
func(c *Cache) Weight() int{
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
	n, _ := c.hashmap.Load(key)
	if n == nil {
		//不存在则新建node放入map
		newnode := node.New(key, value)
		c.weight+=newnode.Weight()
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
	n, _ := c.hashmap.Load(key)
	if n == nil {
		//不存在则新建node放入map
		newnode := node.New(key, value)
		newnode.SetWeight(weight)
		c.weight+=weight
		c.hashmap.Store(key, newnode)
		//fmt.Println("添加node",newnode)
		//写操作后的维护，调整队列或驱逐，后期可改为异步
		c.afterWrite(newnode)
	} else {
		oldnode := n.(*node.Node)
		c.weight-=oldnode.Weight()
		//存在则覆盖nodevalue值
		//fmt.Println("old值：",oldnode.Value())
		oldnode.SetValue(value)
		oldnode.SetWeight(weight)
		c.weight+=weight

		//fmt.Println("new值：",oldnode.Value())
		//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
		c.afterRead(oldnode)
	} //判断权重是否上限
}
func (c *Cache) GetWithWeight(key string) (interface{},int) {
	fmt.Println("========执行Get方法=======")
	n, _ := c.hashmap.Load(key)
	if n != nil {
		oldnode := n.(*node.Node)
		//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
		c.afterRead(oldnode)
		return oldnode.Value(),oldnode.Weight()
	} else {
		return nil,0
	}
}
func (c *Cache) Get(key string) interface{} {
	fmt.Println("========执行Get方法=======")
	n, _ := c.hashmap.Load(key)
	if n != nil {
		oldnode := n.(*node.Node)
		//由于没有新值插入，维护与读操作后的维护相同，后期可改为异步
		c.afterRead(oldnode)
		return oldnode.Value()
	} else {
		return nil
	}
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
}
func (c *Cache) updateTask(nod *node.Node) {
	//执行更新元素位置工作
	//position := nod.Position()
	//switch position {
	//case WINDOW: //处于window
	if (nod.IsBelongsToWindow()) {
		if c.wsize >= nod.Weight() { //当前node值小于weight最大值//后期1改为node权重，代表node不超过window最大大小
			c.onAccess(nod)
		} else if c.windowQ.Contains(nod) { //说明超过最大值，移到队首等待被清除
			c.windowQ.MoveToFront(nod)
		}
	} else if (nod.IsBelongsToProbation()) {
		//case PROBATION: //处于probation
		if c.probationQ.MaxWeight >= nod.Weight() { //后期1改为node权重，代表node不超过probation最大大小
			c.onAccess(nod)
		} else { //说明超过最大值，移除然后放入window队头等待驱逐比较
			c.probationQ.Remove(nod)
			c.windowQ.AddFirst(nod)
			//修改window,probation权重，此时不需要
		}
	} else {
	//case PROTECTED: //处于protected状态
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
	if c.wsize >= nod.Weight() {//说明没有超过最大值，移到队尾部
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
func (c *Cache) evictFromMain(candidates int) {
	victimQueue := node.PROBATION //probation 受害者由probation选出
	victim , _ := c.probationQ.First()
	candidate, _ := c.probationQ.Last()
	for c.weight > c.cap {
		fmt.Println("========执行evictFromMain：最大容量已满=========")
		if candidates == 0 {
			candidate, _ = c.windowQ.Last()
		}
		if candidate == nil && victim == nil { //如果从window和probation里都没有元素，就从protected里面取
			fmt.Println("========执行evictFromMain：candidate == nil && victim == nil=========")
			if victimQueue == node.PROBATION {
				victim, _ = c.protectedQ.First()
				victimQueue = node.PROTECTED //protected
				continue
			} else if victimQueue == node.PROTECTED {
				victim, _ = c.windowQ.First()
				victimQueue = 0 //window
				continue
			}
			break
		}
		//此处权重为0情况，因为计数场景权重都为1，后期根据情况扩展
		//如果只有一个权重立即驱逐
		if victim == nil {
			previous := c.GetPreviousInAccessOrder(candidate)
			evict := candidate
			candidate = previous
			candidates--
			c.evictEntry(evict)
			continue
		} else if candidate == nil {
			evict := victim
			victim = c.GetNextInAccessOrder(victim)
			c.evictEntry(evict)
			continue
		}
		//忽略值引用的情况
		//竞选者本身权重超过最大值直接驱逐，不考虑
		if candidate.Weight()>c.cap{
			evict:=candidate
			if candidates >0{
				candidate=c.GetPreviousInAccessOrder(candidate)
			}else{
				candidate=c.GetNextInAccessOrder(candidate)
			}
			candidates--
			c.evictEntry(evict)
			continue
		}
		//驱逐频率最低的条目
		candidates--
		fmt.Println("========执行evictFromMain：admit（）=========")
		if c.admit(candidate, victim) { //判断是否驱逐受害者，true驱逐受害者
			evict := victim
			victim = c.GetNextInAccessOrder(victim)
			c.evictEntry(evict)
			candidate = c.GetPreviousInAccessOrder(candidate)
		} else { //否，驱逐竞争者
			evict := candidate
			if candidates > 0 {
				candidate = c.GetPreviousInAccessOrder(candidate)
			} else {
				candidate = c.GetNextInAccessOrder(candidate)
			}
			c.evictEntry(evict)
		}
	}
}

func (c *Cache)GetPreviousInAccessOrder(nod *node.Node) *node.Node{
	if nod.IsBelongsToWindow() {
		return c.windowQ.GetNextNodeBy(nod)
	} else if nod.IsBelongsToProbation() {
		return c.probationQ.GetNextNodeBy(nod)
	}
	return c.protectedQ.GetNextNodeBy(nod)
}
func (c *Cache)GetNextInAccessOrder(nod *node.Node) *node.Node{
	if nod.IsBelongsToWindow() {
		return c.windowQ.GetPrevNodeBy(nod)
	} else if nod.IsBelongsToProbation() {
		return c.probationQ.GetPrevNodeBy(nod)
	}
	return c.protectedQ.GetPrevNodeBy(nod)
}

func (c *Cache) evictEntry(nod *node.Node) bool {
	fmt.Println("=========执行evictEntry========")
	//尝试根据给定的删除原因驱逐条目。由于当前只会根据容量驱逐，因此不设其他参数
	//分别判断三种情况导致的驱逐，并更具情况设置是否复活，synchronized执行
	//移除node监听器
	//移除hashmap里的元素
	fmt.Println("=========执行evictEntry:",nod.Key())
	c.hashmap.Delete(nod.Key())
	c.weight-=nod.Weight()
	if nod.IsBelongsToWindow() {
		c.windowQ.Remove(nod)
	} else if nod.IsBelongsToProbation(){
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
	fmt.Printf("========执行admit方法：victim：%s：%d次,candidate:%s：%d=======\n",victim.Key(),victimFreq,candidate.Key(),candidateFreq)
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
		queue.MoveToBack(nod) //todo
	}
}
func (c *Cache) reorderProbation(nod *node.Node) { //从probation队列中移动元素
	if !c.probationQ.Contains(nod) {
		return
	} else if nod.Weight()>c.protectedQ.MaxWeight   {
		//若大小超过protected大小，则放入pb的尾部
		c.reorder(c.probationQ, nod)
		return
	}
	//修改pt的权重大小，但现在不考虑权重故不需要
	c.probationQ.Remove(nod)
	c.protectedQ.Push(nod)
	nod.InProtected()
}


