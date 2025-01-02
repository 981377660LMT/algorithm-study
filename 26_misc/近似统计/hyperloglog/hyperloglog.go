// 估计集合基数（distinct count）。
// HyperLogLog 在大数据场景下以非常低的内存开销，完成近似的去重计数。
// HyperLogLog 原理简述
//
// 1. 将每个元素做哈希，哈希值高位的 p 位用于定位“第几个桶”，然后根据后续低位计算前导零的个数 \(\rho\)；
// 2. 记录在 regs[bucket] 中的值是该桶所见到的最大前导零数；
// 3. 最终通过合并所有桶的估计值来得出全局基数估计。
//
// api:
// - New14() / New16() / NewNoSparse() / New16NoSparse()：创建一个新的 HyperLogLog Sketch
// - Insert([]byte)：插入数据
// - Estimate()：获取基数估计（distinct count）
// - Merge(*Sketch)：合并两个 Sketch
// - MarshalBinary() / UnmarshalBinary()：序列化 / 反序列化

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/bits"
	"sort"
)

func main() {
	// 1) 创建一个新的 HyperLogLog Sketch
	//    - 可以使用 New14() / New16() 等构造函数，根据需求选择精度 p
	//    - 如果要禁用稀疏模式，可用 NewNoSparse() / New16NoSparse()
	s := New16() // p=16 => 2^16=65536个桶，稀疏模式默认开启

	// 2) 插入数据: 可插入任意 []byte
	s.Insert([]byte("apple"))
	s.Insert([]byte("banana"))

	// 模拟大批量插入，比如插入 10000 个 user-id
	for i := 0; i < 10000; i++ {
		// 可以把 int 转成字符串再转成 []byte
		userID := fmt.Sprintf("user-%d", i)
		s.Insert([]byte(userID))
	}

	// 3) 获取基数估计 (distinct count)
	est := s.Estimate()
	fmt.Printf("Estimated distinct count: %d\n", est)

	// 4) 再创建一个新的 Sketch 并插入一些重复/新数据
	s2 := New16()
	s2.Insert([]byte("apple"))  // apple 在 s 中已存在
	s2.Insert([]byte("cherry")) // cherry 是个新元素

	// 5) 合并两个 Sketch
	//    Merge 会把 s2 的信息合并到 s 中，更新去重计数
	if err := s.Merge(s2); err != nil {
		fmt.Println("Merge error:", err)
	}

	// 合并之后再估计一次
	est = s.Estimate()
	fmt.Printf("Estimated distinct count (after merge): %d\n", est)

	// 6) 序列化（Marshal）到二进制数据，以便存到文件或网络传输
	data, err := s.MarshalBinary()
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}
	fmt.Printf("Serialized length = %d bytes\n", len(data))

	// 7) 反序列化（Unmarshal）回来，验证结果一致
	var s3 Sketch
	if err := s3.UnmarshalBinary(data); err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}
	fmt.Printf("Estimated distinct count (unmarshalled): %d\n", s3.Estimate())
}

// #region hyperloglog
const (
	pp      = uint8(25) // 代表了稀疏模式（Sparse Representation）的一些内部处理所需的位宽（25 位）
	mp      = uint32(1) << pp
	version = 2 // 数据结构序列化的版本号
)

type Sketch struct {
	// 精度（precision）。决定了 m = 2^p 个桶。p 越大，HLL 的精度越高，估计越准确，但内存也越大。
	p uint8
	// m：桶数量，m = 1 << p。
	m     uint32
	alpha float64

	// 如果启用了 稀疏模式，前期会把数据以一种“hash -> set/list”的方式存储，等稀疏数据较多时再切换到普通模式
	tmpSet     *set
	sparseList *compressedList

	// 不使用稀疏模式时，直接使用 regs 这个长度为 m 的数组，每个元素记录某个桶的 ρ 值（见后）。
	regs []uint8
}

func New() *Sketch           { return New14() }                     // New returns a HyperLogLog Sketch with 2^14 registers (precision 14)
func New14() *Sketch         { return newSketchNoError(14, true) }  // New14 returns a HyperLogLog Sketch with 2^14 registers (precision 14)
func New16() *Sketch         { return newSketchNoError(16, true) }  // New16 returns a HyperLogLog Sketch with 2^16 registers (precision 16)
func NewNoSparse() *Sketch   { return newSketchNoError(14, false) } // NewNoSparse returns a HyperLogLog Sketch with 2^14 registers (precision 14) that will not use a sparse representation
func New16NoSparse() *Sketch { return newSketchNoError(16, false) } // New16NoSparse returns a HyperLogLog Sketch with 2^16 registers (precision 16) that will not use a sparse representation

func newSketchNoError(precision uint8, sparse bool) *Sketch {
	sk, _ := NewSketch(precision, sparse)
	return sk
}

func NewSketch(precision uint8, sparse bool) (*Sketch, error) {
	if precision < 4 || precision > 18 {
		return nil, fmt.Errorf("p has to be >= 4 and <= 18")
	}
	m := uint32(1) << precision
	s := &Sketch{
		m:     m,
		p:     precision,
		alpha: alpha(float64(m)),
	}
	if sparse {
		s.tmpSet = newSet(0)
		s.sparseList = newCompressedList(0)
	} else {
		s.regs = make([]uint8, m)
	}
	return s, nil
}

func (sk *Sketch) sparse() bool { return sk.sparseList != nil }

// Clone returns a deep copy of sk.
func (sk *Sketch) Clone() *Sketch {
	clone := *sk
	clone.regs = append([]uint8(nil), sk.regs...)
	clone.tmpSet = sk.tmpSet.Clone()
	clone.sparseList = sk.sparseList.Clone()
	return &clone
}

