// https://github.com/FastFilter/xorfilter
// Xor 和二进制融合过滤器是布隆过滤器的更快、更简洁的替代方案。
// 此外，与布隆过滤器不同，xor 和二进制融合过滤器可以使用标准技术（gzip、zstd 等）自然压缩。
// 它们的体积也小于布谷鸟过滤器。它们被用于生产系统。
// Xor8 只保存 8-bit 指纹，而 BinaryFuse 可以针对 8-bit / 16-bit / 32-bit 等等泛型 T Unsigned 做更通用的处理.
//
// Xor8：较为直接、简单，每个 key 只计算三次哈希，对应在 [0, BlockLength-1], [BlockLength, 2*BlockLength-1], [2*BlockLength, 3*BlockLength - 1] 3 个位置；
// BinaryFuse：更优化的结构，分段式设计，让构造及查询可借助更好的缓存局部性；对大规模数据时表现更好。

package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/bits"
	"sort"
)

func main() {
	{ // 1.准备你的静态数据：这里举例 20 个随机 uint64 值
		keys := []uint64{1, 42, 123, 999, 1000, 1234, 5555, 9999, 7777, 55555,
			66666, 88888, 101010, 202020, 303030, 404040, 505050, 606060, 707070, 808080}

		// 2. 构建 Xor8 过滤器
		filter, err := Populate(keys)
		if err != nil {
			log.Fatal("Populate failed:", err)
		}

		// 3. 使用过滤器进行查询
		testKeys := []uint64{42, 999, 99999}
		for _, k := range testKeys {
			mightContain := filter.Contains(k)
			fmt.Printf("Key %d => mightContain=%v\n", k, mightContain)
		}
	}

	{
		// 1. 准备你的数据
		keys := []uint64{1, 42, 123, 999, 1000, 1234 /* ... */, 808080}

		// 2. 使用 BinaryFuse8 来构建过滤器
		fuseFilter, err := PopulateBinaryFuse8(keys)
		if err != nil {
			log.Fatal("PopulateBinaryFuse8 failed:", err)
		}

		// 3. 测试查询
		testKeys := []uint64{42, 999, 99999}
		for _, k := range testKeys {
			mightContain := fuseFilter.Contains(k)
			fmt.Printf("Key %d => mightContain=%v\n", k, mightContain)
		}
	}
}

// #region xorfilter_definitions
// Xor8 offers a 0.3% false-positive probability
type Xor8 struct {
	Seed uint64 // 随机种子，用来混入哈希计算（减少哈希冲突带来的构造失败概率）

	// 块的大小(= capacity/3)
	// 由于采用三重哈希（3-wise hashing），把整个指纹数组分成 3 段，每段大小都是 BlockLength
	BlockLength uint32

	Fingerprints []uint8
}

type xorset struct {
	xormask uint64
	count   uint32
}

type hashes struct {
	h  uint64
	h0 uint32
	h1 uint32
	h2 uint32
}

type keyindex struct {
	hash  uint64
	index uint32
}

// #endregion

// #region xorfilter
func murmur64(h uint64) uint64 {
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	h *= 0xc4ceb9fe1a85ec53
	h ^= h >> 33
	return h
}

