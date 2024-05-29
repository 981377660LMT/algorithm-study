package main

import "fmt"

const INF int = 2e18

func main() {
	// nums1 = [1,2], nums2 = [2,3]
	fmt.Println(minimumXORSum([]int{1, 2}, []int{2, 3})) // 2
}

// 1879. 两个数组最小的异或值之和
// https://leetcode.cn/problems/minimum-xor-sum-of-two-arrays/description/
func minimumXORSum(nums1 []int, nums2 []int) int {
	n := int32(len(nums1))
	START, END := 2*n+1, 2*n+2
	graph := make([][]*CostFlowEdge, END+1)
	addEdge := func(from, to int32, flow, cost int) {
		e1 := NewCostFlowEdge(to, 0, true, cost)
		e2 := NewCostFlowEdge(from, flow, false, -cost)
		e1.rev = e2
		e2.rev = e1
		graph[from] = append(graph[from], e1)
		graph[to] = append(graph[to], e2)
	}
	for i := int32(0); i < n; i++ {
		addEdge(START, i, 1, 0)
		addEdge(i+n, END, 1, 0)
		for j := int32(0); j < n; j++ {
			addEdge(i, j+n, 1, nums1[i]^nums2[j])
		}
	}
	mcmf := NewSpfaMinimumCostFlow()
	_, cost := mcmf.Apply(graph, START, END, INF)
	return cost
}

type SpfaMinimumCostFlow struct {
	queue    []int32
	dists    []int
	inque    []bool
	prev     []*CostFlowEdge
	net      [][]*CostFlowEdge
	callback func(flow, cost int) (keep bool)
}

func NewSpfaMinimumCostFlow() *SpfaMinimumCostFlow {
	return &SpfaMinimumCostFlow{}
}

func (mcf *SpfaMinimumCostFlow) SetCallback(callback func(flow, cost int) (keep bool)) {
	mcf.callback = callback
}

func (mcf *SpfaMinimumCostFlow) Apply(net [][]*CostFlowEdge, s, t int32, send int) (flow, cost int) {
	mcf.prepare(int32(len(net)))
	mcf.net = net
	for flow < send {
		mcf.spfa(t, INF)
		if mcf.dists[s] == INF {
			break
		}
		iter := s
		sent := send - flow
		for mcf.prev[iter] != nil {
			sent = min(sent, mcf.prev[iter].flow)
			iter = mcf.prev[iter].rev.to
		}
		if mcf.callback != nil && !mcf.callback(sent, mcf.dists[s]) {
			break
		}
		iter = s
		for mcf.prev[iter] != nil {
			mcf.send(mcf.prev[iter], -sent)
			iter = mcf.prev[iter].rev.to
		}
		cost += sent * mcf.dists[s]
		flow += sent
	}
	return
}

func (mcf *SpfaMinimumCostFlow) prepare(vertexNum int32) {
	if mcf.dists == nil || len(mcf.dists) < int(vertexNum) {
		mcf.queue = make([]int32, 0, vertexNum)
		mcf.dists = make([]int, vertexNum)
		mcf.inque = make([]bool, vertexNum)
		mcf.prev = make([]*CostFlowEdge, vertexNum)
	}
}

func (mcf *SpfaMinimumCostFlow) send(e *CostFlowEdge, flow int) {
	e.flow += flow
	e.rev.flow -= flow
}

func (mcf *SpfaMinimumCostFlow) spfa(s int32, inf int) {
	mcf.queue = mcf.queue[:0]
	for i := range mcf.net {
		mcf.dists[i] = inf
		mcf.inque[i] = false
	}
	mcf.dists[s] = 0
	mcf.prev[s] = nil
	mcf.queue = append(mcf.queue, s)
	for len(mcf.queue) > 0 {
		head := mcf.queue[0]
		mcf.queue = mcf.queue[1:]
		mcf.inque[head] = false
		for _, e := range mcf.net[head] {
			if e.flow > 0 && mcf.dists[e.to] > mcf.dists[head]-e.cost {
				mcf.dists[e.to] = mcf.dists[head] - e.cost
				mcf.prev[e.to] = e
				if !mcf.inque[e.to] {
					mcf.inque[e.to] = true
					mcf.queue = append(mcf.queue, e.to)
				}
			}
		}
	}
}

type CostFlowEdge struct {
	to   int32
	flow int
	real bool
	cost int
	rev  *CostFlowEdge
}

func NewCostFlowEdge(to int32, flow int, real bool, cost int) *CostFlowEdge {
	return &CostFlowEdge{to: to, flow: flow, real: real, cost: cost}
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
