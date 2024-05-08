package goffeine

import (
	"goffeine/internal/utils"
	"math"
)

// FrequencySketch migrate based on
// https://github.com/ben-manes/caffeine/blob/master/caffeine/src/main/java/com/github/benmanes/caffeine/cache/FrequencySketch.java
// This class maintains a 4-bit CountMinSketch [1] with periodic aging to provide the popularity
// history for the TinyLfu admission policy [2]. The time and space efficiency of the sketch
// allows it to cheaply estimate the frequency of an entry in a stream of cache access events.
//
// The counter matrix is represented as a single-dimensional array holding 16 counters per slot. A
// fixed depth of four balances the accuracy and cost, resulting in a width of four times the
// length of the array. To retain an accurate estimation, the array's length equals the maximum
// number of entries in the cache, increased to the closest power-of-two to exploit more efficient
// bit masking. This configuration results in a confidence of 93.75% and an error bound of
// e / width.
//
// To improve hardware efficiency, an item's counters are constrained to a 64-byte block, which is
// the Size of an L1 cache line. This differs from the theoretical ideal where counters are
// uniformly distributed to minimize collisions. In that configuration, the memory accesses are
// not predictable and lack spatial locality, which may cause the pipeline to need to wait for
// four memory loads. Instead, the items are uniformly distributed to blocks, and each counter is
// uniformly selected from a distinct 16-byte segment. While the runtime memory layout may result
// in the blocks not being cache-aligned, the L2 spatial prefetcher tries to load aligned pairs of
// cache lines, so the typical cost is only one memory access.
//
// The frequency of all entries is aged periodically using a sampling window based on the maximum
// number of entries in the cache. This is referred to as the reset operation by TinyLfu and keeps
// the sketch fresh by dividing all counters by two and subtracting based on the number of odd
// counters found. The O(n) cost of aging is amortized, ideal for hardware prefetching, and uses
// inexpensive bit manipulations per array location.
//
// [1] An Improved Data Stream Summary: The Count-Min Sketch and its Applications
// http://dimacs.rutgers.edu/~graham/pubs/papers/cm-full.pdf
// [2] TinyLFU: A Highly Efficient Cache Admission Policy
// https://dl.acm.org/citation.cfm?id=3149371
// [3] Hash Function Prospector: Three round functions
// https://github.com/skeeto/hash-prospector#three-round-functions
const (
	ResetMask int64 = 0x7777777777777777
	OneMask   int64 = 0x1111111111111111
	LongF     int64 = 0xf
	Long1     int64 = 1
)

func NewSketch(maximumSize int) *FrequencySketch {
	fs := FrequencySketch{}
	fs.EnsureCapacity(maximumSize)
	return &fs
}

type FrequencySketch struct {
	SampleSize int
	BlockMask  int
	Table      []int64
	Size       int
}

// EnsureCapacity
// Initializes and increases the capacity of this FrequencySketch instance, if necessary,
// to ensure that it can accurately estimate the popularity of elements given the maximum Size of
// the cache. This operation forgets all previous counts when resizing.
func (f *FrequencySketch) EnsureCapacity(maximumSize int) {
	maximum := int(utils.Min(maximumSize, int(uint(math.MaxInt32)>>1)))
	if len(f.Table) > maximum {
		return
	}

	newSize := int(utils.Max(utils.CeilingPowerOfTwo32(maximum), 8))
	f.Table = make([]int64, newSize)

	if maximumSize == 0 {
		f.SampleSize = 10
	} else {
		f.SampleSize = 10 * maximum
	}

	f.BlockMask = int(uint(newSize)>>3) - 1 // 需要多少个block
	if int32(f.SampleSize) <= 0 {
		f.SampleSize = math.MaxInt32
	}
	f.Size = 0
}

