// https://github.com/alecthomas/mph
// StaticHashMap

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"sort"
	"time"
)

func main() {
	var testData = map[string][]byte{
		"apple":  []byte("fruit"),
		"banana": []byte("fruit"),
		"cherry": []byte("fruit"),
		"egg":    []byte("food"),
		"milk":   []byte("drink"),
	}

	buildCHDFromMap := func(data map[string][]byte) (*CHD, error) {
		builder := NewCHDBuilder()

		// 可选：设置一个固定随机种子, 以确保每次构建同样的布局
		// builder.Seed(12345)

		// 添加所有 key -> value
		for k, v := range data {
			builder.Add([]byte(k), v)
		}

		// 调用 Build() 生成 CHD
		chd, err := builder.Build()
		if err != nil {
			return nil, err
		}
		return chd, nil
	}

	// 1) 构建 Minimal Perfect Hash
	chd, err := buildCHDFromMap(testData)
	if err != nil {
		log.Fatal(err)
	}

	// 2) 测试查询
	fmt.Println("Check 'banana':", string(chd.Get([]byte("banana"))))
	fmt.Println("Check 'milk':", string(chd.Get([]byte("milk"))))
	fmt.Println("Check 'notfound':", chd.Get([]byte("notfound"))) // nil

	// 3) 遍历 (key, value)
	it := chd.Iterate()
	for it != nil {
		k, v := it.Get()
		fmt.Printf("key=%s, value=%s\n", k, v)
		it = it.Next()
	}

	serializeCHD := func(chd *CHD) ([]byte, error) {
		var buf bytes.Buffer
		err := chd.Write(&buf)
		return buf.Bytes(), err
	}
	// 4) 序列化到内存 (或文件)
	//    这里用内存 buffer 举例
	//    你可以改为 os.Create("mychd.bin") 写入文件
	buf, err := serializeCHD(chd)
	if err != nil {
		log.Fatal(err)
	}

	// 5) 从序列化数据恢复
	loaded, err := Mmap(buf)
	if err != nil {
		log.Fatal(err)
	}

	// 再次测试
	fmt.Println("Loaded check 'cherry':", string(loaded.Get([]byte("cherry"))))
}

// #region chd_builder

type chdHasher struct {
	r       []uint64
	size    uint64
	buckets uint64
	rand    *rand.Rand
}

type bucket struct {
	index  uint64
	keys   [][]byte
	values [][]byte
}

func (b *bucket) String() string {
	a := "bucket{"
	for _, k := range b.keys {
		a += string(k) + ", "
	}
	return a + "}"
}

// Intermediate data structure storing buckets + outer hash index.
type bucketVector []bucket

func (b bucketVector) Len() int           { return len(b) }
func (b bucketVector) Less(i, j int) bool { return len(b[i].keys) > len(b[j].keys) }
func (b bucketVector) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

// Build a new CDH MPH.
type CHDBuilder struct {
	keys   [][]byte
	values [][]byte
	seed   int64
	seeded bool
}

// Create a new CHD hash table builder.
func NewCHDBuilder() *CHDBuilder {
	return &CHDBuilder{}
}

// Seed the RNG. This can be used to reproducible building.
func (b *CHDBuilder) Seed(seed int64) {
	b.seed = seed
	b.seeded = true
}

// Add a key and value to the hash table.
func (b *CHDBuilder) Add(key []byte, value []byte) {
	b.keys = append(b.keys, key)
	b.values = append(b.values, value)
}

// Try to find a hash function that does not cause collisions with table, when
// applied to the keys in the bucket.
func tryHash(hasher *chdHasher, seen map[uint64]bool, keys [][]byte, values [][]byte, indices []uint16, bucket *bucket, ri uint16, r uint64) bool {
	// Track duplicates within this bucket.
	duplicate := make(map[uint64]bool)
	// Make hashes for each entry in the bucket.
	hashes := make([]uint64, len(bucket.keys))
	for i, k := range bucket.keys {
		h := hasher.Table(r, k)
		hashes[i] = h
		if seen[h] {
			return false
		}
		if duplicate[h] {
			return false
		}
		duplicate[h] = true
	}

	// Update seen hashes
	for _, h := range hashes {
		seen[h] = true
	}

	// Add the hash index.
	indices[bucket.index] = ri

	// Update the the hash table.
	for i, h := range hashes {
		keys[h] = bucket.keys[i]
		values[h] = bucket.values[i]
	}
	return true
}

