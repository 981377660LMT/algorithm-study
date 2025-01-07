// https://github.com/blevesearch/vellum/tree/master/levenshtein

package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"unicode/utf8"
)

func main() {
	// 创建 Levenshtein Automaton，设定编辑距离为3
	query := "dog"
	distance := 1
	builder, err := NewLevenshteinAutomatonBuilder(3, false)
	if err != nil {
		log.Fatalf("Failed to create Levenshtein Automaton: %v", err)
	}

	// 构建 DFA
	dfa, err := builder.BuildDfa(query, uint8(distance))
	if err != nil {
		log.Fatalf("Failed to build DFA: %v", err)
	}

	// 匹配字符串
	input := "dog22"
	matched, ed := dfa.MatchAndDistance(input)
	if matched {
		fmt.Printf("'%s' matched with edit distance %d\n", input, ed)
	} else {
		fmt.Printf("'%s' not matched\n", input)
	}
}

// StateLimit is the maximum number of states allowed
const StateLimit = 10000

// ErrTooManyStates is returned if you attempt to build a Levenshtein
// automaton which requires too many states.
var ErrTooManyStates = fmt.Errorf("dfa contains more than %d states",
	StateLimit)

// LevenshteinAutomatonBuilder wraps a precomputed
// datastructure that allows to produce small (but not minimal) DFA.
type LevenshteinAutomatonBuilder struct {
	pDfa *ParametricDFA
}

// NewLevenshteinAutomatonBuilder creates a
// reusable, threadsafe Levenshtein automaton builder.
// `maxDistance` - maximum distance considered by the automaton.
// `transposition` - assign a distance of 1 for transposition
//
// Building this automaton builder is computationally intensive.
// While it takes only a few milliseconds for `d=2`, it grows
// exponentially with `d`. It is only reasonable to `d <= 5`.
func NewLevenshteinAutomatonBuilder(maxDistance uint8,
	transposition bool) (*LevenshteinAutomatonBuilder, error) {
	lnfa := newLevenshtein(maxDistance, transposition)

	pdfa, err := fromNfa(lnfa)
	if err != nil {
		return nil, err
	}

	return &LevenshteinAutomatonBuilder{pDfa: pdfa}, nil
}

// BuildDfa builds the levenshtein automaton for serving
// queries with a given edit distance.
func (lab *LevenshteinAutomatonBuilder) BuildDfa(query string,
	fuzziness uint8) (*DFA, error) {
	return lab.pDfa.buildDfa(query, fuzziness, false)
}

// MaxDistance returns the MaxEdit distance supported by the
// LevenshteinAutomatonBuilder builder.
func (lab *LevenshteinAutomatonBuilder) MaxDistance() uint8 {
	return lab.pDfa.maxDistance
}

// #region parametricDFA
type ParametricState struct {
	shapeID uint32
	offset  uint32
}

func newParametricState() ParametricState {
	return ParametricState{}
}

func (ps *ParametricState) isDeadEnd() bool {
	return ps.shapeID == 0
}

type Transition struct {
	destShapeID uint32
	deltaOffset uint32
}

func (t *Transition) apply(state ParametricState) ParametricState {
	ps := ParametricState{
		shapeID: t.destShapeID}
	// don't need any offset if we are in the dead state,
	// this ensures we have only one dead state.
	if t.destShapeID != 0 {
		ps.offset = state.offset + t.deltaOffset
	}

	return ps
}

type ParametricStateIndex struct {
	stateIndex []uint32
	stateQueue []ParametricState
	numOffsets uint32
}

func newParametricStateIndex(queryLen,
	numParamState uint32) ParametricStateIndex {
	numOffsets := queryLen + 1
	if numParamState == 0 {
		numParamState = numOffsets
	}
	maxNumStates := numParamState * numOffsets
	psi := ParametricStateIndex{
		stateIndex: make([]uint32, maxNumStates),
		stateQueue: make([]ParametricState, 0, 150),
		numOffsets: numOffsets,
	}

	for i := uint32(0); i < maxNumStates; i++ {
		psi.stateIndex[i] = math.MaxUint32
	}
	return psi
}

