// https://cdn.luogu.com.cn/upload/image_hosting/cn8i6rnu.png
// !对同一段区间，拆分出的区间个数不超过线段树的区间个数.
// eg: [1,8)，线段树拆分成[1,2),[2,4),[4,8)三个区间，而树状数组拆分出[0,8)和[0,1)两个区间.
// 树状数组做到了不同的子区间的总数恰好为n。
// !等价于没有右子树的线段树.
//
//	                          7
//	                       ↙ ↓
//	                     ↙   ↓
//			               ↙      ↓
//	                ↙        ↓
//								↙           ↓
//				       3             x
//	        ↙  ↓         ↙  ↓
//		     1     x       5     x
//		  ↙ ↓   ↙↓    ↙ ↓ 	↙↓
//		0    x  2  x   4   x  6  x
//
// 注意树状数组拆分出的区间不是最少的。最少的拆分参考`DivideIntervalAbel`。
package main

import (
	"fmt"
)

func main() {
	// demo()

}

func demo() {
	N := 20
	D := NewDivideIntervalBIT(N)

	{
		D.EnumerateSegment(1, 6, func(i int, sign bool) {
			fmt.Println(i, sign)
			fmt.Println(D.IdToSegment(i))
		})

		fmt.Println("---")
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			fmt.Println(D.IdToSegment(i))
			fmt.Println("p", D.Parent(i))
		}
	}

	{

		data := make([]int, N)
		for i := 0; i < N; i++ {
			data[i] = i
		}
		D.PushUp(func(parent, child int) {
			data[parent] += data[child]
		})

		getSum := func(start, end int) int {
			sum := 0
			D.EnumerateSegment(start, end, func(id int, sign bool) {
				if sign {
					sum += data[id]
				} else {
					sum -= data[id]
				}
			})
			return sum
		}
		fmt.Println(getSum(0, 10))

		D.PathToRoot(7, func(id int) {
			data[id] += 1
		})
		fmt.Println(getSum(2, 10))
		fmt.Println(getSum(2, 7))
		fmt.Println(getSum(2, 8))
	}
}

type DivideIntervalBIT struct {
	n int
}

// [0,n).
func NewDivideIntervalBIT(n int) *DivideIntervalBIT {
	return &DivideIntervalBIT{n: n}
}

// 树状数组拆分区间.运算满足阿贝尔群.
func (d *DivideIntervalBIT) EnumerateSegment(start, end int, f func(id int, sign bool)) {
	for end > start {
		f(end-1, true)
		end &= end - 1 // 减去lowbit
	}
	for start > end {
		f(start-1, false)
		start &= start - 1
	}
}

func (d *DivideIntervalBIT) IdToSegment(id int) (start, end int) {
	id++
	start = id & (id - 1)
	end = id
	return
}

func (d *DivideIntervalBIT) Parent(id int) int {
	id++
	return id | (id - 1) // 加上lowbit
}

// O(n) 从叶子向根方向push.
func (d *DivideIntervalBIT) PushUp(f func(parent, child int)) {
	for i := 0; i < d.n; i++ {
		if p := d.Parent(i); p < d.n {
			f(p, i)
		}
	}
}

// O(logn) 从叶子向根方向push.
func (d *DivideIntervalBIT) PathToRoot(start int, f func(id int)) {
	for i := start; i < d.n; i = d.Parent(i) {
		f(i)
	}
}

func (d *DivideIntervalBIT) Size() int {
	return d.n
}
