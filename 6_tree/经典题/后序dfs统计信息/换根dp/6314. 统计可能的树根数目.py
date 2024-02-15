from collections import defaultdict
from random import randint, randrange
from typing import List
from Rerooting import Rerooting
from LCA import LCA_HLD


# 2581. 统计可能的树根数目
# https://leetcode.cn/problems/count-number-of-possible-root-nodes/description/
# Alice 有一棵 n 个节点的树，节点编号为 0 到 n - 1 。
# 树用一个长度为 n - 1 的二维整数数组 edges 表示，其中 edges[i] = [ai, bi] ，表示树中节点 ai 和 bi 之间有一条边。
# Alice 想要 Bob 找到这棵树的根。她允许 Bob 对这棵树进行若干次 猜测 。每一次猜测，Bob 做如下事情：
# 选择两个 不相等 的整数 u 和 v ，且树中必须存在边 [u, v] 。
# Bob 猜测树中 u 是 v 的 父节点 。
# !Bob 的猜测用二维整数数组 guesses 表示，其中 guesses[j] = [uj, vj] 表示 Bob 猜 uj 是 vj 的`父节点`。
# !Alice 非常懒，她不想逐个回答 Bob 的猜测，只告诉 Bob 这些猜测里面 至少 有 k 个猜测的结果为 true 。
# 给你二维整数数组 edges ，Bob 的所有猜测和整数 k ，请你返回可能成为树根的 节点数目 。如果没有这样的树，则返回 0。
E = int


class Solution:
    def rootCount(self, edges: List[List[int]], guesses: List[List[int]], k: int) -> int:
        def e(root: int) -> E:
            return 0

        def op(childRes1: E, childRes2: E) -> E:
            return childRes1 + childRes2

        def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            if direction == 0:  # cur -> parent
                return fromRes + counter[(parent, cur)]
            return fromRes + counter[(cur, parent)]  # parent -> cur

        counter = defaultdict(int)
        for a, b in guesses:
            counter[(a, b)] += 1
        n = len(edges) + 1
        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u, v)
        dp = R.rerooting(e, op, composition)
        return sum(x >= k for x in dp)

    def rootCount2(self, edges: List[List[int]], guesses: List[List[int]], k: int) -> List[int]:
        """`猜测父结点`改成`猜测祖先结点`的做法"""
        E = int

        def e(root: int) -> E:
            return 0

        def op(childRes1: E, childRes2: E) -> E:
            return childRes1 + childRes2

        def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            if direction == 0:  # cur -> parent
                return fromRes + counter[parent][cur]  # !有多少个查询 (parent-?) 存在于 (parent-cur) 所在链的子树中
            return fromRes + counter[cur][parent]  # !有多少个查询(cur - ?) 存在于 (cur->parent) 的子树中

        n = len(edges) + 1
        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u, v)
        lca = LCA_HLD(n)
        for u, v in edges:
            lca.addEdge(u, v, 1)
        lca.build(0)
        counter = [defaultdict(int) for _ in range(n)]  # 有多少个查询存在于链 (i -> key) 的子树中
        for a, b in guesses:
            next = lca.jump(a, b, 1)
            counter[a][next] += 1
        dp = R.rerooting(e, op, composition)
        return dp


if __name__ == "__main__":
    from collections import defaultdict
    from typing import Iterable, Mapping, Sequence, Tuple, Union

    AdjList = Sequence[Iterable[int]]
    AdjMap = Mapping[int, Iterable[int]]
    Tree = Union[AdjList, AdjMap]

    class _DFSOrder:
        __slots__ = ("starts", "ends", "_n", "_tree", "_dfsId")

        def __init__(self, n: int, tree: Tree, root=0) -> None:
            """dfs序

            Args:
                n (int): 树节点从0开始,根节点为0
                tree (Tree): 无向图邻接表

            1. 按照dfs序遍历k个结点形成的回路 每条边恰好经过两次
            """
            self.starts = [0] * n
            self.ends = [0] * n

            self._n = n
            self._tree = tree
            self._dfsId = 1

            self._dfs(root, -1)

        def queryRange(self, root: int) -> Tuple[int, int]:
            """求子树映射到的区间

            Args:
                root (int): 根节点
            Returns:
                Tuple[int, int]: [start, end] 1 <= start <= end <= n
            """
            return self.starts[root], self.ends[root]

        def queryId(self, root: int) -> int:
            """求root自身的dfsId

            Args:
                root (int): 根节点
            Returns:
                int: id  1 <= id <= n
            """
            return self.ends[root]

        def isAncestor(self, root: int, child: int) -> bool:
            """判断root是否是child的祖先

            Args:
                root (int): 根节点
                child (int): 子节点

            应用:枚举边时给树的边定向
            ```
            if not D.isAncestor(e[0], e[1]):
                e[0], e[1] = e[1], e[0]
            ```
            """
            left1, right1 = self.starts[root], self.ends[root]
            left2, right2 = self.starts[child], self.ends[child]
            return left1 <= left2 <= right2 <= right1

        def _dfs(self, cur: int, pre: int) -> None:
            self.starts[cur] = self._dfsId
            for next in self._tree[cur]:
                if next == pre:
                    continue
                self._dfs(next, cur)
            self.ends[cur] = self._dfsId
            self._dfsId += 1

    # check 2 with brute force
    def rootCount2BruteForce(edges: List[List[int]], guesses: List[List[int]], k: int) -> List[int]:
        n = len(edges) + 1
        dp = [0] * n  # 以每个结点为根的树满足条件的结点数
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
        for i in range(n):
            D = _DFSOrder(n, adjList, i)
            for a, b in guesses:
                if D.isAncestor(a, b):
                    dp[i] += 1
        return dp

    # 产生一个随机的无根树(0-n-1)
    def randomTree(n: int) -> List[List[int]]:
        if n <= 2:
            g = [[] for _ in range(n)]
            if n == 2:
                g[0].append(1)
                g[1].append(0)
            return g
        prufer = [randint(1, n) for _ in range(n - 2)]
        deg = [0] * (n + 1)
        for p in prufer:
            deg[p] += 1
        prufer = prufer + [n]
        parents = [0] * (n + 1)
        i, j = 0, 1
        while i < n - 1:
            while deg[j] > 0:
                j += 1
            parents[j] = prufer[i]
            while i < n - 2:
                p = prufer[i]
                deg[p] -= 1
                if p > j or deg[p] > 0:
                    break
                parents[p] = prufer[i + 1]
                i += 1
            i += 1
            j += 1
        parents = parents[1:]
        tree = [[] for _ in range(n)]
        for i in range(1, n):
            p = parents[i - 1]
            tree[i - 1].append(p - 1)
            tree[p - 1].append(i - 1)
        return tree

    for _ in range(10):
        n = randrange(2, 500)
        tree = randomTree(n)
        edges = []
        for i in range(n):
            for j in tree[i]:
                if i < j:
                    edges.append([i, j])
        guess = []
        while len(guess) < 3:
            a, b = randrange(n), randrange(n)
            if a != b:
                guess.append([a, b])
        k = randrange(n)
        if not rootCount2BruteForce(edges, guess, k) == Solution().rootCount2(edges, guess, k):
            print(n, tree, guess, k, "aasas")
            print(rootCount2BruteForce(edges, guess, k), Solution().rootCount2(edges, guess, k))
            exit(0)
    print("ok")
