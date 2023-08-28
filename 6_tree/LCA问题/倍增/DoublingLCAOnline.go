// !动态维护k级祖先，可以动态添加叶子节点(在树上做递归/回溯操作时比较有用)
// 作者：cup-pyy
// 链接：https://zhuanlan.zhihu.com/p/650275363
// 来源：知乎
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

// #include<iostream>
// #include<cstring>
// #include<vector>
// using namespace std;
// using LL = long long;
// const int maxn = 1e6 + 5;
// vector<int> g[maxn], query[maxn];
// int fa[20][maxn];
// int ans[maxn], a[maxn], cnt[maxn], sum;

// void dfs(int u){
//     if (u != 1){
//         if (++cnt[a[u]] == 1){
//             sum += 1;
//         }
//     }
//     for(auto x : query[u]) ans[x] = sum;
//     for(auto j : g[u]) dfs(j);
//     if (u != 1){
//         if (--cnt[a[u]] == 0){
//             sum -= 1;
//         }
//     }
// }

// int main(){

//     cin.tie(0);
//     cout.tie(0);
//     ios::sync_with_stdio(0);

//     int n;
//     cin >> n;
//     int tot = 1, cur = 1, qs = 0;
//     vector<int> ops;
//     ops.reserve(n);
//     ops.push_back(cur);
//     for(int i = 0; i < n; i++){
//         char op;
//         cin >> op;
//         if (op == '+'){
//             int x;
//             cin >> x;
//             a[++tot] = x;
//             fa[0][tot] = cur;
//             for(int j = 1; j <= 19; j++)
//                 fa[j][tot] = fa[j - 1][fa[j - 1][tot]];
//             g[cur].push_back(tot);
//             cur = tot;
//             ops.push_back(cur);
//         }
//         else if (op == '-'){
//             int k;
//             cin >> k;
//             for(int j = 19; j >= 0; j--){
//                 if (k >> j & 1){
//                     cur = fa[j][cur];
//                 }
//             }
//             ops.push_back(cur);
//         }
//         else if (op == '!'){
//             ops.pop_back();
//             cur = ops.back();
//         }
//         else{
//             query[cur].push_back(++qs);
//         }
//     }
//     dfs(1);
//     for(int i = 1; i <= qs; i++)
//         cout << ans[i] << '\n';
//
// }

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	operations := make([][2]int, q)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		switch op {
		case "+":
			var x int
			fmt.Fscan(in, &x)
			operations[i] = [2]int{1, x}
		case "-":
			var k int
			fmt.Fscan(in, &k)
			operations[i] = [2]int{2, k}
		case "!":
			operations[i] = [2]int{3, 0}
		case "?":
			operations[i] = [2]int{4, 0}
		}
	}

	res := Rollbacks(operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 给定一个初始时为空的数组nums, 需要实现下面四种类型的操作：
// [1, x]: 将x添加到nums尾部
// [2, k]: 将尾部的k个数删除.保证存在k个数.
// [3, 0]: 撤销上一次操作1或2操作
// [4, 0]: 查询当前nums中有多少个不同的数
//
// 1<=q<=1e6,询问次数不超过1e5
//
// !1. 因为要支持撤销，所以需要`保存版本`或者action, 例如immer.js两种方式都会提供.
// 这里保存版本比较合适。不同版本构成了一棵树.
// !2. 离线方法是建立一棵版本树，+操作加边，-操作通过倍增上跳到对应节点，
// 操作回退到上个节点，?操作记录当前节点需要记录答案，最后dfs整棵树求解.
// !3. 这个倍增也非常巧妙，不预处理而是是动态的
func Rollbacks(operations [][2]int) []int {}

// 不预先给出整棵树,而是动态添加叶子节点,维护树节点的LCA和k级祖先.
// 初始时只有一个根节点0.
type DoublingLCAOnline struct {
	n      int
	bitLen int
	dp     [][]int
	depth  []int
}

func NewDoublingLCAOnline(n int) *DoublingLCAOnline {
	bit := bits.Len(uint(n))
	dp := make([][]int, bit)
	for i := range dp {
		cur := make([]int, n)
		for j := range cur {
			cur[j] = -1
		}
		dp[i] = cur
	}
	depth := make([]int, n)
	return &DoublingLCAOnline{n: n, bitLen: bit, dp: dp, depth: depth}
}

// 在树中添加一条从parent到child的边.要求parent已经存在于树中.
func (lca *DoublingLCAOnline) AddEdge(parent, child int) {
	if parent != 0 && lca.depth[parent] == 0 {
		panic(fmt.Sprintf("parent %d not exists", parent))
	}
	lca.depth[child] = lca.depth[parent] + 1
	lca.dp[0][child] = parent
	for i := 0; i < lca.bitLen-1; i++ {
		pre := lca.dp[i][child]
		if pre == -1 {
			break
		}
		lca.dp[i+1][child] = lca.dp[i][pre]
	}
}

// 查询节点node的第k个祖先(向上跳k步).如果不存在,返回-1.
func (lca *DoublingLCAOnline) KthAncestor(node, k int) int {
	if k > lca.depth[node] {
		return -1
	}
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			node = lca.dp[bit][node]
			if node == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return node
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=depth[v],返回跳到的节点.
func (lca *DoublingLCAOnline) UpToDepth(root, toDepth int) int {
	if toDepth >= lca.depth[root] {
		return root
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if (lca.depth[root]-toDepth)&(1<<i) > 0 {
			root = lca.dp[i][root]
		}
	}
	return root
}

func (lca *DoublingLCAOnline) LCA(root1, root2 int) int {
	if lca.depth[root1] < lca.depth[root2] {
		root1, root2 = root2, root1
	}
	root1 = lca.UpToDepth(root1, lca.depth[root2])
	if root1 == root2 {
		return root1
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root1] != lca.dp[i][root2] {
			root1 = lca.dp[i][root1]
			root2 = lca.dp[i][root2]
		}
	}
	return lca.dp[0][root1]
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *DoublingLCAOnline) Jump(start, target, step int) int {
	lca_ := lca.LCA(start, target)
	dep1, dep2, deplca := lca.depth[start], lca.depth[target], lca.depth[lca_]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return lca.KthAncestor(start, step)
	}
	return lca.KthAncestor(target, dist-step)
}
