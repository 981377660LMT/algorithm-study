// 橋や関節点などを効率的に求める際に有効なアルゴリズム.
// グラフをDFSして各頂点 idx について, ord[idx] := DFS で頂点に訪れた順番,
// low[idx] := 頂点 idxからDFS木の葉方向の辺を 0回以上,
// 後退辺を 1回以下通って到達可能な頂点の ord の最小値 を求める.

// build(): LowLink を構築する.
// !構築後, Articulation には関節点, Bridge には橋が格納される.
// 非連結でも多重辺を含んでいてもOK.

package main

func criticalConnections(n int, connections [][]int) [][]int {
	g := make([][]Edge, n)
	for i := 0; i < len(connections); i++ {
		u, v := connections[i][0], connections[i][1]
		g[u] = append(g[u], Edge{from: u, to: v})
		g[v] = append(g[v], Edge{from: v, to: u})
	}
	lowLink := NewLowLink(g)
	lowLink.Build()
	res := [][]int{}
	for i := 0; i < len(lowLink.Bridge); i++ {
		res = append(res, []int{lowLink.Bridge[i].from, lowLink.Bridge[i].to})
	}
	return res
}

type Edge = struct{ from, to int }
type LowLink struct {
	Articulation []int  // 関節点
	Bridge       []Edge // 橋
	g            [][]Edge
	ord, low     []int
	used         []bool
}

func NewLowLink(g [][]Edge) *LowLink {
	return &LowLink{g: g}
}

func (ll *LowLink) Build() {
	ll.used = make([]bool, len(ll.g))
	ll.ord = make([]int, len(ll.g))
	ll.low = make([]int, len(ll.g))
	k := 0
	for i := 0; i < len(ll.g); i++ {
		if !ll.used[i] {
			k = ll.dfs(i, k, -1)
		}
	}
}

func (ll *LowLink) dfs(idx, k, par int) int {
	ll.used[idx] = true
	ll.ord[idx] = k
	k++
	ll.low[idx] = ll.ord[idx]
	isArticulation := false
	beet := false
	cnt := 0
	for _, e := range ll.g[idx] {
		if e.to == par {
			tmp := beet
			beet = true
			if !tmp {
				continue
			}
		}
		if !ll.used[e.to] {
			cnt++
			k = ll.dfs(e.to, k, idx)
			ll.low[idx] = min(ll.low[idx], ll.low[e.to])
			if par >= 0 && ll.low[e.to] >= ll.ord[idx] {
				isArticulation = true
			}
			if ll.ord[idx] < ll.low[e.to] {
				ll.Bridge = append(ll.Bridge, e)
			}
		} else {
			ll.low[idx] = min(ll.low[idx], ll.ord[e.to])
		}
	}

	if par == -1 && cnt > 1 {
		isArticulation = true
	}
	if isArticulation {
		ll.Articulation = append(ll.Articulation, idx)
	}
	return k
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
