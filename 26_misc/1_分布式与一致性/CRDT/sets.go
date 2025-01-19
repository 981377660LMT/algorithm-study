package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	{
		obj := "dummy-object"
		gset := NewGSet()
		gset.Add(obj)
		// Should always print 'true' as `obj` exists in the g-set.
		fmt.Println(gset.Contains(obj))
	}
}

type Set interface {
	Add(interface{})
	Contains(interface{}) bool
}

// #region gset

type mainSet map[interface{}]struct{}

var (
	// GSet should implement the set interface.
	_ Set = &GSet{}
)

// Gset is a grow-only set.
type GSet struct {
	mainSet mainSet
}

// NewGSet returns an instance of GSet.
func NewGSet() *GSet {
	return &GSet{
		mainSet: mainSet{},
	}
}

// Add lets you add an element to grow-only set.
func (g *GSet) Add(elem interface{}) {
	g.mainSet[elem] = struct{}{}
}

// Contains returns true if an element exists within the
// set or false otherwise.
func (g *GSet) Contains(elem interface{}) bool {
	_, ok := g.mainSet[elem]
	return ok
}

// Len returns the no. of elements present within GSet.
func (g *GSet) Len() int {
	return len(g.mainSet)
}

// Elems returns all the elements present in the set.
func (g *GSet) Elems() []interface{} {
	elems := make([]interface{}, 0, len(g.mainSet))

	for elem := range g.mainSet {
		elems = append(elems, elem)
	}

	return elems
}

type gsetJSON struct {
	T string        `json:"type"`
	E []interface{} `json:"e"`
}

// MarshalJSON will be used to generate a serialized output
// of a given GSet.
func (g *GSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&gsetJSON{
		T: "g-set",
		E: g.Elems(),
	})
}

// #endregion

// #region 2pset
var (
	// TwoPhaseSet should implement Set.
	_ Set = &TwoPhaseSet{}
)

// TwoPhaseSet supports both addition and removal of
// elements to set.
type TwoPhaseSet struct {
	addSet *GSet
	rmSet  *GSet
}

// NewTwoPhaseSet returns a new instance of TwoPhaseSet.
func NewTwoPhaseSet() *TwoPhaseSet {
	return &TwoPhaseSet{
		addSet: NewGSet(),
		rmSet:  NewGSet(),
	}
}

// Add inserts element into the TwoPhaseSet.
func (t *TwoPhaseSet) Add(elem interface{}) {
	t.addSet.Add(elem)
}

// Remove deletes the element from the set.
func (t *TwoPhaseSet) Remove(elem interface{}) {
	t.rmSet.Add(elem)
}

// Contains returns true if the set contains the element.
// The set is said to contain the element if it is present
// in the add-set and not in the remove-set.
func (t *TwoPhaseSet) Contains(elem interface{}) bool {
	return t.addSet.Contains(elem) && !t.rmSet.Contains(elem)
}

type tpsetJSON struct {
	T string        `json:"type"`
	A []interface{} `json:"a"`
	R []interface{} `json:"r"`
}

// MarshalJSON serializes TwoPhaseSet into an JSON array.
func (t *TwoPhaseSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&tpsetJSON{
		T: "2p-set",
		A: t.addSet.Elems(),
		R: t.rmSet.Elems(),
	})
}

// #endregion

// #region orset
type ORSet struct {
	addMap map[interface{}]map[string]struct{}
	rmMap  map[interface{}]map[string]struct{}
}

func NewORSet() *ORSet {
	return &ORSet{
		addMap: make(map[interface{}]map[string]struct{}),
		rmMap:  make(map[interface{}]map[string]struct{}),
	}
}

func (o *ORSet) Add(value interface{}) {
	if m, ok := o.addMap[value]; ok {
		m[NewV4().String()] = struct{}{}
		o.addMap[value] = m
		return
	}

	m := make(map[string]struct{})

	m[NewV4().String()] = struct{}{}
	o.addMap[value] = m
}

func (o *ORSet) Remove(value interface{}) {
	r, ok := o.rmMap[value]
	if !ok {
		r = make(map[string]struct{})
	}

	if m, ok := o.addMap[value]; ok {
		for uid, _ := range m {
			r[uid] = struct{}{}
		}
	}

	o.rmMap[value] = r
}

