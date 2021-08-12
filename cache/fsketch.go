package cache

import (
	"goffeine/utils"
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
// 1、每个元素的Fre最多不超过15，那么可以用 4个bit来表示。从0000到1111。一个int64有64个bit，所以一个int64可以表示(64/4)=16个Fre。
// 2、但为了防止hash不够均匀，一个int64只用来表示4个Fre，所以，实际上：64bit/4 = 16bit，即：用16bit表示一个Fre
type FSketch struct {
	table     []uint64 // 数据表格
	length    int      // table的长度
	counter   int      // 计数器，每次increament就需要+1
	threshold int      // 临界值，当counter到临界值到了后，就要reset了
}

// New一个FrequencySketch
//
// @param n 表示你要存放多少个次数
func NewFSketch(n int) FSketch {
	if n <= 0 {
		n = 1
	}
	// 用16个bit来表示一个Fre。所以，理论上分配的容量最好是16的整数倍。
	// 又因为一个Int64是4个Fre，但是16的整数倍不一定是64的整数倍
	// 所以，分配容量应该是Int64的整数倍最合适。
	capacity := (n * 16 / 64) + 1
	capacity = utils.CeilingPowerOfTwo32(capacity)

	f := FSketch{
		table:     make([]uint64, capacity),
		length:    capacity,
		threshold: 10 * capacity,
		counter:   0,
	}
	return f
}

func (s *FSketch) Frequency(x []byte) int {
	frequency := 0
	hashCode := uint64(s.spread(x))
	start := int((hashCode & 3) << 2)
	for i := 0; i < 4; i++ {
		idx := s.indexOf(hashCode, SEEDS[i])
		count := (int)((s.table[idx]>> ((start + i) << 2)) & LONG_15)
		frequency = int(math.Min(float64(frequency), float64(count)))
	}
	return frequency
}

func (s *FSketch) Increment(x []byte) {
	hashCode := uint64(s.spread(x))
	start := int((hashCode & 3) << 2)
	added := true
	for i := 0; i < 4; i++ {
		idx := s.indexOf(hashCode, SEEDS[i])
		added = added || s.incrementAt(idx, start+i)
	}
	if added {
		s.counter += 1
		if s.counter == s.threshold {
			s.reset()
		}
	}
}

// 给table下标为idx的，第pos个fre加1
// 具体什么算法也没有完全高清楚，是移植自Caffeine的源码
//
// @param i  table的下标
// @param j  第几个fre
func (s *FSketch) incrementAt(i, j int) bool {
	// j永远等于0, 4, 8, 12
	// 所以下面的offset永远等于 0, 16, 32, 48
	var offset uint64 = uint64(j << 2)

	mask := LONG_15 << offset

	if (s.table[i] & mask) != mask {
		s.table[i] += (LONG_1 << offset)
		return true
	}
	return false
}

// 得到指定深度的计数器在table的索引
// 具体什么算法也没有完全高清楚，是移植自Caffeine的源码
//
// @param item the element's hash
// @param i the counter depth
// @return the table index
func (s *FSketch) indexOf(item uint64, seed uint64) int {
	hash := (item + seed) * seed
	hash += hash >> 32
	return int(hash & uint64(s.length - 1))
}

// 散列出一个更加好的hash数值
func (s *FSketch) spread(key []byte) uint32 {
	h := 0
	for _, v := range key {
		h = 31 * h + int(v & 0xff)
	}
	x := uint32(h)
	x = ((x >> 16) ^ x) * 0x45d9f3b // x = ((x >>> 16) ^ x) * 0x45d9f3b;
	x = ((x >> 16) ^ x) * 0x45d9f3b // x = ((x >>> 16) ^ x) * 0x45d9f3b;
	return (x >> 16) ^ x            // return (x >>> 16) ^ x;
}

// 重置
func (s *FSketch) reset() {
	var count uint = 0
	for i := 0; i < s.length; i++ {
		count += uint(bits.OnesCount64(s.table[i] & ONE_MASK))
		s.table[i] = (s.table[i] >> 1) & RESET_MASK
	}
	s.counter = int((uint(s.length) >> 1) - (count >> 2))
}