func (psi *ParametricStateIndex) numStates() int {
	return len(psi.stateQueue)
}

func (psi *ParametricStateIndex) maxNumStates() int {
	return len(psi.stateIndex)
}

func (psi *ParametricStateIndex) get(stateID uint32) ParametricState {
	return psi.stateQueue[stateID]
}

func (psi *ParametricStateIndex) getOrAllocate(ps ParametricState) uint32 {
	bucket := ps.shapeID*psi.numOffsets + ps.offset
	if bucket < uint32(len(psi.stateIndex)) &&
		psi.stateIndex[bucket] != math.MaxUint32 {
		return psi.stateIndex[bucket]
	}
	nState := uint32(len(psi.stateQueue))
	psi.stateQueue = append(psi.stateQueue, ps)

	psi.stateIndex[bucket] = nState
	return nState
}

type ParametricDFA struct {
	distance         []uint8
	transitions      []Transition
	maxDistance      uint8
	transitionStride uint32
	diameter         uint32
}

func (pdfa *ParametricDFA) initialState() ParametricState {
	return ParametricState{shapeID: 1}
}

// Returns true iff whatever characters come afterward,
// we will never reach a shorter distance
func (pdfa *ParametricDFA) isPrefixSink(state ParametricState, queryLen uint32) bool {
	if state.isDeadEnd() {
		return true
	}

	remOffset := queryLen - state.offset
	if remOffset < pdfa.diameter {
		stateDistances := pdfa.distance[pdfa.diameter*state.shapeID:]
		prefixDistance := stateDistances[remOffset]
		if prefixDistance > pdfa.maxDistance {
			return false
		}

		for _, d := range stateDistances {
			if d < prefixDistance {
				return false
			}
		}
		return true
	}
	return false
}

func (pdfa *ParametricDFA) numStates() int {
	return len(pdfa.transitions) / int(pdfa.transitionStride)
}

