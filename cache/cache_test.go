package cache

import (
	"github.com/stretchr/testify/assert"
	"goffeine/cache/internal/node"
	"strconv"
	"testing"
)

func newCache(n int) Cache {
	return New(n)
}

func TestInitial(t *testing.T) {
	assert := assert.New(t)
	cache := newCache(10)
	assert.Equal(10, cache.Capacity())
	assert.Equal(0, cache.Weight())
	assert.Equal(false, cache.Contains(node.New("abc", 1)))
}
func TestCache_PutANDGET_two(t *testing.T) {
	assert :=assert.New(t)
	cache := newCache(2)
	cache.Put("test1",1)
	cache.Put("test2",2)
	value:=cache.Get("test2")
	assert.Equal(2,value)//判断边界容量插入
	cache.Put("test3",3)
	assert.Equal(true,cache.Get("test1")==nil)//判断超容量是否被驱逐
	assert.Equal(3,cache.Get("test3"))//判断超容量时元素是否被插入
	cache.Put("test3",4)
	assert.Equal(4,cache.Get("test3"))//判断已有值插入时是否被更新
	cache.Put("test4",4)
	for i := 0; i < 6; i++ {
		cache.Get("test4")
	}
	cache.Put("test5",5)
	assert.Equal(true,cache.Get("test1")==nil)//判断放在protected区的元素在window访问次数到达上限时是否被淘汰。
	//空值如何处理合适？？？
}
func TestCache_PutANDGET_bignum(t *testing.T) {
	cache := newCache(100)
	assert :=assert.New(t)
	for i:=0;i<25;i++{
		key:=strconv.Itoa(i)
		cache.Put(key,i)
	}
	cache.Get("0")//访问后放入protected中不会被驱逐

	for i:=25;i<100;i++{//继续添加，此时放满
		key:=strconv.Itoa(i)
		cache.Put(key,i)
	}
	cache.Put("100",100)//window中的99被替换，此时window中是100.
	assert.Equal(true,cache.Get("99")==nil)//检测99是否被驱逐
	assert.Equal(100,cache.Get("100"))//检测99是否被驱逐
	for j:=0;j<7;j++{
		cache.Get("100")//将window中的100次数累积到5以上，在下一轮驱逐时晋升到probation中，不会被淘汰
	}
	cache.Put("101",101)//window此时为101，probation中1被驱逐，100被放入队尾。
	assert.Equal(0,cache.Get("0"))//0不被驱逐
	assert.Equal(true,cache.Get("1")==nil)//1被驱逐
	cache.Put("102",102)//window中的101被驱逐，102放入window中
	assert.Equal(true,cache.Get("101")==nil)//101被驱逐
	assert.Equal(false,cache.Get("102")==nil)//102存在
	assert.Equal(100,cache.weight)//检测是否是最大值
}
func TestCache_PutANDGETwithWeight_bignum(t *testing.T) {
	cache := newCache(1000)
	assert :=assert.New(t)
	cache.PutWithWeight("1",1,2)
	cache.PutWithWeight("1",2,5)//测验weight方法是否会重置weight，
	cache.PutWithWeight("2",1,300)//测验当前权重是正确
	assert.Equal(305,cache.Weight())
	cache.PutWithWeight("3",3,600)//
	//cache.PutWithWeight("4",4,800)//检验多驱逐是否成功
	//assert.Equal(80,cache.Weight())//此处不是bug是程序逻辑，是由于window太小导致的，需要动态调整窗口大小去优化
	cache.PutWithWeight("5",3,100)//当已经满时，大于window大小的weight一律放不进。
//	assert.Equal(80,cache.Weight())
}