func (o *ORSet) Contains(value interface{}) bool {
	addMap, ok := o.addMap[value]
	if !ok {
		return false
	}

	rmMap, ok := o.rmMap[value]
	if !ok {
		return true
	}

	for uid, _ := range addMap {
		if _, ok := rmMap[uid]; !ok {
			return true
		}
	}

	return false
}

func (o *ORSet) Merge(r *ORSet) {
	for value, m := range r.addMap {
		addMap, ok := o.addMap[value]
		if ok {
			for uid, _ := range m {
				addMap[uid] = struct{}{}
			}

			continue
		}

		o.addMap[value] = m
	}

	for value, m := range r.rmMap {
		rmMap, ok := o.rmMap[value]
		if ok {
			for uid, _ := range m {
				rmMap[uid] = struct{}{}
			}

			continue
		}

		o.rmMap[value] = m
	}
}

// #endregion

// #region lww_e_set
// Last-write-wins element set

type LWWSet struct {
	addMap map[interface{}]time.Time
	rmMap  map[interface{}]time.Time

	bias BiasType

	clock Clock
}

type BiasType string

const (
	BiasAdd    BiasType = "a"
	BiasRemove BiasType = "r"
)

var (
	ErrNoSuchBias = errors.New("no such bias found")
)

func NewLWWSet() (*LWWSet, error) {
	return NewLWWSetWithBias(BiasAdd)
}

func NewLWWSetWithBias(bias BiasType) (*LWWSet, error) {
	if bias != BiasAdd && bias != BiasRemove {
		return nil, ErrNoSuchBias
	}

	return &LWWSet{
		addMap: make(map[interface{}]time.Time),
		rmMap:  make(map[interface{}]time.Time),
		bias:   bias,
		clock:  NewClock(),
	}, nil
}

func (s *LWWSet) Add(value interface{}) {
	s.addMap[value] = s.clock.Now()
}

func (s *LWWSet) Remove(value interface{}) {
	s.rmMap[value] = s.clock.Now()
}

func (s *LWWSet) Contains(value interface{}) bool {
	addTime, addOk := s.addMap[value]

	// If a value is not present in added set then
	// always return false, irrespective of whether
	// it is present in the removed set.
	if !addOk {
		return false
	}

	rmTime, rmOk := s.rmMap[value]

	// If a value is present in added set and not in remove
	// we should always return true.
	if !rmOk {
		return true
	}

	switch s.bias {
	case BiasAdd:
		return !addTime.Before(rmTime)

	case BiasRemove:
		return rmTime.Before(addTime)
	}

	// This case will almost always never be hit. Usually
	// if an invalid Bias value is provided, it is called
	// at a higher level.
	return false
}

func (s *LWWSet) Merge(r *LWWSet) {
	for value, ts := range r.addMap {
		if t, ok := s.addMap[value]; ok && t.Before(ts) {
			s.addMap[value] = ts
		} else {
			if t.Before(ts) {
				s.addMap[value] = ts
			} else {
				s.addMap[value] = t
			}
		}
	}

	for value, ts := range r.rmMap {
		if t, ok := s.rmMap[value]; ok && t.Before(ts) {
			s.rmMap[value] = ts
		} else {
			if t.Before(ts) {
				s.rmMap[value] = ts
			} else {
				s.rmMap[value] = t
			}
		}
	}
}

type lwwesetJSON struct {
	T string   `json:"type"`
	B string   `json:"bias"`
	E []elJSON `json:"e"`
}

type elJSON struct {
	Elem interface{} `json:"el"`
	TAdd int64       `json:"ta,omitempty"`
	TDel int64       `json:"td,omitempty"`
}

