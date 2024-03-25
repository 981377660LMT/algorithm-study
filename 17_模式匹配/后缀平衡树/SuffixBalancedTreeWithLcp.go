// package template.string;

// import java.util.ArrayDeque;
// import java.util.Deque;

// public class SuffixBalancedTreeLcp {
//     private static final double FACTOR = 0.75;
//     private static Node[] stk = new Node[0];
//     private static int tail;
//     public Node root;
//     private ObjectHolder<Node> objectHolder = new ObjectHolder<>();
//     private Deque<Node> dq;

//     private static class ObjectHolder<V> {
//         V data;

//         public void clear() {
//             data = null;
//         }
//     }

//     public SuffixBalancedTreeLcp(int cap) {
//         dq = new ArrayDeque<>(cap + 1);
//         root = Node.NIL;
//         Node dummy = new Node(Integer.MIN_VALUE);
//         dummy.next = dummy;
//         dummy.occur = 0;
//         dummy.offsetToTail = -1;
//         dummy.weight = 0;
//         dq.addFirst(dummy);
//     }

//     private boolean check() {
//         collect(root);
//         for (int i = 1; i < tail; i++) {
//             if (stk[i - 1].weight >= stk[i].weight) {
//                 return false;
//             }
//             if (compare(stk[i - 1], stk[i]) >= 0) {
//                 return false;
//             }
//         }
//         for (int i = 0; i < tail; i++) {
//             if (stk[i].occur < 0 || stk[i].occur > 1) {
//                 return false;
//             }
//         }

//         if (root.aliveSize + 1 != dq.size()) {
//             return false;
//         }

//         return true;
//     }

//     public int lcp(Node a, Node b) {
//         if (a.weight > b.weight) {
//             Node tmp = a;
//             a = b;
//             b = tmp;
//         }
//         return rangeLCPExcludeL(root, 0, 1, a.weight, b.weight);
//     }

//     private int considerLcp(Node a, Node b) {
//         if (a.key != b.key) {
//             return 0;
//         }
//         return 1 + lcp(a.next, b.next);
//     }

//     private void recalcRightLcp(Node prev, Node next) {
//         if (next == Node.NIL) {
//             return;
//         }
//         next.prev = prev;
//         int lcp = considerLcp(prev, next);
//         updateLcp(root, next, lcp);
//     }

//     public Node addPrefix(int x) {
//         objectHolder.clear();
//         root = insert(root, x, dq.peekFirst(), objectHolder, 0, 1);
//         Node node = objectHolder.data;
//         int rank = rank(node);

//         //fix lcp
//         Node prev = rank == 1 ? Node.NIL : kth(root, rank - 1);
//         Node next = rank == root.aliveSize ? Node.NIL : kth(root, rank + 1);
//         recalcRightLcp(prev, node);
//         recalcRightLcp(node, next);

//         dq.addFirst(node);
//         // assert check();
//         return node;
//     }

//     public void removePrefix() {
//         assert dq.size() > 1;
//         Node deleted = dq.removeFirst();
//         int rank = rank(deleted);
//         Node next = rank == root.aliveSize ? Node.NIL : kth(root, rank + 1);

//         //fix lcp
//         if (next != Node.NIL) {
//             int nextLcp = Math.min(next.lcp, deleted.lcp);
//             next.prev = deleted.prev;
//             updateLcp(root, next, nextLcp);
//         }

//         delete(root, deleted);
//         // assert check();

//         //clean or not
//         if (root.aliveSize * 2 < root.size) {
//             collect(root);
//             int wpos = 0;
//             for (int i = 0; i < tail; i++) {
//                 if (stk[i].occur == 0) {
//                     continue;
//                 }
//                 stk[wpos++] = stk[i];
//             }
//             root = refactor(0, wpos - 1, 0, 1);
//         }
//     }

//     public int rank(Node node) {
//         return rank(root, node);
//     }

