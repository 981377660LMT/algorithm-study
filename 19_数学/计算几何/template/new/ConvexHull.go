// x#pragma once

// #include "geo/base.hpp"

// template <typename T>
// vector<int> ConvexHull(vector<Point<T>>& XY, string mode = "full",
//                        bool sorted = false) {
//   assert(mode == "full" || mode == "lower" || mode == "upper");
//   ll N = XY.size();
//   if (N == 1) return {0};
//   if (N == 2) {
//     if (XY[0] < XY[1]) return {0, 1};
//     if (XY[1] < XY[0]) return {1, 0};
//     return {0};
//   }
//   vc<int> I(N);
//   if (sorted) {
//     FOR(i, N) I[i] = i;
//   } else {
//     I = argsort(XY);
//   }

//   auto check = [&](ll i, ll j, ll k) -> bool {
//     return (XY[j] - XY[i]).det(XY[k] - XY[i]) > 0;
//   };

//   auto calc = [&]() {
//     vector<int> P;
//     for (auto&& k: I) {
//       while (P.size() > 1) {
//         auto i = P[P.size() - 2];
//         auto j = P[P.size() - 1];
//         if (check(i, j, k)) break;
//         P.pop_back();
//       }
//       P.eb(k);
//     }
//     return P;
//   };

//   vc<int> P;
//   if (mode == "full" || mode == "lower") {
//     vc<int> Q = calc();
//     P.insert(P.end(), all(Q));
//   }
//   if (mode == "full" || mode == "upper") {
//     if (!P.empty()) P.pop_back();
//     reverse(all(I));
//     vc<int> Q = calc();
//     P.insert(P.end(), all(Q));
//   }
//   if (mode == "upper") reverse(all(P));
//   while (len(P) >= 2 && XY[P[0]] == XY[P.back()]) P.pop_back();
//   return P;
// }

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/furthest_pair
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func() {
		var n int32
		fmt.Fscan(in, &n)
		points := make([][2]int, n)
		for i := int32(0); i < n; i++ {
			fmt.Fscan(in, &points[i][0], &points[i][1])
		}
		i, j := FurthestPair(points)
		fmt.Fprintln(out, i, j)
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		solve()
	}
}

type Mode uint8

const (
	Full Mode = iota
	Lower
	Upper
)

// (凸包/上凸包/下凸包).
func ConvexHull(points [][2]int, mode Mode, isPointsSorted bool) []int32 {
	n := len(points)
	if n == 1 {
		return []int32{0}
	}

	compare := func(i, j int32) int8 {
		x1, y1 := points[i][0], points[i][1]
		x2, y2 := points[j][0], points[j][1]
		if x1 < x2 || (x1 == x2 && y1 < y2) {
			return -1
		}
		if x1 == x2 && y1 == y2 {
			return 0
		}
		return 1
	}

	if n == 2 {
		res := compare(0, 1)
		if res == -1 {
			return []int32{0, 1}
		}
		if res == 1 {
			return []int32{1, 0}
		}
		return []int32{0}
	}

	order := make([]int32, n)
	for i := int32(0); i < int32(n); i++ {
		order[i] = i
	}
	if !isPointsSorted {
		sort.Slice(order, func(i, j int) bool { return compare(order[i], order[j]) == -1 })
	}

	check := func(i, j, k int32) bool {
		x1, y1 := points[j][0]-points[i][0], points[j][1]-points[i][1]
		x2, y2 := points[k][0]-points[i][0], points[k][1]-points[i][1]
		return x1*y2 > x2*y1
	}

	calc := func() []int32 {
		var p []int32
		for _, k := range order {
			for len(p) > 1 {
				i, j := p[len(p)-2], p[len(p)-1]
				if check(i, j, k) {
					break
				}
				p = p[:len(p)-1]
			}
			p = append(p, k)
		}
		return p
	}

	var p []int32
	if mode == Full || mode == Lower {
		p = append(p, calc()...)
	}
	if mode == Full || mode == Upper {
		if len(p) > 0 {
			p = p[:len(p)-1]
		}
		reverse(order)
		p = append(p, calc()...)
	}
	if mode == Upper {
		reverse(p)
	}
	for len(p) >= 2 && points[p[0]] == points[p[len(p)-1]] {
		p = p[:len(p)-1]
	}
	return p
}

// 平面最远点对(旋转卡壳).
func FurthestPair(points [][2]int) (int32, int32) {
	best := -1
	resI, resJ := int32(-1), int32(-1)

	update := func(i, j int32) {
		x1, y1 := points[i][0]-points[j][0], points[i][1]-points[j][1]
		d := x1*x1 + y1*y1
		if d > best {
			best = d
			resI, resJ = i, j
		}
	}

	update(0, 1)

	order := ConvexHull(points, Full, false)
	n := int32(len(order))
	if n == 1 {
		return resI, resJ
	}
	if n == 2 {
		return order[0], order[1]
	}
	order = append(order, order...)

	newPoints := reArrage(points, order)
	check := func(i, j int32) bool {
		x1, y1 := newPoints[i+1][0]-newPoints[i][0], newPoints[i+1][1]-newPoints[i][1]
		x2, y2 := newPoints[j+1][0]-newPoints[j][0], newPoints[j+1][1]-newPoints[j][1]
		return x1*y2 > x2*y1
	}

	j := int32(1)
	for i := int32(0); i < n; i++ {
		if i >= j {
			j = i
		}
		for j < 2*n && check(i, j) {
			j++
		}
		update(order[i], order[j])
	}
	return resI, resJ
}

func argSort(nums []int) []int32 {
	order := make([]int32, len(nums))
	for i := int32(0); i < int32(len(nums)); i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func reverse[T any](nums []T) {
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func reArrage[T any](nums []T, order []int32) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}