func (b *CHDBuilder) Build() (*CHD, error) {
	n := uint64(len(b.keys))
	m := n / 2
	if m == 0 {
		m = 1
	}

	keys := make([][]byte, n)
	values := make([][]byte, n)
	hasher := newCHDHasher(n, m, b.seed, b.seeded)
	buckets := make(bucketVector, m)
	indices := make([]uint16, m)
	// An extra check to make sure we don't use an invalid index
	for i := range indices {
		indices[i] = ^uint16(0)
	}
	// Have we seen a hash before?
	seen := make(map[uint64]bool)
	// Used to ensure there are no duplicate keys.
	duplicates := make(map[string]bool)

	for i := range b.keys {
		key := b.keys[i]
		value := b.values[i]
		k := string(key)
		if duplicates[k] {
			return nil, errors.New("duplicate key " + k)
		}
		duplicates[k] = true
		oh := hasher.HashIndexFromKey(key)

		buckets[oh].index = oh
		buckets[oh].keys = append(buckets[oh].keys, key)
		buckets[oh].values = append(buckets[oh].values, value)
	}

	// Order buckets by size (retaining the hash index)
	collisions := 0
	sort.Sort(buckets)
nextBucket:
	for i, bucket := range buckets {
		if len(bucket.keys) == 0 {
			continue
		}

		// Check existing hash functions.
		for ri, r := range hasher.r {
			if tryHash(hasher, seen, keys, values, indices, &bucket, uint16(ri), r) {
				continue nextBucket
			}
		}

		// Keep trying new functions until we get one that does not collide.
		// The number of retries here is very high to allow a very high
		// probability of not getting collisions.
		for i := 0; i < 10000000; i++ {
			if i > collisions {
				collisions = i
			}
			ri, r := hasher.Generate()
			if tryHash(hasher, seen, keys, values, indices, &bucket, ri, r) {
				hasher.Add(r)
				continue nextBucket
			}
		}

		// Failed to find a hash function with no collisions.
		return nil, fmt.Errorf(
			"failed to find a collision-free hash function after ~10000000 attempts, for bucket %d/%d with %d entries: %s",
			i, len(buckets), len(bucket.keys), &bucket)
	}

	// println("max bucket collisions:", collisions)
	// println("keys:", len(table))
	// println("hash functions:", len(hasher.r))

	keylist := make([]dataSlice, len(b.keys))
	valuelist := make([]dataSlice, len(b.values))
	var buf bytes.Buffer
	for i, k := range keys {
		keylist[i].start = uint64(buf.Len())
		buf.Write(k)
		keylist[i].end = uint64(buf.Len())
		valuelist[i].start = uint64(buf.Len())
		buf.Write(values[i])
		valuelist[i].end = uint64(buf.Len())
	}

	return &CHD{
		r:       hasher.r,
		indices: indices,
		mmap:    buf.Bytes(),
		keys:    keylist,
		values:  valuelist,
	}, nil
}

func newCHDHasher(size, buckets uint64, seed int64, seeded bool) *chdHasher {
	if !seeded {
		seed = time.Now().UnixNano()
	}
	rs := rand.NewSource(seed)
	c := &chdHasher{size: size, buckets: buckets, rand: rand.New(rs)}
	c.Add(c.rand.Uint64())
	return c
}

// Hash index from key.
func (h *chdHasher) HashIndexFromKey(b []byte) uint64 {
	return (hasher(b) ^ h.r[0]) % h.buckets
}

// Table hash from random value and key. Generate() returns these random values.
func (h *chdHasher) Table(r uint64, b []byte) uint64 {
	return (hasher(b) ^ h.r[0] ^ r) % h.size
}

func (c *chdHasher) Generate() (uint16, uint64) {
	return c.Len(), c.rand.Uint64()
}

// Add a random value generated by Generate().
func (c *chdHasher) Add(r uint64) {
	c.r = append(c.r, r)
}

func (c *chdHasher) Len() uint16 {
	return uint16(len(c.r))
}

func (h *chdHasher) String() string {
	return fmt.Sprintf("chdHasher{size: %d, buckets: %d, r: %v}", h.size, h.buckets, h.r)
}

// #endregion

// #region chd

// CHD hash table lookup.
type CHD struct {
	// Random hash function table.
	r []uint64
	// Array of indices into hash function table r. We assume there aren't
	// more than 2^16 hash functions O_o
	indices []uint16
	// Final table of values.
	mmap   []byte
	keys   []dataSlice
	values []dataSlice
}

type dataSlice struct {
	start uint64
	end   uint64
}

func hasher(data []byte) uint64 {
	var hash uint64 = 14695981039346656037
	for _, c := range data {
		hash ^= uint64(c)
		hash *= 1099511628211
	}
	return hash
}

