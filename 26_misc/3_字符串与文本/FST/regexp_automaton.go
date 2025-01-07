package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"regexp/syntax"
	"unicode"

	"unicode/utf8"
	unicode_utf8 "unicode/utf8"
)

func main() {
	pattern := "cat.*"
	regexpAutomaton, err := NewRegexp(pattern)
	if err != nil {
		log.Fatalf("Failed to create RegexpAutomaton: %v", err)
	}

	// 匹配字符串
	input := "catdog"
	if regexpAutomaton.MatchesRegex(input) {
		fmt.Printf("'%s' matches the pattern '%s'\n", input, pattern)
	} else {
		fmt.Printf("'%s' does not match the pattern '%s'\n", input, pattern)
	}
}

// #region regexp

// ErrNoEmpty returned when "zero width assertions" are used
var ErrNoEmpty = fmt.Errorf("zero width assertions not allowed")

// ErrNoWordBoundary returned when word boundaries are used
var ErrNoWordBoundary = fmt.Errorf("word boundaries are not allowed")

// ErrNoBytes returned when byte literals are used
var ErrNoBytes = fmt.Errorf("byte literals are not allowed")

// ErrNoLazy returned when lazy quantifiers are used
var ErrNoLazy = fmt.Errorf("lazy quantifiers are not allowed")

// ErrCompiledTooBig returned when regular expression parses into
// too many instructions
var ErrCompiledTooBig = fmt.Errorf("too many instructions")

var DefaultLimit = uint(10 * (1 << 20))

// Regexp implements the vellum.Automaton interface for matcing a user
// specified regular expression.
type Regexp struct {
	orig string
	dfa  *dfa
}

// NewRegexp creates a new Regular Expression automaton with the specified
// expression.  By default it is limited to approximately 10MB for the
// compiled finite state automaton.  If this size is exceeded,
// ErrCompiledTooBig will be returned.
func NewRegexp(expr string) (*Regexp, error) {
	return NewRegexpWithLimit(expr, DefaultLimit)
}

// NewRegexpWithLimit creates a new Regular Expression automaton with
// the specified expression.  The size of the compiled finite state
// automaton exceeds the user specified size,  ErrCompiledTooBig will be
// returned.
func NewRegexpWithLimit(expr string, size uint) (*Regexp, error) {
	parsed, err := syntax.Parse(expr, syntax.Perl)
	if err != nil {
		return nil, err
	}
	return NewRegexpParsedWithLimit(expr, parsed, size)
}

func NewRegexpParsedWithLimit(expr string, parsed *syntax.Regexp, size uint) (*Regexp, error) {
	compiler := newCompiler(size)
	insts, err := compiler.compile(parsed)
	if err != nil {
		return nil, err
	}
	dfaBuilder := newDfaBuilder(insts)
	dfa, err := dfaBuilder.build()
	if err != nil {
		return nil, err
	}
	return &Regexp{
		orig: expr,
		dfa:  dfa,
	}, nil
}

// Start returns the start state of this automaton.
func (r *Regexp) Start() int {
	return 1
}

// IsMatch returns if the specified state is a matching state.
func (r *Regexp) IsMatch(s int) bool {
	if s < len(r.dfa.states) {
		return r.dfa.states[s].match
	}
	return false
}

// CanMatch returns if the specified state can ever transition to a matching
// state.
func (r *Regexp) CanMatch(s int) bool {
	if s < len(r.dfa.states) && s > 0 {
		return true
	}
	return false
}

// WillAlwaysMatch returns if the specified state will always end in a
// matching state.
func (r *Regexp) WillAlwaysMatch(int) bool {
	return false
}

// Accept returns the new state, resulting from the transition byte b
// when currently in the state s.
func (r *Regexp) Accept(s int, b byte) int {
	if s < len(r.dfa.states) {
		return r.dfa.states[s].next[b]
	}
	return 0
}

func (r *Regexp) MatchesRegex(input string) bool {
	currentState := r.Start()
	index := 0
	// Traverse the DFA while characters can still match
	for r.CanMatch(currentState) && index < len(input) {
		currentState = r.Accept(currentState, input[index])
		index++
	}
	return index == len(input) && r.IsMatch(currentState)
}

// #endregion

// #region dfa

// StateLimit is the maximum number of states allowed
const StateLimit = 10000

// ErrTooManyStates is returned if you attempt to build a Levenshtein
// automaton which requires too many states.
var ErrTooManyStates = fmt.Errorf("dfa contains more than %d states",
	StateLimit)