// Returns the estimated number of occurrences of an element, up to the maximum (15).
// @param e the element to count occurrences of
// @return the estimated number of occurrences of the element; possibly zero but never negative
func (f *FrequencySketch) Frequency(e string) int {
	count := make([]int, 4)
	blockHash := spread(hashCode(e))
	counterHash := rehash(blockHash)
	block := (blockHash & f.BlockMask) << 3

	for i := 0; i < 4; i++ {
		h := int(uint(counterHash) >> (i << 3))
		index := int(uint(h)>>1) & 15
		offset := h & 1
		count[i] = int(int64(uint64(f.Table[offset+block+(i<<1)])>>(index<<2)) & LongF)
	}
	return int(utils.Min(utils.Min(count[0], count[1]), utils.Min(count[2], count[3])))
}

// Increment
// Increments the popularity of the element if it does not exceed the maximum (15). The popularity
// of all elements will be periodically down sampled when the observed events exceed a threshold.
// This process provides a frequency aging to allow expired long term entries to fade away.
func (f *FrequencySketch) Increment(e string) {
	index := make([]int, 8)
	blockHash := spread(hashCode(e))
	counterHash := rehash(blockHash)
	block := (blockHash & f.BlockMask) << 3
	for i := 0; i < 4; i++ {
		h := int(uint(counterHash) >> (i << 3))
		index[i] = int(uint(h)>>1) & 15
		offset := h & 1
		index[i+4] = block + offset + (i << 1)
	}

	added := f.incrementAt(index[4], index[0])
	added = f.incrementAt(index[5], index[1]) || added
	added = f.incrementAt(index[6], index[2]) || added
	added = f.incrementAt(index[7], index[3]) || added

	f.Size += 1
	if added && f.Size == f.SampleSize {
		f.Reset()
	}
}

// Increments the specified counter by 1 if it is not already at the maximum value (15).
//
// @param i the Table index (16 counters)
// @param j the counter to increment
// @return if incremented
func (f *FrequencySketch) incrementAt(i int, j int) bool {
	offset := j << 2
	mask := LongF << offset
	if (f.Table[i] & mask) != mask {
		f.Table[i] += (Long1 << offset)
		return true
	}
	return false
}

// Reset Reduces every counter by half of its original value.
func (f *FrequencySketch) Reset() {
	count := 0
	for i := 0; i < len(f.Table); i++ {
		count += BitCount64(f.Table[i] & OneMask)
		f.Table[i] = int64(uint64(f.Table[i])>>1) & ResetMask
	}
	f.Size = f.Size - int(uint(count)>>2)
	f.Size = int(uint(f.Size) >> 1)
}

// hashCode is a simple hash function that returns an int.
func hashCode(s string) int {
	h := 0
	for _, v := range []byte(s) {
		h = 31*h + int(v&0xff)
	}
	return h
}

// spread Applies a supplemental hash function to defend against a poor quality hash.
func spread(x int) int {
	x ^= int(uint(x) >> 17)
	x *= 0xed5ad4bb
	x ^= int(uint(x) >> 11)
	x *= 0xac4c1b51
	x ^= int(uint(x) >> 15)
	return x
}

// rehash Applies another round of hashing for additional randomization.
func rehash(x int) int {
	x *= 0x31848bab
	x ^= int(uint(x) >> 14)
	return x
}

const (
	Long5Mask  int64 = 0x5555555555555555
	Long3Mask  int64 = 0x3333333333333333
	Long0FMask int64 = 0x0f0f0f0f0f0f0f0f
)

func BitCount64(i int64) int {
	// HD, Figure 5-2
	i = i - (int64(uint64(i)>>1) & Long5Mask)
	i = (i & Long3Mask) + (int64(uint64(i)>>2) & Long3Mask)
	i = (i + int64(uint64(i)>>4)) & Long0FMask
	i = i + int64(uint64(i)>>8)
	i = i + int64(uint64(i)>>16)
	i = i + int64(uint64(i)>>32)
	return int(i) & 0x7f
}
