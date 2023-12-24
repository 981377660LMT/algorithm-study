// G - Christmas Color Grid 2-可撤销并查集维护联通分量个数
// https://atcoder.jp/contests/abc334/editorial/8995

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// G - Christmas Color Grid 2
// https://atcoder.jp/contests/abc334/tasks/abc334_g
// 给定一个网格图，由红/绿两种点组成.
// !询问：将每个绿色点变为红色(删除这个点)后，剩余的绿色点连通块的大小.
// ROW,COL<=1000.
//
// !反向思考，问题等价于：开始时所有点为红色，第i次操作将[0,i)+[i+1,n)的点全部涂成绿色，求此时绿色点的连通块个数.
// 分治删点即可.
// O(ROW*COL*lg(ROW*COL)^2)
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const GREEN byte = '#'
	const RED byte = '.'
	const MOD int = 998244353
	var DIR4 = [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	pow := func(base, exp, mod int) int {
		base %= mod
		res := 1 % mod
		for ; exp > 0; exp >>= 1 {
			if exp&1 == 1 {
				res = res * base % mod
			}
			base = base * base % mod
		}
		return res
	}

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)
	grid := make([][]byte, ROW)
	for i := 0; i < ROW; i++ {
		grid[i] = make([]byte, COL)
		fmt.Fscan(in, &grid[i])
	}

	uf := NewUnionFindArrayWithUndo(ROW * COL)
	greenPos := [][2]int{}
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			if grid[i][j] == GREEN {
				greenPos = append(greenPos, [2]int{i, j})
			}
		}
	}

	res := 0
	history := [][2]int{} // 保存(并查集的state,上一次修改绿色点的index)
	colors := make([][]byte, ROW)
	for i := 0; i < ROW; i++ {
		colors[i] = make([]byte, COL)
		for j := 0; j < COL; j++ {
			colors[i][j] = RED
		}
	}

	uf.SetPart(0) // 初始时绿色点的联通分量数为0
	MutateWithoutOneUndo(
		0, len(greenPos),
		func(index int) {
			curX, curY := greenPos[index][0], greenPos[index][1]
			history = append(history, [2]int{uf.GetState(), index})
			colors[curX][curY] = GREEN
			uf.Part++

			for _, d := range DIR4 {
				nextX, nextY := curX+d[0], curY+d[1]
				if 0 <= nextX && nextX < ROW && 0 <= nextY && nextY < COL && colors[nextX][nextY] == GREEN {
					uf.Union(uf.Find(curX*COL+curY), uf.Find(nextX*COL+nextY))
				}
			}
		},
		func() {
			pre := history[len(history)-1]
			preState, preIndex := pre[0], pre[1]
			history = history[:len(history)-1]

			uf.Rollback(preState)
			colors[greenPos[preIndex][0]][greenPos[preIndex][1]] = RED
			uf.Part--
		},
		func(_ int) {
			res = (res + uf.Part) % MOD
		},
	)

	fmt.Println(res * pow(len(greenPos), MOD-2, MOD) % MOD)
}

// 线段树分治的特殊情形.
// 调用 `query` 时，`state` 为对除了 `index` 以外所有点均调用过了 `mutate` 的状态。但不保证调用 `mutate` 的顺序。
// 总计会调用 $O(NlgN)$ 次的 `mutate` 和 `undo`, 以及 $O(N)$ 次的 `query`.
func MutateWithoutOneUndo(
	start, end int,
	/** 这里的 index 也就是 time. */
	mutate func(index int),
	undo func(),
	query func(index int),
) {
	var dfs func(curStart, curEnd int)
	dfs = func(curStart, curEnd int) {
		if curEnd == curStart+1 {
			query(curStart)
			return
		}

		mid := (curStart + curEnd) >> 1
		for i := curStart; i < mid; i++ {
			mutate(i)
		}
		dfs(mid, curEnd)
		for i := curStart; i < mid; i++ {
			undo()
		}

		for i := mid; i < curEnd; i++ {
			mutate(i)
		}
		dfs(curStart, mid)
		for i := mid; i < curEnd; i++ {
			undo()
		}
	}

	dfs(start, end)
}

type historyItem struct{ a, b int32 }

type UnionFindArrayWithUndo struct {
	Part      int
	n         int32
	innerSnap int
	data      []int32
	history   []historyItem // (root,data)
}

func NewUnionFindArrayWithUndo(n int) *UnionFindArrayWithUndo {
	data := make([]int32, n)
	for i := range data {
		data[i] = -1
	}
	return &UnionFindArrayWithUndo{Part: n, n: int32(n), data: data}
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *UnionFindArrayWithUndo) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}
	small, smallData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	big, bigData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	uf.data[small] = smallData
	uf.data[big] = bigData
	if big != small {
		uf.Part++
	}
	return true
}

// 保存并查集当前的状态.
//
//	!Snapshot() 之后可以调用 Rollback(-1) 回滚到这个状态.
func (uf *UnionFindArrayWithUndo) Snapshot() {
	uf.innerSnap = len(uf.history) >> 1
}

// 回滚到指定的状态.
//
//	state 为 -1 表示回滚到上一次 `SnapShot` 时保存的状态.
//	其他值表示回滚到状态id为state时的状态.
func (uf *UnionFindArrayWithUndo) Rollback(state int) bool {
	if state == -1 {
		state = uf.innerSnap
	}
	state <<= 1
	if state < 0 || state > len(uf.history) {
		return false
	}
	for state < len(uf.history) {
		uf.Undo()
	}
	return true
}

// 获取当前并查集的状态id.
//
//	也就是当前合并(Union)被调用的次数.
func (uf *UnionFindArrayWithUndo) GetState() int {
	return len(uf.history) >> 1
}

func (uf *UnionFindArrayWithUndo) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndo) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	uf.history = append(uf.history, historyItem{int32(x), uf.data[x]})
	uf.history = append(uf.history, historyItem{int32(y), uf.data[y]})
	if x == y {
		return false
	}
	if uf.data[x] > uf.data[y] {
		x ^= y
		y ^= x
		x ^= y
	}
	uf.data[x] += uf.data[y]
	uf.data[y] = int32(x)
	uf.Part--
	return true
}

func (uf *UnionFindArrayWithUndo) Find(x int) int {
	cur := int32(x)
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return int(cur)
}

func (uf *UnionFindArrayWithUndo) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }

func (uf *UnionFindArrayWithUndo) GetSize(x int) int { return int(-uf.data[uf.Find(x)]) }

func (ufa *UnionFindArrayWithUndo) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < int(ufa.n); i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArrayWithUndo) SetPart(part int) { ufa.Part = part }

func (ufa *UnionFindArrayWithUndo) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