//     public int leq(IntSequence seq) {
//         return rank(root, seq);
//     }

//     public Node sa(int k) {
//         k++;
//         return kth(root, k);
//     }

//     public int[] sa() {
//         collect(root);
//         int[] sa = new int[size()];
//         int wpos = 0;
//         for (int i = 0; i < tail; i++) {
//             if (stk[i].occur == 0) {
//                 continue;
//             }
//             sa[wpos++] = stk[i].offsetToTail;
//         }
//         return sa;
//     }

//     public int size() {
//         return root.aliveSize;
//     }

//     private static void ensureSpace(int n) {
//         if (stk.length >= n) {
//             return;
//         }
//         int nextSize = Math.max(1 << 16, stk.length);
//         while (nextSize < n) {
//             nextSize += nextSize;
//         }
//         stk = new Node[nextSize];
//     }

//     private static void updateLcp(Node root, Node target, int lcp) {
//         root.pushDown();
//         if (root == target) {
//             root.lcp = lcp;
//         } else {
//             if (root.weight > target.weight) {
//                 updateLcp(root.left, target, lcp);
//             } else {
//                 updateLcp(root.right, target, lcp);
//             }
//         }
//         root.pushUp();
//     }

//     private static int insertCompare(Node a, int key, Node next) {
//         if (a.key != key) {
//             return Integer.compare(a.key, key);
//         }
//         return Double.compare(a.next.weight, next.weight);
//     }

//     private static int compare(Node root, IntSequence seq) {
//         int len = seq.length();
//         for (int i = 0; i < len; i++, root = root.next) {
//             if (seq.get(i) != root.key) {
//                 return Integer.compare(root.key, seq.get(i));
//             }
//         }
//         return 0;
//     }

//     private static int compare(Node a, Node b) {
//         for (int i = 0; a != b; i++, a = a.next, b = b.next) {
//             if (a.key != b.key) {
//                 return Integer.compare(a.key, b.key);
//             }
//         }
//         return 0;
//     }

//     private static int rangeLCPExcludeL(Node root, double L, double R, double l, double r) {
//         if (root == Node.NIL || R <= l || L > r) {
//             return Integer.MAX_VALUE;
//         }
//         if (L > l && R <= r) {
//             return root.rangeMinLCP;
//         }
//         root.pushDown();
//         int ans = Math.min(rangeLCPExcludeL(root.left, L, root.weight, l, r), rangeLCPExcludeL(root.right, root.weight, R, l, r));
//         if (root.occur > 0 && l < root.weight && root.weight <= r) {
//             ans = Math.min(ans, root.lcp);
//         }
//         return ans;
//     }

//     private static Node kth(Node root, int k) {
//         if (root == Node.NIL) {
//             return root;
//         }
//         root.pushDown();
//         Node ans;
//         if (root.left.aliveSize >= k) {
//             ans = kth(root.left, k);
//         } else if (root.left.aliveSize + root.occur >= k) {
//             ans = root;
//         } else {
//             ans = kth(root.right, k - root.left.aliveSize - root.occur);
//         }
//         //push up for calc purpose
//         root.pushUp();
//         return ans;
//     }

//     private static int rank(Node root, IntSequence seq) {
//         if (root == Node.NIL) {
//             return 0;
//         }

//         int ans = 0;
// //        root = refactor(root, L, R);
//         root.pushDown();
//         int compRes = compare(root, seq);
//         if (compRes > 0) {
//             ans += rank(root.left, seq);
//         } else {
//             ans += root.aliveSize - root.right.aliveSize;
//             ans += rank(root.right, seq);
//         }
// //        root.pushUp();
//         return ans;
//     }

//     private static int rank(Node root, Node node) {
//         if (root == Node.NIL) {
//             return 0;
//         }
// //        root = refactor(root, L, R);
//         root.pushDown();
//         int ans = 0;
//         if (root == node) {
//             ans += root.aliveSize - root.right.aliveSize;
//         } else {
//             int compRes = root.compareTo(node);
//             if (compRes > 0) {
//                 ans += rank(root.left, node);
//             } else {
//                 ans += root.aliveSize - root.right.aliveSize;
//                 ans += rank(root.right, node);
//             }
//         }
// //        root.pushUp();
//         return ans;
//     }