func min(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

func (pdfa *ParametricDFA) transition(state ParametricState,
	chi uint32) Transition {
	return pdfa.transitions[pdfa.transitionStride*state.shapeID+chi]
}

func (pdfa *ParametricDFA) getDistance(state ParametricState,
	qLen uint32) Distance {
	remainingOffset := qLen - state.offset
	if state.isDeadEnd() || remainingOffset >= pdfa.diameter {
		return Atleast{d: pdfa.maxDistance + 1}
	}
	dist := pdfa.distance[int(pdfa.diameter*state.shapeID)+int(remainingOffset)]
	if dist > pdfa.maxDistance {
		return Atleast{d: dist}
	}
	return Exact{d: dist}
}

func (pdfa *ParametricDFA) computeDistance(left, right string) Distance {
	state := pdfa.initialState()
	leftChars := []rune(left)
	for _, chr := range []rune(right) {
		start := state.offset
		stop := min(start+pdfa.diameter, uint32(len(leftChars)))
		chi := characteristicVector(leftChars[start:stop], chr)
		transition := pdfa.transition(state, uint32(chi))
		state = transition.apply(state)
		if state.isDeadEnd() {
			return Atleast{d: pdfa.maxDistance + 1}
		}
	}
	return pdfa.getDistance(state, uint32(len(left)))
}

func (pdfa *ParametricDFA) buildDfa(query string, distance uint8,
	prefix bool) (*DFA, error) {
	qLen := uint32(len([]rune(query)))
	alphabet := queryChars(query)

	psi := newParametricStateIndex(qLen, uint32(pdfa.numStates()))
	maxNumStates := psi.maxNumStates()
	deadEndStateID := psi.getOrAllocate(newParametricState())
	if deadEndStateID != 0 {
		return nil, fmt.Errorf("Invalid dead end state")
	}

	initialStateID := psi.getOrAllocate(pdfa.initialState())
	dfaBuilder := withMaxStates(uint32(maxNumStates))
	mask := uint32((1 << pdfa.diameter) - 1)

	var stateID int
	for stateID = 0; stateID < StateLimit; stateID++ {
		if stateID == psi.numStates() {
			break
		}
		state := psi.get(uint32(stateID))
		if prefix && pdfa.isPrefixSink(state, qLen) {
			distance := pdfa.getDistance(state, qLen)
			dfaBuilder.addState(uint32(stateID), uint32(stateID), distance)
		} else {
			transition := pdfa.transition(state, 0)
			defSuccessor := transition.apply(state)
			defSuccessorID := psi.getOrAllocate(defSuccessor)
			distance := pdfa.getDistance(state, qLen)
			stateBuilder, err := dfaBuilder.addState(uint32(stateID), defSuccessorID, distance)

			if err != nil {
				return nil, fmt.Errorf("parametric_dfa: buildDfa, err: %v", err)
			}

			alphabet.resetNext()
			chr, cv, err := alphabet.next()
			for err == nil {
				chi := cv.shiftAndMask(state.offset, mask)

				transition := pdfa.transition(state, chi)

				destState := transition.apply(state)

				destStateID := psi.getOrAllocate(destState)

				stateBuilder.addTransition(chr, destStateID)

				chr, cv, err = alphabet.next()
			}
		}
	}

	if stateID == StateLimit {
		return nil, ErrTooManyStates
	}

	dfaBuilder.setInitialState(initialStateID)
	return dfaBuilder.build(distance), nil
}

func fromNfa(nfa *LevenshteinNFA) (*ParametricDFA, error) {
	lookUp := newHash()
	lookUp.getOrAllocate(*newMultiState())
	initialState := nfa.initialStates()
	lookUp.getOrAllocate(*initialState)

	maxDistance := nfa.maxDistance()
	msDiameter := nfa.msDiameter()

	numChi := 1 << msDiameter
	chiValues := make([]uint64, numChi)
	for i := 0; i < numChi; i++ {
		chiValues[i] = uint64(i)
	}

	transitions := make([]Transition, 0, numChi*int(msDiameter))
	var stateID int
	for stateID = 0; stateID < StateLimit; stateID++ {
		if stateID == len(lookUp.items) {
			break
		}

		for _, chi := range chiValues {
			destMs := newMultiState()

			ms := lookUp.getFromID(stateID)

			nfa.transition(ms, destMs, chi)

			translation := destMs.normalize()

			destID := lookUp.getOrAllocate(*destMs)

			transitions = append(transitions, Transition{
				destShapeID: uint32(destID),
				deltaOffset: translation,
			})
		}
	}

	if stateID == StateLimit {
		return nil, ErrTooManyStates
	}

	ns := len(lookUp.items)
	diameter := int(msDiameter)

	distances := make([]uint8, 0, diameter*ns)
	for stateID := 0; stateID < ns; stateID++ {
		ms := lookUp.getFromID(stateID)
		for offset := 0; offset < diameter; offset++ {
			dist := nfa.multistateDistance(ms, uint32(offset))
			distances = append(distances, dist.distance())
		}
	}

	return &ParametricDFA{
		diameter:         uint32(msDiameter),
		transitions:      transitions,
		maxDistance:      maxDistance,
		transitionStride: uint32(numChi),
		distance:         distances,
	}, nil
}

type hash struct {
	index map[[16]byte]int
	items []MultiState
}

func newHash() *hash {
	return &hash{
		index: make(map[[16]byte]int, 100),
		items: make([]MultiState, 0, 100),
	}
}

func (h *hash) getOrAllocate(m MultiState) int {
	size := len(h.items)
	var exists bool
	var pos int
	md5 := getHash(&m)
	if pos, exists = h.index[md5]; !exists {
		h.index[md5] = size
		pos = size
		h.items = append(h.items, m)
	}
	return pos
}

func (h *hash) getFromID(id int) *MultiState {
	return &h.items[id]
}

func getHash(ms *MultiState) [16]byte {
	msBytes := []byte{}
	for _, state := range ms.states {
		jsonBytes, _ := json.Marshal(&state)
		msBytes = append(msBytes, jsonBytes...)
	}
	return md5.Sum(msBytes)
}

// #endregion

// #region levenshtein_nfa

// / Levenshtein Distance computed by a Levenshtein Automaton.
// /
// / Levenshtein automata can only compute the exact Levenshtein distance
// / up to a given `max_distance`.
// /
// / Over this distance, the automaton will invariably
// / return `Distance::AtLeast(max_distance + 1)`.
type Distance interface {
	distance() uint8
}

type Exact struct {
	d uint8
}

func (e Exact) distance() uint8 {
	return e.d
}

type Atleast struct {
	d uint8
}

func (a Atleast) distance() uint8 {
	return a.d
}

func characteristicVector(query []rune, c rune) uint64 {
	chi := uint64(0)
	for i := 0; i < len(query); i++ {
		if query[i] == c {
			chi |= 1 << uint64(i)
		}
	}
	return chi
}

type NFAState struct {
	Offset      uint32
	Distance    uint8
	InTranspose bool
}

type NFAStates []NFAState

func (ns NFAStates) Len() int {
	return len(ns)
}

func (ns NFAStates) Less(i, j int) bool {
	if ns[i].Offset != ns[j].Offset {
		return ns[i].Offset < ns[j].Offset
	}

	if ns[i].Distance != ns[j].Distance {
		return ns[i].Distance < ns[j].Distance
	}

	return !ns[i].InTranspose && ns[j].InTranspose
}

func (ns NFAStates) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns *NFAState) imply(other NFAState) bool {
	transposeImply := ns.InTranspose
	if !other.InTranspose {
		transposeImply = !other.InTranspose
	}

	deltaOffset := ns.Offset - other.Offset
	if ns.Offset < other.Offset {
		deltaOffset = other.Offset - ns.Offset
	}

	if transposeImply {
		return uint32(other.Distance) >= (uint32(ns.Distance) + deltaOffset)
	}

	return uint32(other.Distance) > (uint32(ns.Distance) + deltaOffset)
}

