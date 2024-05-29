// package template.problem;

// import template.graph.Graph;
// import template.primitve.generated.graph.LongAugmentCallback;
// import template.primitve.generated.graph.LongAugmentMinimumCostFlow;
// import template.primitve.generated.graph.LongCostFlowEdge;
// import template.primitve.generated.graph.LongFlow;

// import java.util.ArrayList;
// import java.util.List;

// /**
//  * <p>
//  * Given directed graph. You can increase the length of some edges with each unit cost and you can't
//  * make such operation more than limit times. Find the minimum possible distance in
//  * O(FE\log_2V+q\log_2F) while F=\sum_e cost(e).
//  * </p>
//  */
// public class MakeDistanceFarthest {
//   private static long inf = Long.MAX_VALUE / 4;

//   private static class LinearFunction {
//     long l;
//     long a;
//     long b;

//     public LinearFunction(long l, long a, long b) {
//       this.l = l;
//       this.a = a;
//       this.b = b;
//     }

//     long getL() {
//       return a * l + b;
//     }

//     double inverse(double y) {
//       return (y - b) / a;
//     }

//     long apply(long x) {
//       return a * x + b;
//     }
//   }

//   private List<LongCostFlowEdge>[] g;
//   LinearFunction[] fs;

//   public MakeDistanceFarthest(int n) {
//     g = Graph.createGraph(n);
//   }

//   public void addEdge(int u, int v, long len, long cost) {
//     LongFlow.addCostFlowEdge(g, u, v, cost, len);
//   }

//   public void addLimitedEdge(int u, int v, long len, long cost, long limit) {
//     LongFlow.addCostFlowEdge(g, u, v, cost, len);
//     LongFlow.addCostFlowEdge(g, u, v, inf, len + limit);
//   }

//   /**
//    * ensure there is a path from s to t
//    *
//    * @param mcf
//    * @param s
//    * @param t
//    */
//   public void solve(LongAugmentMinimumCostFlow mcf, int s, int t, long budgeLimit, long distLimit,
//       long flowLimit) {
//     List<LinearFunction> list = new ArrayList<>();
//     LongAugmentCallback callback = new LongAugmentCallback() {
//       long sumFlow = 0;
//       long sumCost = 0;

//       @Override
//       public boolean callback(long flow, long pathCost) {
//         sumFlow += flow;
//         sumCost += flow * pathCost;

//         if (!list.isEmpty() && list.get(list.size() - 1).l == pathCost) {
//           list.remove(list.size() - 1);
//         }
//         LinearFunction func = new LinearFunction(pathCost, sumFlow, -sumCost);
//         list.add(func);
//         return func.getL() <= budgeLimit && func.l <= distLimit;
//       }
//     };
//     mcf.setCallback(callback);
//     mcf.apply(g, s, t, flowLimit);
//     fs = list.toArray(new LinearFunction[0]);
//   }

//   /**
//    * get maximum distance with no more than x expense in O(\log_2F).
//    *
//    * @param x
//    * @return
//    */
//   public double queryByExpense(long x) {
//     int l = 0;
//     int r = fs.length - 1;
//     while (l < r) {
//       int mid = (l + r) / 2;
//       boolean valid = mid + 1 >= fs.length || fs[mid + 1].getL() > x;
//       if (valid) {
//         r = mid;
//       } else {
//         l = mid + 1;
//       }
//     }
//     return fs[l].inverse(x);
//   }

//   /**
//    * <pre>
//    * get minimum expense that make distance greater than or equal to x.
//    * time complexity: O(\log_2n)
//    * </pre>
//    */
//   public long queryByShortestPath(long x) {
//     int l = 0;
//     int r = fs.length - 1;
//     while (l < r) {
//       int mid = (l + r + 1) / 2;
//       boolean valid = fs[mid].l <= x;
//       if (valid) {
//         l = mid;
//       } else {
//         r = mid - 1;
//       }
//     }
//     return fs[l].apply(x);
//   }
// }