//     private static void init(int key, Node root, Node next, double weight) {
//         root.key = key;
//         root.weight = weight;
//         root.next = next;
//         root.occur++;
//         root.offsetToTail = next.offsetToTail + 1;
//         root.lcp = Integer.MAX_VALUE;
//         root.prev = Node.NIL;
//         root.pushUp();
//     }

//     private static Node newNode(int key, Node next, double weight) {
//         Node root = new Node();
//         init(key, root, next, weight);
//         return root;
//     }

//     private static Node insert(Node root, int key, Node next, ObjectHolder<Node> insertNode, double L, double R) {
//         if (root == Node.NIL) {
//             root = newNode(key, next, (L + R) / 2);
//             insertNode.data = root;
//             return root;
//         }
//         root.pushDown();
//         int cmpRes = insertCompare(root, key, next);
//         if (cmpRes == 0) {
//             insertNode.data = root;
//             init(key, root, next, root.weight);
//         } else if (cmpRes > 0) {
//             root.left = insert(root.left, key, next, insertNode, L, root.weight);
//         } else {
//             root.right = insert(root.right, key, next, insertNode, root.weight, R);
//         }
//         root.pushUp();
//         root = refactor(root, L, R);
//         return root;
//     }

//     private static void delete(Node root, Node node) {
//         assert root != Node.NIL;
//         root.pushDown();
//         if (root == node) {
//             root.occur--;
//         } else {
//             int compRes = root.compareTo(node);
//             if (compRes > 0) {
//                 delete(root.left, node);
//             } else {
//                 delete(root.right, node);
//             }
//         }
//         root.pushUp();
//     }

//     private static void collect(Node root) {
//         ensureSpace(root.size);
//         tail = 0;
//         _collect(root);
//         assert tail == root.size;
//     }

//     private static Node refactor(Node root, double L, double R) {
//         double threshold = root.size * FACTOR;
//         if (root.left.size > threshold || root.right.size > threshold) {
//             collect(root);
//             root = refactor(0, tail - 1, L, R);
//         }
//         return root;
//     }

//     private static void _collect(Node root) {
//         if (root == Node.NIL) {
//             return;
//         }
//         root.pushDown();
//         _collect(root.left);
//         stk[tail++] = root;
//         _collect(root.right);
//     }

//     private static Node refactor(int l, int r, double L, double R) {
//         if (l > r) {
//             return Node.NIL;
//         }
//         int m = (l + r) / 2;
//         Node root = stk[m];
//         root.weight = (L + R) / 2;
//         root.left = refactor(l, m - 1, L, root.weight);
//         root.right = refactor(m + 1, r, root.weight, R);
//         root.pushUp();
//         return root;
//     }

//     @Override
//     public String toString() {
//         collect(root);
//         StringBuilder ans = new StringBuilder("{");
//         for (int i = 0; i < tail; i++) {
//             ans.append(stk[i]).append(',');
//         }
//         if (ans.length() > 1) {
//             ans.setLength(ans.length() - 1);
//         }
//         ans.append("}");
//         return ans.toString();
//     }

//     public static class Node implements Cloneable, Comparable<Node> {
//         public static final Node NIL = new Node();

//         Node left = NIL;
//         Node right = NIL;
//         int size;
//         int aliveSize;
//         int key;
//         byte occur;
//         public double weight;
//         Node next;
//         //prev means the floor node
//         Node prev = Node.NIL;
//         public int offsetToTail;
//         public int lcp = Integer.MAX_VALUE;
//         private int rangeMinLCP = Integer.MAX_VALUE;