type MultiState struct {
	states []NFAState
}

func (ms *MultiState) States() []NFAState {
	return ms.states
}

func (ms *MultiState) Clear() {
	ms.states = ms.states[:0]
}

func newMultiState() *MultiState {
	return &MultiState{states: make([]NFAState, 0)}
}

func (ms *MultiState) normalize() uint32 {
	minOffset := uint32(math.MaxUint32)

	for _, s := range ms.states {
		if s.Offset < minOffset {
			minOffset = s.Offset
		}
	}
	if minOffset == uint32(math.MaxUint32) {
		minOffset = 0
	}

	for i := 0; i < len(ms.states); i++ {
		ms.states[i].Offset -= minOffset
	}

	sort.Sort(NFAStates(ms.states))

	return minOffset
}

func (ms *MultiState) addStates(nState NFAState) {

	for _, s := range ms.states {
		if s.imply(nState) {
			return
		}
	}

	i := 0
	for i < len(ms.states) {
		if nState.imply(ms.states[i]) {
			ms.states = append(ms.states[:i], ms.states[i+1:]...)
		} else {
			i++
		}
	}
	ms.states = append(ms.states, nState)

}

func extractBit(bitset uint64, pos uint8) bool {
	shift := bitset >> pos
	bit := shift & 1
	return bit == uint64(1)
}

func dist(left, right uint32) uint32 {
	if left > right {
		return left - right
	}
	return right - left
}

type LevenshteinNFA struct {
	mDistance uint8
	damerau   bool // 是否支持转置
}

func newLevenshtein(maxD uint8, transposition bool) *LevenshteinNFA {
	return &LevenshteinNFA{mDistance: maxD,
		damerau: transposition,
	}
}

func (la *LevenshteinNFA) maxDistance() uint8 {
	return la.mDistance
}

func (la *LevenshteinNFA) msDiameter() uint8 {
	return 2*la.mDistance + 1
}

func (la *LevenshteinNFA) initialStates() *MultiState {
	ms := MultiState{}
	nfaState := NFAState{}
	ms.addStates(nfaState)
	return &ms
}

func (la *LevenshteinNFA) multistateDistance(ms *MultiState,
	queryLen uint32) Distance {
	minDistance := Atleast{d: la.mDistance + 1}
	for _, s := range ms.states {
		t := s.Distance + uint8(dist(queryLen, s.Offset))
		if t <= uint8(la.mDistance) {
			if minDistance.distance() > t {
				minDistance.d = t
			}
		}
	}

	if minDistance.distance() == la.mDistance+1 {
		return Atleast{d: la.mDistance + 1}
	}

	return minDistance
}

