// 非公平组合游戏（Partizan Game）与公平组合游戏的区别在于在非公平组合游戏中，游戏者在某一确定状态可以做出的决策集合与游戏者有关。
// 非公平博弈

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc209_e()
}

// 成语接龙.
func abc209_e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	D := NewDictionary[string]()
	for i := int32(0); i < n; i++ {
		a := words[i][:3]
		b := words[i][len(words[i])-3:]
		D.Id(a)
		D.Id(b)
	}

	game := NewGraphGame(D.Size())
	for i := int32(0); i < n; i++ {
		a := words[i][:3]
		b := words[i][len(words[i])-3:]
		ida, idb := D.Id(a), D.Id(b)
		game.AddEdge(ida, idb)
	}
	game.Build()

	f := func(win, lose bool) string {
		if win {
			return "Aoki"
		}
		if lose {
			return "Takahashi"
		}
		return "Draw"
	}

	for i := int32(0); i < n; i++ {
		j := game.Edges[i][1]
		res := f(game.Win[j], game.Lose[j])
		fmt.Fprintln(out, res)
	}

}

const INF32 int32 = 1e9 + 10

// 给定一张有向图，玩家轮流行动，无法行动的玩家输.
type GraphGame struct {
	n            int32
	Win          []bool
	Lose         []bool
	Edges        [][2]int32
	endTurn      []int32
	bestStrategy []int32
}

func NewGraphGame(n int32) *GraphGame {
	return &GraphGame{n: n}
}

func (g *GraphGame) AddEdge(u, v int32) {
	g.Edges = append(g.Edges, [2]int32{u, v})
}

func (g *GraphGame) Build() {
	graph, revGraph := make([][]int32, g.n), make([][]int32, g.n)
	for _, e := range g.Edges {
		u, v := e[0], e[1]
		graph[u] = append(graph[u], v)
		revGraph[v] = append(revGraph[v], u)
	}
	n := g.n
	indeg, outdeg := make([]int32, n), make([]int32, n)
	for cur, nexts := range graph {
		for _, next := range nexts {
			indeg[next]++
			outdeg[cur]++
		}
	}
	win, lose := make([]bool, n), make([]bool, n)
	endTurn := make([]int32, n)
	bestStrategy := make([]int32, n)
	for v := int32(0); v < n; v++ {
		endTurn[v] = INF32
		bestStrategy[v] = -1
	}
	queue := make([]int32, 0)
	for v := int32(0); v < n; v++ {
		if outdeg[v] == 0 {
			queue = append(queue, v)
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if win[cur] || lose[cur] {
			continue
		}
		lose[cur] = true
		for _, next := range graph[cur] {
			if lose[next] {
				win[cur] = true
			}
			if !win[next] {
				lose[cur] = false
			}
		}
		if win[cur] {
			for _, next := range graph[cur] {
				if endTurn[next] > endTurn[cur]+1 {
					endTurn[next] = endTurn[cur] + 1
					bestStrategy[next] = cur
				}
			}
		}
		if lose[cur] {
			endTurn[cur] = 0
			for _, next := range graph[cur] {
				if endTurn[next] < endTurn[cur]+1 {
					endTurn[cur] = endTurn[next] + 1
					bestStrategy[cur] = next
				}
			}
		}
		for _, next := range revGraph[cur] {
			outdeg[next]--
			if lose[cur] || outdeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	g.Win, g.Lose, g.endTurn, g.bestStrategy = win, lose, endTurn, bestStrategy
}

// 有向图反图.
func reverseGraph(graph [][]int32) [][]int32 {
	n := len(graph)
	rg := make([][]int32, n)
	for cur, nexts := range graph {
		for _, next := range nexts {
			rg[next] = append(rg[next], int32(cur))
		}
	}
	return rg
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

type Dictionary[V comparable] struct {
	_idToValue []V
	_valueToId map[V]int32
}

// A dictionary that maps values to unique ids.
func NewDictionary[V comparable]() *Dictionary[V] {
	return &Dictionary[V]{
		_valueToId: map[V]int32{},
	}
}
func (d *Dictionary[V]) Id(value V) int32 {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := int32(len(d._idToValue))
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary[V]) Value(id int32) V {
	return d._idToValue[id]
}
func (d *Dictionary[V]) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary[V]) Size() int32 {
	return int32(len(d._idToValue))
}