func (s *LWWSet) MarshalJSON() ([]byte, error) {
	l := &lwwesetJSON{
		T: "lww-e-set",
		B: string(s.bias),
		E: make([]elJSON, 0, len(s.addMap)),
	}

	for e, t := range s.addMap {
		el := elJSON{Elem: e, TAdd: t.Unix()}
		if td, ok := s.rmMap[e]; ok {
			el.TDel = td.Unix()
		}

		l.E = append(l.E, el)
	}

	for e, t := range s.rmMap {
		if _, ok := s.addMap[e]; ok {
			continue
		}

		l.E = append(l.E, elJSON{Elem: e, TDel: t.Unix()})
	}

	return json.Marshal(l)
}

// #endregion

// #region clock
// Re-export of time.Duration
type Duration = time.Duration

// Clock represents an interface to the functions in the standard library time
// package. Two implementations are available in the clock package. The first
// is a real-time clock which simply wraps the time package's functions. The
// second is a mock clock which will only change when
// programmatically adjusted.
type Clock interface {
	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) *Timer
	Now() time.Time
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	Ticker(d time.Duration) *Ticker
	Timer(d time.Duration) *Timer
	WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc)
	WithTimeout(parent context.Context, t time.Duration) (context.Context, context.CancelFunc)
}

// NewClock returns an instance of a real-time clock.
func NewClock() Clock {
	return &clock{}
}

// clock implements a real-time clock by simply wrapping the time package functions.
type clock struct{}

func (c *clock) After(d time.Duration) <-chan time.Time { return time.After(d) }

func (c *clock) AfterFunc(d time.Duration, f func()) *Timer {
	return &Timer{timer: time.AfterFunc(d, f)}
}

func (c *clock) Now() time.Time { return time.Now() }

func (c *clock) Since(t time.Time) time.Duration { return time.Since(t) }

func (c *clock) Until(t time.Time) time.Duration { return time.Until(t) }

func (c *clock) Sleep(d time.Duration) { time.Sleep(d) }

func (c *clock) Tick(d time.Duration) <-chan time.Time { return time.Tick(d) }

func (c *clock) Ticker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{C: t.C, ticker: t}
}

func (c *clock) Timer(d time.Duration) *Timer {
	t := time.NewTimer(d)
	return &Timer{C: t.C, timer: t}
}

func (c *clock) WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, d)
}

func (c *clock) WithTimeout(parent context.Context, t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, t)
}

// Mock represents a mock clock that only moves forward programmically.
// It can be preferable to a real-time clock when testing time-based functionality.
type Mock struct {
	// mu protects all other fields in this struct, and the data that they
	// point to.
	mu sync.Mutex

	now    time.Time   // current time
	timers clockTimers // tickers & timers
}

// NewMock returns an instance of a mock clock.
// The current time of the mock clock on initialization is the Unix epoch.
func NewMock() *Mock {
	return &Mock{now: time.Unix(0, 0)}
}

// Add moves the current time of the mock clock forward by the specified duration.
// This should only be called from a single goroutine at a time.
func (m *Mock) Add(d time.Duration) {
	// Calculate the final current time.
	m.mu.Lock()
	t := m.now.Add(d)
	m.mu.Unlock()

	// Continue to execute timers until there are no more before the new time.
	for {
		if !m.runNextTimer(t) {
			break
		}
	}

	// Ensure that we end with the new time.
	m.mu.Lock()
	m.now = t
	m.mu.Unlock()

	// Give a small buffer to make sure that other goroutines get handled.
	gosched()
}

// Set sets the current time of the mock clock to a specific one.
// This should only be called from a single goroutine at a time.
func (m *Mock) Set(t time.Time) {
	// Continue to execute timers until there are no more before the new time.
	for {
		if !m.runNextTimer(t) {
			break
		}
	}

	// Ensure that we end with the new time.
	m.mu.Lock()
	m.now = t
	m.mu.Unlock()

	// Give a small buffer to make sure that other goroutines get handled.
	gosched()
}

// WaitForAllTimers sets the clock until all timers are expired
func (m *Mock) WaitForAllTimers() time.Time {
	// Continue to execute timers until there are no more
	for {
		m.mu.Lock()
		if len(m.timers) == 0 {
			m.mu.Unlock()
			return m.Now()
		}

		sort.Sort(m.timers)
		next := m.timers[len(m.timers)-1].Next()
		m.mu.Unlock()
		m.Set(next)
	}
}

