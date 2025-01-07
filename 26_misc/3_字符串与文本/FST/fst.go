// https://github.com/blevesearch/vellum

/*
Package vellum 是一个用于构建、序列化并执行 FST（有限状态转换器）的库。

该库的使用分为两个阶段：构建 FST 和使用 FST。

在构建 FST 的过程中，你需要按字典序依次插入键（[]byte 类型）及其对应的值（uint64）。
在插入的同时，数据会被流式写入到底层的 Writer。
完成构建后，必须调用 builder 的 Close() 方法以结束构建过程。

构建完成后，如果你将 FST 序列化到了磁盘上，可以使用 Open() 方法打开它；
如果相应的字节数据已经在内存中，则可以使用 Load() 方法直接加载。
默认情况下，Open() 会使用 mmap，以避免将整个文件全部加载到内存中。

当 FST 准备就绪后，可以使用 Contains() 方法判断某个键是否存在于 FST 中；
使用 Get() 方法不仅能判断键是否存在，还能获取键对应的值；
同时，你也可以使用 Iterator() 方法遍历指定范围内的键/值对。
*/

// builder：实现了构造 FST 的主要逻辑；
// fst：FST 的核心结构、查询操作；
// encoder / decoder：实现 FST 在字节流中的序列化/反序列化；
// iterator：实现区间/顺序遍历 (key-range iteration)；
// merge_iterator：提供将多个迭代器的结果合并生成一个新的 FST 的能力；
// 其它辅助数据结构：registry, builderNode, automaton 等

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	var buf bytes.Buffer
	builder, err := NewFSTBuilder(&buf, nil)
	if err != nil {
		log.Fatal(err)
	}
	{
		// 将 keys 按字典序插入
		err = builder.Insert([]byte("cat"), 1)
		if err != nil {
			log.Fatal(err)
		}
		err = builder.Insert([]byte("dog"), 2)
		if err != nil {
			log.Fatal(err)
		}
		err = builder.Insert([]byte("fish"), 3)
		if err != nil {
			log.Fatal(err)
		}
		err = builder.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	fst, err := Load(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	{
		val, exists, err := fst.Get([]byte("dog"))
		if err != nil {
			log.Fatal(err)
		}
		if exists {
			fmt.Printf("contains dog with val: %d\n", val)
		} else {
			fmt.Printf("does not contain dog")
		}

		iter, err := fst.Iterator(nil, nil)
		if err != nil {
			log.Fatal(err)
		}

		for {
			k, v := iter.Current()
			fmt.Printf("key: %s, val: %d\n", string(k), v)
			err = iter.Next()
			if err != nil {
				break
			}
		}
	}

}

// nommap.
func open(path string) (*FST, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return new(data, nil)
}

// #region vellum
// ErrOutOfOrder is returned when values are not inserted in
// lexicographic order.
var ErrOutOfOrder = errors.New("values not inserted in lexicographic order")

// ErrIteratorDone is returned by Iterator/Next/Seek methods when the
// Current() value pointed to by the iterator is greater than the last
// key in this FST, or outside the configured startKeyInclusive/endKeyExclusive
// range of the Iterator.
var ErrIteratorDone = errors.New("iterator-done")

type BuilderOpts struct {
	Encoder           int // 指定用于序列化 FST 的编码器版本
	RegistryTableSize int // 定义注册表（registry）的哈希表大小
	RegistryMRUSize   int // 每个哈希桶中缓存的最近使用（MRU，Most Recently Used）节点数量
}

// 构建器会将 FST 的底层表示流式写入到提供的 io.Writer 中。
// 这意味着在插入键值对时，数据会即时写入，而不是先全部保存在内存中再一次性序列化。
func NewFSTBuilder(w io.Writer, opts *BuilderOpts) (*Builder, error) {
	return newBuilder(w, opts)
}

// 直接从文件路径加载 FST，适用于文件存储的场景。
// 可利用 mmap 等机制高效访问大文件，节省内存。
func Open(path string) (*FST, error) {
	return open(path)
}

// 从内存中的字节数据加载 FST，适用于数据已在内存中的场景，如网络传输或内存缓存。
func Load(data []byte) (*FST, error) {
	return new(data, nil)
}

// 合并多个迭代器中的键值对，处理重复键，并构建一个新的 FST 写入到指定的 io.Writer.
func Merge(w io.Writer, opts *BuilderOpts, itrs []Iterator, f MergeFunc) error {
	builder, err := NewFSTBuilder(w, opts)
	if err != nil {
		return err
	}

	itr, err := NewMergeIterator(itrs, f)
	for err == nil {
		k, v := itr.Current()
		err = builder.Insert(k, v)
		if err != nil {
			return err
		}
		err = itr.Next()
	}

	if err != nil && err != ErrIteratorDone {
		return err
	}

	err = itr.Close()
	if err != nil {
		return err
	}

	err = builder.Close()
	if err != nil {
		return err
	}

	return nil
}

// #endregion

// #region fst

// FST is an in-memory representation of a finite state transducer,
// capable of returning the uint64 value associated with
// each []byte key stored, as well as enumerating all of the keys
// in order.
type FST struct {
	f       io.Closer
	ver     int
	len     int
	typ     int
	data    []byte
	decoder decoder
}

func new(data []byte, f io.Closer) (rv *FST, err error) {
	rv = &FST{
		data: data,
		f:    f,
	}

	rv.ver, rv.typ, err = decodeHeader(data)
	if err != nil {
		return nil, err
	}

	rv.decoder, err = loadDecoder(rv.ver, rv.data)
	if err != nil {
		return nil, err
	}

	rv.len = rv.decoder.getLen()

	return rv, nil
}

// Contains returns true if this FST contains the specified key.
func (f *FST) Contains(val []byte) (bool, error) {
	_, exists, err := f.Get(val)
	return exists, err
}

// Get returns the value associated with the key.  NOTE: a value of zero
// does not imply the key does not exist, you must consult the second
// return value as well.
func (f *FST) Get(input []byte) (uint64, bool, error) {
	return f.get(input, nil)
}

func (f *FST) get(input []byte, prealloc fstState) (uint64, bool, error) {
	var total uint64
	curr := f.decoder.getRoot()
	state, err := f.decoder.stateAt(curr, prealloc)
	if err != nil {
		return 0, false, err
	}
	for _, c := range input {
		_, curr, output := state.TransitionFor(c)
		if curr == noneAddr {
			return 0, false, nil
		}

		state, err = f.decoder.stateAt(curr, state)
		if err != nil {
			return 0, false, err
		}

		total += output
	}

	if state.Final() {
		total += state.FinalOutput()
		return total, true, nil
	}
	return 0, false, nil
}

// Version returns the encoding version used by this FST instance.
func (f *FST) Version() int {
	return f.ver
}

// Len returns the number of entries in this FST instance.
func (f *FST) Len() int {
	return f.len
}

// Type returns the type of this FST instance.
func (f *FST) Type() int {
	return f.typ
}

// Close will unmap any mmap'd data (if managed by vellum) and it will close
// the backing file (if managed by vellum).  You MUST call Close() for any
// FST instance that is created.
func (f *FST) Close() error {
	if f.f != nil {
		err := f.f.Close()
		if err != nil {
			return err
		}
	}
	f.data = nil
	f.decoder = nil
	return nil
}

// Start returns the start state of this Automaton
func (f *FST) Start() int {
	return f.decoder.getRoot()
}

// IsMatch returns if this state is a matching state in this Automaton
func (f *FST) IsMatch(addr int) bool {
	match, _ := f.IsMatchWithVal(addr)
	return match
}

// CanMatch returns if this state can ever transition to a matching state
// in this Automaton
func (f *FST) CanMatch(addr int) bool {
	if addr == noneAddr {
		return false
	}
	return true
}

// WillAlwaysMatch returns if from this state the Automaton will always
// be in a matching state
func (f *FST) WillAlwaysMatch(int) bool {
	return false
}

// Accept returns the next state for this Automaton on input of byte b
func (f *FST) Accept(addr int, b byte) int {
	next, _ := f.AcceptWithVal(addr, b)
	return next
}

// IsMatchWithVal returns if this state is a matching state in this Automaton
// and also returns the final output value for this state
func (f *FST) IsMatchWithVal(addr int) (bool, uint64) {
	s, err := f.decoder.stateAt(addr, nil)
	if err != nil {
		return false, 0
	}
	return s.Final(), s.FinalOutput()
}

// AcceptWithVal returns the next state for this Automaton on input of byte b
// and also returns the output value for the transition
func (f *FST) AcceptWithVal(addr int, b byte) (int, uint64) {
	s, err := f.decoder.stateAt(addr, nil)
	if err != nil {
		return noneAddr, 0
	}
	_, next, output := s.TransitionFor(b)
	return next, output
}

// Iterator returns a new Iterator capable of enumerating the key/value pairs
// between the provided startKeyInclusive and endKeyExclusive.
func (f *FST) Iterator(startKeyInclusive, endKeyExclusive []byte) (*FSTIterator, error) {
	return newIterator(f, startKeyInclusive, endKeyExclusive, nil)
}

// Search returns a new Iterator capable of enumerating the key/value pairs
// between the provided startKeyInclusive and endKeyExclusive that also
// satisfy the provided automaton.
func (f *FST) Search(aut Automaton, startKeyInclusive, endKeyExclusive []byte) (*FSTIterator, error) {
	return newIterator(f, startKeyInclusive, endKeyExclusive, aut)
}

// Debug is only intended for debug purposes, it simply asks the underlying
// decoder visit each state, and pass it to the provided callback.
func (f *FST) Debug(callback func(int, interface{}) error) error {

	addr := f.decoder.getRoot()
	set := NewBitSet(uint(addr))
	stack := addrStack{addr}

	stateNumber := 0
	stack, addr = stack[:len(stack)-1], stack[len(stack)-1]
	for addr != noneAddr {
		if set.Test(uint(addr)) {
			stack, addr = stack.Pop()
			continue
		}
		set.Set(uint(addr))
		state, err := f.decoder.stateAt(addr, nil)
		if err != nil {
			return err
		}
		err = callback(stateNumber, state)
		if err != nil {
			return err
		}
		for i := 0; i < state.NumTransitions(); i++ {
			tchar := state.TransitionAt(i)
			_, dest, _ := state.TransitionFor(tchar)
			stack = append(stack, dest)
		}
		stateNumber++
		stack, addr = stack.Pop()
	}

	return nil
}

type addrStack []int

func (a addrStack) Pop() (addrStack, int) {
	l := len(a)
	if l < 1 {
		return a, noneAddr
	}
	return a[:l-1], a[l-1]
}

// Reader() returns a Reader instance that a single thread may use to
// retrieve data from the FST
func (f *FST) Reader() (*Reader, error) {
	return &Reader{f: f}, nil
}

func (f *FST) GetMinKey() ([]byte, error) {
	var rv []byte

	curr := f.decoder.getRoot()
	state, err := f.decoder.stateAt(curr, nil)
	if err != nil {
		return nil, err
	}

	for !state.Final() {
		nextTrans := state.TransitionAt(0)
		_, curr, _ = state.TransitionFor(nextTrans)
		state, err = f.decoder.stateAt(curr, state)
		if err != nil {
			return nil, err
		}

		rv = append(rv, nextTrans)
	}

	return rv, nil
}

func (f *FST) GetMaxKey() ([]byte, error) {
	var rv []byte

	curr := f.decoder.getRoot()
	state, err := f.decoder.stateAt(curr, nil)
	if err != nil {
		return nil, err
	}

	for state.NumTransitions() > 0 {
		nextTrans := state.TransitionAt(state.NumTransitions() - 1)
		_, curr, _ = state.TransitionFor(nextTrans)
		state, err = f.decoder.stateAt(curr, state)
		if err != nil {
			return nil, err
		}

		rv = append(rv, nextTrans)
	}

	return rv, nil
}

// A Reader is meant for a single threaded use
type Reader struct {
	f        *FST
	prealloc fstStateV1
}

func (r *Reader) Get(input []byte) (uint64, bool, error) {
	return r.f.get(input, &r.prealloc)
}

// #endregion

// #region builder

var defaultBuilderOpts = &BuilderOpts{
	Encoder:           1,
	RegistryTableSize: 10000,
	RegistryMRUSize:   2,
}

// A Builder is used to build a new FST.  When possible data is
// streamed out to the underlying Writer as soon as possible.
type Builder struct {
	unfinished *unfinishedNodes // 暂存尚未编译成“冻结”节点的那部分前缀路径
	registry   *registry        // 用来去重/合并已经“冻结”的节点——因为多个分支可能合并成相同后缀
	last       []byte
	len        int

	lastAddr int

	encoder encoder
	opts    *BuilderOpts

	builderNodePool *builderNodePool
}

const noneAddr = 1
const emptyAddr = 0

// NewBuilder returns a new Builder which will stream out the
// underlying representation to the provided Writer as the set is built.
func newBuilder(w io.Writer, opts *BuilderOpts) (*Builder, error) {
	if opts == nil {
		opts = defaultBuilderOpts
	}
	builderNodePool := &builderNodePool{}
	rv := &Builder{
		unfinished:      newUnfinishedNodes(builderNodePool),
		registry:        newRegistry(builderNodePool, opts.RegistryTableSize, opts.RegistryMRUSize),
		builderNodePool: builderNodePool,
		opts:            opts,
		lastAddr:        noneAddr,
	}

	var err error
	rv.encoder, err = loadEncoder(opts.Encoder, w)
	if err != nil {
		return nil, err
	}
	err = rv.encoder.start()
	if err != nil {
		return nil, err
	}
	return rv, nil
}

func (b *Builder) Reset(w io.Writer) error {
	b.unfinished.Reset()
	b.registry.Reset()
	b.lastAddr = noneAddr
	b.encoder.reset(w)
	b.last = nil
	b.len = 0

	err := b.encoder.start()
	if err != nil {
		return err
	}
	return nil
}

// Insert the provided value to the set being built.
// NOTE: values must be inserted in lexicographical order.
func (b *Builder) Insert(key []byte, val uint64) error {
	// ensure items are added in lexicographic order
	if bytes.Compare(key, b.last) < 0 {
		return ErrOutOfOrder
	}
	if len(key) == 0 {
		b.len = 1
		b.unfinished.setRootOutput(val)
		return nil
	}

	prefixLen, out := b.unfinished.findCommonPrefixAndSetOutput(key, val)
	b.len++
	err := b.compileFrom(prefixLen)
	if err != nil {
		return err
	}
	b.copyLastKey(key)
	b.unfinished.addSuffix(key[prefixLen:], out)

	return nil
}

func (b *Builder) copyLastKey(key []byte) {
	if b.last == nil {
		b.last = make([]byte, 0, 64)
	} else {
		b.last = b.last[:0]
	}
	b.last = append(b.last, key...)
}

// Close MUST be called after inserting all values.
func (b *Builder) Close() error {
	err := b.compileFrom(0)
	if err != nil {
		return err
	}
	root := b.unfinished.popRoot()
	rootAddr, err := b.compile(root)
	if err != nil {
		return err
	}
	return b.encoder.finish(b.len, rootAddr)
}

func (b *Builder) compileFrom(iState int) error {
	addr := noneAddr
	for iState+1 < len(b.unfinished.stack) {
		var node *builderNode
		if addr == noneAddr {
			node = b.unfinished.popEmpty()
		} else {
			node = b.unfinished.popFreeze(addr)
		}
		var err error
		addr, err = b.compile(node)
		if err != nil {
			return nil
		}
	}
	b.unfinished.topLastFreeze(addr)
	return nil
}

func (b *Builder) compile(node *builderNode) (int, error) {
	if node.final && len(node.trans) == 0 &&
		node.finalOutput == 0 {
		return 0, nil
	}
	found, addr, entry := b.registry.entry(node)
	if found {
		return addr, nil
	}
	addr, err := b.encoder.encodeState(node, b.lastAddr)
	if err != nil {
		return 0, err
	}

	b.lastAddr = addr
	entry.addr = addr
	return addr, nil
}

type unfinishedNodes struct {
	stack []*builderNodeUnfinished

	// cache allocates a reasonable number of builderNodeUnfinished
	// objects up front and tries to keep reusing them
	// because the main data structure is a stack, we assume the
	// same access pattern, and don't track items separately
	// this means calls get() and pushXYZ() must be paired,
	// as well as calls put() and popXYZ()
	cache []builderNodeUnfinished

	builderNodePool *builderNodePool
}

func (u *unfinishedNodes) Reset() {
	u.stack = u.stack[:0]
	for i := 0; i < len(u.cache); i++ {
		u.cache[i] = builderNodeUnfinished{}
	}
	u.pushEmpty(false)
}

func newUnfinishedNodes(p *builderNodePool) *unfinishedNodes {
	rv := &unfinishedNodes{
		stack:           make([]*builderNodeUnfinished, 0, 64),
		cache:           make([]builderNodeUnfinished, 64),
		builderNodePool: p,
	}
	rv.pushEmpty(false)
	return rv
}

// get new builderNodeUnfinished, reusing cache if possible
func (u *unfinishedNodes) get() *builderNodeUnfinished {
	if len(u.stack) < len(u.cache) {
		return &u.cache[len(u.stack)]
	}
	// full now allocate a new one
	return &builderNodeUnfinished{}
}

// return builderNodeUnfinished, clearing it for reuse
func (u *unfinishedNodes) put() {
	if len(u.stack) >= len(u.cache) {
		return
		// do nothing, not part of cache
	}
	u.cache[len(u.stack)] = builderNodeUnfinished{}
}

func (u *unfinishedNodes) findCommonPrefixAndSetOutput(key []byte,
	out uint64) (int, uint64) {
	var i int
	for i < len(key) {
		if i >= len(u.stack) {
			break
		}
		var addPrefix uint64
		if !u.stack[i].hasLastT {
			break
		}
		if u.stack[i].lastIn == key[i] {
			commonPre := outputPrefix(u.stack[i].lastOut, out)
			addPrefix = outputSub(u.stack[i].lastOut, commonPre)
			out = outputSub(out, commonPre)
			u.stack[i].lastOut = commonPre
			i++
		} else {
			break
		}

		if addPrefix != 0 {
			u.stack[i].addOutputPrefix(addPrefix)
		}
	}

	return i, out
}

func (u *unfinishedNodes) pushEmpty(final bool) {
	next := u.get()
	next.node = u.builderNodePool.Get()
	next.node.final = final
	u.stack = append(u.stack, next)
}

func (u *unfinishedNodes) popRoot() *builderNode {
	l := len(u.stack)
	var unfinished *builderNodeUnfinished
	u.stack, unfinished = u.stack[:l-1], u.stack[l-1]
	rv := unfinished.node
	u.put()
	return rv
}

func (u *unfinishedNodes) popFreeze(addr int) *builderNode {
	l := len(u.stack)
	var unfinished *builderNodeUnfinished
	u.stack, unfinished = u.stack[:l-1], u.stack[l-1]
	unfinished.lastCompiled(addr)
	rv := unfinished.node
	u.put()
	return rv
}

func (u *unfinishedNodes) popEmpty() *builderNode {
	l := len(u.stack)
	var unfinished *builderNodeUnfinished
	u.stack, unfinished = u.stack[:l-1], u.stack[l-1]
	rv := unfinished.node
	u.put()
	return rv
}

func (u *unfinishedNodes) setRootOutput(out uint64) {
	u.stack[0].node.final = true
	u.stack[0].node.finalOutput = out
}

func (u *unfinishedNodes) topLastFreeze(addr int) {
	last := len(u.stack) - 1
	u.stack[last].lastCompiled(addr)
}

func (u *unfinishedNodes) addSuffix(bs []byte, out uint64) {
	if len(bs) == 0 {
		return
	}
	last := len(u.stack) - 1
	u.stack[last].hasLastT = true
	u.stack[last].lastIn = bs[0]
	u.stack[last].lastOut = out
	for _, b := range bs[1:] {
		next := u.get()
		next.node = u.builderNodePool.Get()
		next.hasLastT = true
		next.lastIn = b
		next.lastOut = 0
		u.stack = append(u.stack, next)
	}
	u.pushEmpty(true)
}

type builderNodeUnfinished struct {
	node     *builderNode
	lastOut  uint64
	lastIn   byte
	hasLastT bool
}

func (b *builderNodeUnfinished) lastCompiled(addr int) {
	if b.hasLastT {
		transIn := b.lastIn
		transOut := b.lastOut
		b.hasLastT = false
		b.lastOut = 0
		b.node.trans = append(b.node.trans, transition{
			in:   transIn,
			out:  transOut,
			addr: addr,
		})
	}
}

func (b *builderNodeUnfinished) addOutputPrefix(prefix uint64) {
	if b.node.final {
		b.node.finalOutput = outputCat(prefix, b.node.finalOutput)
	}
	for i := range b.node.trans {
		b.node.trans[i].out = outputCat(prefix, b.node.trans[i].out)
	}
	if b.hasLastT {
		b.lastOut = outputCat(prefix, b.lastOut)
	}
}

type builderNode struct {
	finalOutput uint64
	trans       []transition
	final       bool

	// intrusive linked list
	next *builderNode
}

// reset resets the receiver builderNode to a re-usable state.
func (n *builderNode) reset() {
	n.final = false
	n.finalOutput = 0
	n.trans = n.trans[:0]
	n.next = nil
}

func (n *builderNode) equiv(o *builderNode) bool {
	if n.final != o.final {
		return false
	}
	if n.finalOutput != o.finalOutput {
		return false
	}
	if len(n.trans) != len(o.trans) {
		return false
	}
	for i, ntrans := range n.trans {
		otrans := o.trans[i]
		if ntrans.in != otrans.in {
			return false
		}
		if ntrans.addr != otrans.addr {
			return false
		}
		if ntrans.out != otrans.out {
			return false
		}
	}
	return true
}

type transition struct {
	out  uint64
	addr int
	in   byte
}

func outputPrefix(l, r uint64) uint64 {
	if l < r {
		return l
	}
	return r
}

func outputSub(l, r uint64) uint64 {
	return l - r
}

func outputCat(l, r uint64) uint64 {
	return l + r
}

// builderNodePool pools builderNodes using a singly linked list.
//
// NB: builderNode lifecylce is described by the following interactions -
// +------------------------+                            +----------------------+
// |    Unfinished Nodes    |      Transfer once         |        Registry      |
// |(not frozen builderNode)|-----builderNode is ------->| (frozen builderNode) |
// +------------------------+      marked frozen         +----------------------+
//
//	^                                                     |
//	|                                                     |
//	|                                                   Put()
//	| Get() on        +-------------------+             when
//	+-new char--------| builderNode Pool  |<-----------evicted
//	                  +-------------------+
type builderNodePool struct {
	head *builderNode
}

func (p *builderNodePool) Get() *builderNode {
	if p.head == nil {
		return &builderNode{}
	}
	head := p.head
	p.head = p.head.next
	return head
}

func (p *builderNodePool) Put(v *builderNode) {
	if v == nil {
		return
	}
	v.reset()
	v.next = p.head
	p.head = v
}

// #endregion

// #region common
const maxCommon = 1<<6 - 1

func encodeCommon(in byte) byte {
	val := byte((int(commonInputs[in]) + 1) % 256)
	if val > maxCommon {
		return 0
	}
	return val
}

func decodeCommon(in byte) byte {
	return commonInputsInv[in-1]
}

var commonInputs = []byte{
	84,  // '\x00'
	85,  // '\x01'
	86,  // '\x02'
	87,  // '\x03'
	88,  // '\x04'
	89,  // '\x05'
	90,  // '\x06'
	91,  // '\x07'
	92,  // '\x08'
	93,  // '\t'
	94,  // '\n'
	95,  // '\x0b'
	96,  // '\x0c'
	97,  // '\r'
	98,  // '\x0e'
	99,  // '\x0f'
	100, // '\x10'
	101, // '\x11'
	102, // '\x12'
	103, // '\x13'
	104, // '\x14'
	105, // '\x15'
	106, // '\x16'
	107, // '\x17'
	108, // '\x18'
	109, // '\x19'
	110, // '\x1a'
	111, // '\x1b'
	112, // '\x1c'
	113, // '\x1d'
	114, // '\x1e'
	115, // '\x1f'
	116, // ' '
	80,  // '!'
	117, // '"'
	118, // '#'
	79,  // '$'
	39,  // '%'
	30,  // '&'
	81,  // "'"
	75,  // '('
	74,  // ')'
	82,  // '*'
	57,  // '+'
	66,  // ','
	16,  // '-'
	12,  // '.'
	2,   // '/'
	19,  // '0'
	20,  // '1'
	21,  // '2'
	27,  // '3'
	32,  // '4'
	29,  // '5'
	35,  // '6'
	36,  // '7'
	37,  // '8'
	34,  // '9'
	24,  // ':'
	73,  // ';'
	119, // '<'
	23,  // '='
	120, // '>'
	40,  // '?'
	83,  // '@'
	44,  // 'A'
	48,  // 'B'
	42,  // 'C'
	43,  // 'D'
	49,  // 'E'
	46,  // 'F'
	62,  // 'G'
	61,  // 'H'
	47,  // 'I'
	69,  // 'J'
	68,  // 'K'
	58,  // 'L'
	56,  // 'M'
	55,  // 'N'
	59,  // 'O'
	51,  // 'P'
	72,  // 'Q'
	54,  // 'R'
	45,  // 'S'
	52,  // 'T'
	64,  // 'U'
	65,  // 'V'
	63,  // 'W'
	71,  // 'X'
	67,  // 'Y'
	70,  // 'Z'
	77,  // '['
	121, // '\\'
	78,  // ']'
	122, // '^'
	31,  // '_'
	123, // '`'
	4,   // 'a'
	25,  // 'b'
	9,   // 'c'
	17,  // 'd'
	1,   // 'e'
	26,  // 'f'
	22,  // 'g'
	13,  // 'h'
	7,   // 'i'
	50,  // 'j'
	38,  // 'k'
	14,  // 'l'
	15,  // 'm'
	10,  // 'n'
	3,   // 'o'
	8,   // 'p'
	60,  // 'q'
	6,   // 'r'
	5,   // 's'
	0,   // 't'
	18,  // 'u'
	33,  // 'v'
	11,  // 'w'
	41,  // 'x'
	28,  // 'y'
	53,  // 'z'
	124, // '{'
	125, // '|'
	126, // '}'
	76,  // '~'
	127, // '\x7f'
	128, // '\x80'
	129, // '\x81'
	130, // '\x82'
	131, // '\x83'
	132, // '\x84'
	133, // '\x85'
	134, // '\x86'
	135, // '\x87'
	136, // '\x88'
	137, // '\x89'
	138, // '\x8a'
	139, // '\x8b'
	140, // '\x8c'
	141, // '\x8d'
	142, // '\x8e'
	143, // '\x8f'
	144, // '\x90'
	145, // '\x91'
	146, // '\x92'
	147, // '\x93'
	148, // '\x94'
	149, // '\x95'
	150, // '\x96'
	151, // '\x97'
	152, // '\x98'
	153, // '\x99'
	154, // '\x9a'
	155, // '\x9b'
	156, // '\x9c'
	157, // '\x9d'
	158, // '\x9e'
	159, // '\x9f'
	160, // '\xa0'
	161, // '¡'
	162, // '¢'
	163, // '£'
	164, // '¤'
	165, // '¥'
	166, // '¦'
	167, // '§'
	168, // '¨'
	169, // '©'
	170, // 'ª'
	171, // '«'
	172, // '¬'
	173, // '\xad'
	174, // '®'
	175, // '¯'
	176, // '°'
	177, // '±'
	178, // '²'
	179, // '³'
	180, // '´'
	181, // 'µ'
	182, // '¶'
	183, // '·'
	184, // '¸'
	185, // '¹'
	186, // 'º'
	187, // '»'
	188, // '¼'
	189, // '½'
	190, // '¾'
	191, // '¿'
	192, // 'À'
	193, // 'Á'
	194, // 'Â'
	195, // 'Ã'
	196, // 'Ä'
	197, // 'Å'
	198, // 'Æ'
	199, // 'Ç'
	200, // 'È'
	201, // 'É'
	202, // 'Ê'
	203, // 'Ë'
	204, // 'Ì'
	205, // 'Í'
	206, // 'Î'
	207, // 'Ï'
	208, // 'Ð'
	209, // 'Ñ'
	210, // 'Ò'
	211, // 'Ó'
	212, // 'Ô'
	213, // 'Õ'
	214, // 'Ö'
	215, // '×'
	216, // 'Ø'
	217, // 'Ù'
	218, // 'Ú'
	219, // 'Û'
	220, // 'Ü'
	221, // 'Ý'
	222, // 'Þ'
	223, // 'ß'
	224, // 'à'
	225, // 'á'
	226, // 'â'
	227, // 'ã'
	228, // 'ä'
	229, // 'å'
	230, // 'æ'
	231, // 'ç'
	232, // 'è'
	233, // 'é'
	234, // 'ê'
	235, // 'ë'
	236, // 'ì'
	237, // 'í'
	238, // 'î'
	239, // 'ï'
	240, // 'ð'
	241, // 'ñ'
	242, // 'ò'
	243, // 'ó'
	244, // 'ô'
	245, // 'õ'
	246, // 'ö'
	247, // '÷'
	248, // 'ø'
	249, // 'ù'
	250, // 'ú'
	251, // 'û'
	252, // 'ü'
	253, // 'ý'
	254, // 'þ'
	255, // 'ÿ'
}

var commonInputsInv = []byte{
	't',
	'e',
	'/',
	'o',
	'a',
	's',
	'r',
	'i',
	'p',
	'c',
	'n',
	'w',
	'.',
	'h',
	'l',
	'm',
	'-',
	'd',
	'u',
	'0',
	'1',
	'2',
	'g',
	'=',
	':',
	'b',
	'f',
	'3',
	'y',
	'5',
	'&',
	'_',
	'4',
	'v',
	'9',
	'6',
	'7',
	'8',
	'k',
	'%',
	'?',
	'x',
	'C',
	'D',
	'A',
	'S',
	'F',
	'I',
	'B',
	'E',
	'j',
	'P',
	'T',
	'z',
	'R',
	'N',
	'M',
	'+',
	'L',
	'O',
	'q',
	'H',
	'G',
	'W',
	'U',
	'V',
	',',
	'Y',
	'K',
	'J',
	'Z',
	'X',
	'Q',
	';',
	')',
	'(',
	'~',
	'[',
	']',
	'$',
	'!',
	'\'',
	'*',
	'@',
	'\x00',
	'\x01',
	'\x02',
	'\x03',
	'\x04',
	'\x05',
	'\x06',
	'\x07',
	'\x08',
	'\t',
	'\n',
	'\x0b',
	'\x0c',
	'\r',
	'\x0e',
	'\x0f',
	'\x10',
	'\x11',
	'\x12',
	'\x13',
	'\x14',
	'\x15',
	'\x16',
	'\x17',
	'\x18',
	'\x19',
	'\x1a',
	'\x1b',
	'\x1c',
	'\x1d',
	'\x1e',
	'\x1f',
	' ',
	'"',
	'#',
	'<',
	'>',
	'\\',
	'^',
	'`',
	'{',
	'|',
	'}',
	'\x7f',
	'\x80',
	'\x81',
	'\x82',
	'\x83',
	'\x84',
	'\x85',
	'\x86',
	'\x87',
	'\x88',
	'\x89',
	'\x8a',
	'\x8b',
	'\x8c',
	'\x8d',
	'\x8e',
	'\x8f',
	'\x90',
	'\x91',
	'\x92',
	'\x93',
	'\x94',
	'\x95',
	'\x96',
	'\x97',
	'\x98',
	'\x99',
	'\x9a',
	'\x9b',
	'\x9c',
	'\x9d',
	'\x9e',
	'\x9f',
	'\xa0',
	'\xa1',
	'\xa2',
	'\xa3',
	'\xa4',
	'\xa5',
	'\xa6',
	'\xa7',
	'\xa8',
	'\xa9',
	'\xaa',
	'\xab',
	'\xac',
	'\xad',
	'\xae',
	'\xaf',
	'\xb0',
	'\xb1',
	'\xb2',
	'\xb3',
	'\xb4',
	'\xb5',
	'\xb6',
	'\xb7',
	'\xb8',
	'\xb9',
	'\xba',
	'\xbb',
	'\xbc',
	'\xbd',
	'\xbe',
	'\xbf',
	'\xc0',
	'\xc1',
	'\xc2',
	'\xc3',
	'\xc4',
	'\xc5',
	'\xc6',
	'\xc7',
	'\xc8',
	'\xc9',
	'\xca',
	'\xcb',
	'\xcc',
	'\xcd',
	'\xce',
	'\xcf',
	'\xd0',
	'\xd1',
	'\xd2',
	'\xd3',
	'\xd4',
	'\xd5',
	'\xd6',
	'\xd7',
	'\xd8',
	'\xd9',
	'\xda',
	'\xdb',
	'\xdc',
	'\xdd',
	'\xde',
	'\xdf',
	'\xe0',
	'\xe1',
	'\xe2',
	'\xe3',
	'\xe4',
	'\xe5',
	'\xe6',
	'\xe7',
	'\xe8',
	'\xe9',
	'\xea',
	'\xeb',
	'\xec',
	'\xed',
	'\xee',
	'\xef',
	'\xf0',
	'\xf1',
	'\xf2',
	'\xf3',
	'\xf4',
	'\xf5',
	'\xf6',
	'\xf7',
	'\xf8',
	'\xf9',
	'\xfa',
	'\xfb',
	'\xfc',
	'\xfd',
	'\xfe',
	'\xff',
}

// #endregion

// #region decoder

func init() {
	registerDecoder(versionV1, func(data []byte) decoder {
		return newDecoderV1(data)
	})
}

type decoderV1 struct {
	data []byte
}

func newDecoderV1(data []byte) *decoderV1 {
	return &decoderV1{
		data: data,
	}
}

func (d *decoderV1) getRoot() int {
	if len(d.data) < footerSizeV1 {
		return noneAddr
	}
	footer := d.data[len(d.data)-footerSizeV1:]
	root := binary.LittleEndian.Uint64(footer[8:])
	return int(root)
}

func (d *decoderV1) getLen() int {
	if len(d.data) < footerSizeV1 {
		return 0
	}
	footer := d.data[len(d.data)-footerSizeV1:]
	dlen := binary.LittleEndian.Uint64(footer)
	return int(dlen)
}

func (d *decoderV1) stateAt(addr int, prealloc fstState) (fstState, error) {
	state, ok := prealloc.(*fstStateV1)
	if ok && state != nil {
		*state = fstStateV1{} // clear the struct
	} else {
		state = &fstStateV1{}
	}
	err := state.at(d.data, addr)
	if err != nil {
		return nil, err
	}
	return state, nil
}

type fstStateV1 struct {
	data     []byte
	top      int
	bottom   int
	numTrans int

	// single trans only
	singleTransChar byte
	singleTransNext bool
	singleTransAddr uint64
	singleTransOut  uint64

	// shared
	transSize int
	outSize   int

	// multiple trans only
	final       bool
	transTop    int
	transBottom int
	destTop     int
	destBottom  int
	outTop      int
	outBottom   int
	outFinal    int
}

func (f *fstStateV1) isEncodedSingle() bool {
	if f.data[f.top]>>7 > 0 {
		return true
	}
	return false
}

func (f *fstStateV1) at(data []byte, addr int) error {
	f.data = data
	if addr == emptyAddr {
		return f.atZero()
	} else if addr == noneAddr {
		return f.atNone()
	}
	if addr > len(data) || addr < 16 {
		return fmt.Errorf("invalid address %d/%d", addr, len(data))
	}
	f.top = addr
	f.bottom = addr
	if f.isEncodedSingle() {
		return f.atSingle(data, addr)
	}
	return f.atMulti(data, addr)
}

func (f *fstStateV1) atZero() error {
	f.top = 0
	f.bottom = 1
	f.numTrans = 0
	f.final = true
	f.outFinal = 0
	return nil
}

func (f *fstStateV1) atNone() error {
	f.top = 0
	f.bottom = 1
	f.numTrans = 0
	f.final = false
	f.outFinal = 0
	return nil
}

func (f *fstStateV1) atSingle(data []byte, addr int) error {
	// handle single transition case
	f.numTrans = 1
	f.singleTransNext = data[f.top]&transitionNext > 0
	f.singleTransChar = data[f.top] & maxCommon
	if f.singleTransChar == 0 {
		f.bottom-- // extra byte for uncommon
		f.singleTransChar = data[f.bottom]
	} else {
		f.singleTransChar = decodeCommon(f.singleTransChar)
	}
	if f.singleTransNext {
		// now we know the bottom, can compute next addr
		f.singleTransAddr = uint64(f.bottom - 1)
		f.singleTransOut = 0
	} else {
		f.bottom-- // extra byte with pack sizes
		f.transSize, f.outSize = decodePackSize(data[f.bottom])
		f.bottom -= f.transSize // exactly one trans
		f.singleTransAddr = readPackedUint(data[f.bottom : f.bottom+f.transSize])
		if f.outSize > 0 {
			f.bottom -= f.outSize // exactly one out (could be length 0 though)
			f.singleTransOut = readPackedUint(data[f.bottom : f.bottom+f.outSize])
		} else {
			f.singleTransOut = 0
		}
		// need to wait till we know bottom
		if f.singleTransAddr != 0 {
			f.singleTransAddr = uint64(f.bottom) - f.singleTransAddr
		}
	}
	return nil
}

func (f *fstStateV1) atMulti(data []byte, addr int) error {
	// handle multiple transitions case
	f.final = data[f.top]&stateFinal > 0
	f.numTrans = int(data[f.top] & maxNumTrans)
	if f.numTrans == 0 {
		f.bottom-- // extra byte for number of trans
		f.numTrans = int(data[f.bottom])
		if f.numTrans == 1 {
			// can't really be 1 here, this is special case that means 256
			f.numTrans = 256
		}
	}
	f.bottom-- // extra byte with pack sizes
	f.transSize, f.outSize = decodePackSize(data[f.bottom])

	f.transTop = f.bottom
	f.bottom -= f.numTrans // one byte for each transition
	f.transBottom = f.bottom

	f.destTop = f.bottom
	f.bottom -= f.numTrans * f.transSize
	f.destBottom = f.bottom

	if f.outSize > 0 {
		f.outTop = f.bottom
		f.bottom -= f.numTrans * f.outSize
		f.outBottom = f.bottom
		if f.final {
			f.bottom -= f.outSize
			f.outFinal = f.bottom
		}
	}
	return nil
}

func (f *fstStateV1) Address() int {
	return f.top
}

func (f *fstStateV1) Final() bool {
	return f.final
}

func (f *fstStateV1) FinalOutput() uint64 {
	if f.final && f.outSize > 0 {
		return readPackedUint(f.data[f.outFinal : f.outFinal+f.outSize])
	}
	return 0
}

func (f *fstStateV1) NumTransitions() int {
	return f.numTrans
}

func (f *fstStateV1) TransitionAt(i int) byte {
	if f.isEncodedSingle() {
		return f.singleTransChar
	}
	transitionKeys := f.data[f.transBottom:f.transTop]
	return transitionKeys[f.numTrans-i-1]
}

func (f *fstStateV1) TransitionFor(b byte) (int, int, uint64) {
	if f.isEncodedSingle() {
		if f.singleTransChar == b {
			return 0, int(f.singleTransAddr), f.singleTransOut
		}
		return -1, noneAddr, 0
	}
	transitionKeys := f.data[f.transBottom:f.transTop]
	pos := bytes.IndexByte(transitionKeys, b)
	if pos < 0 {
		return -1, noneAddr, 0
	}
	transDests := f.data[f.destBottom:f.destTop]
	dest := int(readPackedUint(transDests[pos*f.transSize : pos*f.transSize+f.transSize]))
	if dest > 0 {
		// convert delta
		dest = f.bottom - dest
	}
	transVals := f.data[f.outBottom:f.outTop]
	var out uint64
	if f.outSize > 0 {
		out = readPackedUint(transVals[pos*f.outSize : pos*f.outSize+f.outSize])
	}
	return f.numTrans - pos - 1, dest, out
}

func (f *fstStateV1) String() string {
	rv := ""
	rv += fmt.Sprintf("State: %d (%#x)", f.top, f.top)
	if f.final {
		rv += " final"
		fout := f.FinalOutput()
		if fout != 0 {
			rv += fmt.Sprintf(" (%d)", fout)
		}
	}
	rv += "\n"
	rv += fmt.Sprintf("Data: % x\n", f.data[f.bottom:f.top+1])

	for i := 0; i < f.numTrans; i++ {
		transChar := f.TransitionAt(i)
		_, transDest, transOut := f.TransitionFor(transChar)
		rv += fmt.Sprintf(" - %d (%#x) '%s' ---> %d (%#x)  with output: %d", transChar, transChar, string(transChar), transDest, transDest, transOut)
		rv += "\n"
	}
	if f.numTrans == 0 {
		rv += "\n"
	}
	return rv
}

func (f *fstStateV1) DotString(num int) string {
	rv := ""
	label := fmt.Sprintf("%d", num)
	final := ""
	if f.final {
		final = ",peripheries=2"
	}
	rv += fmt.Sprintf("    %d [label=\"%s\"%s];\n", f.top, label, final)

	for i := 0; i < f.numTrans; i++ {
		transChar := f.TransitionAt(i)
		_, transDest, transOut := f.TransitionFor(transChar)
		out := ""
		if transOut != 0 {
			out = fmt.Sprintf("/%d", transOut)
		}
		rv += fmt.Sprintf("    %d -> %d [label=\"%s%s\"];\n", f.top, transDest, escapeInput(transChar), out)
	}

	return rv
}

func escapeInput(b byte) string {
	x := strconv.AppendQuoteRune(nil, rune(b))
	return string(x[1:(len(x) - 1)])
}

// #endregion

// #region encoder

const versionV1 = 1
const oneTransition = 1 << 7
const transitionNext = 1 << 6
const stateFinal = 1 << 6
const footerSizeV1 = 16

func init() {
	registerEncoder(versionV1, func(w io.Writer) encoder {
		return newEncoderV1(w)
	})
}

type encoderV1 struct {
	bw *writer
}

func newEncoderV1(w io.Writer) *encoderV1 {
	return &encoderV1{
		bw: newWriter(w),
	}
}

func (e *encoderV1) reset(w io.Writer) {
	e.bw.Reset(w)
}

func (e *encoderV1) start() error {
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint64(header, versionV1)
	binary.LittleEndian.PutUint64(header[8:], uint64(0)) // type
	n, err := e.bw.Write(header)
	if err != nil {
		return err
	}
	if n != headerSize {
		return fmt.Errorf("short write of header %d/%d", n, headerSize)
	}
	return nil
}

func (e *encoderV1) encodeState(s *builderNode, lastAddr int) (int, error) {
	if len(s.trans) == 0 && s.final && s.finalOutput == 0 {
		return 0, nil
	} else if len(s.trans) != 1 || s.final {
		return e.encodeStateMany(s)
	} else if !s.final && s.trans[0].out == 0 && s.trans[0].addr == lastAddr {
		return e.encodeStateOneFinish(s, transitionNext)
	}
	return e.encodeStateOne(s)
}

func (e *encoderV1) encodeStateOne(s *builderNode) (int, error) {
	start := uint64(e.bw.counter)
	outPackSize := 0
	if s.trans[0].out != 0 {
		outPackSize = packedSize(s.trans[0].out)
		err := e.bw.WritePackedUintIn(s.trans[0].out, outPackSize)
		if err != nil {
			return 0, err
		}
	}
	delta := deltaAddr(start, uint64(s.trans[0].addr))
	transPackSize := packedSize(delta)
	err := e.bw.WritePackedUintIn(delta, transPackSize)
	if err != nil {
		return 0, err
	}

	packSize := encodePackSize(transPackSize, outPackSize)
	err = e.bw.WriteByte(packSize)
	if err != nil {
		return 0, err
	}

	return e.encodeStateOneFinish(s, 0)
}

func (e *encoderV1) encodeStateOneFinish(s *builderNode, next byte) (int, error) {
	enc := encodeCommon(s.trans[0].in)

	// not a common input
	if enc == 0 {
		err := e.bw.WriteByte(s.trans[0].in)
		if err != nil {
			return 0, err
		}
	}
	err := e.bw.WriteByte(oneTransition | next | enc)
	if err != nil {
		return 0, err
	}

	return e.bw.counter - 1, nil
}

func (e *encoderV1) encodeStateMany(s *builderNode) (int, error) {
	start := uint64(e.bw.counter)
	transPackSize := 0
	outPackSize := packedSize(s.finalOutput)
	anyOutputs := s.finalOutput != 0
	for i := range s.trans {
		delta := deltaAddr(start, uint64(s.trans[i].addr))
		tsize := packedSize(delta)
		if tsize > transPackSize {
			transPackSize = tsize
		}
		osize := packedSize(s.trans[i].out)
		if osize > outPackSize {
			outPackSize = osize
		}
		anyOutputs = anyOutputs || s.trans[i].out != 0
	}
	if !anyOutputs {
		outPackSize = 0
	}

	if anyOutputs {
		// output final value
		if s.final {
			err := e.bw.WritePackedUintIn(s.finalOutput, outPackSize)
			if err != nil {
				return 0, err
			}
		}
		// output transition values (in reverse)
		for j := len(s.trans) - 1; j >= 0; j-- {
			err := e.bw.WritePackedUintIn(s.trans[j].out, outPackSize)
			if err != nil {
				return 0, err
			}
		}
	}

	// output transition dests (in reverse)
	for j := len(s.trans) - 1; j >= 0; j-- {
		delta := deltaAddr(start, uint64(s.trans[j].addr))
		err := e.bw.WritePackedUintIn(delta, transPackSize)
		if err != nil {
			return 0, err
		}
	}

	// output transition keys (in reverse)
	for j := len(s.trans) - 1; j >= 0; j-- {
		err := e.bw.WriteByte(s.trans[j].in)
		if err != nil {
			return 0, err
		}
	}

	packSize := encodePackSize(transPackSize, outPackSize)
	err := e.bw.WriteByte(packSize)
	if err != nil {
		return 0, err
	}

	numTrans := encodeNumTrans(len(s.trans))

	// if number of transitions wont fit in edge header byte
	// write out separately
	if numTrans == 0 {
		if len(s.trans) == 256 {
			// this wouldn't fit in single byte, but reuse value 1
			// which would have always fit in the edge header instead
			err = e.bw.WriteByte(1)
			if err != nil {
				return 0, err
			}
		} else {
			err = e.bw.WriteByte(byte(len(s.trans)))
			if err != nil {
				return 0, err
			}
		}
	}

	// finally write edge header
	if s.final {
		numTrans |= stateFinal
	}
	err = e.bw.WriteByte(numTrans)
	if err != nil {
		return 0, err
	}

	return e.bw.counter - 1, nil
}

func (e *encoderV1) finish(count, rootAddr int) error {
	footer := make([]byte, footerSizeV1)
	binary.LittleEndian.PutUint64(footer, uint64(count))        // root addr
	binary.LittleEndian.PutUint64(footer[8:], uint64(rootAddr)) // root addr
	n, err := e.bw.Write(footer)
	if err != nil {
		return err
	}
	if n != footerSizeV1 {
		return fmt.Errorf("short write of footer %d/%d", n, footerSizeV1)
	}
	err = e.bw.Flush()
	if err != nil {
		return err
	}
	return nil
}

// #endregion

// #region encoding
const headerSize = 16

type encoderConstructor func(w io.Writer) encoder
type decoderConstructor func([]byte) decoder

var encoders = map[int]encoderConstructor{}
var decoders = map[int]decoderConstructor{}

type encoder interface {
	start() error
	encodeState(s *builderNode, addr int) (int, error)
	finish(count, rootAddr int) error
	reset(w io.Writer)
}

func loadEncoder(ver int, w io.Writer) (encoder, error) {
	if cons, ok := encoders[ver]; ok {
		return cons(w), nil
	}
	return nil, fmt.Errorf("no encoder for version %d registered", ver)
}

func registerEncoder(ver int, cons encoderConstructor) {
	encoders[ver] = cons
}

type decoder interface {
	getRoot() int
	getLen() int
	stateAt(addr int, prealloc fstState) (fstState, error)
}

func loadDecoder(ver int, data []byte) (decoder, error) {
	if cons, ok := decoders[ver]; ok {
		return cons(data), nil
	}
	return nil, fmt.Errorf("no decoder for version %d registered", ver)
}

func registerDecoder(ver int, cons decoderConstructor) {
	decoders[ver] = cons
}

func decodeHeader(header []byte) (ver int, typ int, err error) {
	if len(header) < headerSize {
		err = fmt.Errorf("invalid header < 16 bytes")
		return
	}
	ver = int(binary.LittleEndian.Uint64(header[0:8]))
	typ = int(binary.LittleEndian.Uint64(header[8:16]))
	return
}

// fstState represents a state inside the FTS runtime
// It is the main contract between the FST impl and the decoder
// The FST impl should work only with this interface, while only the decoder
// impl knows the physical representation.
type fstState interface {
	Address() int
	Final() bool
	FinalOutput() uint64
	NumTransitions() int
	TransitionFor(b byte) (int, int, uint64)
	TransitionAt(i int) byte
}

// #endregion

// #region automaton
// Automaton represents the general contract of a byte-based finite automaton
type Automaton interface {

	// Start returns the start state
	Start() int

	// IsMatch returns true if and only if the state is a match
	IsMatch(int) bool

	// CanMatch returns true if and only if it is possible to reach a match
	// in zero or more steps
	CanMatch(int) bool

	// WillAlwaysMatch returns true if and only if the current state matches
	// and will always match no matter what steps are taken
	WillAlwaysMatch(int) bool

	// Accept returns the next state given the input to the specified state
	Accept(int, byte) int
}

// AutomatonContains implements an generic Contains() method which works
// on any implementation of Automaton
func AutomatonContains(a Automaton, k []byte) bool {
	i := 0
	curr := a.Start()
	for a.CanMatch(curr) && i < len(k) {
		curr = a.Accept(curr, k[i])
		if curr == noneAddr {
			break
		}
		i++
	}
	if i != len(k) {
		return false
	}
	return a.IsMatch(curr)
}

// AlwaysMatch is an Automaton implementation which always matches
type AlwaysMatch struct{}

// Start returns the AlwaysMatch start state
func (m *AlwaysMatch) Start() int {
	return 0
}

// IsMatch always returns true
func (m *AlwaysMatch) IsMatch(int) bool {
	return true
}

// CanMatch always returns true
func (m *AlwaysMatch) CanMatch(int) bool {
	return true
}

// WillAlwaysMatch always returns true
func (m *AlwaysMatch) WillAlwaysMatch(int) bool {
	return true
}

// Accept returns the next AlwaysMatch state
func (m *AlwaysMatch) Accept(int, byte) int {
	return 0
}

// creating an alwaysMatchAutomaton to avoid unnecessary repeated allocations.
var alwaysMatchAutomaton = &AlwaysMatch{}

type FuzzyAutomaton interface {
	Automaton
	EditDistance(int) uint8
	MatchAndDistance(input string) (bool, uint8)
}

// #endregion

// #region iterator
// Iterator represents a means of visiting key/value pairs in order.
type Iterator interface {

	// Current() returns the key/value pair currently pointed to.
	// The []byte of the key is ONLY guaranteed to be valid until
	// another call to Next/Seek/Close.  If you need it beyond that
	// point you MUST make a copy.
	Current() ([]byte, uint64)

	// Next() advances the iterator to the next key/value pair.
	// If no more key/value pairs exist, ErrIteratorDone is returned.
	Next() error

	// Seek() advances the iterator the specified key, or the next key
	// if it does not exist.
	// If no keys exist after that point, ErrIteratorDone is returned.
	Seek(key []byte) error

	// Reset resets the Iterator' internal state to allow for iterator
	// reuse (e.g. pooling).
	Reset(f *FST, startKeyInclusive, endKeyExclusive []byte, aut Automaton) error

	// Close() frees any resources held by this iterator.
	Close() error
}

type FuzzyIterator interface {
	Iterator
	EditDistance() uint8
}

// FSTIterator is a structure for iterating key/value pairs in this FST in
// lexicographic order.  Iterators should be constructed with the FSTIterator
// method on the parent FST structure.
type FSTIterator struct {
	f   *FST
	aut Automaton

	startKeyInclusive []byte
	endKeyExclusive   []byte

	statesStack    []fstState
	keysStack      []byte
	keysPosStack   []int
	valsStack      []uint64
	autStatesStack []int

	nextStart []byte

	editDistance uint8
}

func newIterator(f *FST, startKeyInclusive, endKeyExclusive []byte,
	aut Automaton) (*FSTIterator, error) {

	rv := &FSTIterator{}
	err := rv.Reset(f, startKeyInclusive, endKeyExclusive, aut)
	if err != nil {
		return nil, err
	}
	return rv, nil
}

func (i *FSTIterator) EditDistance() uint8 {
	return i.editDistance
}

// Reset resets the Iterator' internal state to allow for iterator
// reuse (e.g. pooling).
func (i *FSTIterator) Reset(f *FST,
	startKeyInclusive, endKeyExclusive []byte, aut Automaton) error {
	if aut == nil {
		aut = alwaysMatchAutomaton
	}

	i.f = f
	i.startKeyInclusive = startKeyInclusive
	i.endKeyExclusive = endKeyExclusive
	i.aut = aut

	return i.pointTo(startKeyInclusive)
}

// pointTo attempts to point us to the specified location
func (i *FSTIterator) pointTo(key []byte) error {
	// tried to seek before start
	if bytes.Compare(key, i.startKeyInclusive) < 0 {
		key = i.startKeyInclusive
	}

	// tried to see past end
	if i.endKeyExclusive != nil &&
		bytes.Compare(key, i.endKeyExclusive) > 0 {
		key = i.endKeyExclusive
	}

	// reset any state, pointTo always starts over
	i.statesStack = i.statesStack[:0]
	i.keysStack = i.keysStack[:0]
	i.keysPosStack = i.keysPosStack[:0]
	i.valsStack = i.valsStack[:0]
	i.autStatesStack = i.autStatesStack[:0]

	root, err := i.f.decoder.stateAt(i.f.decoder.getRoot(), nil)
	if err != nil {
		return err
	}

	autStart := i.aut.Start()

	maxQ := -1
	// root is always part of the path
	i.statesStack = append(i.statesStack, root)
	i.autStatesStack = append(i.autStatesStack, autStart)
	for j := 0; j < len(key); j++ {
		keyJ := key[j]
		curr := i.statesStack[len(i.statesStack)-1]
		autCurr := i.autStatesStack[len(i.autStatesStack)-1]

		pos, nextAddr, nextVal := curr.TransitionFor(keyJ)
		if nextAddr == noneAddr {
			// needed transition doesn't exist
			// find last trans before the one we needed
			for q := curr.NumTransitions() - 1; q >= 0; q-- {
				if curr.TransitionAt(q) < keyJ {
					maxQ = q
					break
				}
			}
			break
		}
		autNext := i.aut.Accept(autCurr, keyJ)

		next, err := i.f.decoder.stateAt(nextAddr, nil)
		if err != nil {
			return err
		}

		i.statesStack = append(i.statesStack, next)
		i.keysStack = append(i.keysStack, keyJ)
		i.keysPosStack = append(i.keysPosStack, pos)
		i.valsStack = append(i.valsStack, nextVal)
		i.autStatesStack = append(i.autStatesStack, autNext)
		continue
	}

	if !i.statesStack[len(i.statesStack)-1].Final() ||
		!i.aut.IsMatch(i.autStatesStack[len(i.autStatesStack)-1]) ||
		bytes.Compare(i.keysStack, key) < 0 {
		return i.next(maxQ)
	}

	return nil
}

// Current returns the key and value currently pointed to by the iterator.
// If the iterator is not pointing at a valid value (because Iterator/Next/Seek)
// returned an error previously, it may return nil,0.
func (i *FSTIterator) Current() ([]byte, uint64) {
	curr := i.statesStack[len(i.statesStack)-1]
	if curr.Final() {
		var total uint64
		for _, v := range i.valsStack {
			total += v
		}
		total += curr.FinalOutput()
		return i.keysStack, total
	}
	return nil, 0
}

// Next advances this iterator to the next key/value pair.  If there is none
// or the advancement goes beyond the configured endKeyExclusive, then
// ErrIteratorDone is returned.
func (i *FSTIterator) Next() error {
	return i.next(-1)
}

func (i *FSTIterator) next(lastOffset int) error {
	// remember where we started with keysStack in this next() call
	i.nextStart = append(i.nextStart[:0], i.keysStack...)

	nextOffset := lastOffset + 1
	allowCompare := false

OUTER:
	for true {
		curr := i.statesStack[len(i.statesStack)-1]
		autCurr := i.autStatesStack[len(i.autStatesStack)-1]

		if curr.Final() && i.aut.IsMatch(autCurr) && allowCompare {
			// check to see if new keystack might have gone too far
			if i.endKeyExclusive != nil &&
				bytes.Compare(i.keysStack, i.endKeyExclusive) >= 0 {
				return ErrIteratorDone
			}

			cmp := bytes.Compare(i.keysStack, i.nextStart)
			if cmp > 0 {
				if fa, ok := i.aut.(FuzzyAutomaton); ok {
					i.editDistance = fa.EditDistance(autCurr)
				}
				// in final state greater than start key
				return nil
			}
		}

		numTrans := curr.NumTransitions()

	INNER:
		for nextOffset < numTrans {
			t := curr.TransitionAt(nextOffset)

			autNext := i.aut.Accept(autCurr, t)
			if !i.aut.CanMatch(autNext) {
				// TODO: potential optimization to skip nextOffset
				// forwards more directly to something that the
				// automaton likes rather than a linear scan?
				nextOffset += 1
				continue INNER
			}

			pos, nextAddr, v := curr.TransitionFor(t)

			// the next slot in the statesStack might have an
			// fstState instance that we can reuse
			var nextPrealloc fstState
			if len(i.statesStack) < cap(i.statesStack) {
				nextPrealloc = i.statesStack[0:cap(i.statesStack)][len(i.statesStack)]
			}

			// push onto stack
			next, err := i.f.decoder.stateAt(nextAddr, nextPrealloc)
			if err != nil {
				return err
			}

			i.statesStack = append(i.statesStack, next)
			i.keysStack = append(i.keysStack, t)
			i.keysPosStack = append(i.keysPosStack, pos)
			i.valsStack = append(i.valsStack, v)
			i.autStatesStack = append(i.autStatesStack, autNext)

			nextOffset = 0
			allowCompare = true

			continue OUTER
		}

		// no more transitions, so need to backtrack and stack pop
		if len(i.statesStack) <= 1 {
			// stack len is 1 (root), can't go back further, we're done
			break
		}

		// if the top of the stack represents a linear chain of states
		// (i.e., a suffix of nodes linked by single transitions),
		// then optimize by popping the suffix in one shot without
		// going back all the way to the OUTER loop
		var popNum int
		for j := len(i.statesStack) - 1; j > 0; j-- {
			if j == 1 || i.statesStack[j].NumTransitions() != 1 {
				popNum = len(i.statesStack) - 1 - j
				break
			}
		}
		if popNum < 1 { // always pop at least 1 entry from the stacks
			popNum = 1
		}

		nextOffset = i.keysPosStack[len(i.keysPosStack)-popNum] + 1
		allowCompare = false

		i.statesStack = i.statesStack[:len(i.statesStack)-popNum]
		i.keysStack = i.keysStack[:len(i.keysStack)-popNum]
		i.keysPosStack = i.keysPosStack[:len(i.keysPosStack)-popNum]
		i.valsStack = i.valsStack[:len(i.valsStack)-popNum]
		i.autStatesStack = i.autStatesStack[:len(i.autStatesStack)-popNum]
	}

	return ErrIteratorDone
}

// Seek advances this iterator to the specified key/value pair.  If this key
// is not in the FST, Current() will return the next largest key.  If this
// seek operation would go past the last key, or outside the configured
// startKeyInclusive/endKeyExclusive then ErrIteratorDone is returned.
func (i *FSTIterator) Seek(key []byte) error {
	return i.pointTo(key)
}

// Close will free any resources held by this iterator.
func (i *FSTIterator) Close() error {
	// at the moment we don't do anything,
	// but wanted this for API completeness
	return nil
}

// #endregion

// #region merge_iterator

// MergeFunc is used to choose the new value for a key when merging a slice
// of iterators, and the same key is observed with multiple values.
// Values presented to the MergeFunc will be in the same order as the
// original slice creating the MergeIterator.  This allows some MergeFunc
// implementations to prioritize one iterator over another.
type MergeFunc func([]uint64) uint64

// MergeIterator implements the Iterator interface by traversing a slice
// of iterators and merging the contents of them.  If the same key exists
// in mulitipe underlying iterators, a user-provided MergeFunc will be
// invoked to choose the new value.
type MergeIterator struct {
	itrs   []Iterator
	f      MergeFunc
	currKs [][]byte
	currVs []uint64

	lowK    []byte
	lowV    uint64
	lowIdxs []int

	mergeV []uint64
}

// NewMergeIterator creates a new MergeIterator over the provided slice of
// Iterators and with the specified MergeFunc to resolve duplicate keys.
func NewMergeIterator(itrs []Iterator, f MergeFunc) (*MergeIterator, error) {
	rv := &MergeIterator{
		itrs:    itrs,
		f:       f,
		currKs:  make([][]byte, len(itrs)),
		currVs:  make([]uint64, len(itrs)),
		lowIdxs: make([]int, 0, len(itrs)),
		mergeV:  make([]uint64, 0, len(itrs)),
	}
	rv.init()
	if rv.lowK == nil {
		return rv, ErrIteratorDone
	}
	return rv, nil
}

func (m *MergeIterator) init() {
	for i, itr := range m.itrs {
		m.currKs[i], m.currVs[i] = itr.Current()
	}
	m.updateMatches()
}

func (m *MergeIterator) updateMatches() {
	if len(m.itrs) < 1 {
		return
	}
	m.lowK = m.currKs[0]
	m.lowIdxs = m.lowIdxs[:0]
	m.lowIdxs = append(m.lowIdxs, 0)
	for i := 1; i < len(m.itrs); i++ {
		if m.currKs[i] == nil {
			continue
		}
		cmp := bytes.Compare(m.currKs[i], m.lowK)
		if m.lowK == nil || cmp < 0 {
			// reached a new low
			m.lowK = m.currKs[i]
			m.lowIdxs = m.lowIdxs[:0]
			m.lowIdxs = append(m.lowIdxs, i)
		} else if cmp == 0 {
			m.lowIdxs = append(m.lowIdxs, i)
		}
	}
	if len(m.lowIdxs) > 1 {
		// merge multiple values
		m.mergeV = m.mergeV[:0]
		for _, vi := range m.lowIdxs {
			m.mergeV = append(m.mergeV, m.currVs[vi])
		}
		m.lowV = m.f(m.mergeV)
	} else if len(m.lowIdxs) == 1 {
		m.lowV = m.currVs[m.lowIdxs[0]]
	}
}

// Current returns the key and value currently pointed to by this iterator.
// If the iterator is not pointing at a valid value (because Iterator/Next/Seek)
// returned an error previously, it may return nil,0.
func (m *MergeIterator) Current() ([]byte, uint64) {
	return m.lowK, m.lowV
}

// Next advances this iterator to the next key/value pair.  If there is none,
// then ErrIteratorDone is returned.
func (m *MergeIterator) Next() error {
	// move all the current low iterators to next
	for _, vi := range m.lowIdxs {
		err := m.itrs[vi].Next()
		if err != nil && err != ErrIteratorDone {
			return err
		}
		m.currKs[vi], m.currVs[vi] = m.itrs[vi].Current()
	}
	m.updateMatches()
	if m.lowK == nil {
		return ErrIteratorDone
	}
	return nil
}

// Seek advances this iterator to the specified key/value pair.  If this key
// is not in the FST, Current() will return the next largest key.  If this
// seek operation would go past the last key, then ErrIteratorDone is returned.
func (m *MergeIterator) Seek(key []byte) error {
	for i := range m.itrs {
		err := m.itrs[i].Seek(key)
		if err != nil && err != ErrIteratorDone {
			return err
		}
	}
	m.updateMatches()
	if m.lowK == nil {
		return ErrIteratorDone
	}
	return nil
}

// Close will attempt to close all the underlying Iterators.  If any errors
// are encountered, the first will be returned.
func (m *MergeIterator) Close() error {
	var rv error
	for i := range m.itrs {
		// close all iterators, return first error if any
		err := m.itrs[i].Close()
		if rv == nil {
			rv = err
		}
	}
	return rv
}

// MergeMin chooses the minimum value
func MergeMin(vals []uint64) uint64 {
	rv := vals[0]
	for _, v := range vals[1:] {
		if v < rv {
			rv = v
		}
	}
	return rv
}

// MergeMax chooses the maximum value
func MergeMax(vals []uint64) uint64 {
	rv := vals[0]
	for _, v := range vals[1:] {
		if v > rv {
			rv = v
		}
	}
	return rv
}

// MergeSum sums the values
func MergeSum(vals []uint64) uint64 {
	rv := vals[0]
	for _, v := range vals[1:] {
		rv += v
	}
	return rv
}

// #endregion

// #region pack
func deltaAddr(base, trans uint64) uint64 {
	// transition dest of 0 is special case
	if trans == 0 {
		return 0
	}
	return base - trans
}

const packOutMask = 1<<4 - 1

func encodePackSize(transSize, outSize int) byte {
	var rv byte
	rv = byte(transSize << 4)
	rv |= byte(outSize)
	return rv
}

func decodePackSize(pack byte) (transSize int, packSize int) {
	transSize = int(pack >> 4)
	packSize = int(pack & packOutMask)
	return
}

const maxNumTrans = 1<<6 - 1

func encodeNumTrans(n int) byte {
	if n <= maxNumTrans {
		return byte(n)
	}
	return 0
}

func readPackedUint(data []byte) (rv uint64) {
	for i := range data {
		shifted := uint64(data[i]) << uint(i*8)
		rv |= shifted
	}
	return
}

// #endregion

// #region register
type registryCell struct {
	addr int
	node *builderNode
}

type registry struct {
	builderNodePool *builderNodePool
	table           []registryCell
	tableSize       uint
	mruSize         uint
}

func newRegistry(p *builderNodePool, tableSize, mruSize int) *registry {
	nsize := tableSize * mruSize
	rv := &registry{
		builderNodePool: p,
		table:           make([]registryCell, nsize),
		tableSize:       uint(tableSize),
		mruSize:         uint(mruSize),
	}
	return rv
}

func (r *registry) Reset() {
	var empty registryCell
	for i := range r.table {
		r.builderNodePool.Put(r.table[i].node)
		r.table[i] = empty
	}
}

func (r *registry) entry(node *builderNode) (bool, int, *registryCell) {
	if len(r.table) == 0 {
		return false, 0, nil
	}
	bucket := r.hash(node)
	start := r.mruSize * uint(bucket)
	end := start + r.mruSize
	rc := registryCache(r.table[start:end])
	return rc.entry(node, r.builderNodePool)
}

const fnvPrime = 1099511628211

func (r *registry) hash(b *builderNode) int {
	var final uint64
	if b.final {
		final = 1
	}

	var h uint64 = 14695981039346656037
	h = (h ^ final) * fnvPrime
	h = (h ^ b.finalOutput) * fnvPrime
	for _, t := range b.trans {
		h = (h ^ uint64(t.in)) * fnvPrime
		h = (h ^ t.out) * fnvPrime
		h = (h ^ uint64(t.addr)) * fnvPrime
	}
	return int(h % uint64(r.tableSize))
}

type registryCache []registryCell

func (r registryCache) entry(node *builderNode, pool *builderNodePool) (bool, int, *registryCell) {
	if len(r) == 1 {
		if r[0].node != nil && r[0].node.equiv(node) {
			return true, r[0].addr, nil
		}
		pool.Put(r[0].node)
		r[0].node = node
		return false, 0, &r[0]
	}
	for i := range r {
		if r[i].node != nil && r[i].node.equiv(node) {
			addr := r[i].addr
			r.promote(i)
			return true, addr, nil
		}
	}
	// no match
	last := len(r) - 1
	pool.Put(r[last].node)
	r[last].node = node // discard LRU
	r.promote(last)
	return false, 0, &r[0]

}

func (r registryCache) promote(i int) {
	for i > 0 {
		r.swap(i-1, i)
		i--
	}
}

func (r registryCache) swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// #endregion

// #region transducer
// Transducer represents the general contract of a byte-based finite transducer
type Transducer interface {

	// all transducers are also automatons
	Automaton

	// IsMatchWithValue returns true if and only if the state is a match
	// additionally it returns a states final value (if any)
	IsMatchWithVal(int) (bool, uint64)

	// Accept returns the next state given the input to the specified state
	// additionally it returns the value associated with the transition
	AcceptWithVal(int, byte) (int, uint64)
}

// TransducerGet implements an generic Get() method which works
// on any implementation of Transducer
// The caller MUST check the boolean return value for a match.
// Zero is a valid value regardless of match status,
// and if it is NOT a match, the value collected so far is returned.
func TransducerGet(t Transducer, k []byte) (bool, uint64) {
	var total uint64
	i := 0
	curr := t.Start()
	for t.CanMatch(curr) && i < len(k) {
		var transVal uint64
		curr, transVal = t.AcceptWithVal(curr, k[i])
		if curr == noneAddr {
			break
		}
		total += transVal
		i++
	}
	if i != len(k) {
		return false, total
	}
	match, finalVal := t.IsMatchWithVal(curr)
	return match, total + finalVal
}

// #endregion
// #region writer

// A writer is a buffered writer used by vellum. It counts how many bytes have
// been written and has some convenience methods used for encoding the data.
type writer struct {
	w       *bufio.Writer
	counter int
}

func newWriter(w io.Writer) *writer {
	return &writer{
		w: bufio.NewWriter(w),
	}
}

func (w *writer) Reset(newWriter io.Writer) {
	w.w.Reset(newWriter)
	w.counter = 0
}

func (w *writer) WriteByte(c byte) error {
	err := w.w.WriteByte(c)
	if err != nil {
		return err
	}
	w.counter++
	return nil
}

func (w *writer) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.counter += n
	return n, err
}