func (sk *Sketch) maybeToNormal() {
	if uint32(sk.tmpSet.Len())*100 > sk.m {
		sk.mergeSparse()
		if uint32(sk.sparseList.Len()) > sk.m {
			sk.toNormal()
		}
	}
}

func (sk *Sketch) Merge(other *Sketch) error {
	if other == nil {
		return nil
	}
	if sk.p != other.p {
		return errors.New("precisions must be equal")
	}

	if sk.sparse() && other.sparse() {
		sk.mergeSparseSketch(other)
	} else {
		sk.mergeDenseSketch(other)
	}
	return nil
}

func (sk *Sketch) mergeSparseSketch(other *Sketch) {
	sk.tmpSet.Merge(other.tmpSet)
	for iter := other.sparseList.Iter(); iter.HasNext(); {
		sk.tmpSet.add(iter.Next())
	}
	sk.maybeToNormal()
}

func (sk *Sketch) mergeDenseSketch(other *Sketch) {
	if sk.sparse() {
		sk.toNormal()
	}

	if other.sparse() {
		other.tmpSet.ForEach(func(k uint32) {
			i, r := decodeHash(k, other.p, pp)
			sk.insert(i, r)
		})
		for iter := other.sparseList.Iter(); iter.HasNext(); {
			i, r := decodeHash(iter.Next(), other.p, pp)
			sk.insert(i, r)
		}
	} else {
		for i, v := range other.regs {
			if v > sk.regs[i] {
				sk.regs[i] = v
			}
		}
	}
}

func (sk *Sketch) toNormal() {
	if sk.tmpSet.Len() > 0 {
		sk.mergeSparse()
	}

	sk.regs = make([]uint8, sk.m)
	for iter := sk.sparseList.Iter(); iter.HasNext(); {
		i, r := decodeHash(iter.Next(), sk.p, pp)
		sk.insert(i, r)
	}

	sk.tmpSet = nil
	sk.sparseList = nil
}

func (sk *Sketch) insert(i uint32, r uint8) { sk.regs[i] = max(r, sk.regs[i]) }
func (sk *Sketch) Insert(e []byte)          { sk.InsertHash(hash(e)) }

func (sk *Sketch) InsertHash(x uint64) {
	if sk.sparse() {
		if sk.tmpSet.add(encodeHash(x, sk.p, pp)) {
			sk.maybeToNormal()
		}
		return
	}
	i, r := getPosVal(x, sk.p)
	sk.insert(uint32(i), r)
}

func (sk *Sketch) Estimate() uint64 {
	if sk.sparse() {
		sk.mergeSparse()
		return uint64(linearCount(mp, mp-sk.sparseList.count))
	}

	sum, ez := sumAndZeros(sk.regs)
	m := float64(sk.m)

	est := sk.alpha * m * (m - ez) / (sum + beta(sk.p, ez))
	return uint64(est + 0.5)
}

func (sk *Sketch) mergeSparse() {
	if sk.tmpSet.Len() == 0 {
		return
	}

	keys := make(uint64Slice, 0, sk.tmpSet.Len())
	sk.tmpSet.ForEach(func(k uint32) {
		keys = append(keys, k)
	})
	sort.Sort(keys)

	newList := newCompressedList(4*sk.tmpSet.Len() + sk.sparseList.Len())
	for iter, i := sk.sparseList.Iter(), 0; iter.HasNext() || i < len(keys); {
		if !iter.HasNext() {
			newList.Append(keys[i])
			i++
			continue
		}

		if i >= len(keys) {
			newList.Append(iter.Next())
			continue
		}

		x1, x2 := iter.Peek(), keys[i]
		if x1 == x2 {
			newList.Append(iter.Next())
			i++
		} else if x1 > x2 {
			newList.Append(x2)
			i++
		} else {
			newList.Append(iter.Next())
		}
	}

	sk.sparseList = newList
	sk.tmpSet = newSet(0)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (sk *Sketch) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 0, 8+len(sk.regs))
	// Marshal a version marker.
	data = append(data, version)
	// Marshal p.
	data = append(data, sk.p)
	// Marshal b
	data = append(data, 0)

	if sk.sparse() {
		// It's using the sparse Sketch.
		data = append(data, byte(1))

		// Add the tmp_set
		tsdata, err := sk.tmpSet.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, tsdata...)

		// Add the sparse Sketch
		sdata, err := sk.sparseList.MarshalBinary()
		if err != nil {
			return nil, err
		}
		return append(data, sdata...), nil
	}

	// It's using the dense Sketch.
	data = append(data, byte(0))

	// Add the dense sketch Sketch.
	sz := len(sk.regs)
	data = append(data, []byte{
		byte(sz >> 24),
		byte(sz >> 16),
		byte(sz >> 8),
		byte(sz),
	}...)

	// Marshal each element in the list.
	for _, v := range sk.regs {
		data = append(data, byte(v))
	}

	return data, nil
}

// ErrorTooShort is an error that UnmarshalBinary try to parse too short
// binary.
var ErrorTooShort = errors.New("too short binary")

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (sk *Sketch) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return ErrorTooShort
	}

	// Unmarshal version. We may need this in the future if we make
	// non-compatible changes.
	v := data[0]

	// Unmarshal p.
	p := data[1]

	// Unmarshal b.
	b := data[2]

	// Determine if we need a sparse Sketch
	sparse := data[3] == byte(1)

	// Make a newSketch Sketch if the precision doesn't match or if the Sketch was used
	if sk.p != p || sk.regs != nil || sk.tmpSet.Len() > 0 || (sk.sparseList != nil && sk.sparseList.Len() > 0) {
		newh, err := NewSketch(p, sparse)
		if err != nil {
			return err
		}
		*sk = *newh
	}

	// h is now initialised with the correct p. We just need to fill the
	// rest of the details out.
	if sparse {
		// Using the sparse Sketch.

		// Unmarshal the tmp_set.
		tssz := binary.BigEndian.Uint32(data[4:8])
		sk.tmpSet = newSet(int(tssz))

		// We need to unmarshal tssz values in total, and each value requires us
		// to read 4 bytes.
		tsLastByte := int((tssz * 4) + 8)
		for i := 8; i < tsLastByte; i += 4 {
			k := binary.BigEndian.Uint32(data[i : i+4])
			sk.tmpSet.add(k)
		}

		// Unmarshal the sparse Sketch.
		return sk.sparseList.UnmarshalBinary(data[tsLastByte:])
	}

	// Using the dense Sketch.
	sk.sparseList = nil
	sk.tmpSet = nil

	if v == 1 {
		return sk.unmarshalBinaryV1(data[8:], b)
	}
	return sk.unmarshalBinaryV2(data)
}

