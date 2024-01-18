package sketch

import (
	"goffeine/internal/node"
	"goffeine/internal/utils"
	"math"
	"math/bits"
)

var SEEDS = [4]uint64{0xc3a5c85c97cb3127, 0xb492b66fbe98f273, 0x9ae16a3b2f90404f, 0xcbf29ce484222325} // 来自 FNV-1a, CityHash, and Murmur3 的种子数
const (
	LONG_15    uint64 = 0xf
	LONG_1     uint64 = 1
	RESET_MASK uint64 = 0x7777777777777777
	ONE_MASK   uint64 = 0x1111111111111111
)

// 完整名字：FrequencySketch
// 用一个uint64表示所有的次数，下面详细解释一下：
// 1、每个元素的Fre最多不超过15，那么可以用 4个bit来表示。从0000到1111。
// 2、一个int64有64个bit，所以一个int64可以表示(64/4)=16个Fre。
type FrequencySketch struct {
	table     []uint64 // 数据表格
	length    int      // table的长度
	counter   int      // 计数器，每次increament就需要+1
	threshold int      // 临界值，当counter到临界值到了后，就要reset了
}

// New一个FrequencySketch
//
// @param n 表示你要存放多少个次数
func New(n int) *FrequencySketch {
	if n <= 0 {
		n = 1
	}
	// 用4个bit来表示一个Fre。所以，理论上分配的容量最好是4的整数倍。
	// 但是4的整数倍不一定是64的整数倍，所以，分配容量应该是Int64的整数倍最合适。
	//length := n * 4 / 64
	length := n

	if length < 4 {
		length = 4
	}

	// 分配容量优化都是按照2的n次方个Byte来分配
	// 所以如果按照2的n次方个int64来分配
	length = utils.CeilingPowerOfTwo32(length)
	f := FrequencySketch{
		table:     make([]uint64, length),
		length:    length,
		threshold: 10 * length,
		counter:   0,
	}
	return &f
}

func (s *FrequencySketch) Frequency(pNode *node.Node) int {
	var x []byte = pNode.KeyHash
	frequency := math.MaxInt32
	hashCode := uint64(s.spread(x))
	start := int((hashCode & 3) << 2)
	for i := 0; i < 4; i++ {
		idx := s.indexOf(hashCode, i)
		offset := (start + i) << 2
		count := (int)((s.table[idx] >> offset) & LONG_15)
		frequency = int(math.Min(float64(frequency), float64(count)))
	}
	return frequency
}

func (s *FrequencySketch) Increment(pNode *node.Node) {
	var x []byte = pNode.KeyHash
	hashCode := uint64(s.spread(x))
	start := int((hashCode & 3) << 2)
	added := 0
	for i := 0; i < 4; i++ {
		idx := s.indexOf(hashCode, i)
		added |= s.incrementAt(idx, start+i)
	}
	if added == 1 {
		s.counter += 1
		if s.counter == s.threshold {
			s.reset()
		}
	}
}

// 给table下标为idx的，第pos个fre加1
//
// @param i  table的下标
// @param j  第几个fre
func (s *FrequencySketch) incrementAt(i, j int) int {
	// j永远是下面的值：
	// 0, 1, 2, 3
	// 4, 5, 6, 7
	// 8, 9, 10, 11
	// 12, 13, 14, 15

	// 所以下面的offset永远是：
	// 0, 4, 8, 12
	// 16, 20, 24, 28,
	// 32, 36, 40, 44,
	// 48, 52, 56, 60
	var offset uint64 = uint64(j << 2)

	mask := LONG_15 << offset

	if (s.table[i] & mask) != mask {
		s.table[i] += (LONG_1 << offset)
		return 1
	}
	return 0
}

// 得到指定深度的计数器在table的索引
// 具体什么算法也没有完全高清楚，是移植自Caffeine的源码
//
// @param item the element's hash
// @param i the counter depth
// @return the table index
func (s *FrequencySketch) indexOf(item uint64, i int) int {
	seed := SEEDS[i]
	hash := (item + seed) * seed
	hash += hash >> 32
	return int(hash & uint64(s.length-1))
}

// 散列出一个更加好的hash数值
func (s *FrequencySketch) spread(key []byte) uint32 {
	// h算法移植于 Java 的 StringLatin1.hashCode()
	h := 0
	for _, v := range key {
		h = 31*h + int(v&0xff)
	}
	// 怕hashCode不够散列，再来一次
	// 算法移植于 Caffeine
	x := uint32(h)
	x = ((x >> 16) ^ x) * 0x45d9f3b // x = ((x >>> 16) ^ x) * 0x45d9f3b;
	x = ((x >> 16) ^ x) * 0x45d9f3b // x = ((x >>> 16) ^ x) * 0x45d9f3b;
	return (x >> 16) ^ x            // return (x >>> 16) ^ x;
}

// 重置
func (s *FrequencySketch) reset() {
	var count uint = 0
	for i := 0; i < s.length; i++ {
		count += uint(bits.OnesCount64(s.table[i] & ONE_MASK))
		s.table[i] = (s.table[i] >> 1) & RESET_MASK
	}
	s.counter = int((uint(s.length) >> 1) - (count >> 2))
}