// runNextTimer executes the next timer in chronological order and moves the
// current time to the timer's next tick time. The next time is not executed if
// its next time is after the max time. Returns true if a timer was executed.
func (m *Mock) runNextTimer(max time.Time) bool {
	m.mu.Lock()

	// Sort timers by time.
	sort.Sort(m.timers)

	// If we have no more timers then exit.
	if len(m.timers) == 0 {
		m.mu.Unlock()
		return false
	}

	// Retrieve next timer. Exit if next tick is after new time.
	t := m.timers[0]
	if t.Next().After(max) {
		m.mu.Unlock()
		return false
	}

	// Move "now" forward and unlock clock.
	m.now = t.Next()
	now := m.now
	m.mu.Unlock()

	// Execute timer.
	t.Tick(now)
	return true
}

// After waits for the duration to elapse and then sends the current time on the returned channel.
func (m *Mock) After(d time.Duration) <-chan time.Time {
	return m.Timer(d).C
}

// AfterFunc waits for the duration to elapse and then executes a function in its own goroutine.
// A Timer is returned that can be stopped.
func (m *Mock) AfterFunc(d time.Duration, f func()) *Timer {
	m.mu.Lock()
	defer m.mu.Unlock()
	ch := make(chan time.Time, 1)
	t := &Timer{
		c:       ch,
		fn:      f,
		mock:    m,
		next:    m.now.Add(d),
		stopped: false,
	}
	m.timers = append(m.timers, (*internalTimer)(t))
	return t
}

// Now returns the current wall time on the mock clock.
func (m *Mock) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.now
}

// Since returns time since `t` using the mock clock's wall time.
func (m *Mock) Since(t time.Time) time.Duration {
	return m.Now().Sub(t)
}

// Until returns time until `t` using the mock clock's wall time.
func (m *Mock) Until(t time.Time) time.Duration {
	return t.Sub(m.Now())
}

// Sleep pauses the goroutine for the given duration on the mock clock.
// The clock must be moved forward in a separate goroutine.
func (m *Mock) Sleep(d time.Duration) {
	<-m.After(d)
}

// Tick is a convenience function for Ticker().
// It will return a ticker channel that cannot be stopped.
func (m *Mock) Tick(d time.Duration) <-chan time.Time {
	return m.Ticker(d).C
}

// Ticker creates a new instance of Ticker.
func (m *Mock) Ticker(d time.Duration) *Ticker {
	m.mu.Lock()
	defer m.mu.Unlock()
	ch := make(chan time.Time, 1)
	t := &Ticker{
		C:    ch,
		c:    ch,
		mock: m,
		d:    d,
		next: m.now.Add(d),
	}
	m.timers = append(m.timers, (*internalTicker)(t))
	return t
}

// Timer creates a new instance of Timer.
func (m *Mock) Timer(d time.Duration) *Timer {
	m.mu.Lock()
	ch := make(chan time.Time, 1)
	t := &Timer{
		C:       ch,
		c:       ch,
		mock:    m,
		next:    m.now.Add(d),
		stopped: false,
	}
	m.timers = append(m.timers, (*internalTimer)(t))
	now := m.now
	m.mu.Unlock()
	m.runNextTimer(now)
	return t
}

// removeClockTimer removes a timer from m.timers. m.mu MUST be held
// when this method is called.
func (m *Mock) removeClockTimer(t clockTimer) {
	for i, timer := range m.timers {
		if timer == t {
			copy(m.timers[i:], m.timers[i+1:])
			m.timers[len(m.timers)-1] = nil
			m.timers = m.timers[:len(m.timers)-1]
			break
		}
	}
	sort.Sort(m.timers)
}

// clockTimer represents an object with an associated start time.
type clockTimer interface {
	Next() time.Time
	Tick(time.Time)
}

// clockTimers represents a list of sortable timers.
type clockTimers []clockTimer

func (a clockTimers) Len() int           { return len(a) }
func (a clockTimers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a clockTimers) Less(i, j int) bool { return a[i].Next().Before(a[j].Next()) }