func sumAndZeros(regs []uint8) (res, ez float64) {
	for _, v := range regs {
		if v == 0 {
			ez++
		}
		res += 1.0 / math.Pow(2.0, float64(v))
	}
	return res, ez
}

func (sk *Sketch) unmarshalBinaryV1(data []byte, b uint8) error {
	sk.regs = make([]uint8, len(data)*2)
	for i, v := range data {
		sk.regs[i*2] = uint8((v >> 4)) + b
		sk.regs[i*2+1] = uint8((v<<4)>>4) + b
	}
	return nil
}

func (sk *Sketch) unmarshalBinaryV2(data []byte) error {
	sk.regs = data[8:]
	return nil
}

// #endregion

// #region sparse
func getIndex(k uint32, p, pp uint8) uint32 {
	if k&1 == 1 {
		return bextr32(k, 32-p, p)
	}
	return bextr32(k, pp-p+1, p)
}

// 这里把一个 64 位 hash 拆分并压缩进 32 位 uint32，作为稀疏结构的标记.
// Encode a hash to be used in the sparse representation.
func encodeHash(x uint64, p, pp uint8) uint32 {
	idx := uint32(bextr(x, 64-pp, pp))
	if bextr(x, 64-pp, pp-p) == 0 {
		zeros := bits.LeadingZeros64((bextr(x, 0, 64-pp)<<pp)|(1<<pp-1)) + 1
		return idx<<7 | uint32(zeros<<1) | 1
	}
	return idx << 1
}

// Decode a hash from the sparse representation.
func decodeHash(k uint32, p, pp uint8) (uint32, uint8) {
	var r uint8
	if k&1 == 1 {
		r = uint8(bextr32(k, 1, 6)) + pp - p
	} else {
		// We can use the 64bit clz implementation and reduce the result
		// by 32 to get a clz for a 32bit word.
		r = uint8(bits.LeadingZeros64(uint64(k<<(32-pp+p-1))) - 31) // -32 + 1
	}
	return getIndex(k, p, pp), r
}

type set struct {
	m *Set[uint32]
}

func newSet(size int) *set {
	return &set{m: NewSet[uint32](size)}
}

func (s *set) ForEach(fn func(v uint32)) {
	s.m.ForEach(func(v uint32) bool {
		fn(v)
		return true
	})
}

func (s *set) Merge(other *set) {
	other.m.ForEach(func(v uint32) bool {
		s.m.Add(v)
		return true
	})
}

func (s *set) Len() int {
	return s.m.Len()
}

func (s *set) add(v uint32) bool {
	if s.m.Has(v) {
		return false
	}
	s.m.Add(v)
	return true
}

func (s *set) Clone() *set {
	if s == nil {
		return nil
	}

	newS := NewSet[uint32](s.m.Len())
	s.m.ForEach(func(v uint32) bool {
		newS.Add(v)
		return true
	})
	return &set{m: newS}
}

func (s *set) MarshalBinary() (data []byte, err error) {
	// 4 bytes for the size of the set, and 4 bytes for each key.
	// list.
	data = make([]byte, 0, 4+(4*s.m.Len()))

	// Length of the set. We only need 32 bits because the size of the set
	// couldn't exceed that on 32 bit architectures.
	sl := s.m.Len()
	data = append(data, []byte{
		byte(sl >> 24),
		byte(sl >> 16),
		byte(sl >> 8),
		byte(sl),
	}...)

	// Marshal each element in the set.
	s.m.ForEach(func(k uint32) bool {
		data = append(data, []byte{
			byte(k >> 24),
			byte(k >> 16),
			byte(k >> 8),
			byte(k),
		}...)
		return true
	})

	return data, nil
}

type uint64Slice []uint32

func (p uint64Slice) Len() int           { return len(p) }
func (p uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// #endregion

// #region utils
var hash = hashFunc

func alpha(m float64) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	}
	return 0.7213 / (1 + 1.079/m)
}

func getPosVal(x uint64, p uint8) (uint64, uint8) {
	i := bextr(x, 64-p, p) // {x63,...,x64-p}
	w := x<<p | 1<<(p-1)   // {x63-p,...,x0}
	rho := uint8(bits.LeadingZeros64(w)) + 1
	return i, rho
}

func linearCount(m uint32, v uint32) float64 {
	fm := float64(m)
	return fm * math.Log(fm/float64(v))
}

func bextr(v uint64, start, length uint8) uint64 {
	return (v >> start) & ((1 << length) - 1)
}

func bextr32(v uint32, start, length uint8) uint32 {
	return (v >> start) & ((1 << length) - 1)
}

func hashFunc(e []byte) uint64 {
	return Hash64(e, 1337)
}

// #endregion

// #region beta
var betaMap = map[uint8]func(float64) float64{
	4:  beta4,
	5:  beta5,
	6:  beta6,
	7:  beta7,
	8:  beta8,
	9:  beta9,
	10: beta10,
	11: beta11,
	12: beta12,
	13: beta13,
	14: beta14,
	15: beta15,
	16: beta16,
	17: beta17,
	18: beta18,
}