type dfaBuilder struct {
	dfa    *dfa
	cache  map[string]int
	keyBuf []byte
}

func newDfaBuilder(insts prog) *dfaBuilder {
	d := &dfaBuilder{
		dfa: &dfa{
			insts:  insts,
			states: make([]state, 0, 16),
		},
		cache: make(map[string]int, 1024),
	}
	// add 0 state that is invalid
	d.dfa.states = append(d.dfa.states, state{
		next:  make([]int, 256),
		match: false,
	})
	return d
}

func (d *dfaBuilder) build() (*dfa, error) {
	cur := newSparseSet(uint(len(d.dfa.insts)))
	next := newSparseSet(uint(len(d.dfa.insts)))

	d.dfa.add(cur, 0)
	ns, instsReuse := d.cachedState(cur, nil)
	states := intStack{ns}
	seen := make(map[int]struct{})
	var s int
	states, s = states.Pop()
	for s != 0 {
		for b := 0; b < 256; b++ {
			var ns int
			ns, instsReuse = d.runState(cur, next, s, byte(b), instsReuse)
			if ns != 0 {
				if _, ok := seen[ns]; !ok {
					seen[ns] = struct{}{}
					states = states.Push(ns)
				}
			}
			if len(d.dfa.states) > StateLimit {
				return nil, ErrTooManyStates
			}
		}
		states, s = states.Pop()
	}
	return d.dfa, nil
}

func (d *dfaBuilder) runState(cur, next *sparseSet, state int, b byte, instsReuse []uint) (
	int, []uint) {
	cur.Clear()
	for _, ip := range d.dfa.states[state].insts {
		cur.Add(ip)
	}
	d.dfa.run(cur, next, b)
	var nextState int
	nextState, instsReuse = d.cachedState(next, instsReuse)
	d.dfa.states[state].next[b] = nextState
	return nextState, instsReuse
}

func instsKey(insts []uint, buf []byte) []byte {
	if cap(buf) < 8*len(insts) {
		buf = make([]byte, 8*len(insts))
	} else {
		buf = buf[0 : 8*len(insts)]
	}
	for i, inst := range insts {
		binary.LittleEndian.PutUint64(buf[i*8:], uint64(inst))
	}
	return buf
}

func (d *dfaBuilder) cachedState(set *sparseSet,
	instsReuse []uint) (int, []uint) {
	insts := instsReuse[:0]
	if cap(insts) == 0 {
		insts = make([]uint, 0, set.Len())
	}
	var isMatch bool
	for i := uint(0); i < uint(set.Len()); i++ {
		ip := set.Get(i)
		switch d.dfa.insts[ip].op {
		case OpRange:
			insts = append(insts, ip)
		case OpMatch:
			isMatch = true
			insts = append(insts, ip)
		}
	}
	if len(insts) == 0 {
		return 0, insts
	}
	d.keyBuf = instsKey(insts, d.keyBuf)
	v, ok := d.cache[string(d.keyBuf)]
	if ok {
		return v, insts
	}
	d.dfa.states = append(d.dfa.states, state{
		insts: insts,
		next:  make([]int, 256),
		match: isMatch,
	})
	newV := len(d.dfa.states) - 1
	d.cache[string(d.keyBuf)] = newV
	return newV, nil
}

type dfa struct {
	insts  prog
	states []state
}

func (d *dfa) add(set *sparseSet, ip uint) {
	if set.Contains(ip) {
		return
	}
	set.Add(ip)
	switch d.insts[ip].op {
	case OpJmp:
		d.add(set, d.insts[ip].to)
	case OpSplit:
		d.add(set, d.insts[ip].splitA)
		d.add(set, d.insts[ip].splitB)
	}
}

func (d *dfa) run(from, to *sparseSet, b byte) bool {
	to.Clear()
	var isMatch bool
	for i := uint(0); i < uint(from.Len()); i++ {
		ip := from.Get(i)
		switch d.insts[ip].op {
		case OpMatch:
			isMatch = true
		case OpRange:
			if d.insts[ip].rangeStart <= b &&
				b <= d.insts[ip].rangeEnd {
				d.add(to, ip+1)
			}
		}
	}
	return isMatch
}

type state struct {
	insts []uint
	next  []int
	match bool
}

type intStack []int

func (s intStack) Push(v int) intStack {
	return append(s, v)
}

func (s intStack) Pop() (intStack, int) {
	l := len(s)
	if l < 1 {
		return s, 0
	}
	return s[:l-1], s[l-1]
}

