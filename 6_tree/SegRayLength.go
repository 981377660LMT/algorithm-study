// https://github.com/old-yan/CP-template/blob/a07b6fe0092e9ee890a0e35ada6ea1bb2c83ba05/TREE/SegRayLengthHelper_vector.md#L3
// 最长线段射线助手。

// seg:线段,这里指的是某条路径.
// ray:射线,这里指的是以某个树节点为起点的路径.
// 对于根结点，所有的连通部分都在自己下方。
// 对于非根结点，一个连通部分在自己上方，剩余连通部分都在自己下方。
// 在预处理后，
// !  `ray[i]` 存储了 以 `i` 作为一端的线段的最长长度`前三名`。
// `seg[i]` 存储了不经过 `i` 的线段的最长长度前两名。
// `downRay[i]` 存储了以 `i` 作为一端的向下的线段的最长长度。
// `downSeg[i]` 存储了 `i` 下方的不越过 `i` 的线段的线段最长长度。注意，不越过，但是可以经过。
// `upRay[i]` 存储了 以 `i` 作为一端的向上的线段的最长长度。
// `upSeg[i]` 存储了 `i` 上方的不越过 `i` 的线段的线段最长长度。注意，不越过，但是可以经过。

// !注意：此处的前三名、前两名不允许某个连通部分(子树)重复占用。
// 比如结点 `i` 有两个邻居，从第一个连通部分可以找到长度为 `10` 的 `ray` ，也可以找到长度为 `9` 的 `ray` ；
// 从第二个联通部分可以找到长度为 `8` 的 `ray` ，
// 那么 `m_ray[i][0]==10` ，`m_ray[i][1]==8` ，`m_ray[i][2]==0` ，也就是没有第三名。

package main

import "fmt"

func main() {
	tree := NewSegRayLength(5)
	tree.AddEdge(2, 0, 1)

	tree.AddEdge(1, 3, 1)
	tree.AddEdge(4, 0, 1)
	tree.AddEdge(3, 0, 1)
	tree.Build(3)

	fmt.Println(tree.Ray)
	fmt.Println(tree.Seg)
	fmt.Println(tree.DownRay)
	fmt.Println(tree.DownSeg)
	fmt.Println(tree.UpRay)
	fmt.Println(tree.UpSeg)

	fmt.Println(tree.QueryMaxRayAndSeg(0, 2))
}

func treeDiameter(edges [][]int) int {
	SR := NewSegRayLength(len(edges) + 1)
	for _, e := range edges {
		SR.AddEdge(e[0], e[1], 1)
	}
	SR.Build(0)
	res := 0
	for _, r := range SR.Ray {
		if r[0] > res {
			res = r[0]
		}
	}
	return res
}

// 用于处理树中最长射线/线段的问题.
//  https://github.com/old-yan/CP-template/blob/main/TREE/SegRayLengthHelper_vector.h
type SegRayLength struct {
	Ray     [][3]int // 以i为端点的线段长度前三大的值
	Seg     [][2]int // 不经过i的线段长度前两大的值
	DownRay []int    // 以i为端点的向下的线段长度最大值
	DownSeg []int    // i下方不越过i的线段长度最大值
	UpRay   []int    // 以i为端点的向上的线段长度最大值
	UpSeg   []int    // i上方不越过i的线段长度最大值

	tree          [][][2]int
	depthWeighted []int
	lid           []int
	top           []int
	parent        []int
	heavySon      []int
	dfn           int
}

func NewSegRayLength(n int) *SegRayLength {
	res := &SegRayLength{
		Ray:           make([][3]int, n),
		Seg:           make([][2]int, n),
		DownRay:       make([]int, n),
		DownSeg:       make([]int, n),
		UpRay:         make([]int, n),
		UpSeg:         make([]int, n),
		tree:          make([][][2]int, n),
		depthWeighted: make([]int, n),
		lid:           make([]int, n),
		top:           make([]int, n),
		parent:        make([]int, n),
		heavySon:      make([]int, n),
	}
	for i := 0; i < n; i++ {
		res.parent[i] = -1
	}
	return res
}

func (s *SegRayLength) AddEdge(u, v, w int) {
	s.tree[u] = append(s.tree[u], [2]int{v, w})
	s.tree[v] = append(s.tree[v], [2]int{u, w})
}

func (s *SegRayLength) AddDirectedEdge(u, v, w int) {
	s.tree[u] = append(s.tree[u], [2]int{v, w})
}

func (s *SegRayLength) Build(root int) {
	s._dfs1(root, -1, 0)
	s._dfs2(root, 0, 0)
	s._dfs3(root, root)
}

// 查询最长射线ray,最长线段seg.
//  u:树节点.
//  ignoreRoot:屏蔽掉以ignoreRoot为根的子树.需要保证ignoreRoot在u的子树中.
//  -1:不屏蔽任何子树.
//  返回:树中剩余部分从u出发的最长射线,树中剩余部分的最长线段.
func (s *SegRayLength) QueryMaxRayAndSeg(u, ignoreRoot int) (maxRay, maxSeg int) {
	if ignoreRoot == -1 {
		return s._maxRaySeg(u, -1, -1)
	}
	return s._maxRaySeg(u, s.DownRay[ignoreRoot]+s._weightedDist(u, ignoreRoot), s.DownSeg[ignoreRoot])
}