func beta(p uint8, ez float64) float64 {
	f, ok := betaMap[p]
	if !ok {
		panic(fmt.Sprintf("invalid precision %d", p))
	}
	return f(ez)
}

/*
p=4
[-0.582581413904517,-1.935300357560050,11.07932375 8035073,-22.131357446444323,22.505391846630037,-12 .000723834917984,3.220579408194167,-0.342225302271 235]
*/
func beta4(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.582581413904517*ez +
		-1.935300357560050*zl +
		11.079323758035073*math.Pow(zl, 2) +
		-22.131357446444323*math.Pow(zl, 3) +
		22.505391846630037*math.Pow(zl, 4) +
		-12.000723834917984*math.Pow(zl, 5) +
		3.220579408194167*math.Pow(zl, 6) +
		-0.342225302271235*math.Pow(zl, 7)
}

/*
p=5
[-0.7518999460733967,-0.9590030077748760,5.5997371 322141607,-8.2097636999765520,6.5091254894472037,- 2.6830293734323729,0.5612891113138221,-0.046333162 2196545]
*/
func beta5(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.7518999460733967*ez +
		-0.9590030077748760*zl +
		5.5997371322141607*math.Pow(zl, 2) +
		-8.2097636999765520*math.Pow(zl, 3) +
		6.5091254894472037*math.Pow(zl, 4) +
		-2.6830293734323729*math.Pow(zl, 5) +
		0.5612891113138221*math.Pow(zl, 6) +
		-0.0463331622196545*math.Pow(zl, 7)
}

/*
p=6
[29.8257900969619634,-31.3287083337725925,-10.5942 523036582283,-11.5720125689099618,3.81887543739074 92,-2.4160130328530811,0.4542208940970826,-0.05751 55452020420]
*/
func beta6(ez float64) float64 {
	zl := math.Log(ez + 1)
	return 29.8257900969619634*ez +
		-31.3287083337725925*zl +
		-10.5942523036582283*math.Pow(zl, 2) +
		-11.5720125689099618*math.Pow(zl, 3) +
		3.8188754373907492*math.Pow(zl, 4) +
		-2.4160130328530811*math.Pow(zl, 5) +
		0.4542208940970826*math.Pow(zl, 6) +
		-0.0575155452020420*math.Pow(zl, 7)
}

/*
p=7
[2.8102921290820060,-3.9780498518175995,1.31626800 41351582,-3.9252486335805901,2.0080835753946471,-0 .7527151937556955,0.1265569894242751,-0.0109946438726240]
*/
func beta7(ez float64) float64 {
	zl := math.Log(ez + 1)
	return 2.8102921290820060*ez +
		-3.9780498518175995*zl +
		1.3162680041351582*math.Pow(zl, 2) +
		-3.9252486335805901*math.Pow(zl, 3) +
		2.0080835753946471*math.Pow(zl, 4) +
		-0.7527151937556955*math.Pow(zl, 5) +
		0.1265569894242751*math.Pow(zl, 6) +
		-0.0109946438726240*math.Pow(zl, 7)
}

/*
p=8
[1.00633544887550519,-2.00580666405112407,1.643697 49366514117,-2.70560809940566172,1.392099802442225 98,-0.46470374272183190,0.07384282377269775,-0.00578554885254223]
*/
func beta8(ez float64) float64 {
	zl := math.Log(ez + 1)
	return 1.00633544887550519*ez +
		-2.00580666405112407*zl +
		1.64369749366514117*math.Pow(zl, 2) +
		-2.70560809940566172*math.Pow(zl, 3) +
		1.39209980244222598*math.Pow(zl, 4) +
		-0.46470374272183190*math.Pow(zl, 5) +
		0.07384282377269775*math.Pow(zl, 6) +
		-0.00578554885254223*math.Pow(zl, 7)
}

/*
p=9
[-0.09415657458167959,-0.78130975924550528,1.71514 946750712460,-1.73711250406516338,0.86441508489048 924,-0.23819027465047218,0.03343448400269076,-0.00 207858528178157]
*/
func beta9(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.09415657458167959*ez +
		-0.78130975924550528*zl +
		1.71514946750712460*math.Pow(zl, 2) +
		-1.73711250406516338*math.Pow(zl, 3) +
		0.86441508489048924*math.Pow(zl, 4) +
		-0.23819027465047218*math.Pow(zl, 5) +
		0.03343448400269076*math.Pow(zl, 6) +
		-0.00207858528178157*math.Pow(zl, 7)
}

/*
p=10
[-0.25935400670790054,-0.52598301999805808,1.48933 034925876839,-1.29642714084993571,0.62284756217221615,-0.15672326770251041,0.02054415903878563,-0.00 112488483925502]
*/
func beta10(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.25935400670790054*ez +
		-0.52598301999805808*zl +
		1.48933034925876839*math.Pow(zl, 2) +
		-1.29642714084993571*math.Pow(zl, 3) +
		0.62284756217221615*math.Pow(zl, 4) +
		-0.15672326770251041*math.Pow(zl, 5) +
		0.02054415903878563*math.Pow(zl, 6) +
		-0.00112488483925502*math.Pow(zl, 7)
}

/*
p=11
[-4.32325553856025e-01,-1.08450736399632e-01,6.091 56550741120e-01,-1.65687801845180e-02,-7.958293410 87617e-02,4.71830602102918e-02,-7.81372902346934e- 03,5.84268708489995e-04]
*/
func beta11(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.432325553856025*ez +
		-0.108450736399632*zl +
		0.609156550741120*math.Pow(zl, 2) +
		-0.0165687801845180*math.Pow(zl, 3) +
		-0.0795829341087617*math.Pow(zl, 4) +
		0.0471830602102918*math.Pow(zl, 5) +
		-0.00781372902346934*math.Pow(zl, 6) +
		0.000584268708489995*math.Pow(zl, 7)
}

