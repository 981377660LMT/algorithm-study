package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// int dx[] = {1, 0, -1, 0, 1, 1, -1, -1};
// int dy[] = {0, 1, 0, -1, 1, -1, 1, -1};

// int LIM = 200100;

// struct T {
//   int idx;
//   ll y1, y2;
// };
// void solve() {
//   LL(N);
//   VEC(pi, stone, N);
//   for (auto& [x, y]: stone) ++x, ++y;
//   vvc<int> XtoY(LIM);
//   for (auto& [x, y]: stone) XtoY[x].eb(y);

//   // X -> [y1,y2]
//   vvc<T> seg(LIM);
//   int n = 0;
//   FOR(x, LIM) {
//     auto Y = XtoY[x];
//     sort(all(Y));
//     ll a = 0;
//     for (auto& y: Y) {
//       // [a,y-1]
//       if (a < y) seg[x].eb(T{n++, a, y - 1});
//       a = y + 1;
//     }
//     seg[x].eb(T{n++, a, LIM - 1});
//   }

//   UnionFind uf(n);
//   FOR(x, LIM - 1) {
//     auto A = seg[x];
//     auto B = seg[x + 1];
//     while (len(A) && len(B)) {
//       auto& [i, y1, y2] = A.back();
//       auto& [j, y3, y4] = B.back();
//       if (y1 > y4) {
//         POP(A);
//         continue;
//       }
//       if (y3 > y2) {
//         POP(B);
//         continue;
//       }
//       assert(max(y1, y3) <= min(y2, y4));
//       uf.merge(i, j);
//       if (y1 < y3) {
//         POP(B);
//       } else {
//         POP(A);
//       }
//     }
//   }
//   ll ANS = 0;
//   FOR(x, LIM) {
//     for (auto& [i, y1, y2]: seg[x]) {
//       if (uf[i] != uf[0]) { ANS += y2 - y1 + 1; }
//     }
//   }
//   print(ANS);
// }

const INF int = 1e18

// G - Go Territory (围棋盘上不能到达原点的点)
// https://atcoder.jp/contests/abc361/tasks/abc361_g
//
// 二维平面，有障碍物，可以上下左右走。
// 问有多少个点，不可以走到(-1,-1).
//
// 1.横纵坐标+1，将点变为格子处理；
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	stones := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &stones[i][0], &stones[i][1])
		stones[i][0]++
		stones[i][1]++
	}

	xToY := make(map[int][]int)
	for i := 0; i < n; i++ {
		x, y := stones[i][0], stones[i][1]
		xToY[x] = append(xToY[x], y)
	}
	keys := make([]int, 0, len(xToY))
	for k := range xToY {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// !扫描线，维护每个x坐标的y坐标区间(矩形)
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *UnionFindArraySimple32) Size(key int32) int32 {
	return -u.data[u.Find(key)]
}