func (s *SegRayLength) _dfs1(cur, pre, dist int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	downRay, downSeg := s.DownRay, s.DownSeg
	for _, e := range s.tree[cur] {
		next, weight := e[0], e[1]
		if next != pre {
			nextSize := s._dfs1(next, cur, dist+weight)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
			s._addDownRay(cur, downRay[next]+weight)
			s._addDownSeg(cur, downSeg[next])
		}
	}
	len1, len2 := s.Ray[cur][0], s.Ray[cur][1]
	cand := len1 + len2
	if cand > downSeg[cur] {
		downSeg[cur] = cand
	}
	s.depthWeighted[cur] = dist
	s.parent[cur] = pre
	s.heavySon[cur] = heavySon
	return subSize
}

func (s *SegRayLength) _dfs2(cur, upRay, upSeg int) {
	s._setUpRay(cur, upRay)
	s._setUpSeg(cur, upSeg)
	downRay, downSeg := s.DownRay, s.DownSeg
	for _, e := range s.tree[cur] {
		next, weight := e[0], e[1]
		if next != s.parent[cur] {
			ray, seg := s._maxRaySeg(cur, downRay[next]+weight, downSeg[next])
			s._addSeg(next, seg)
			ray += weight
			if ray > seg {
				seg = ray
			}
			s._dfs2(next, ray, seg)
		}
	}
}

func (s *SegRayLength) _dfs3(cur, top int) int {
	s.top[cur] = top
	s.lid[cur] = s.dfn
	s.dfn++
	heavySon := s.heavySon[cur]
	if heavySon != -1 {
		s._dfs3(heavySon, top)
		for _, e := range s.tree[cur] {
			next := e[0]
			if next != heavySon && next != s.parent[cur] {
				s._dfs3(next, next)
			}
		}
	}
	return s.lid[cur]
}

func (s *SegRayLength) _addRay(i, ray int) {
	rayI := &s.Ray[i]
	if ray > rayI[0] {
		rayI[2] = rayI[1]
		rayI[1] = rayI[0]
		rayI[0] = ray
	} else if ray > rayI[1] {
		rayI[2] = rayI[1]
		rayI[1] = ray
	} else if ray > rayI[2] {
		rayI[2] = ray
	}
}

func (s *SegRayLength) _addSeg(i, seg int) {
	segI := &s.Seg[i]
	if seg > segI[0] {
		segI[1] = segI[0]
		segI[0] = seg
	} else if seg > segI[1] {
		segI[1] = seg
	}
}

func (s *SegRayLength) _addDownRay(i, ray int) {
	if ray > s.DownRay[i] {
		s.DownRay[i] = ray
	}
	s._addRay(i, ray)
}

func (s *SegRayLength) _addDownSeg(i, seg int) {
	if seg > s.DownSeg[i] {
		s.DownSeg[i] = seg
	}
	s._addSeg(i, seg)
}

func (s *SegRayLength) _setUpRay(i, ray int) {
	s.UpRay[i] = ray
	s._addRay(i, ray)
}

func (s *SegRayLength) _setUpSeg(i, seg int) {
	s.UpSeg[i] = seg
}

// 查询树中某部分的最长射线和线段.
//  u:树中某点.
//  ignoreRay:屏蔽掉的部分提供的最长射线.
//  ignorSeg:屏蔽掉的部分提供的最长线段.
//  返回:从u出发的最长射线和树中剩余部分的最长线段.
func (s *SegRayLength) _maxRaySeg(u, ignoreRay, ignorSeg int) (maxRay, maxSeg int) {
	r0, r1, r2 := s.Ray[u][0], s.Ray[u][1], s.Ray[u][2]
	s0, s1 := s.Seg[u][0], s.Seg[u][1]
	maxRay = r0
	if ignoreRay == r0 {
		maxRay = r1
	}
	maxSeg = s0
	if ignorSeg == s0 {
		maxSeg = s1
	}
	twoRay := r0 + r1
	if ignoreRay == r0 {
		twoRay = r1 + r2
	} else if ignoreRay == r1 {
		twoRay = r0 + r2
	}
	if maxSeg < twoRay {
		maxSeg = twoRay
	}
	return
}

func (s *SegRayLength) _lca(u, v int) int {
	for {
		if s.lid[u] > s.lid[v] {
			u, v = v, u
		}
		if s.top[u] == s.top[v] {
			return u
		}
		v = s.parent[s.top[v]]
	}
}

func (s *SegRayLength) _weightedDist(u, v int) int {
	return s.depthWeighted[u] + s.depthWeighted[v] - 2*s.depthWeighted[s._lca(u, v)]
}