/*
p=12
[-3.84979202588598e-01,1.83162233114364e-01,1.3039 6688841854e-01,7.04838927629266e-02,-8.95893971464 453e-03,1.13010036741605e-02,-1.94285569591290e-03 ,2.25435774024964e-04]
*/
func beta12(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.384979202588598*ez +
		0.183162233114364*zl +
		0.130396688841854*math.Pow(zl, 2) +
		0.0704838927629266*math.Pow(zl, 3) +
		-0.0089589397146453*math.Pow(zl, 4) +
		0.0113010036741605*math.Pow(zl, 5) +
		-0.00194285569591290*math.Pow(zl, 6) +
		0.000225435774024964*math.Pow(zl, 7)
}

/*
p=13
[-0.41655270946462997,-0.22146677040685156,0.38862 131236999947,0.45340979746062371,-0.36264738324476 375,0.12304650053558529,-0.01701540384555510,0.001 02750367080838]
*/
func beta13(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.41655270946462997*ez +
		-0.22146677040685156*zl +
		0.38862131236999947*math.Pow(zl, 2) +
		0.45340979746062371*math.Pow(zl, 3) +
		-0.36264738324476375*math.Pow(zl, 4) +
		0.12304650053558529*math.Pow(zl, 5) +
		-0.01701540384555510*math.Pow(zl, 6) +
		0.00102750367080838*math.Pow(zl, 7)
}

/*
p=14
[-3.71009760230692e-01,9.78811941207509e-03,1.8579 6293324165e-01,2.03015527328432e-01,-1.16710521803 686e-01,4.31106699492820e-02,-5.99583540511831e-03 ,4.49704299509437e-04]
*/

func beta14(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.371009760230692*ez +
		0.00978811941207509*zl +
		0.185796293324165*math.Pow(zl, 2) +
		0.203015527328432*math.Pow(zl, 3) +
		-0.116710521803686*math.Pow(zl, 4) +
		0.0431106699492820*math.Pow(zl, 5) +
		-0.00599583540511831*math.Pow(zl, 6) +
		0.000449704299509437*math.Pow(zl, 7)
}

/*
p=15
[-0.38215145543875273,-0.89069400536090837,0.37602 335774678869,0.99335977440682377,-0.65577441638318 956,0.18332342129703610,-0.02241529633062872,0.001 21399789330194]
*/
func beta15(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.38215145543875273*ez +
		-0.89069400536090837*zl +
		0.37602335774678869*math.Pow(zl, 2) +
		0.99335977440682377*math.Pow(zl, 3) +
		-0.65577441638318956*math.Pow(zl, 4) +
		0.18332342129703610*math.Pow(zl, 5) +
		-0.02241529633062872*math.Pow(zl, 6) +
		0.00121399789330194*math.Pow(zl, 7)
}

/*
p=16
[-0.37331876643753059,-1.41704077448122989,0.407291 84796612533,1.56152033906584164,-0.99242233534286128,0.26064681399483092,-0.03053811369682807,0.00155770210179105]
*/
func beta16(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.37331876643753059*ez +
		-1.41704077448122989*zl +
		0.40729184796612533*math.Pow(zl, 2) +
		1.56152033906584164*math.Pow(zl, 3) +
		-0.99242233534286128*math.Pow(zl, 4) +
		0.26064681399483092*math.Pow(zl, 5) +
		-0.03053811369682807*math.Pow(zl, 6) +
		0.00155770210179105*math.Pow(zl, 7)
}

/*
p=17
[-0.36775502299404605,0.53831422351377967,0.769702 89278767923,0.55002583586450560,-0.745755882611469 41,0.25711835785821952,-0.03437902606864149,0.0018 5949146371616]
*/
func beta17(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.36775502299404605*ez +
		0.53831422351377967*zl +
		0.76970289278767923*math.Pow(zl, 2) +
		0.55002583586450560*math.Pow(zl, 3) +
		-0.74575588261146941*math.Pow(zl, 4) +
		0.25711835785821952*math.Pow(zl, 5) +
		-0.03437902606864149*math.Pow(zl, 6) +
		0.00185949146371616*math.Pow(zl, 7)
}

/*
p=18
[-0.36479623325960542,0.99730412328635032,1.553543 86230081221,1.25932677198028919,-1.533259482091101 63,0.47801042200056593,-0.05951025172951174,0.0029 1076804642205]
*/
func beta18(ez float64) float64 {
	zl := math.Log(ez + 1)
	return -0.36479623325960542*ez +
		0.99730412328635032*zl +
		1.55354386230081221*math.Pow(zl, 2) +
		1.25932677198028919*math.Pow(zl, 3) +
		-1.53325948209110163*math.Pow(zl, 4) +
		0.47801042200056593*math.Pow(zl, 5) +
		-0.05951025172951174*math.Pow(zl, 6) +
		0.00291076804642205*math.Pow(zl, 7)
}

// #endregion

// #region compressed
// Original author of this file is github.com/clarkduvall/hyperloglog
type iterable interface {
	decode(i int, last uint32) (uint32, int)
	Len() int
	Iter() *iterator
}

type iterator struct {
	i    int
	last uint32
	v    iterable
}

func (iter *iterator) Next() uint32 {
	n, i := iter.v.decode(iter.i, iter.last)
	iter.last = n
	iter.i = i
	return n
}

func (iter *iterator) Peek() uint32 {
	n, _ := iter.v.decode(iter.i, iter.last)
	return n
}

func (iter iterator) HasNext() bool {
	return iter.i < iter.v.Len()
}

// compressedList 则是一个有序存储的结构，里面用“前缀差分 + 可变长编码”来节省空间.
type compressedList struct {
	count uint32
	last  uint32
	b     variableLengthList
}