func (w *writer) Flush() error {
	return w.w.Flush()
}

func (w *writer) WritePackedUintIn(v uint64, n int) error {
	for shift := uint(0); shift < uint(n*8); shift += 8 {
		err := w.WriteByte(byte(v >> shift))
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *writer) WritePackedUint(v uint64) error {
	n := packedSize(v)
	return w.WritePackedUintIn(v, n)
}

func packedSize(n uint64) int {
	if n < 1<<8 {
		return 1
	} else if n < 1<<16 {
		return 2
	} else if n < 1<<24 {
		return 3
	} else if n < 1<<32 {
		return 4
	} else if n < 1<<40 {
		return 5
	} else if n < 1<<48 {
		return 6
	} else if n < 1<<56 {
		return 7
	}
	return 8
}

// #endregion

// #region bitset

// the wordSize of a bit set
const wordSize = 64

// the wordSize of a bit set in bytes
const wordBytes = wordSize / 8

// wordMask is wordSize-1, used for bit indexing in a word
const wordMask = wordSize - 1

// log2WordSize is lg(wordSize)
const log2WordSize = 6

// A BitSet is a set of bits. The zero value of a BitSet is an empty set of length 0.
type BitSet struct {
	length uint
	set    []uint64
}

// wordsNeeded calculates the number of words needed for i bits
func wordsNeeded(i uint) int {
	if i > (Cap() - wordMask) {
		return int(Cap() >> log2WordSize)
	}
	return int((i + wordMask) >> log2WordSize)
}

// wordsIndex calculates the index of words in a `uint64`
func wordsIndex(i uint) uint {
	return i & wordMask
}

// New creates a new BitSet with a hint that length bits will be required.
// The memory usage is at least length/8 bytes.
// In case of allocation failure, the function will return a BitSet with zero
// capacity.
func NewBitSet(length uint) (bset *BitSet) {
	defer func() {
		if r := recover(); r != nil {
			bset = &BitSet{
				0,
				make([]uint64, 0),
			}
		}
	}()

	bset = &BitSet{
		length,
		make([]uint64, wordsNeeded(length)),
	}

	return bset
}

// Test whether bit i is set.
func (b *BitSet) Test(i uint) bool {
	if i >= b.length {
		return false
	}
	return b.set[i>>log2WordSize]&(1<<wordsIndex(i)) != 0
}

func (b *BitSet) Set(i uint) *BitSet {
	if i >= b.length { // if we need more bits, make 'em
		b.extendSet(i)
	}
	b.set[i>>log2WordSize] |= 1 << wordsIndex(i)
	return b
}

// extendSet adds additional words to incorporate new bits if needed
func (b *BitSet) extendSet(i uint) {
	if i >= Cap() {
		panic("You are exceeding the capacity")
	}
	nsize := wordsNeeded(i + 1)
	if b.set == nil {
		b.set = make([]uint64, nsize)
	} else if cap(b.set) >= nsize {
		b.set = b.set[:nsize] // fast resize
	} else if len(b.set) < nsize {
		newset := make([]uint64, nsize, 2*nsize) // increase capacity 2x
		copy(newset, b.set)
		b.set = newset
	}
	b.length = i + 1
}

func Cap() uint {
	return ^uint(0)
}

// #endregion