// Timer represents a single event.
// The current time will be sent on C, unless the timer was created by AfterFunc.
type Timer struct {
	C       <-chan time.Time
	c       chan time.Time
	timer   *time.Timer // realtime impl, if set
	next    time.Time   // next tick time
	mock    *Mock       // mock clock, if set
	fn      func()      // AfterFunc function, if set
	stopped bool        // True if stopped, false if running
}

// Stop turns off the ticker.
func (t *Timer) Stop() bool {
	if t.timer != nil {
		return t.timer.Stop()
	}

	t.mock.mu.Lock()
	registered := !t.stopped
	t.mock.removeClockTimer((*internalTimer)(t))
	t.stopped = true
	t.mock.mu.Unlock()
	return registered
}

// Reset changes the expiry time of the timer
func (t *Timer) Reset(d time.Duration) bool {
	if t.timer != nil {
		return t.timer.Reset(d)
	}

	t.mock.mu.Lock()
	t.next = t.mock.now.Add(d)
	defer t.mock.mu.Unlock()

	registered := !t.stopped
	if t.stopped {
		t.mock.timers = append(t.mock.timers, (*internalTimer)(t))
	}

	t.stopped = false
	return registered
}

type internalTimer Timer

func (t *internalTimer) Next() time.Time { return t.next }
func (t *internalTimer) Tick(now time.Time) {
	// a gosched() after ticking, to allow any consequences of the
	// tick to complete
	defer gosched()

	t.mock.mu.Lock()
	if t.fn != nil {
		// defer function execution until the lock is released, and
		defer func() { go t.fn() }()
	} else {
		t.c <- now
	}
	t.mock.removeClockTimer((*internalTimer)(t))
	t.stopped = true
	t.mock.mu.Unlock()
}

// Ticker holds a channel that receives "ticks" at regular intervals.
type Ticker struct {
	C       <-chan time.Time
	c       chan time.Time
	ticker  *time.Ticker  // realtime impl, if set
	next    time.Time     // next tick time
	mock    *Mock         // mock clock, if set
	d       time.Duration // time between ticks
	stopped bool          // True if stopped, false if running
}

// Stop turns off the ticker.
func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
	} else {
		t.mock.mu.Lock()
		t.mock.removeClockTimer((*internalTicker)(t))
		t.stopped = true
		t.mock.mu.Unlock()
	}
}

// Reset resets the ticker to a new duration.
func (t *Ticker) Reset(dur time.Duration) {
	if t.ticker != nil {
		t.ticker.Reset(dur)
		return
	}

	t.mock.mu.Lock()
	defer t.mock.mu.Unlock()

	if t.stopped {
		t.mock.timers = append(t.mock.timers, (*internalTicker)(t))
		t.stopped = false
	}

	t.d = dur
	t.next = t.mock.now.Add(dur)
}

type internalTicker Ticker

func (t *internalTicker) Next() time.Time { return t.next }
func (t *internalTicker) Tick(now time.Time) {
	select {
	case t.c <- now:
	default:
	}
	t.mock.mu.Lock()
	t.next = now.Add(t.d)
	t.mock.mu.Unlock()
	gosched()
}

// Sleep momentarily so that other goroutines can process.
func gosched() { time.Sleep(1 * time.Millisecond) }

var (
	// type checking
	_ Clock = &Mock{}
)

func (m *Mock) WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return m.WithDeadline(parent, m.Now().Add(timeout))
}

func (m *Mock) WithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	if cur, ok := parent.Deadline(); ok && cur.Before(deadline) {
		// The current deadline is already sooner than the new one.
		return context.WithCancel(parent)
	}
	ctx := &timerCtx{clock: m, parent: parent, deadline: deadline, done: make(chan struct{})}
	propagateCancel(parent, ctx)
	dur := m.Until(deadline)
	if dur <= 0 {
		ctx.cancel(context.DeadlineExceeded) // deadline has already passed
		return ctx, func() {}
	}
	ctx.Lock()
	defer ctx.Unlock()
	if ctx.err == nil {
		ctx.timer = m.AfterFunc(dur, func() {
			ctx.cancel(context.DeadlineExceeded)
		})
	}
	return ctx, func() { ctx.cancel(context.Canceled) }
}

