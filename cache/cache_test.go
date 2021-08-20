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
	assert.Equal(0, cache.Len())
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
	for i:=0;i<150;i++{
		key:=strconv.Itoa(i)
		cache.Put(key,i)
	}
	for j:=0;j<50;j++{
		cache.Get("149")
	}
	for i:=150;i<200;i++{
		key:=strconv.Itoa(i)
		cache.Put(key,i)
	}
	assert.Equal(100,cache.weight)
	assert.Equal(149,cache.Get("149"))
	assert.Equal(true,cache.Get("79")==nil)
}
