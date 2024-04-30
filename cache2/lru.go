package cache2

type LRU struct{}

func (lru *LRU) add(k, v string) {
	//if Map中存有这条Key {
	//	替换Map中的Value值
	//	将链表中的对应节点移到最前面
	//} else {
	//  if 已经达到缓存容量上限 {
	//  	获取链表尾部节点的Key，并从Map中删除
	//      移除链表尾部的Node
	//  }
	//	创建要插入的新节点
	//	将新节点插入到链表头部
	//	放入Map中
	//}
}

func (lru *LRU) get(k string) string {
	//if Map中存有这条Key {
	//	返回查询到的Value
	//	将对应节点移动到链表头部
	//} else {
	//	返回 空
	//}
	return ""
}