// #endregion

// #region compile

type compiler struct {
	sizeLimit uint
	insts     prog
	instsPool []inst

	sequences  Sequences
	rangeStack RangeStack
	startBytes []byte
	endBytes   []byte
}

func newCompiler(sizeLimit uint) *compiler {
	return &compiler{
		sizeLimit:  sizeLimit,
		startBytes: make([]byte, unicode_utf8.UTFMax),
		endBytes:   make([]byte, unicode_utf8.UTFMax),
	}
}

func (c *compiler) compile(ast *syntax.Regexp) (prog, error) {
	err := c.c(ast)
	if err != nil {
		return nil, err
	}
	inst := c.allocInst()
	inst.op = OpMatch
	c.insts = append(c.insts, inst)
	return c.insts, nil
}

func (c *compiler) c(ast *syntax.Regexp) (err error) {
	if ast.Flags&syntax.NonGreedy > 1 {
		return ErrNoLazy
	}

	switch ast.Op {
	case syntax.OpEndLine, syntax.OpBeginLine,
		syntax.OpBeginText, syntax.OpEndText:
		return ErrNoEmpty
	case syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		return ErrNoWordBoundary
	case syntax.OpEmptyMatch:
		return nil
	case syntax.OpLiteral:
		for _, r := range ast.Rune {
			if ast.Flags&syntax.FoldCase > 0 {
				next := syntax.Regexp{
					Op:    syntax.OpCharClass,
					Flags: ast.Flags & syntax.FoldCase,
					Rune0: [2]rune{r, r},
				}
				next.Rune = next.Rune0[0:2]
				// try to find more folded runes
				for r1 := unicode.SimpleFold(r); r1 != r; r1 = unicode.SimpleFold(r1) {
					next.Rune = append(next.Rune, r1, r1)
				}
				err = c.c(&next)
				if err != nil {
					return err
				}
			} else {
				c.sequences, c.rangeStack, err = NewSequencesPrealloc(
					r, r, c.sequences, c.rangeStack, c.startBytes, c.endBytes)
				if err != nil {
					return err
				}
				for _, seq := range c.sequences {
					c.compileUtf8Ranges(seq)
				}
			}
		}
	case syntax.OpAnyChar:
		next := syntax.Regexp{
			Op:    syntax.OpCharClass,
			Flags: ast.Flags & syntax.FoldCase,
			Rune0: [2]rune{0, unicode.MaxRune},
		}
		next.Rune = next.Rune0[:2]
		return c.c(&next)
	case syntax.OpAnyCharNotNL:
		next := syntax.Regexp{
			Op:    syntax.OpCharClass,
			Flags: ast.Flags & syntax.FoldCase,
			Rune:  []rune{0, 0x09, 0x0B, unicode.MaxRune},
		}
		return c.c(&next)
	case syntax.OpCharClass:
		return c.compileClass(ast)
	case syntax.OpCapture:
		return c.c(ast.Sub[0])
	case syntax.OpConcat:
		for _, sub := range ast.Sub {
			err := c.c(sub)
			if err != nil {
				return err
			}
		}
		return nil
	case syntax.OpAlternate:
		if len(ast.Sub) == 0 {
			return nil
		}
		jmpsToEnd := make([]uint, 0, len(ast.Sub)-1)
		// does not handle last entry
		for i := 0; i < len(ast.Sub)-1; i++ {
			sub := ast.Sub[i]
			split := c.emptySplit()
			j1 := c.top()
			err := c.c(sub)
			if err != nil {
				return err
			}
			jmpsToEnd = append(jmpsToEnd, c.emptyJump())
			j2 := c.top()
			c.setSplit(split, j1, j2)
		}
		// handle last entry
		err := c.c(ast.Sub[len(ast.Sub)-1])
		if err != nil {
			return err
		}
		end := uint(len(c.insts))
		for _, jmpToEnd := range jmpsToEnd {
			c.setJump(jmpToEnd, end)
		}
	case syntax.OpQuest:
		split := c.emptySplit()
		j1 := c.top()
		err := c.c(ast.Sub[0])
		if err != nil {
			return err
		}
		j2 := c.top()
		c.setSplit(split, j1, j2)

	case syntax.OpStar:
		j1 := c.top()
		split := c.emptySplit()
		j2 := c.top()
		err := c.c(ast.Sub[0])
		if err != nil {
			return err
		}
		jmp := c.emptyJump()
		j3 := uint(len(c.insts))

		c.setJump(jmp, j1)
		c.setSplit(split, j2, j3)

	case syntax.OpPlus:
		j1 := c.top()
		err := c.c(ast.Sub[0])
		if err != nil {
			return err
		}
		split := c.emptySplit()
		j2 := c.top()
		c.setSplit(split, j1, j2)

	case syntax.OpRepeat:
		if ast.Max == -1 {
			for i := 0; i < ast.Min; i++ {
				err := c.c(ast.Sub[0])
				if err != nil {
					return err
				}
			}
			next := syntax.Regexp{
				Op:    syntax.OpStar,
				Flags: ast.Flags,
				Sub:   ast.Sub,
				Sub0:  ast.Sub0,
				Rune:  ast.Rune,
				Rune0: ast.Rune0,
			}
			return c.c(&next)
		}
		for i := 0; i < ast.Min; i++ {
			err := c.c(ast.Sub[0])
			if err != nil {
				return err
			}
		}
		splits := make([]uint, 0, ast.Max-ast.Min)
		starts := make([]uint, 0, ast.Max-ast.Min)
		for i := ast.Min; i < ast.Max; i++ {
			splits = append(splits, c.emptySplit())
			starts = append(starts, uint(len(c.insts)))
			err := c.c(ast.Sub[0])
			if err != nil {
				return err
			}
		}
		end := uint(len(c.insts))
		for i := 0; i < len(splits); i++ {
			c.setSplit(splits[i], starts[i], end)
		}

	}

	return c.checkSize()
}