func (v *compressedList) Clone() *compressedList {
	if v == nil {
		return nil
	}

	newV := &compressedList{
		count: v.count,
		last:  v.last,
	}

	newV.b = make(variableLengthList, len(v.b))
	copy(newV.b, v.b)
	return newV
}

func (v *compressedList) MarshalBinary() (data []byte, err error) {
	// Marshal the variableLengthList
	bdata, err := v.b.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// At least 4 bytes for the two fixed sized values plus the size of bdata.
	data = make([]byte, 0, 4+4+len(bdata))

	// Marshal the count and last values.
	data = append(data, []byte{
		// Number of items in the list.
		byte(v.count >> 24),
		byte(v.count >> 16),
		byte(v.count >> 8),
		byte(v.count),
		// The last item in the list.
		byte(v.last >> 24),
		byte(v.last >> 16),
		byte(v.last >> 8),
		byte(v.last),
	}...)

	// Append the list
	return append(data, bdata...), nil
}

func (v *compressedList) UnmarshalBinary(data []byte) error {
	if len(data) < 12 {
		return ErrorTooShort
	}

	// Set the count.
	v.count, data = binary.BigEndian.Uint32(data[:4]), data[4:]

	// Set the last value.
	v.last, data = binary.BigEndian.Uint32(data[:4]), data[4:]

	// Set the list.
	sz, data := binary.BigEndian.Uint32(data[:4]), data[4:]
	v.b = make([]uint8, sz)
	if uint32(len(data)) < sz {
		return ErrorTooShort
	}
	for i := uint32(0); i < sz; i++ {
		v.b[i] = data[i]
	}
	return nil
}

func newCompressedList(capacity int) *compressedList {
	v := &compressedList{}
	v.b = make(variableLengthList, 0, capacity)
	return v
}

func (v *compressedList) Len() int {
	return len(v.b)
}

func (v *compressedList) decode(i int, last uint32) (uint32, int) {
	n, i := v.b.decode(i, last)
	return n + last, i
}

func (v *compressedList) Append(x uint32) {
	v.count++
	v.b = v.b.Append(x - v.last)
	v.last = x
}

func (v *compressedList) Iter() *iterator {
	return &iterator{0, 0, v}
}

type variableLengthList []uint8

func (v variableLengthList) MarshalBinary() (data []byte, err error) {
	// 4 bytes for the size of the list, and a byte for each element in the
	// list.
	data = make([]byte, 0, 4+v.Len())

	// Length of the list. We only need 32 bits because the size of the set
	// couldn't exceed that on 32 bit architectures.
	sz := v.Len()
	data = append(data, []byte{
		byte(sz >> 24),
		byte(sz >> 16),
		byte(sz >> 8),
		byte(sz),
	}...)

	// Marshal each element in the list.
	for i := 0; i < sz; i++ {
		data = append(data, v[i])
	}

	return data, nil
}

func (v variableLengthList) Len() int {
	return len(v)
}

func (v *variableLengthList) Iter() *iterator {
	return &iterator{0, 0, v}
}

func (v variableLengthList) decode(i int, last uint32) (uint32, int) {
	var x uint32
	j := i
	for ; v[j]&0x80 != 0; j++ {
		x |= uint32(v[j]&0x7f) << (uint(j-i) * 7)
	}
	x |= uint32(v[j]) << (uint(j-i) * 7)
	return x, j + 1
}

func (v variableLengthList) Append(x uint32) variableLengthList {
	for x&0xffffff80 != 0 {
		v = append(v, uint8((x&0x7f)|0x80))
		x >>= 7
	}
	return append(v, uint8(x&0x7f))
}

// #endregion

// #region intmap
// Set is a specialization of Map modelling a set of integers.
// Like Map, methods that read from the set are valid on the nil Set.
// This include Has, Len, and ForEach.
type Set[K IntKey] Map[K, struct{}]

// NewSet creates a new Set with a given initial capacity.
func NewSet[K IntKey](capacity int) *Set[K] {
	return (*Set[K])(NewMap[K, struct{}](capacity))
}

// Add an element to the set. Returns true if the element was not already present.
func (s *Set[K]) Add(k K) bool {
	_, found := (*Map[K, struct{}])(s).PutIfNotExists(k, struct{}{})
	return found
}

// Del deletes a key, returning true iff the key was found
func (s *Set[K]) Del(k K) bool {
	return (*Map[K, struct{}])(s).Del(k)
}

// Clear removes all items from the Set, but keeps the internal buffers for reuse.
func (s *Set[K]) Clear() {
	(*Map[K, struct{}])(s).Clear()
}

// Has returns true if the key is in the set.
// If the set is nil this method always return false.
func (s *Set[K]) Has(k K) bool {
	return (*Map[K, struct{}])(s).Has(k)
}

// Len returns the number of elements in the set.
// If the set is nil this method return 0.
func (s *Set[K]) Len() int {
	return (*Map[K, struct{}])(s).Len()
}

// ForEach iterates over the elements in the set while the visit function returns true.
// This method returns immediately if the set is nil.
//
// The iteration order of a Set is not defined, so please avoid relying on it.
func (s *Set[K]) ForEach(visit func(k K) bool) {
	(*Map[K, struct{}])(s).ForEach(func(k K, _ struct{}) bool {
		return visit(k)
	})
}

// IntKey is a type constraint for values that can be used as keys in Map
type IntKey interface {
	~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8 | ~uintptr
}

type pair[K IntKey, V any] struct {
	K K
	V V
}

const fillFactor64 = 0.7

func phiMix64(x int) int {
	h := int64(x) * int64(0x9E3779B9)
	return int(h ^ (h >> 16))
}