func (la *LevenshteinNFA) simpleTransition(state NFAState,
	symbol uint64, ms *MultiState) {

	if state.Distance < la.mDistance {
		// insertion
		ms.addStates(NFAState{Offset: state.Offset,
			Distance:    state.Distance + 1,
			InTranspose: false})

		// substitution
		ms.addStates(NFAState{Offset: state.Offset + 1,
			Distance:    state.Distance + 1,
			InTranspose: false})

		n := la.mDistance + 1 - state.Distance
		for d := uint8(1); d < n; d++ {
			if extractBit(symbol, d) {
				//  for d > 0, as many deletion and character match
				ms.addStates(NFAState{Offset: state.Offset + 1 + uint32(d),
					Distance:    state.Distance + d,
					InTranspose: false})
			}
		}

		if la.damerau && extractBit(symbol, 1) {
			ms.addStates(NFAState{
				Offset:      state.Offset,
				Distance:    state.Distance + 1,
				InTranspose: true})
		}

	}

	if extractBit(symbol, 0) {
		ms.addStates(NFAState{Offset: state.Offset + 1,
			Distance:    state.Distance,
			InTranspose: false})
	}

	if state.InTranspose && extractBit(symbol, 0) {
		ms.addStates(NFAState{Offset: state.Offset + 2,
			Distance:    state.Distance,
			InTranspose: false})
	}

}

func (la *LevenshteinNFA) transition(cState *MultiState,
	dState *MultiState, scv uint64) {
	dState.Clear()
	mask := (uint64(1) << la.msDiameter()) - uint64(1)

	for _, state := range cState.states {
		cv := (scv >> state.Offset) & mask
		la.simpleTransition(state, cv, dState)
	}

	sort.Sort(NFAStates(dState.states))
}

func (la *LevenshteinNFA) computeDistance(query, other []rune) Distance {
	cState := la.initialStates()
	nState := newMultiState()

	for _, i := range other {
		nState.Clear()
		chi := characteristicVector(query, i)
		la.transition(cState, nState, chi)
		cState, nState = nState, cState
	}

	return la.multistateDistance(cState, uint32(len(query)))
}

// #endregion
// #region dfa
const SinkState = uint32(0)

type DFA struct {
	transitions [][256]uint32
	distances   []Distance
	initState   int
	ed          uint8
}

// Returns the initial state
func (d *DFA) initialState() int {
	return d.initState
}

// Returns the Levenshtein distance associated to the
// current state.
func (d *DFA) distance(stateId int) Distance {
	return d.distances[stateId]
}

func (d *DFA) EditDistance(stateId int) uint8 {
	return d.distances[stateId].distance()
}

func (d *DFA) MatchAndDistance(input string) (bool, uint8) {
	currentState := d.Start()
	index := 0
	// Traverse the DFA while characters can still match
	for d.CanMatch(currentState) && index < len(input) {
		currentState = d.Accept(currentState, input[index])
		if currentState == int(SinkState) {
			break
		}
		index++
	}
	// Ensure we've processed the entire input and check if the current state is a match
	if index == len(input) && d.IsMatch(currentState) {
		return true, d.EditDistance(currentState)
	}
	return false, 0
}

// Returns the number of states in the `DFA`.
func (d *DFA) numStates() int {
	return len(d.transitions)
}

// Returns the destination state reached after consuming a given byte.
func (d *DFA) transition(fromState int, b uint8) int {
	return int(d.transitions[fromState][b])
}

func (d *DFA) eval(bytes []uint8) Distance {
	state := d.initialState()

	for _, b := range bytes {
		state = d.transition(state, b)
	}

	return d.distance(state)
}

func (d *DFA) Start() int {
	return int(d.initialState())
}

func (d *DFA) IsMatch(state int) bool {
	if _, ok := d.distance(state).(Exact); ok {
		return true
	}
	return false
}

func (d *DFA) CanMatch(state int) bool {
	return state > 0 && state < d.numStates()
}