// propagateCancel arranges for child to be canceled when parent is.
func propagateCancel(parent context.Context, child *timerCtx) {
	if parent.Done() == nil {
		return // parent is never canceled
	}
	go func() {
		select {
		case <-parent.Done():
			child.cancel(parent.Err())
		case <-child.Done():
		}
	}()
}

type timerCtx struct {
	sync.Mutex

	clock    Clock
	parent   context.Context
	deadline time.Time
	done     chan struct{}

	err   error
	timer *Timer
}

func (c *timerCtx) cancel(err error) {
	c.Lock()
	defer c.Unlock()
	if c.err != nil {
		return // already canceled
	}
	c.err = err
	close(c.done)
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) { return c.deadline, true }

func (c *timerCtx) Done() <-chan struct{} { return c.done }

func (c *timerCtx) Err() error { return c.err }

func (c *timerCtx) Value(key interface{}) interface{} { return c.parent.Value(key) }

func (c *timerCtx) String() string {
	return fmt.Sprintf("clock.WithDeadline(%s [%s])", c.deadline, c.deadline.Sub(c.clock.Now()))
}

// #endregion

// #region uuid

// #region uuid
// Size of a UUID in bytes.
const Size = 16

// UUID representation compliant with specification
// described in RFC 4122.
type UUID [Size]byte

// UUID versions
const (
	_ byte = iota
	V1
	V2
	V3
	V4
	V5
)

// UUID layout variants.
const (
	VariantNCS byte = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// UUID DCE domains.
const (
	DomainPerson = iota
	DomainGroup
	DomainOrg
)

// String parse helpers.
var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
)

// Nil is special form of UUID that is specified to have all
// 128 bits set to zero.
var Nil = UUID{}

// Predefined namespace UUIDs.
var (
	NamespaceDNS  = Must(FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceURL  = Must(FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceOID  = Must(FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceX500 = Must(FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
)

// Equal returns true if u1 and u2 equals, otherwise returns false.
func Equal(u1 UUID, u2 UUID) bool {
	return bytes.Equal(u1[:], u2[:])
}

// Version returns algorithm version used to generate UUID.
func (u UUID) Version() byte {
	return u[6] >> 4
}

// Variant returns UUID layout variant.
func (u UUID) Variant() byte {
	switch {
	case (u[8] >> 7) == 0x00:
		return VariantNCS
	case (u[8] >> 6) == 0x02:
		return VariantRFC4122
	case (u[8] >> 5) == 0x06:
		return VariantMicrosoft
	case (u[8] >> 5) == 0x07:
		fallthrough
	default:
		return VariantFuture
	}
}

// Bytes returns bytes slice representation of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

// SetVersion sets version bits.
func (u *UUID) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits.
func (u *UUID) SetVariant(v byte) {
	switch v {
	case VariantNCS:
		u[8] = (u[8]&(0xff>>1) | (0x00 << 7))
	case VariantRFC4122:
		u[8] = (u[8]&(0xff>>2) | (0x02 << 6))
	case VariantMicrosoft:
		u[8] = (u[8]&(0xff>>3) | (0x06 << 5))
	case VariantFuture:
		fallthrough
	default:
		u[8] = (u[8]&(0xff>>3) | (0x07 << 5))
	}
}

// Must is a helper that wraps a call to a function returning (UUID, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//
//	var packageUUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426655440000"));
func Must(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}

// FromString returns UUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

// FromBytes returns UUID converted from raw byte slice input.
// It will return error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (u UUID, err error) {
	err = u.UnmarshalBinary(input)
	return
}

// FromBytesOrNil returns UUID converted from raw byte slice input.
// Same behavior as FromBytes, but returns a Nil UUID on error.
func FromBytesOrNil(input []byte) UUID {
	uuid, err := FromBytes(input)
	if err != nil {
		return Nil
	}
	return uuid
}

// FromStringOrNil returns UUID parsed from string input.
// Same behavior as FromString, but returns a Nil UUID on error.
func FromStringOrNil(input string) UUID {
	uuid, err := FromString(input)
	if err != nil {
		return Nil
	}
	return uuid
}

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (u UUID) MarshalText() (text []byte, err error) {
	text = []byte(u.String())
	return
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//
//	"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//	"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//	"urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//	"6ba7b8109dad11d180b400c04fd430c8"
//
// ABNF for supported UUID text representation follows:
//
//	uuid := canonical | hashlike | braced | urn
//	plain := canonical | hashlike
//	canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
//	hashlike := 12hexoct
//	braced := '{' plain '}'
//	urn := URN ':' UUID-NID ':' plain
//	URN := 'urn'
//	UUID-NID := 'uuid'
//	12hexoct := 6hexoct 6hexoct
//	6hexoct := 4hexoct 2hexoct
//	4hexoct := 2hexoct 2hexoct
//	2hexoct := hexoct hexoct
//	hexoct := hexdig hexdig
//	hexdig := '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' |
//	          'a' | 'b' | 'c' | 'd' | 'e' | 'f' |
//	          'A' | 'B' | 'C' | 'D' | 'E' | 'F'
func (u *UUID) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	case 32:
		return u.decodeHashLike(text)
	case 36:
		return u.decodeCanonical(text)
	case 38:
		return u.decodeBraced(text)
	case 41:
		fallthrough
	case 45:
		return u.decodeURN(text)
	default:
		return fmt.Errorf("uuid: incorrect UUID length: %s", text)
	}
}

// decodeCanonical decodes UUID string in format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8".
func (u *UUID) decodeCanonical(t []byte) (err error) {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	src := t[:]
	dst := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err = hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	return
}

// decodeHashLike decodes UUID string in format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeHashLike(t []byte) (err error) {
	src := t[:]
	dst := u[:]

	if _, err = hex.Decode(dst, src); err != nil {
		return err
	}
	return
}

// decodeBraced decodes UUID string in format
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}" or in format
// "{6ba7b8109dad11d180b400c04fd430c8}".
func (u *UUID) decodeBraced(t []byte) (err error) {
	l := len(t)

	if t[0] != '{' || t[l-1] != '}' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	return u.decodePlain(t[1 : l-1])
}