// Read a serialized CHD.
func Read(r io.Reader) (*CHD, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Mmap(b)
}

// Mmap creates a new CHD aliasing the CHD structure over an existing byte region (typically mmapped).
func Mmap(b []byte) (*CHD, error) {
	c := &CHD{mmap: b}

	bi := &sliceReader{b: b}

	// Read vector of hash functions.
	rl := bi.ReadInt()
	c.r = bi.ReadUint64Array(rl)

	// Read hash function indices.
	il := bi.ReadInt()
	c.indices = bi.ReadUint16Array(il)

	el := bi.ReadInt()

	c.keys = make([]dataSlice, el)
	c.values = make([]dataSlice, el)

	for i := uint64(0); i < el; i++ {
		kl := bi.ReadInt()
		vl := bi.ReadInt()
		c.keys[i].start = bi.pos
		bi.pos += kl
		c.keys[i].end = bi.pos
		c.values[i].start = bi.pos
		bi.pos += vl
		c.values[i].end = bi.pos
	}

	return c, nil
}

func (c *CHD) slice(s dataSlice) []byte {
	return c.mmap[s.start:s.end]
}

// Get an entry from the hash table.
func (c *CHD) Get(key []byte) []byte {
	r0 := c.r[0]
	h := hasher(key) ^ r0
	i := h % uint64(len(c.indices))
	ri := c.indices[i]
	// This can occur if there were unassigned slots in the hash table.
	if ri >= uint16(len(c.r)) {
		return nil
	}
	r := c.r[ri]
	ti := (h ^ r) % uint64(len(c.keys))
	// fmt.Printf("r[0]=%d, h=%d, i=%d, ri=%d, r=%d, ti=%d\n", c.r[0], h, i, ri, r, ti)
	k := c.keys[ti]
	if bytes.Compare(c.slice(k), key) != 0 {
		return nil
	}
	v := c.values[ti]
	return c.slice(v)
}

func (c *CHD) Len() int {
	return len(c.keys)
}

// Iterate over entries in the hash table.
func (c *CHD) Iterate() *Iterator {
	if len(c.keys) == 0 {
		return nil
	}
	return &Iterator{c: c}
}

// Serialize the CHD. The serialized form is conducive to mmapped access. See
// the Mmap function for details.
func (c *CHD) Write(w io.Writer) error {
	write := func(nd ...interface{}) error {
		for _, d := range nd {
			if err := binary.Write(w, binary.LittleEndian, d); err != nil {
				return err
			}
		}
		return nil
	}

	data := []interface{}{
		uint32(len(c.r)), c.r,
		uint32(len(c.indices)), c.indices,
		uint32(len(c.keys)),
	}

	if err := write(data...); err != nil {
		return err
	}

	for i := range c.keys {
		k, v := c.keys[i], c.values[i]
		if err := write(uint32(k.end-k.start), uint32(v.end-v.start)); err != nil {
			return err
		}
		if _, err := w.Write(c.slice(k)); err != nil {
			return err
		}
		if _, err := w.Write(c.slice(v)); err != nil {
			return err
		}
	}
	return nil
}

type Iterator struct {
	i int
	c *CHD
}

func (c *Iterator) Get() (key []byte, value []byte) {
	return c.c.slice(c.c.keys[c.i]), c.c.slice(c.c.values[c.i])
}

func (c *Iterator) Next() *Iterator {
	c.i++
	if c.i >= len(c.c.keys) {
		return nil
	}
	return c
}

// #endregion

// #region sliceReader

// Read values and typed vectors from a byte slice without copying where possible.
type sliceReader struct {
	b   []byte
	pos uint64
}

func (b *sliceReader) Read(size uint64) []byte {
	start := b.pos
	b.pos += size
	return b.b[start:b.pos]
}

func (b *sliceReader) ReadUint64Array(n uint64) []uint64 {
	buf := b.Read(n * 8)
	out := make([]uint64, n)
	for i := 0; i < len(buf); i += 8 {
		out[i>>3] = binary.LittleEndian.Uint64(buf[i : i+8])
	}
	return out
}

func (b *sliceReader) ReadUint16Array(n uint64) []uint16 {
	buf := b.Read(n * 2)
	out := make([]uint16, n)
	for i := 0; i < len(buf); i += 2 {
		out[i>>1] = binary.LittleEndian.Uint16(buf[i : i+2])
	}
	return out
}

func (b *sliceReader) ReadInt() uint64 {
	return uint64(binary.LittleEndian.Uint32(b.Read(4)))
}

// #endregion