func (d *DFA) Accept(state int, b byte) int {
	return int(d.transition(state, b))
}

// WillAlwaysMatch returns if the specified state will always end in a
// matching state.
func (d *DFA) WillAlwaysMatch(state int) bool {
	return false
}

func fill(dest []uint32, val uint32) {
	for i := range dest {
		dest[i] = val
	}
}

func fillTransitions(dest *[256]uint32, val uint32) {
	for i := range dest {
		dest[i] = val
	}
}

type Utf8DFAStateBuilder struct {
	dfaBuilder       *Utf8DFABuilder
	stateID          uint32
	defaultSuccessor []uint32
}

func (sb *Utf8DFAStateBuilder) addTransitionID(fromStateID uint32, b uint8,
	toStateID uint32) {
	sb.dfaBuilder.transitions[fromStateID][b] = toStateID
}

func (sb *Utf8DFAStateBuilder) addTransition(in rune, toStateID uint32) {
	fromStateID := sb.stateID
	chars := []byte(string(in))
	lastByte := chars[len(chars)-1]

	for i, ch := range chars[:len(chars)-1] {
		remNumBytes := len(chars) - i - 1
		defaultSuccessor := sb.defaultSuccessor[remNumBytes]
		intermediateStateID := sb.dfaBuilder.transitions[fromStateID][ch]

		if intermediateStateID == defaultSuccessor {
			intermediateStateID = sb.dfaBuilder.allocate()
			fillTransitions(&sb.dfaBuilder.transitions[intermediateStateID],
				sb.defaultSuccessor[remNumBytes-1])
		}

		sb.addTransitionID(fromStateID, ch, intermediateStateID)
		fromStateID = intermediateStateID
	}

	toStateIDDecoded := sb.dfaBuilder.getOrAllocate(original(toStateID))
	sb.addTransitionID(fromStateID, lastByte, toStateIDDecoded)
}

type Utf8StateId uint32

func original(stateId uint32) Utf8StateId {
	return predecessor(stateId, 0)
}

func predecessor(stateId uint32, numSteps uint8) Utf8StateId {
	return Utf8StateId(stateId*4 + uint32(numSteps))
}

// Utf8DFABuilder makes it possible to define a DFA
// that takes unicode character, and build a `DFA`
// that operates on utf-8 encoded
type Utf8DFABuilder struct {
	index        []uint32
	distances    []Distance
	transitions  [][256]uint32
	initialState uint32
	numStates    uint32
	maxNumStates uint32
}

func withMaxStates(maxStates uint32) *Utf8DFABuilder {
	rv := &Utf8DFABuilder{
		index:        make([]uint32, maxStates*2+100),
		distances:    make([]Distance, 0, maxStates),
		transitions:  make([][256]uint32, 0, maxStates),
		maxNumStates: maxStates,
	}

	for i := range rv.index {
		rv.index[i] = math.MaxUint32
	}

	return rv
}

func (dfab *Utf8DFABuilder) allocate() uint32 {
	newState := dfab.numStates
	dfab.numStates++

	dfab.distances = append(dfab.distances, Atleast{d: 255})
	dfab.transitions = append(dfab.transitions, [256]uint32{})

	return newState
}

func (dfab *Utf8DFABuilder) getOrAllocate(state Utf8StateId) uint32 {
	if int(state) >= cap(dfab.index) {
		cloneIndex := make([]uint32, int(state)*2)
		copy(cloneIndex, dfab.index)
		dfab.index = cloneIndex
	}
	if dfab.index[state] != math.MaxUint32 {
		return dfab.index[state]
	}

	nstate := dfab.allocate()
	dfab.index[state] = nstate

	return nstate
}

func (dfab *Utf8DFABuilder) setInitialState(iState uint32) {
	decodedID := dfab.getOrAllocate(original(iState))
	dfab.initialState = decodedID
}

func (dfab *Utf8DFABuilder) build(ed uint8) *DFA {
	return &DFA{
		transitions: dfab.transitions,
		distances:   dfab.distances,
		initState:   int(dfab.initialState),
		ed:          ed,
	}
}

