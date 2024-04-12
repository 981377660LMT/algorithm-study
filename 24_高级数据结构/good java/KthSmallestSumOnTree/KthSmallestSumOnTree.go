// package template.graph;

// import java.util.*;

// /**
//  * 给定一颗树，每条边上有一个非负权重，要求从根出发，找到路径权重最小的k条路径
//  */
// public class KthSmallestSumOnTree {
//     /**
//      * O(klog k)
//      *
//      * @param root
//      * @param k
//      */
//     public static List<State> kthSmallestSumOnTree(Vertex root, int k) {
//         List<State> ans = new ArrayList<>(k);
//         PriorityQueue<State> pq = new PriorityQueue<>(2 * k, Comparator.comparingLong(x -> x.sum));
//         pq.add(new State(root, null, 0));
//         while (!pq.isEmpty() && ans.size() < k) {
//             State state = pq.remove();
//             ans.add(state);
//             //child or bro
//             if (state.iterator.hasNext()) {
//                 Edge e = state.iterator.next();
//                 pq.add(new State(e.to, state, state.sum + e.weight));
//             }
//             if (state.parent != null && state.parent.iterator.hasNext()) {
//                 Edge e = state.parent.iterator.next();
//                 pq.add(new State(e.to, state.parent, state.parent.sum + e.weight));
//             }
//         }
//         return ans;
//     }

//     /**
//      * O(klog k)
//      *
//      * @param root
//      * @param k
//      */
//     public List<EdgeState> kthSmallestSumOnTreeWithEdge(Vertex root, int k) {
//         List<EdgeState> ans = new ArrayList<>(k);
//         PriorityQueue<EdgeState> pq = new PriorityQueue<>(2 * k, Comparator.comparingLong(x -> x.sum));
//         pq.add(new EdgeState(root, null, 0, null));
//         while (!pq.isEmpty() && ans.size() < k) {
//             EdgeState state = pq.remove();
//             ans.add(state);
//             //child or bro
//             if (state.iterator.hasNext()) {
//                 Edge e = state.iterator.next();
//                 pq.add(new EdgeState(e.to, state, state.sum + e.weight, e));
//             }
//             if (state.parent != null && state.parent.iterator.hasNext()) {
//                 Edge e = state.parent.iterator.next();
//                 pq.add(new EdgeState(e.to, state.parent, state.parent.sum + e.weight, e));
//             }
//         }
//         return ans;
//     }

//     public static class State {
//         public Vertex v;
//         Iterator<Edge> iterator;
//         public State parent;
//         public long sum;

//         public State(Vertex v, State parent, long sum) {
//             this.v = v;
//             this.parent = parent;
//             this.sum = sum;
//             iterator = v.children();
//         }

//         public State(Vertex v, Iterator<Edge> iterator, State parent, long sum) {
//             this.v = v;
//             this.iterator = iterator;
//             this.parent = parent;
//             this.sum = sum;
//         }
//     }

//     public static class EdgeState extends State {
//         public Edge edge;

//         public EdgeState(Vertex v, State parent, long sum, Edge e) {
//             super(v, parent, sum);
//             this.edge = e;
//         }
//     }

//     public interface Vertex {
//         Iterator<Edge> children();
//     }

//     public static class Edge {
//         public Vertex to;
//         public long weight;
//     }
// }

package main

func main() {
	//    0
	//   / \
	//  1   2
	//     / \
	//    3   4
	//       / \
	//      5   6

	n := int32(7)
	tree := make([]*vertex, n)
}

// 给定一颗树，每条边上有一个非负权重，要求从根出发，找到路径权重最小的k条路径
func KthSmallestSumOnTree(root *vertex, k int32) []*state {
	res := make([]*state, 0, k)
	pq := NewHeap(func(a, b *state) bool { return a.Sum < b.Sum }, nil)
	pq.Push(newState(root, nil, 0))
	for pq.Len() > 0 && int32(len(res)) < k {
		state := pq.Pop()
		res = append(res, state)
		// child
		if state.iterator.HasNext() {
			e := state.iterator.Next()
			pq.Push(newState(e.To, state, state.Sum+int(e.Weight)))
		}
		// bro
		if state.Parent != nil && state.Parent.iterator.HasNext() {
			e := state.Parent.iterator.Next()
			pq.Push(newState(e.To, state.Parent, state.Parent.Sum+int(e.Weight)))
		}
	}
	return res
}

// 遍历的一个状态.
type state struct {
	Cur      *vertex
	Parent   *state
	Sum      int
	iterator *edgeIterator
}

func newState(cur *vertex, parent *state, sum int) *state {
	return &state{Cur: cur, Parent: parent, Sum: sum, iterator: cur.Children()}
}

func newStateWithIterator(cur *vertex, parent *state, sum int, iterator *edgeIterator) *state {
	return &state{Cur: cur, Parent: parent, Sum: sum, iterator: iterator}
}

type edge struct {
	To     *vertex
	Weight int
}

type vertex struct {
	Children func() *edgeIterator
}

type edgeIterator struct {
	Next    func() *edge
	HasNext func() bool
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