func (c *compiler) checkSize() error {
	if uint(len(c.insts)*instSize) > c.sizeLimit {
		return ErrCompiledTooBig
	}
	return nil
}

func (c *compiler) compileClass(ast *syntax.Regexp) error {
	if len(ast.Rune) == 0 {
		return nil
	}
	jmps := make([]uint, 0, len(ast.Rune)-2)
	// does not do last pair
	for i := 0; i < len(ast.Rune)-2; i += 2 {
		rstart := ast.Rune[i]
		rend := ast.Rune[i+1]

		split := c.emptySplit()
		j1 := c.top()
		err := c.compileClassRange(rstart, rend)
		if err != nil {
			return err
		}
		jmps = append(jmps, c.emptyJump())
		j2 := c.top()
		c.setSplit(split, j1, j2)
	}
	// handle last pair
	rstart := ast.Rune[len(ast.Rune)-2]
	rend := ast.Rune[len(ast.Rune)-1]
	err := c.compileClassRange(rstart, rend)
	if err != nil {
		return err
	}
	end := c.top()
	for _, jmp := range jmps {
		c.setJump(jmp, end)
	}
	return nil
}

func (c *compiler) compileClassRange(startR, endR rune) (err error) {
	c.sequences, c.rangeStack, err = NewSequencesPrealloc(
		startR, endR, c.sequences, c.rangeStack, c.startBytes, c.endBytes)
	if err != nil {
		return err
	}
	jmps := make([]uint, 0, len(c.sequences)-1)
	// does not do last entry
	for i := 0; i < len(c.sequences)-1; i++ {
		seq := c.sequences[i]
		split := c.emptySplit()
		j1 := c.top()
		c.compileUtf8Ranges(seq)
		jmps = append(jmps, c.emptyJump())
		j2 := c.top()
		c.setSplit(split, j1, j2)
	}
	// handle last entry
	c.compileUtf8Ranges(c.sequences[len(c.sequences)-1])
	end := c.top()
	for _, jmp := range jmps {
		c.setJump(jmp, end)
	}

	return nil
}

func (c *compiler) compileUtf8Ranges(seq Sequence) {
	for _, r := range seq {
		inst := c.allocInst()
		inst.op = OpRange
		inst.rangeStart = r.Start
		inst.rangeEnd = r.End
		c.insts = append(c.insts, inst)
	}
}

func (c *compiler) emptySplit() uint {
	inst := c.allocInst()
	inst.op = OpSplit
	c.insts = append(c.insts, inst)
	return c.top() - 1
}

func (c *compiler) emptyJump() uint {
	inst := c.allocInst()
	inst.op = OpJmp
	c.insts = append(c.insts, inst)
	return c.top() - 1
}

func (c *compiler) setSplit(i, pc1, pc2 uint) {
	split := c.insts[i]
	split.splitA = pc1
	split.splitB = pc2
}

func (c *compiler) setJump(i, pc uint) {
	jmp := c.insts[i]
	jmp.to = pc
}

func (c *compiler) top() uint {
	return uint(len(c.insts))
}

