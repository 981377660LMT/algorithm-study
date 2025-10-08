// ref: https://nyaannyaan.github.io/library/game/impartial-game.hpp
// FIXME

package main

import (
	"fmt"
	"math/bits"
)

// --- 示例：一个简单的党派游戏 ---

// SimpleGameLogic 实现一个简单的游戏。
// 盘面是一个整数。
// Left 可以将 n 变为 n-2。
// Right 可以将 n 变为 n-1。
// 游戏在 n <= 0 时结束。
type SimpleGameLogic struct{}

func (sgl *SimpleGameLogic) NextStates(g int) (left []int, right []int) {
	// Left's move
	if g-2 > 0 {
		left = append(left, g-2)
	}
	// Right's move
	if g-1 > 0 {
		right = append(right, g-1)
	}
	return
}

func main() {
	logic := &SimpleGameLogic{}
	solver := NewPartisanGameSolver(logic)

	// 计算盘面为 5 时的值
	val := solver.Get(5)
	fmt.Printf("Value of game at state 5 is: %s\n", val) // 预期是 1/2

	val2 := solver.Get(6)
	fmt.Printf("Value of game at state 6 is: %s\n", val2) // 预期是 1
}

// #region PartisanGame

type GameState comparable

type GameLogic[G GameState] interface {
	// NextStates 返回 Left 和 Right 玩家的后继局面。
	NextStates(g G) (left []G, right []G)
}

type PartisanGameSolver[G GameState] struct {
	memo  map[G]SurrealNumber
	logic GameLogic[G]
}

func NewPartisanGameSolver[G GameState](logic GameLogic[G]) *PartisanGameSolver[G] {
	return &PartisanGameSolver[G]{
		memo:  make(map[G]SurrealNumber),
		logic: logic,
	}
}

// Get 计算一个盘面的超现实数值。
func (s *PartisanGameSolver[G]) Get(g G) SurrealNumber {
	if val, ok := s.memo[g]; ok {
		return val
	}
	val := s.calculate(g)
	s.memo[g] = val
	return val
}

func (s *PartisanGameSolver[G]) calculate(g G) SurrealNumber {
	gl, gr := s.logic.NextStates(g)

	if len(gl) == 0 && len(gr) == 0 {
		return NewSurrealNumber(0, 0)
	}

	var leftValues, rightValues []SurrealNumber
	for _, cg := range gl {
		leftValues = append(leftValues, s.Get(cg))
	}
	for _, cg := range gr {
		rightValues = append(rightValues, s.Get(cg))
	}

	var sl, sr SurrealNumber
	if len(leftValues) > 0 {
		sl = leftValues[0]
		for i := 1; i < len(leftValues); i++ {
			if sl.LessThan(leftValues[i]) {
				sl = leftValues[i]
			}
		}
	}
	if len(rightValues) > 0 {
		sr = rightValues[0]
		for i := 1; i < len(rightValues); i++ {
			if rightValues[i].LessThan(sr) {
				sr = rightValues[i]
			}
		}
	}

	if len(rightValues) == 0 {
		return sl.Larger()
	}
	if len(leftValues) == 0 {
		return sr.Smaller()
	}

	if !sl.LessThan(sr) {
		panic("game rule violation: max of left options is not less than min of right options")
	}

	return Reduce(sl, sr)
}

// #endregion

// #region SurrealNumber
// SurrealNumber 表示一个二进分数 p / 2^q。
type SurrealNumber struct {
	p int
	q int
}

func NewSurrealNumber(p int, q int) SurrealNumber {
	return SurrealNumber{p, q}
}

func (s SurrealNumber) Normalize() SurrealNumber {
	if s.p == 0 {
		return SurrealNumber{0, 0}
	}
	tz := bits.TrailingZeros(uint(s.p))
	if tz > s.q {
		tz = s.q
	}
	return SurrealNumber{s.p >> tz, s.q - tz}
}

func (s SurrealNumber) String() string {
	if s.q == 0 {
		return fmt.Sprintf("%d", s.p)
	}
	return fmt.Sprintf("%d / %d", s.p, int(1)<<s.q)
}

func (s SurrealNumber) Add(other SurrealNumber) SurrealNumber {
	cq := s.q
	if other.q > cq {
		cq = other.q
	}
	cp := (s.p << (cq - s.q)) + (other.p << (cq - other.q))
	return SurrealNumber{cp, cq}.Normalize()
}

func (s SurrealNumber) Sub(other SurrealNumber) SurrealNumber {
	cq := s.q
	if other.q > cq {
		cq = other.q
	}
	cp := (s.p << (cq - s.q)) - (other.p << (cq - other.q))
	return SurrealNumber{cp, cq}.Normalize()
}

func (s SurrealNumber) Neg() SurrealNumber {
	return SurrealNumber{-s.p, s.q}
}

func (s SurrealNumber) LessThan(other SurrealNumber) bool {
	return other.Sub(s).p > 0
}

func (s SurrealNumber) EqualTo(other SurrealNumber) bool {
	return s.Sub(other).p == 0
}

func (s SurrealNumber) Child() (SurrealNumber, SurrealNumber) {
	if s.p == 0 {
		return NewSurrealNumber(-1, 0), NewSurrealNumber(1, 0)
	}
	if s.q == 0 && s.p > 0 {
		return NewSurrealNumber(s.p, 0).Add(NewSurrealNumber(-1, 0)),
			NewSurrealNumber(s.p+1, 0)
	}
	if s.q == 0 && s.p < 0 {
		return NewSurrealNumber(s.p-1, 0),
			NewSurrealNumber(s.p, 0).Add(NewSurrealNumber(1, 1))
	}
	return s.Sub(NewSurrealNumber(1, s.q+1)), s.Add(NewSurrealNumber(1, s.q+1))
}

func (s SurrealNumber) Larger() SurrealNumber {
	root := NewSurrealNumber(0, 0)
	for !s.LessThan(root) {
		_, rr := root.Child()
		root = rr
	}
	return root
}

func (s SurrealNumber) Smaller() SurrealNumber {
	root := NewSurrealNumber(0, 0)
	for !root.LessThan(s) {
		lr, _ := root.Child()
		root = lr
	}
	return root
}

func Reduce(l, r SurrealNumber) SurrealNumber {
	if !l.LessThan(r) {
		panic("l must be less than r for Reduce")
	}
	root := NewSurrealNumber(0, 0)
	for !l.LessThan(root) || !root.LessThan(r) {
		lr, rr := root.Child()
		if !root.LessThan(r) {
			root = lr
		} else {
			root = rr
		}
	}
	return root
}

// #endregion