// Map is a hashmap where the keys are some any integer type.
// It is valid to call methods that read a nil map, similar to a standard Go map.
// Methods valid on a nil map are Has, Get, Len, and ForEach.
type Map[K IntKey, V any] struct {
	data []pair[K, V] // key-value pairs
	size int

	zeroVal    V    // value of 'zero' key
	hasZeroKey bool // do we have 'zero' key in the map?
}

// New creates a new map with keys being any integer subtype.
// The map can store up to the given capacity before reallocation and rehashing occurs.
func NewMap[K IntKey, V any](capacity int) *Map[K, V] {
	return &Map[K, V]{
		data: make([]pair[K, V], arraySize(capacity, fillFactor64)),
	}
}

// Has checks if the given key exists in the map.
// Calling this method on a nil map will return false.
func (m *Map[K, V]) Has(key K) bool {
	if m == nil {
		return false
	}

	if key == K(0) {
		return m.hasZeroKey
	}

	idx := m.startIndex(key)
	p := m.data[idx]

	if p.K == K(0) { // end of chain already
		return false
	}
	if p.K == key { // we check zero prior to this call
		return true
	}

	// hash collision, seek next hash match, bailing on first empty
	for {
		idx = m.nextIndex(idx)
		p = m.data[idx]
		if p.K == K(0) {
			return false
		}
		if p.K == key {
			return true
		}
	}
}

// Get returns the value if the key is found.
// If you just need to check for existence it is easier to use Has.
// Calling this method on a nil map will return the zero value for V and false.
func (m *Map[K, V]) Get(key K) (V, bool) {
	if m == nil {
		var zero V
		return zero, false
	}

	if key == K(0) {
		if m.hasZeroKey {
			return m.zeroVal, true
		}
		var zero V
		return zero, false
	}

	idx := m.startIndex(key)
	p := m.data[idx]

	if p.K == K(0) { // end of chain already
		var zero V
		return zero, false
	}
	if p.K == key { // we check zero prior to this call
		return p.V, true
	}

	// hash collision, seek next hash match, bailing on first empty
	for {
		idx = m.nextIndex(idx)
		p = m.data[idx]
		if p.K == K(0) {
			var zero V
			return zero, false
		}
		if p.K == key {
			return p.V, true
		}
	}
}

// Put adds or updates key with value val.
func (m *Map[K, V]) Put(key K, val V) {
	if key == K(0) {
		if !m.hasZeroKey {
			m.size++
		}
		m.zeroVal = val
		m.hasZeroKey = true
		return
	}

	idx := m.startIndex(key)
	p := &m.data[idx]

	if p.K == K(0) { // end of chain already
		p.K = key
		p.V = val
		if m.size >= m.sizeThreshold() {
			m.rehash()
		} else {
			m.size++
		}
		return
	} else if p.K == key { // overwrite existing value
		p.V = val
		return
	}

	// hash collision, seek next empty or key match
	for {
		idx = m.nextIndex(idx)
		p = &m.data[idx]

		if p.K == K(0) {
			p.K = key
			p.V = val
			if m.size >= m.sizeThreshold() {
				m.rehash()
			} else {
				m.size++
			}
			return
		} else if p.K == key {
			p.V = val
			return
		}
	}
}

// PutIfNotExists adds the key-value pair only if the key does not already exist
// in the map, and returns the current value associated with the key and a boolean
// indicating whether the value was newly added or not.
func (m *Map[K, V]) PutIfNotExists(key K, val V) (V, bool) {
	if key == K(0) {
		if m.hasZeroKey {
			return m.zeroVal, false
		}
		m.zeroVal = val
		m.hasZeroKey = true
		m.size++
		return val, true
	}

	idx := m.startIndex(key)
	p := &m.data[idx]

	if p.K == K(0) { // end of chain already
		p.K = key
		p.V = val
		m.size++
		if m.size >= m.sizeThreshold() {
			m.rehash()
		}
		return val, true
	} else if p.K == key {
		return p.V, false
	}

	// hash collision, seek next hash match, bailing on first empty
	for {
		idx = m.nextIndex(idx)
		p = &m.data[idx]

		if p.K == K(0) {
			p.K = key
			p.V = val
			m.size++
			if m.size >= m.sizeThreshold() {
				m.rehash()
			}
			return val, true
		} else if p.K == key {
			return p.V, false
		}
	}
}

// ForEach iterates through key-value pairs in the map while the function f returns true.
// This method returns immediately if invoked on a nil map.
//
// The iteration order of a Map is not defined, so please avoid relying on it.
func (m *Map[K, V]) ForEach(f func(K, V) bool) {
	if m == nil {
		return
	}

	if m.hasZeroKey && !f(K(0), m.zeroVal) {
		return
	}
	forEach64(m.data, f)
}

// Clear removes all items from the map, but keeps the internal buffers for reuse.
func (m *Map[K, V]) Clear() {
	var zero V
	m.hasZeroKey = false
	m.zeroVal = zero

	// compiles down to runtime.memclr()
	for i := range m.data {
		m.data[i] = pair[K, V]{}
	}

	m.size = 0
}

func (m *Map[K, V]) rehash() {
	oldData := m.data
	m.data = make([]pair[K, V], 2*len(m.data))

	// reset size
	if m.hasZeroKey {
		m.size = 1
	} else {
		m.size = 0
	}

	forEach64(oldData, func(k K, v V) bool {
		m.Put(k, v)
		return true
	})
}

// Len returns the number of elements in the map.
// The length of a nil map is defined to be zero.
func (m *Map[K, V]) Len() int {
	if m == nil {
		return 0
	}

	return m.size
}

func (m *Map[K, V]) sizeThreshold() int {
	return int(math.Floor(float64(len(m.data)) * fillFactor64))
}

func (m *Map[K, V]) startIndex(key K) int {
	return phiMix64(int(key)) & (len(m.data) - 1)
}