// decodeURN decodes UUID string in format
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in format
// "urn:uuid:6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeURN(t []byte) (err error) {
	total := len(t)

	urn_uuid_prefix := t[:9]

	if !bytes.Equal(urn_uuid_prefix, urnPrefix) {
		return fmt.Errorf("uuid: incorrect UUID format: %s", t)
	}

	return u.decodePlain(t[9:total])
}

// decodePlain decodes UUID string in canonical format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in hash-like format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodePlain(t []byte) (err error) {
	switch len(t) {
	case 32:
		return u.decodeHashLike(t)
	case 36:
		return u.decodeCanonical(t)
	default:
		return fmt.Errorf("uuid: incorrrect UUID length: %s", t)
	}
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != Size {
		err = fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
		return
	}
	copy(u[:], data)

	return
}

// Difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and Unix epoch (January 1, 1970).
const epochStart = 122192928000000000

type epochFunc func() time.Time
type hwAddrFunc func() (net.HardwareAddr, error)

var (
	global = newRFC4122Generator()

	posixUID = uint32(os.Getuid())
	posixGID = uint32(os.Getgid())
)

// NewV1 returns UUID based on current timestamp and MAC address.
func NewV1() (UUID, error) {
	return global.NewV1()
}

// NewV3 returns UUID based on MD5 hash of namespace UUID and name.
func NewV3(ns UUID, name string) UUID {
	return global.NewV3(ns, name)
}

// NewV4 returns random generated UUID.
func NewV4() UUID {
	res, _ := global.NewV4()
	return res
}

// NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
func NewV5(ns UUID, name string) UUID {
	return global.NewV5(ns, name)
}

// Generator provides interface for generating UUIDs.
type Generator interface {
	NewV1() (UUID, error)
	NewV2(domain byte) (UUID, error)
	NewV3(ns UUID, name string) UUID
	NewV4() (UUID, error)
	NewV5(ns UUID, name string) UUID
}

// Default generator implementation.
type rfc4122Generator struct {
	clockSequenceOnce sync.Once
	hardwareAddrOnce  sync.Once
	storageMutex      sync.Mutex

	rand io.Reader

	epochFunc     epochFunc
	hwAddrFunc    hwAddrFunc
	lastTime      uint64
	clockSequence uint16
	hardwareAddr  [6]byte
}