func (c *compiler) allocInst() *inst {
	if len(c.instsPool) <= 0 {
		c.instsPool = make([]inst, 16)
	}
	inst := &c.instsPool[0]
	c.instsPool = c.instsPool[1:]
	return inst
}

// #endregion

// #region inst

// instOp represents a instruction operation
type instOp int

// the enumeration of operations
const (
	OpMatch instOp = iota
	OpJmp
	OpSplit
	OpRange
)

// instSize is the approximate size of the an inst struct in bytes
const instSize = 40

type inst struct {
	op         instOp
	to         uint
	splitA     uint
	splitB     uint
	rangeStart byte
	rangeEnd   byte
}

func (i *inst) String() string {
	switch i.op {
	case OpJmp:
		return fmt.Sprintf("JMP: %d", i.to)
	case OpSplit:
		return fmt.Sprintf("SPLIT: %d - %d", i.splitA, i.splitB)
	case OpRange:
		return fmt.Sprintf("RANGE: %x - %x", i.rangeStart, i.rangeEnd)
	}
	return "MATCH"
}

type prog []*inst

func (p prog) String() string {
	rv := "\n"
	for i, pi := range p {
		rv += fmt.Sprintf("%d %v\n", i, pi)
	}
	return rv
}

// #endregion

// #region sparse

type sparseSet struct {
	dense  []uint
	sparse []uint
	size   uint
}

func newSparseSet(size uint) *sparseSet {
	return &sparseSet{
		dense:  make([]uint, size),
		sparse: make([]uint, size),
		size:   0,
	}
}

func (s *sparseSet) Len() int {
	return int(s.size)
}

func (s *sparseSet) Add(ip uint) uint {
	i := s.size
	s.dense[i] = ip
	s.sparse[ip] = i
	s.size++
	return i
}

func (s *sparseSet) Get(i uint) uint {
	return s.dense[i]
}

func (s *sparseSet) Contains(ip uint) bool {
	i := s.sparse[ip]
	return i < s.size && s.dense[i] == ip
}

func (s *sparseSet) Clear() {
	s.size = 0
}

// #endregion

// #region utf8

// Sequences is a collection of Sequence
type Sequences []Sequence

// NewSequences constructs a collection of Sequence which describe the
// byte ranges covered between the start and end runes.
func NewSequences(start, end rune) (Sequences, error) {
	rv, _, err := NewSequencesPrealloc(start, end, nil, nil, nil, nil)
	return rv, err
}

func NewSequencesPrealloc(start, end rune,
	preallocSequences Sequences,
	preallocRangeStack RangeStack,
	preallocStartBytes, preallocEndBytes []byte) (Sequences, RangeStack, error) {
	rv := preallocSequences[:0]

	startBytes := preallocStartBytes
	if cap(startBytes) < utf8.UTFMax {
		startBytes = make([]byte, utf8.UTFMax)
	}
	startBytes = startBytes[:utf8.UTFMax]

	endBytes := preallocEndBytes
	if cap(endBytes) < utf8.UTFMax {
		endBytes = make([]byte, utf8.UTFMax)
	}
	endBytes = endBytes[:utf8.UTFMax]

	rangeStack := preallocRangeStack[:0]
	rangeStack = rangeStack.Push(scalarRange{start, end})

	rangeStack, r := rangeStack.Pop()
TOP:
	for r != nilScalarRange {
	INNER:
		for {
			r1, r2 := r.split()
			if r1 != nilScalarRange {
				rangeStack = rangeStack.Push(scalarRange{r2.start, r2.end})
				r.start = r1.start
				r.end = r1.end
				continue INNER
			}
			if !r.valid() {
				rangeStack, r = rangeStack.Pop()
				continue TOP
			}
			for i := 1; i < utf8.UTFMax; i++ {
				max := maxScalarValue(i)
				if r.start <= max && max < r.end {
					rangeStack = rangeStack.Push(scalarRange{max + 1, r.end})
					r.end = max
					continue INNER
				}
			}
			asciiRange := r.ascii()
			if asciiRange != nilRange {
				rv = append(rv, Sequence{
					asciiRange,
				})
				rangeStack, r = rangeStack.Pop()
				continue TOP
			}
			for i := uint(1); i < utf8.UTFMax; i++ {
				m := rune((1 << (6 * i)) - 1)
				if (r.start & ^m) != (r.end & ^m) {
					if (r.start & m) != 0 {
						rangeStack = rangeStack.Push(scalarRange{(r.start | m) + 1, r.end})
						r.end = r.start | m
						continue INNER
					}
					if (r.end & m) != m {
						rangeStack = rangeStack.Push(scalarRange{r.end & ^m, r.end})
						r.end = (r.end & ^m) - 1
						continue INNER
					}
				}
			}
			n, m := r.encode(startBytes, endBytes)
			seq, err := SequenceFromEncodedRange(startBytes[0:n], endBytes[0:m])
			if err != nil {
				return nil, nil, err
			}
			rv = append(rv, seq)
			rangeStack, r = rangeStack.Pop()
			continue TOP
		}
	}

	return rv, rangeStack, nil
}