func (m *Map[K, V]) nextIndex(idx int) int {
	return (idx + 1) & (len(m.data) - 1)
}

func forEach64[K IntKey, V any](pairs []pair[K, V], f func(k K, v V) bool) {
	for _, p := range pairs {
		if p.K != K(0) && !f(p.K, p.V) {
			return
		}
	}
}

// Del deletes a key and its value, returning true iff the key was found
func (m *Map[K, V]) Del(key K) bool {
	if key == K(0) {
		if m.hasZeroKey {
			m.hasZeroKey = false
			m.size--
			return true
		}
		return false
	}

	idx := m.startIndex(key)
	p := m.data[idx]

	if p.K == key {
		// any keys that were pushed back needs to be shifted nack into the empty slot
		// to avoid breaking the chain
		m.shiftKeys(idx)
		m.size--
		return true
	} else if p.K == K(0) { // end of chain already
		return false
	}

	for {
		idx = m.nextIndex(idx)
		p = m.data[idx]

		if p.K == key {
			// any keys that were pushed back needs to be shifted nack into the empty slot
			// to avoid breaking the chain
			m.shiftKeys(idx)
			m.size--
			return true
		} else if p.K == K(0) {
			return false
		}

	}
}

func (m *Map[K, V]) shiftKeys(idx int) int {
	// Shift entries with the same hash.
	// We need to do this on deletion to ensure we don't have zeroes in the hash chain
	for {
		var p pair[K, V]
		lastIdx := idx
		idx = m.nextIndex(idx)
		for {
			p = m.data[idx]
			if p.K == K(0) {
				m.data[lastIdx] = pair[K, V]{}
				return lastIdx
			}

			slot := m.startIndex(p.K)
			if lastIdx <= idx {
				if lastIdx >= slot || slot > idx {
					break
				}
			} else {
				if lastIdx >= slot && slot > idx {
					break
				}
			}
			idx = m.nextIndex(idx)
		}
		m.data[lastIdx] = p
	}
}

func nextPowerOf2(x uint32) uint32 {
	if x == math.MaxUint32 {
		return x
	}

	if x == 0 {
		return 1
	}

	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16

	return x + 1
}

func arraySize(exp int, fill float64) int {
	s := nextPowerOf2(uint32(math.Ceil(float64(exp) / fill)))
	if s < 2 {
		s = 2
	}
	return int(s)
}

// #endregion

// #region metro

var (
	altHash = [256]uint{}
	masks   = [65]uint{}
)

func init() {
	for i := 0; i < 256; i++ {
		altHash[i] = (uint(Hash64([]byte{byte(i)}, 1337)))
	}
	for i := uint(0); i <= 64; i++ {
		masks[i] = (1 << i) - 1
	}
}

func Hash64(buffer []byte, seed uint64) uint64 {
	const (
		k0 = 0xD6D018F5
		k1 = 0xA2AA033B
		k2 = 0x62992FC1
		k3 = 0x30BC5B29
	)

	ptr := buffer

	hash := (seed + k2) * k0

	if len(ptr) >= 32 {
		v0, v1, v2, v3 := hash, hash, hash, hash

		for len(ptr) >= 32 {
			v0 += binary.LittleEndian.Uint64(ptr[:8]) * k0
			v0 = bits.RotateLeft64(v0, -29) + v2
			v1 += binary.LittleEndian.Uint64(ptr[8:16]) * k1
			v1 = bits.RotateLeft64(v1, -29) + v3
			v2 += binary.LittleEndian.Uint64(ptr[16:24]) * k2
			v2 = bits.RotateLeft64(v2, -29) + v0
			v3 += binary.LittleEndian.Uint64(ptr[24:32]) * k3
			v3 = bits.RotateLeft64(v3, -29) + v1
			ptr = ptr[32:]
		}

		v2 ^= bits.RotateLeft64(((v0+v3)*k0)+v1, -37) * k1
		v3 ^= bits.RotateLeft64(((v1+v2)*k1)+v0, -37) * k0
		v0 ^= bits.RotateLeft64(((v0+v2)*k0)+v3, -37) * k1
		v1 ^= bits.RotateLeft64(((v1+v3)*k1)+v2, -37) * k0
		hash += v0 ^ v1
	}

	if len(ptr) >= 16 {
		v0 := hash + (binary.LittleEndian.Uint64(ptr[:8]) * k2)
		v0 = bits.RotateLeft64(v0, -29) * k3
		v1 := hash + (binary.LittleEndian.Uint64(ptr[8:16]) * k2)
		v1 = bits.RotateLeft64(v1, -29) * k3
		v0 ^= bits.RotateLeft64(v0*k0, -21) + v1
		v1 ^= bits.RotateLeft64(v1*k3, -21) + v0
		hash += v1
		ptr = ptr[16:]
	}

	if len(ptr) >= 8 {
		hash += binary.LittleEndian.Uint64(ptr[:8]) * k3
		ptr = ptr[8:]
		hash ^= bits.RotateLeft64(hash, -55) * k1
	}

	if len(ptr) >= 4 {
		hash += uint64(binary.LittleEndian.Uint32(ptr[:4])) * k3
		hash ^= bits.RotateLeft64(hash, -26) * k1
		ptr = ptr[4:]
	}

	if len(ptr) >= 2 {
		hash += uint64(binary.LittleEndian.Uint16(ptr[:2])) * k3
		ptr = ptr[2:]
		hash ^= bits.RotateLeft64(hash, -48) * k1
	}

	if len(ptr) >= 1 {
		hash += uint64(ptr[0]) * k3
		hash ^= bits.RotateLeft64(hash, -37) * k1
	}

	hash ^= bits.RotateLeft64(hash, -28)
	hash *= k0
	hash ^= bits.RotateLeft64(hash, -29)

	return hash
}

// #endregion