//         static {
//             NIL.left = NIL.right = NIL;
//             NIL.size = NIL.aliveSize = 0;
//             NIL.key = -1;
//             NIL.offsetToTail = -1;
//         }

//         @Override
//         public int compareTo(Node o) {
//             return Double.compare(weight, o.weight);
//         }

//         public void pushUp() {
//             if (this == NIL) {
//                 return;
//             }
//             size = left.size + right.size + 1;
//             aliveSize = left.aliveSize + right.aliveSize + occur;
//             rangeMinLCP = Math.min(left.rangeMinLCP, right.rangeMinLCP);
//             if (occur > 0) {
//                 rangeMinLCP = Math.min(rangeMinLCP, lcp);
//             }
//         }

//         public void pushDown() {

//         }

//         private Node() {
//         }

//         private Node(int key) {
//             this.key = key;
//             pushUp();
//         }

//         @Override
//         public String toString() {
//             StringBuilder ans = new StringBuilder("[");
//             int remain = 10;
//             Node node = this;
//             for (; node != null && remain > 0; node = node.next == node ? null : node.next, remain--) {
//                 ans.append(node.key).append(',');
//             }
//             if (node != null) {
//                 ans.append(",...,");
//             }
//             if (ans.length() > 1) {
//                 ans.setLength(ans.length() - 1);
//             }
//             ans.append("/").append(this.lcp);
//             ans.append("]");
//             return ans.toString();
//         }
//     }
// }

package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	sbt := NewSuffixBalancedTreeLcp(10)
	sbt.AddPrefix('a')
	sbt.AddPrefix('c')
	sbt.AddPrefix('b')
	sbt.AddPrefix('a')
	fmt.Println(sbt.SaAll()) // abca
}

const ALPHA_NUM int32 = 4
const ALPHA_DENO int32 = 5

type SbtNode struct {
	Weight             float64
	OffsetToTail       int32
	Lcp                int32
	left, right        *SbtNode
	aliveSize, allSize int32
	key                int32
	alive              bool
	prev, next         *SbtNode
	rangeMinLcp        int32
}

func CreateNILNode() *SbtNode {
	res := &SbtNode{}
	res.Lcp = math.MaxInt32
	res.OffsetToTail = -1
	res.left = res
	res.right = res
	res.key = -1
	res.prev = res
	res.rangeMinLcp = math.MaxInt32
	return res
}

var NIL = CreateNILNode()

func NewSbtNode() *SbtNode {
	res := &SbtNode{
		Lcp:  math.MaxInt32,
		left: NIL, right: NIL,
		prev:        NIL,
		rangeMinLcp: math.MaxInt32,
	}
	return res
}

func NewSbtNodeWithKey(key int32) *SbtNode {
	res := &SbtNode{
		Lcp:  math.MaxInt32,
		left: NIL, right: NIL,
		prev:        NIL,
		key:         key,
		rangeMinLcp: math.MaxInt32,
	}
	res.PushUp()
	return res
}

func (n *SbtNode) PushUp() {
	if n == NIL {
		return
	}
	n.allSize = n.left.allSize + n.right.allSize + 1
	n.aliveSize = n.left.aliveSize + n.right.aliveSize
	if n.alive {
		n.aliveSize++
	}
	n.rangeMinLcp = min32(n.left.rangeMinLcp, n.right.rangeMinLcp)
	if n.alive {
		n.rangeMinLcp = min32(n.rangeMinLcp, n.Lcp)
	}
}

func (n *SbtNode) PushDown() {}

func (n *SbtNode) CompareTo(o *SbtNode) int8 {
	if n.Weight < o.Weight {
		return -1
	}
	if n.Weight > o.Weight {
		return 1
	}
	return 0
}

func (n *SbtNode) String() string {
	res := strings.Builder{}
	res.WriteString("[")
	remain := 10
	node := n
	for node != nil && remain > 0 {
		res.WriteString(string(node.key))
		res.WriteString(",")
		node = node.next
		if node == n {
			node = nil
		}
		remain--
	}
	if node != nil {
		res.WriteString(",...,")
	}
	res.WriteString("/")
	res.WriteString(string(n.Lcp))
	res.WriteString("]")
	return res.String()
}

