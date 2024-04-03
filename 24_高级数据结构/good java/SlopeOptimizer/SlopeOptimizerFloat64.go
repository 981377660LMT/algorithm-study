package main

import (
	"fmt"
	"strings"
)

func main() {
	S := NewSlopeOptimizerFloat64(16, false)

	S.Add(1, 2, 0)
	S.Add(1, 2, 1)
	S.Add(2, 3, 2)
	S.Add(3, 4, 3)
	S.Add(6, 8, 4)
	S.Add(7, 9, 5)
	S.Add(8, 10, 6)
	S.Add(9, 11, 7)
	fmt.Println(S.GetBestChoice(1))
	fmt.Println(S.deque)
}

// 维护凸包，常用于斜率优化.
type SlopeOptimizerFloat64 struct {
	deque           *Deque
	upperConvexHull bool // 是否维护上凸包

}

func NewSlopeOptimizerFloat64(initCapacity int, upperConvexHull bool) *SlopeOptimizerFloat64 {
	initCapacity = max(16, initCapacity)
	return &SlopeOptimizerFloat64{
		deque:           NewDeque2(initCapacity),
		upperConvexHull: upperConvexHull,
	}
}

// id需要按照递增顺序加入.
func (so *SlopeOptimizerFloat64) Add(x, y float64, id int32) *Point {
	t1 := &Point{x: x, y: y, id: id}
	if so.upperConvexHull {
		for so.deque.Len() >= 2 {
			t2 := so.deque.Pop()
			t3 := so.deque.Back()
			if so._slope(t3, t2) > so._slope(t2, t1) {
				so.deque.Append(t2)
				break
			}
		}
	} else {
		for so.deque.Len() >= 2 {
			t2 := so.deque.Pop()
			t3 := so.deque.Back()
			if so._slope(t3, t2) < so._slope(t2, t1) {
				so.deque.Append(t2)
				break
			}
		}
	}
	so.deque.Append(t1)
	return t1
}

// 保留下凸包中id大于等于给定值的点.
func (so *SlopeOptimizerFloat64) Since(id int32) {
	for !so.deque.Empty() && so.deque.Front().id < id {
		so.deque.PopLeft()
	}
}

func (so *SlopeOptimizerFloat64) GetBestChoice(s float64) int32 {
	if so.upperConvexHull {
		for so.deque.Len() >= 2 {
			h1 := so.deque.PopLeft()
			h2 := so.deque.Front()
			if so._slope(h2, h1) < s {
				so.deque.AppendLeft(h1)
				break
			}
		}
	} else {
		for so.deque.Len() >= 2 {
			h1 := so.deque.PopLeft()
			h2 := so.deque.Front()
			if so._slope(h2, h1) > s {
				so.deque.AppendLeft(h1)
				break
			}
		}
	}
	return so.deque.Front().id
}

func (so *SlopeOptimizerFloat64) Empty() bool { return so.deque.Empty() }
func (so *SlopeOptimizerFloat64) Clear()      { so.deque.Clear() }
func (so *SlopeOptimizerFloat64) Len() int    { return so.deque.Len() }

func (so *SlopeOptimizerFloat64) _slope(a, b *Point) float64 {
	if b.x == a.x {
		if b.y == a.y {
			return 0
		} else if b.y > a.y {
			return 1e50
		} else {
			return 1e-50
		}
	}
	return (b.y - a.y) / (b.x - a.x)
}

type Point struct {
	x, y float64
	id   int32
}

func (p *Point) String() string {
	return fmt.Sprintf("{x: %f, y: %f, id: %d}", p.x, p.y, p.id)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type D = *Point
type Deque struct{ l, r []D }

func NewDeque2(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q *Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q *Deque) Len() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q *Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q *Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q *Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}

func (q *Deque) Clear() {
	q.l = q.l[:0]
	q.r = q.r[:0]
}

func (q *Deque) String() string {
	var sb strings.Builder
	sb.WriteString("Deque{")
	for i := 0; i < q.Len(); i++ {
		sb.WriteString(fmt.Sprintf("%v", q.At(i)))
		if i < q.Len()-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}