// Sequence is a collection of Range
type Sequence []Range

// SequenceFromEncodedRange creates sequence from the encoded bytes
func SequenceFromEncodedRange(start, end []byte) (Sequence, error) {
	if len(start) != len(end) {
		return nil, fmt.Errorf("byte slices must be the same length")
	}
	switch len(start) {
	case 2:
		return Sequence{
			Range{start[0], end[0]},
			Range{start[1], end[1]},
		}, nil
	case 3:
		return Sequence{
			Range{start[0], end[0]},
			Range{start[1], end[1]},
			Range{start[2], end[2]},
		}, nil
	case 4:
		return Sequence{
			Range{start[0], end[0]},
			Range{start[1], end[1]},
			Range{start[2], end[2]},
			Range{start[3], end[3]},
		}, nil
	}

	return nil, fmt.Errorf("invalid encoded byte length")
}

// Matches checks to see if the provided byte slice matches the Sequence
func (u Sequence) Matches(bytes []byte) bool {
	if len(bytes) < len(u) {
		return false
	}
	for i := 0; i < len(u); i++ {
		if !u[i].matches(bytes[i]) {
			return false
		}
	}
	return true
}

func (u Sequence) String() string {
	switch len(u) {
	case 1:
		return fmt.Sprintf("%v", u[0])
	case 2:
		return fmt.Sprintf("%v%v", u[0], u[1])
	case 3:
		return fmt.Sprintf("%v%v%v", u[0], u[1], u[2])
	case 4:
		return fmt.Sprintf("%v%v%v%v", u[0], u[1], u[2], u[3])
	default:
		return fmt.Sprintf("invalid utf8 sequence")
	}
}

// Range describes a single range of byte values
type Range struct {
	Start byte
	End   byte
}

var nilRange = Range{0xff, 0}

func (u Range) matches(b byte) bool {
	if u.Start <= b && b <= u.End {
		return true
	}
	return false
}

func (u Range) String() string {
	if u.Start == u.End {
		return fmt.Sprintf("[%X]", u.Start)
	}
	return fmt.Sprintf("[%X-%X]", u.Start, u.End)
}

type scalarRange struct {
	start rune
	end   rune
}

var nilScalarRange = scalarRange{0xffff, 0}

func (s *scalarRange) String() string {
	return fmt.Sprintf("ScalarRange(%d,%d)", s.start, s.end)
}

// split this scalar range if it overlaps with a surrogate codepoint
func (s *scalarRange) split() (scalarRange, scalarRange) {
	if s.start < 0xe000 && s.end > 0xd7ff {
		return scalarRange{
				start: s.start,
				end:   0xd7ff,
			},
			scalarRange{
				start: 0xe000,
				end:   s.end,
			}
	}
	return nilScalarRange, nilScalarRange
}

func (s *scalarRange) valid() bool {
	return s.start <= s.end
}

func (s *scalarRange) ascii() Range {
	if s.valid() && s.end <= 0x7f {
		return Range{
			Start: byte(s.start),
			End:   byte(s.end),
		}
	}
	return nilRange
}

// start and end MUST have capacity for utf8.UTFMax bytes
func (s *scalarRange) encode(start, end []byte) (int, int) {
	n := utf8.EncodeRune(start, s.start)
	m := utf8.EncodeRune(end, s.end)
	return n, m
}

type RangeStack []scalarRange

func (s RangeStack) Push(v scalarRange) RangeStack {
	return append(s, v)
}

func (s RangeStack) Pop() (RangeStack, scalarRange) {
	l := len(s)
	if l < 1 {
		return s, nilScalarRange
	}
	return s[:l-1], s[l-1]
}

func maxScalarValue(nbytes int) rune {
	switch nbytes {
	case 1:
		return 0x007f
	case 2:
		return 0x07FF
	case 3:
		return 0xFFFF
	default:
		return 0x10FFFF
	}
}

// #endregion