type SuffixBalancedTreeLcp struct {
	Root  *SbtNode
	nodes []*SbtNode

	objectHolder **SbtNode
	collector    []*SbtNode
}

func NewSuffixBalancedTreeLcp(cap int32) *SuffixBalancedTreeLcp {
	nodes := make([]*SbtNode, 0, max32(cap+1, 16))
	root := NIL
	dummy := NewSbtNodeWithKey(math.MinInt32)
	dummy.next = dummy
	dummy.OffsetToTail = -1
	nodes = append(nodes, dummy)
	return &SuffixBalancedTreeLcp{
		Root:         root,
		nodes:        nodes,
		objectHolder: new(*SbtNode),
	}
}

func (sbt *SuffixBalancedTreeLcp) AddPrefix(x int32) *SbtNode {
	sbt.Root = sbt._insert(sbt.Root, x, sbt.nodes[len(sbt.nodes)-1], sbt.objectHolder, 0, 1)
	node := *sbt.objectHolder
	rank := sbt.Rank(node)

	// fix lcp
	var prev, next *SbtNode
	if rank == 1 {
		prev = NIL
	} else {
		prev = sbt._kth(sbt.Root, rank-1)
	}
	if rank == sbt.Root.aliveSize {
		next = NIL
	} else {
		next = sbt._kth(sbt.Root, rank+1)
	}
	sbt._recalcRightLcp(prev, node)
	sbt._recalcRightLcp(node, next)

	sbt.nodes = append(sbt.nodes, node)
	return node
}

func (sbt *SuffixBalancedTreeLcp) RemovePrefix() {
	deleted := sbt.nodes[len(sbt.nodes)-1]
	sbt.nodes = sbt.nodes[:len(sbt.nodes)-1]
	rank := sbt.Rank(deleted)
	var next *SbtNode
	if rank == sbt.Root.aliveSize {
		next = NIL
	} else {
		next = sbt._kth(sbt.Root, rank+1)
	}

	// fix lcp
	if next != NIL {
		nextLcp := min32(next.Lcp, deleted.Lcp)
		next.prev = deleted.prev
		sbt._updateLcp(sbt.Root, next, nextLcp)
	}

	sbt._delete(sbt.Root, deleted)

	// clean or not
	if sbt.Root.aliveSize*2 < sbt.Root.allSize {
		sbt._collect(sbt.Root)
		ptr := int32(0)
		for _, node := range sbt.collector {
			if node.alive {
				sbt.collector[ptr] = node
				ptr++
			}
		}
		sbt.Root = sbt._rebuild(0, ptr-1, 0, 1)
	}
}

func (sbt *SuffixBalancedTreeLcp) Lcp(a, b *SbtNode) int32 {
	if a.Weight > b.Weight {
		a, b = b, a
	}
	return sbt._rangeLcpExcludeL(sbt.Root, 0, 1, a.Weight, b.Weight)
}

func (sbt *SuffixBalancedTreeLcp) Sa(k int32) *SbtNode {
	k++
	return sbt._kth(sbt.Root, k)
}

func (sbt *SuffixBalancedTreeLcp) Rank(node *SbtNode) int32 {
	return sbt._rank(sbt.Root, node)
}

// <=
func (sbt *SuffixBalancedTreeLcp) Leq(n int32, f func(i int32) int32) int32 {
	return sbt._rankSequence(sbt.Root, n, f)
}

func (sbt *SuffixBalancedTreeLcp) SaAll() []int32 {
	sbt._collect(sbt.Root)
	res := make([]int32, sbt.Size())
	ptr := 0
	for _, node := range sbt.collector {
		if node.alive {
			res[ptr] = node.OffsetToTail
			ptr++
		}
	}
	return res
}