func newRFC4122Generator() Generator {
	return &rfc4122Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}
}

// NewV1 returns UUID based on current timestamp and MAC address.
func (g *rfc4122Generator) NewV1() (UUID, error) {
	u := UUID{}

	timeNow, clockSeq, err := g.getClockSequence()
	if err != nil {
		return Nil, err
	}
	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	hardwareAddr, err := g.getHardwareAddr()
	if err != nil {
		return Nil, err
	}
	copy(u[10:], hardwareAddr)

	u.SetVersion(V1)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

// NewV2 returns DCE Security UUID based on POSIX UID/GID.
func (g *rfc4122Generator) NewV2(domain byte) (UUID, error) {
	u, err := g.NewV1()
	if err != nil {
		return Nil, err
	}

	switch domain {
	case DomainPerson:
		binary.BigEndian.PutUint32(u[:], posixUID)
	case DomainGroup:
		binary.BigEndian.PutUint32(u[:], posixGID)
	}

	u[9] = domain

	u.SetVersion(V2)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

// NewV3 returns UUID based on MD5 hash of namespace UUID and name.
func (g *rfc4122Generator) NewV3(ns UUID, name string) UUID {
	u := newFromHash(md5.New(), ns, name)
	u.SetVersion(V3)
	u.SetVariant(VariantRFC4122)

	return u
}

// NewV4 returns random generated UUID.
func (g *rfc4122Generator) NewV4() (UUID, error) {
	u := UUID{}
	if _, err := io.ReadFull(g.rand, u[:]); err != nil {
		return Nil, err
	}
	u.SetVersion(V4)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

// NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
func (g *rfc4122Generator) NewV5(ns UUID, name string) UUID {
	u := newFromHash(sha1.New(), ns, name)
	u.SetVersion(V5)
	u.SetVariant(VariantRFC4122)

	return u
}

// Returns epoch and clock sequence.
func (g *rfc4122Generator) getClockSequence() (uint64, uint16, error) {
	var err error
	g.clockSequenceOnce.Do(func() {
		buf := make([]byte, 2)
		if _, err = io.ReadFull(g.rand, buf); err != nil {
			return
		}
		g.clockSequence = binary.BigEndian.Uint16(buf)
	})
	if err != nil {
		return 0, 0, err
	}

	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	timeNow := g.getEpoch()
	// Clock didn't change since last UUID generation.
	// Should increase clock sequence.
	if timeNow <= g.lastTime {
		g.clockSequence++
	}
	g.lastTime = timeNow

	return timeNow, g.clockSequence, nil
}

// Returns hardware address.
func (g *rfc4122Generator) getHardwareAddr() ([]byte, error) {
	var err error
	g.hardwareAddrOnce.Do(func() {
		if hwAddr, err := g.hwAddrFunc(); err == nil {
			copy(g.hardwareAddr[:], hwAddr)
			return
		}

		// Initialize hardwareAddr randomly in case
		// of real network interfaces absence.
		if _, err = io.ReadFull(g.rand, g.hardwareAddr[:]); err != nil {
			return
		}
		// Set multicast bit as recommended by RFC 4122
		g.hardwareAddr[0] |= 0x01
	})
	if err != nil {
		return []byte{}, err
	}
	return g.hardwareAddr[:], nil
}

// Returns difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and current time.
func (g *rfc4122Generator) getEpoch() uint64 {
	return epochStart + uint64(g.epochFunc().UnixNano()/100)
}

// Returns UUID based on hashing of namespace UUID and name.
func newFromHash(h hash.Hash, ns UUID, name string) UUID {
	u := UUID{}
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))

	return u
}

// Returns hardware address.
func defaultHWAddrFunc() (net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return []byte{}, err
	}
	for _, iface := range ifaces {
		if len(iface.HardwareAddr) >= 6 {
			return iface.HardwareAddr, nil
		}
	}
	return []byte{}, fmt.Errorf("uuid: no HW address found")
}

// #endregion