// returns random number, modifies the seed
func splitmix64(seed *uint64) uint64 {
	*seed = *seed + 0x9E3779B97F4A7C15
	z := *seed
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

func mixsplit(key, seed uint64) uint64 {
	return murmur64(key + seed)
}

func rotl64(n uint64, c int) uint64 {
	return (n << uint(c&63)) | (n >> uint((-c)&63))
}

func reduce(hash, n uint32) uint32 {
	// http://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
	return uint32((uint64(hash) * uint64(n)) >> 32)
}

func fingerprint(hash uint64) uint64 {
	return hash ^ (hash >> 32)
}

// Contains tell you whether the key is likely part of the set
func (filter *Xor8) Contains(key uint64) bool {
	hash := mixsplit(key, filter.Seed)
	f := uint8(fingerprint(hash))
	// !计算三重哈希：h0, h1, h2，映射到三个段中对应的位置
	r0 := uint32(hash)
	r1 := uint32(rotl64(hash, 21))
	r2 := uint32(rotl64(hash, 42))
	h0 := reduce(r0, filter.BlockLength)
	h1 := reduce(r1, filter.BlockLength) + filter.BlockLength
	h2 := reduce(r2, filter.BlockLength) + 2*filter.BlockLength
	// !对比与 key 的“fingerprint”是否相同
	return f == (filter.Fingerprints[h0] ^ filter.Fingerprints[h1] ^ filter.Fingerprints[h2])
}

func (filter *Xor8) geth0h1h2(k uint64) hashes {
	hash := mixsplit(k, filter.Seed)
	answer := hashes{}
	answer.h = hash
	r0 := uint32(hash)
	r1 := uint32(rotl64(hash, 21))
	r2 := uint32(rotl64(hash, 42))

	answer.h0 = reduce(r0, filter.BlockLength)
	answer.h1 = reduce(r1, filter.BlockLength)
	answer.h2 = reduce(r2, filter.BlockLength)
	return answer
}

func (filter *Xor8) geth0(hash uint64) uint32 {
	r0 := uint32(hash)
	return reduce(r0, filter.BlockLength)
}

func (filter *Xor8) geth1(hash uint64) uint32 {
	r1 := uint32(rotl64(hash, 21))
	return reduce(r1, filter.BlockLength)
}

func (filter *Xor8) geth2(hash uint64) uint32 {
	r2 := uint32(rotl64(hash, 42))
	return reduce(r2, filter.BlockLength)
}

// scan for values with a count of one
func scanCount(Qi []keyindex, setsi []xorset) ([]keyindex, int) {
	QiSize := 0

	// len(setsi) = filter.BlockLength
	for i := uint32(0); i < uint32(len(setsi)); i++ {
		if setsi[i].count == 1 {
			Qi[QiSize].index = i
			Qi[QiSize].hash = setsi[i].xormask
			QiSize++
		}
	}

	return Qi, QiSize
}

// fill setsi to xorset{0, 0}
func resetSets(setsi []xorset) []xorset {
	for i := range setsi {
		setsi[i] = xorset{0, 0}
	}
	return setsi
}

// The maximum  number of iterations allowed before the populate function returns an error
var MaxIterations = 1024

// !一次性把所有要存的 key（元素）填充进过滤器，构造成功后就不再支持动态插入或删除（是静态结构）.
// Populate fills the filter with provided keys. For best results,
// !the caller should avoid having too many duplicated keys.
// The function may return an error if the set is empty.
func Populate(keys []uint64) (*Xor8, error) {
	size := len(keys)
	if size == 0 {
		return nil, errors.New("provide a non-empty set")
	}
	capacity := 32 + uint32(math.Ceil(1.23*float64(size)))
	capacity = capacity / 3 * 3 // round it down to a multiple of 3

	filter := &Xor8{}
	var rngcounter uint64 = 1
	filter.Seed = splitmix64(&rngcounter)
	filter.BlockLength = capacity / 3

	// slice capacity defaults to length
	filter.Fingerprints = make([]uint8, capacity)

	stack := make([]keyindex, size)
	Q0 := make([]keyindex, filter.BlockLength)
	Q1 := make([]keyindex, filter.BlockLength)
	Q2 := make([]keyindex, filter.BlockLength)
	sets0 := make([]xorset, filter.BlockLength)
	sets1 := make([]xorset, filter.BlockLength)
	sets2 := make([]xorset, filter.BlockLength)
	iterations := 0

	for {
		iterations += 1
		if iterations > MaxIterations {
			// The probability of this happening is lower than the
			// the cosmic-ray probability (i.e., a cosmic ray corrupts your system).
			return nil, errors.New("too many iterations")
		}

		for i := 0; i < size; i++ {
			key := keys[i]
			hs := filter.geth0h1h2(key)
			sets0[hs.h0].xormask ^= hs.h
			sets0[hs.h0].count++
			sets1[hs.h1].xormask ^= hs.h
			sets1[hs.h1].count++
			sets2[hs.h2].xormask ^= hs.h
			sets2[hs.h2].count++
		}

		// scan for values with a count of one
		Q0, Q0size := scanCount(Q0, sets0)
		Q1, Q1size := scanCount(Q1, sets1)
		Q2, Q2size := scanCount(Q2, sets2)

		stacksize := 0
		for Q0size+Q1size+Q2size > 0 {
			for Q0size > 0 {
				Q0size--
				keyindexvar := Q0[Q0size]
				index := keyindexvar.index
				if sets0[index].count == 0 {
					continue // not actually possible after the initial scan.
				}
				hash := keyindexvar.hash
				h1 := filter.geth1(hash)
				h2 := filter.geth2(hash)
				stack[stacksize] = keyindexvar
				stacksize++
				sets1[h1].xormask ^= hash

				sets1[h1].count--
				if sets1[h1].count == 1 {
					Q1[Q1size].index = h1
					Q1[Q1size].hash = sets1[h1].xormask
					Q1size++
				}
				sets2[h2].xormask ^= hash
				sets2[h2].count--
				if sets2[h2].count == 1 {
					Q2[Q2size].index = h2
					Q2[Q2size].hash = sets2[h2].xormask
					Q2size++
				}
			}
			for Q1size > 0 {
				Q1size--
				keyindexvar := Q1[Q1size]
				index := keyindexvar.index
				if sets1[index].count == 0 {
					continue
				}
				hash := keyindexvar.hash
				h0 := filter.geth0(hash)
				h2 := filter.geth2(hash)
				keyindexvar.index += filter.BlockLength
				stack[stacksize] = keyindexvar
				stacksize++
				sets0[h0].xormask ^= hash
				sets0[h0].count--
				if sets0[h0].count == 1 {
					Q0[Q0size].index = h0
					Q0[Q0size].hash = sets0[h0].xormask
					Q0size++
				}
				sets2[h2].xormask ^= hash
				sets2[h2].count--
				if sets2[h2].count == 1 {
					Q2[Q2size].index = h2
					Q2[Q2size].hash = sets2[h2].xormask
					Q2size++
				}
			}
			for Q2size > 0 {
				Q2size--
				keyindexvar := Q2[Q2size]
				index := keyindexvar.index
				if sets2[index].count == 0 {
					continue
				}
				hash := keyindexvar.hash
				h0 := filter.geth0(hash)
				h1 := filter.geth1(hash)
				keyindexvar.index += 2 * filter.BlockLength

				stack[stacksize] = keyindexvar
				stacksize++
				sets0[h0].xormask ^= hash
				sets0[h0].count--
				if sets0[h0].count == 1 {
					Q0[Q0size].index = h0
					Q0[Q0size].hash = sets0[h0].xormask
					Q0size++
				}
				sets1[h1].xormask ^= hash
				sets1[h1].count--
				if sets1[h1].count == 1 {
					Q1[Q1size].index = h1
					Q1[Q1size].hash = sets1[h1].xormask
					Q1size++
				}

			}
		}

		if stacksize == size {
			// success
			break
		}

		if iterations == 10 {
			keys = pruneDuplicates(keys)
			size = len(keys)
		}

		sets0 = resetSets(sets0)
		sets1 = resetSets(sets1)
		sets2 = resetSets(sets2)

		filter.Seed = splitmix64(&rngcounter)
	}

	stacksize := size
	for stacksize > 0 {
		stacksize--
		ki := stack[stacksize]
		val := uint8(fingerprint(ki.hash))
		if ki.index < filter.BlockLength {
			val ^= filter.Fingerprints[filter.geth1(ki.hash)+filter.BlockLength] ^ filter.Fingerprints[filter.geth2(ki.hash)+2*filter.BlockLength]
		} else if ki.index < 2*filter.BlockLength {
			val ^= filter.Fingerprints[filter.geth0(ki.hash)] ^ filter.Fingerprints[filter.geth2(ki.hash)+2*filter.BlockLength]
		} else {
			val ^= filter.Fingerprints[filter.geth0(ki.hash)] ^ filter.Fingerprints[filter.geth1(ki.hash)+filter.BlockLength]
		}
		filter.Fingerprints[ki.index] = val
	}
	return filter, nil
}

func pruneDuplicates(array []uint64) []uint64 {
	sort.Slice(array, func(i, j int) bool {
		return array[i] < array[j]
	})
	pos := 0
	for i := 1; i < len(array); i++ {
		if array[i] != array[pos] {
			array[pos+1] = array[i]
			pos += 1
		}
	}
	return array[:pos+1]
}

// #endregion

// #region binaryfusefilter
type Unsigned interface {
	~uint8 | ~uint16 | ~uint32
}

// T：可以是 uint8, uint16, or uint32，代表指纹占多少位.
type BinaryFuse[T Unsigned] struct {
	Seed               uint64
	SegmentLength      uint32
	SegmentLengthMask  uint32
	SegmentCount       uint32
	SegmentCountLength uint32

	Fingerprints []T
}

// NewBinaryFuse fills the filter with provided keys. For best results,
// the caller should avoid having too many duplicated keys.
// The function may return an error if the set is empty.
func NewBinaryFuse[T Unsigned](keys []uint64) (*BinaryFuse[T], error) {
	size := uint32(len(keys))
	filter := &BinaryFuse[T]{}
	filter.initializeParameters(size)
	rngcounter := uint64(1)
	filter.Seed = splitmix64(&rngcounter)
	capacity := uint32(len(filter.Fingerprints))

	alone := make([]uint32, capacity)
	// the lowest 2 bits are the h index (0, 1, or 2)
	// so we only have 6 bits for counting;
	// but that's sufficient
	t2count := make([]T, capacity)
	reverseH := make([]T, size)

	t2hash := make([]uint64, capacity)
	reverseOrder := make([]uint64, size+1)
	reverseOrder[size] = 1

	// the array h0, h1, h2, h0, h1, h2
	var h012 [6]uint32
	// this could be used to compute the mod3
	// tabmod3 := [5]uint8{0,1,2,0,1}
	iterations := 0
	for {
		iterations += 1
		if iterations > MaxIterations {
			// The probability of this happening is lower than the
			// the cosmic-ray probability (i.e., a cosmic ray corrupts your system).
			return nil, errors.New("too many iterations")
		}

		blockBits := 1
		for (1 << blockBits) < filter.SegmentCount {
			blockBits += 1
		}
		startPos := make([]uint, 1<<blockBits)
		for i := range startPos {
			// important: we do not want i * size to overflow!!!
			startPos[i] = uint((uint64(i) * uint64(size)) >> blockBits)
		}
		for _, key := range keys {
			hash := mixsplit(key, filter.Seed)
			segment_index := hash >> (64 - blockBits)
			for reverseOrder[startPos[segment_index]] != 0 {
				segment_index++
				segment_index &= (1 << blockBits) - 1
			}
			reverseOrder[startPos[segment_index]] = hash
			startPos[segment_index] += 1
		}
		error := 0
		duplicates := uint32(0)

		for i := uint32(0); i < size; i++ {
			hash := reverseOrder[i]
			index1, index2, index3 := filter.getHashFromHash(hash)
			t2count[index1] += 4
			// t2count[index1] ^= 0 // noop
			t2hash[index1] ^= hash
			t2count[index2] += 4
			t2count[index2] ^= 1
			t2hash[index2] ^= hash
			t2count[index3] += 4
			t2count[index3] ^= 2
			t2hash[index3] ^= hash
			// If we have duplicated hash values, then it is likely that
			// the next comparison is true
			if t2hash[index1]&t2hash[index2]&t2hash[index3] == 0 {
				// next we do the actual test
				if ((t2hash[index1] == 0) && (t2count[index1] == 8)) || ((t2hash[index2] == 0) && (t2count[index2] == 8)) || ((t2hash[index3] == 0) && (t2count[index3] == 8)) {
					duplicates += 1
					t2count[index1] -= 4
					t2hash[index1] ^= hash
					t2count[index2] -= 4
					t2count[index2] ^= 1
					t2hash[index2] ^= hash
					t2count[index3] -= 4
					t2count[index3] ^= 2
					t2hash[index3] ^= hash
				}
			}
			if t2count[index1] < 4 {
				error = 1
			}
			if t2count[index2] < 4 {
				error = 1
			}
			if t2count[index3] < 4 {
				error = 1
			}
		}
		if error == 1 {
			for i := uint32(0); i < size; i++ {
				reverseOrder[i] = 0
			}
			for i := uint32(0); i < capacity; i++ {
				t2count[i] = 0
				t2hash[i] = 0
			}
			filter.Seed = splitmix64(&rngcounter)
			continue
		}

		// End of key addition

		Qsize := 0
		// Add sets with one key to the queue.
		for i := uint32(0); i < capacity; i++ {
			alone[Qsize] = i
			if (t2count[i] >> 2) == 1 {
				Qsize++
			}
		}
		stacksize := uint32(0)
		for Qsize > 0 {
			Qsize--
			index := alone[Qsize]
			if (t2count[index] >> 2) == 1 {
				hash := t2hash[index]
				found := t2count[index] & 3
				reverseH[stacksize] = found
				reverseOrder[stacksize] = hash
				stacksize++

				index1, index2, index3 := filter.getHashFromHash(hash)

				h012[1] = index2
				h012[2] = index3
				h012[3] = index1
				h012[4] = h012[1]

				other_index1 := h012[found+1]
				alone[Qsize] = other_index1
				if (t2count[other_index1] >> 2) == 2 {
					Qsize++
				}
				t2count[other_index1] -= 4
				t2count[other_index1] ^= filter.mod3(found + 1) // could use this instead: tabmod3[found+1]
				t2hash[other_index1] ^= hash

				other_index2 := h012[found+2]
				alone[Qsize] = other_index2
				if (t2count[other_index2] >> 2) == 2 {
					Qsize++
				}
				t2count[other_index2] -= 4
				t2count[other_index2] ^= filter.mod3(found + 2) // could use this instead: tabmod3[found+2]
				t2hash[other_index2] ^= hash
			}
		}

		if stacksize+duplicates == size {
			// Success
			size = stacksize
			break
		} else if duplicates > 0 {
			// Duplicates were found, but we did not
			// manage to remove them all. We may simply sort the key to
			// solve the issue. This will run in time O(n log n) and it
			// mutates the input.
			keys = pruneDuplicates(keys)
		}
		for i := uint32(0); i < size; i++ {
			reverseOrder[i] = 0
		}
		for i := uint32(0); i < capacity; i++ {
			t2count[i] = 0
			t2hash[i] = 0
		}
		filter.Seed = splitmix64(&rngcounter)
	}
	if size == 0 {
		return filter, nil
	}

	for i := int(size - 1); i >= 0; i-- {
		// the hash of the key we insert next
		hash := reverseOrder[i]
		xor2 := T(fingerprint(hash))
		index1, index2, index3 := filter.getHashFromHash(hash)
		found := reverseH[i]
		h012[0] = index1
		h012[1] = index2
		h012[2] = index3
		h012[3] = h012[0]
		h012[4] = h012[1]
		filter.Fingerprints[h012[found]] = xor2 ^ filter.Fingerprints[h012[found+1]] ^ filter.Fingerprints[h012[found+2]]
	}

	return filter, nil
}

func (filter *BinaryFuse[T]) initializeParameters(size uint32) {
	arity := uint32(3)
	filter.SegmentLength = calculateSegmentLength(arity, size)
	if filter.SegmentLength > 262144 {
		filter.SegmentLength = 262144
	}
	filter.SegmentLengthMask = filter.SegmentLength - 1
	sizeFactor := calculateSizeFactor(arity, size)
	capacity := uint32(0)
	if size > 1 {
		capacity = uint32(math.Round(float64(size) * sizeFactor))
	}
	initSegmentCount := (capacity+filter.SegmentLength-1)/filter.SegmentLength - (arity - 1)
	arrayLength := (initSegmentCount + arity - 1) * filter.SegmentLength
	filter.SegmentCount = (arrayLength + filter.SegmentLength - 1) / filter.SegmentLength
	if filter.SegmentCount <= arity-1 {
		filter.SegmentCount = 1
	} else {
		filter.SegmentCount = filter.SegmentCount - (arity - 1)
	}
	arrayLength = (filter.SegmentCount + arity - 1) * filter.SegmentLength
	filter.SegmentCountLength = filter.SegmentCount * filter.SegmentLength
	filter.Fingerprints = make([]T, arrayLength)
}

func (filter *BinaryFuse[T]) mod3(x T) T {
	if x > 2 {
		x -= 3
	}

	return x
}

func (filter *BinaryFuse[T]) getHashFromHash(hash uint64) (uint32, uint32, uint32) {
	hi, _ := bits.Mul64(hash, uint64(filter.SegmentCountLength))
	h0 := uint32(hi)
	h1 := h0 + filter.SegmentLength
	h2 := h1 + filter.SegmentLength
	h1 ^= uint32(hash>>18) & filter.SegmentLengthMask
	h2 ^= uint32(hash) & filter.SegmentLengthMask
	return h0, h1, h2
}

// Contains returns `true` if key is part of the set with a false positive probability.
func (filter *BinaryFuse[T]) Contains(key uint64) bool {
	hash := mixsplit(key, filter.Seed)
	f := T(fingerprint(hash))
	h0, h1, h2 := filter.getHashFromHash(hash)
	f ^= filter.Fingerprints[h0] ^ filter.Fingerprints[h1] ^ filter.Fingerprints[h2]
	return f == 0
}

func calculateSegmentLength(arity uint32, size uint32) uint32 {
	// These parameters are very sensitive. Replacing 'floor' by 'round' can
	// substantially affect the construction time.
	if size == 0 {
		return 4
	}
	if arity == 3 {
		return uint32(1) << int(math.Floor(math.Log(float64(size))/math.Log(3.33)+2.25))
	} else if arity == 4 {
		return uint32(1) << int(math.Floor(math.Log(float64(size))/math.Log(2.91)-0.5))
	} else {
		return 65536
	}
}

func calculateSizeFactor(arity uint32, size uint32) float64 {
	if arity == 3 {
		return math.Max(1.125, 0.875+0.25*math.Log(1000000)/math.Log(float64(size)))
	} else if arity == 4 {
		return math.Max(1.075, 0.77+0.305*math.Log(600000)/math.Log(float64(size)))
	} else {
		return 2.0
	}
}

// #endregion

// #region binaryfusefilter8
type BinaryFuse8 BinaryFuse[uint8]

// !构造一个 8-bit 指纹的 BinaryFuse filter.
// PopulateBinaryFuse8 fills the filter with provided keys. For best results,
// the caller should avoid having too many duplicated keys.
// The function may return an error if the set is empty.
func PopulateBinaryFuse8(keys []uint64) (*BinaryFuse8, error) {
	filter, err := NewBinaryFuse[uint8](keys)
	if err != nil {
		return nil, err
	}

	return (*BinaryFuse8)(filter), nil
}

// Contains returns `true` if key is part of the set with a false positive probability of <0.4%.
func (filter *BinaryFuse8) Contains(key uint64) bool {
	return (*BinaryFuse[uint8])(filter).Contains(key)
}

// #endregion