func (sbt *SuffixBalancedTreeLcp) Size() int32 {
	return sbt.Root.aliveSize
}

func (sbt *SuffixBalancedTreeLcp) _insert(root *SbtNode, key int32, next *SbtNode, insertNode **SbtNode, L, R float64) *SbtNode {
	if root == NIL {
		root = sbt._newNode(key, next, (L+R)/2)
		*insertNode = root
		return root
	}
	root.PushDown()
	compareRes := sbt._insertCompare(root, key, next)
	if compareRes == 0 {
		*insertNode = root
		sbt._init(key, root, next, root.Weight)
	} else if compareRes > 0 {
		root.left = sbt._insert(root.left, key, next, insertNode, L, root.Weight)
	} else {
		root.right = sbt._insert(root.right, key, next, insertNode, root.Weight, R)
	}
	root.PushUp()
	root = sbt._tryRebuild(root, L, R)
	return root
}

func (sbt *SuffixBalancedTreeLcp) _delete(root *SbtNode, node *SbtNode) {
	root.PushDown()
	if root == node {
		root.alive = false
	} else {
		compareRes := root.CompareTo(node)
		if compareRes > 0 {
			sbt._delete(root.left, node)
		} else {
			sbt._delete(root.right, node)
		}
	}
	root.PushUp()
}

func (sbt *SuffixBalancedTreeLcp) _updateLcp(root *SbtNode, target *SbtNode, lcp int32) {
	root.PushDown()
	if root == target {
		root.Lcp = lcp
	} else {
		if root.Weight > target.Weight {
			sbt._updateLcp(root.left, target, lcp)
		} else {
			sbt._updateLcp(root.right, target, lcp)
		}
	}
	root.PushUp()
}

func (sbt *SuffixBalancedTreeLcp) _rangeLcpExcludeL(root *SbtNode, L, R float64, l, r float64) int32 {
	if root == NIL || R <= l || L > r {
		return math.MaxInt32
	}
	if L > l && R <= r {
		return root.rangeMinLcp
	}
	root.PushDown()
	res := min32(sbt._rangeLcpExcludeL(root.left, L, root.Weight, l, r), sbt._rangeLcpExcludeL(root.right, root.Weight, R, l, r))
	if root.alive && l < root.Weight && root.Weight <= r {
		res = min32(res, root.Lcp)
	}
	return res
}

func (sbt *SuffixBalancedTreeLcp) _considerLcp(a, b *SbtNode) int32 {
	if a.key != b.key {
		return 0
	}
	return 1 + sbt.Lcp(a.next, b.next)
}

func (sbt *SuffixBalancedTreeLcp) _recalcRightLcp(prev, next *SbtNode) {
	if next == NIL {
		return
	}
	next.prev = prev
	lcp := sbt._considerLcp(prev, next)
	sbt._updateLcp(sbt.Root, next, lcp)
}

func (sbt *SuffixBalancedTreeLcp) _kth(root *SbtNode, k int32) (res *SbtNode) {
	if root == NIL {
		return NIL
	}
	root.PushDown()
	if root.left.aliveSize >= k {
		res = sbt._kth(root.left, k)
	} else {
		count := root.left.aliveSize
		if root.alive {
			count++
		}
		if count >= k {
			res = root
		} else {
			res = sbt._kth(root.right, k-count)
		}
	}
	root.PushUp()
	return
}

func (sbt *SuffixBalancedTreeLcp) _rank(root, node *SbtNode) int32 {
	if root == NIL {
		return 0
	}
	root.PushDown()
	if root == node {
		return root.aliveSize - root.right.aliveSize
	} else {
		compareRes := root.CompareTo(node)
		if compareRes > 0 {
			return sbt._rank(root.left, node)
		} else {
			return root.aliveSize - root.right.aliveSize + sbt._rank(root.right, node)
		}
	}
}