// 给定一个有向图.
// 你可以增加某些边的长度，每增加一单位长度需要花费一定的代价，但是你不能增加的代价超过一定的限制.
// 请找出可能的最小距离.

package main

func main() {

}

const INF int = 2e18

type MakeDistanceFarthest struct {
	g  [][]*CostFlowEdge
	fs []*LinearFunction
}

func NewMakeDistanceFarthest(n int) *MakeDistanceFarthest {
	return &MakeDistanceFarthest{g: make([][]*CostFlowEdge, n)}
}

func (mdf *MakeDistanceFarthest) AddEdge(u, v int32, len, cost int) {
	mdf.addCostFlowEdge(mdf.g, u, v, cost, len)
}

func (mdf *MakeDistanceFarthest) AddLimitedEdge(u, v int32, len, cost, limit int) {
	mdf.addCostFlowEdge(mdf.g, u, v, cost, len)
	mdf.addCostFlowEdge(mdf.g, u, v, INF, len+limit)
}

func (mdf *MakeDistanceFarthest) Solve(flow *SpfaMinimumCostFlow, s, t int32, budgeLimit, distLimit, flowLimit int) {
	var list []*LinearFunction
	sumFlow, sumCost := 0, 0
	callback := func(flow, pathCost int) bool {
		sumFlow += flow
		sumCost += flow * pathCost
		if len(list) > 0 && list[len(list)-1].l == pathCost {
			list = list[:len(list)-1]
		}
		f := NewLinearFunction(pathCost, sumFlow, -sumCost)
		list = append(list, f)
		return f.GetL() <= budgeLimit && f.l <= distLimit
	}
	flow.SetCallback(callback)
	flow.Apply(mdf.g, s, t, flowLimit)
	mdf.fs = list
}

// 花费不超过 x 时的最大距离.
func (mdf *MakeDistanceFarthest) QueryByExpense(expense int) float64 {
	l, r := int32(0), int32(len(mdf.fs)-1)
	for l < r {
		mid := (l + r) >> 1
		valid := (mid+1 >= int32(len(mdf.fs))) || (mdf.fs[mid+1].GetL() > expense)
		if valid {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return mdf.fs[l].Inverse(float64(expense))
}

// 距离不小于 x 时的最小花费.
func (mdf *MakeDistanceFarthest) QueryByShortestPath(dist int) int {
	l, r := int32(0), int32(len(mdf.fs)-1)
	for l < r {
		mid := (l + r + 1) >> 1
		valid := mdf.fs[mid].l <= dist
		if valid {
			l = mid
		} else {
			r = mid - 1
		}
	}
	return mdf.fs[l].Apply(dist)
}

func (mdf *MakeDistanceFarthest) addCostFlowEdge(g [][]*CostFlowEdge, s, t int32, cap, cost int) *CostFlowEdge {
	real := NewCostFlowEdge(t, 0, true, cost)
	virtual := NewCostFlowEdge(s, cap, false, -cost)
	real.rev = virtual
	virtual.rev = real
	g[s] = append(g[s], real)
	g[t] = append(g[t], virtual)
	return real
}

type LinearFunction struct{ l, a, b int }

func NewLinearFunction(l, a, b int) *LinearFunction  { return &LinearFunction{l: l, a: a, b: b} }
func (lf *LinearFunction) GetL() int                 { return lf.a*lf.l + lf.b }
func (lf *LinearFunction) Inverse(y float64) float64 { return (y - float64(lf.b)) / float64(lf.a) }
func (lf *LinearFunction) Apply(x int) int           { return lf.a*x + lf.b }

type SpfaMinimumCostFlow struct {
	queue    []int32
	dists    []int
	inque    []bool
	prev     []*CostFlowEdge
	net      [][]*CostFlowEdge
	callback func(int, int) bool
}

func NewSpfaMinimumCostFlow() *SpfaMinimumCostFlow {
	return &SpfaMinimumCostFlow{}
}

func (mcf *SpfaMinimumCostFlow) SetCallback(callback func(int, int) bool) {
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