func (dfab *Utf8DFABuilder) addState(state, default_suc_orig uint32,
	distance Distance) (*Utf8DFAStateBuilder, error) {
	if state > dfab.maxNumStates {
		return nil, fmt.Errorf("State id is larger than maxNumStates")
	}

	stateID := dfab.getOrAllocate(original(state))
	dfab.distances[stateID] = distance

	defaultSuccID := dfab.getOrAllocate(original(default_suc_orig))
	// creates a chain of states of predecessors of `default_suc_orig`.
	// Accepting k-bytes (whatever the bytes are) from `predecessor_states[k-1]`
	// leads to the `default_suc_orig` state.
	predecessorStates := []uint32{defaultSuccID,
		defaultSuccID,
		defaultSuccID,
		defaultSuccID}

	for numBytes := uint8(1); numBytes < 4; numBytes++ {
		predecessorState := predecessor(default_suc_orig, numBytes)
		predecessorStateID := dfab.getOrAllocate(predecessorState)
		predecessorStates[numBytes] = predecessorStateID
		succ := predecessorStates[numBytes-1]
		fillTransitions(&dfab.transitions[predecessorStateID], succ)
	}

	// 1-byte encoded chars.
	fill(dfab.transitions[stateID][0:192], predecessorStates[0])
	// 2-bytes encoded chars.
	fill(dfab.transitions[stateID][192:224], predecessorStates[1])
	// 3-bytes encoded chars.
	fill(dfab.transitions[stateID][224:240], predecessorStates[2])
	// 4-bytes encoded chars.
	fill(dfab.transitions[stateID][240:256], predecessorStates[3])

	return &Utf8DFAStateBuilder{
		dfaBuilder:       dfab,
		stateID:          stateID,
		defaultSuccessor: predecessorStates}, nil
}

// #endregion

// #region alphabet
type FullCharacteristicVector []uint32

func (fcv FullCharacteristicVector) shiftAndMask(offset, mask uint32) uint32 {
	bucketID := offset / 32
	align := offset - bucketID*32
	if align == 0 {
		return fcv[bucketID] & mask
	}
	left := fcv[bucketID] >> align
	right := fcv[bucketID+1] << (32 - align)
	return (left | right) & mask
}

type tuple struct {
	char rune
	fcv  FullCharacteristicVector
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func sortRune(r []rune) []rune {
	sort.Sort(sortRunes(r))
	return r
}

type Alphabet struct {
	charset []tuple
	index   uint32
}

func (a *Alphabet) resetNext() {
	a.index = 0
}

func (a *Alphabet) next() (rune, FullCharacteristicVector, error) {
	if int(a.index) >= len(a.charset) {
		return 0, nil, fmt.Errorf("eof")
	}

	rv := a.charset[a.index]
	a.index++
	return rv.char, rv.fcv, nil
}

func dedupe(in string) string {
	lookUp := make(map[rune]struct{}, len(in))
	var rv string
	for len(in) > 0 {
		r, size := utf8.DecodeRuneInString(in)
		in = in[size:]
		if _, ok := lookUp[r]; !ok {
			rv += string(r)
			lookUp[r] = struct{}{}
		}
	}
	return rv
}

func queryChars(qChars string) Alphabet {
	chars := dedupe(qChars)
	inChars := sortRune([]rune(chars))
	charsets := make([]tuple, 0, len(inChars))

	for _, c := range inChars {
		tempChars := qChars
		var bits []uint32
		for len(tempChars) > 0 {
			var chunk string
			if len(tempChars) > 32 {
				chunk = tempChars[0:32]
				tempChars = tempChars[32:]
			} else {
				chunk = tempChars
				tempChars = tempChars[:0]
			}

			chunkBits := uint32(0)
			bit := uint32(1)
			for _, chr := range chunk {
				if chr == c {
					chunkBits |= bit
				}
				bit <<= 1
			}
			bits = append(bits, chunkBits)
		}
		bits = append(bits, 0)
		charsets = append(charsets, tuple{char: c, fcv: FullCharacteristicVector(bits)})
	}
	return Alphabet{charset: charsets}
}

// #endregion