func (sbt *SuffixBalancedTreeLcp) _rankSequence(root *SbtNode, n int32, f func(i int32) int32) int32 {
	if root == NIL {
		return 0
	}
	root.PushDown()
	compareRes := root._sequenceCompare(root, n, f)
	if compareRes > 0 {
		return sbt._rankSequence(root.left, n, f)
	} else {
		return root.aliveSize - root.right.aliveSize + sbt._rankSequence(root.right, n, f)
	}
}

func (sbt *SuffixBalancedTreeLcp) _tryRebuild(root *SbtNode, L, R float64) *SbtNode {
	if sbt._isUnbalanced(root) {
		sbt._collect(root)
		root = sbt._rebuild(0, int32(len(sbt.collector)-1), L, R)
	}
	return root
}

func (sbt *SuffixBalancedTreeLcp) _rebuild(l, r int32, L, R float64) *SbtNode {
	if l > r {
		return NIL
	}
	m := (l + r) >> 1
	root := sbt.collector[m]
	root.Weight = (L + R) / 2
	root.left = sbt._rebuild(l, m-1, L, root.Weight)
	root.right = sbt._rebuild(m+1, r, root.Weight, R)
	root.PushUp()
	return root
}

func (sbt *SuffixBalancedTreeLcp) _isUnbalanced(node *SbtNode) bool {
	left, right := node.left, node.right
	// +5，避免不必要的重构
	threshold := node.allSize*ALPHA_NUM + 5*ALPHA_DENO
	return (left.allSize*ALPHA_DENO > threshold) || (right.allSize*ALPHA_DENO > threshold)
}

func (sbt *SuffixBalancedTreeLcp) _collect(node *SbtNode) {
	sbt.collector = sbt.collector[:0]
	sbt._doCollect(node)
}

func (sbt *SuffixBalancedTreeLcp) _doCollect(root *SbtNode) {
	if root == NIL {
		return
	}
	root.PushDown()
	sbt._doCollect(root.left)
	sbt.collector = append(sbt.collector, root)
	sbt._doCollect(root.right)
}

func (sbt *SuffixBalancedTreeLcp) _newNode(key int32, next *SbtNode, weight float64) *SbtNode {
	root := NewSbtNode()
	sbt._init(key, root, next, weight)
	return root
}

func (sbt *SuffixBalancedTreeLcp) _init(key int32, root, next *SbtNode, weight float64) {
	root.key = key
	root.Weight = weight
	root.next = next
	root.alive = true
	root.OffsetToTail = next.OffsetToTail + 1
	root.Lcp = math.MaxInt32
	root.prev = NIL
	root.PushUp()
}

func (sbt *SuffixBalancedTreeLcp) _insertCompare(a *SbtNode, key int32, next *SbtNode) int8 {
	if a.key != key {
		if a.key < key {
			return -1
		}
		return 1
	}
	if a.next.Weight < next.Weight {
		return -1
	}
	if a.next.Weight > next.Weight {
		return 1
	}
	return 0

}
func (sbt *SbtNode) _sequenceCompare(root *SbtNode, n int32, f func(i int32) int32) int8 {
	for i := int32(0); i < n; i++ {
		v := f(i)
		if root.key != v {
			if root.key < v {
				return -1
			}
			return 1
		}
		root = root.next
	}
	return 0
}

func (sbt *SuffixBalancedTreeLcp) _nodeCompare(a, b *SbtNode) int8 {
	for a != b {
		if a.key != b.key {
			if a.key < b.key {
				return -1
			}
			return 1
		}
		a = a.next
		b = b.next
	}
	return 0
}

func (sbt *SuffixBalancedTreeLcp) String() string {
	sbt._collect(sbt.Root)
	res := strings.Builder{}
	res.WriteString("{")
	for i, node := range sbt.collector {
		res.WriteString(node.String())
		if i < len(sbt.collector)-1 {
			res.WriteString(",")
		}
	}
	res.WriteString("}")
	return res.String()
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}
