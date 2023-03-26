# 给你 n 个城市，编号为从 1 到 n 。同时给你一个大小为 n-1 的数组 edges ，
# 其中 edges[i] = [ui, vi] 表示城市 ui 和 vi 之间有一条双向边。
# 题目保证任意城市之间只有唯一的一条路径。换句话说，所有城市形成了一棵 树 。
# 对于 d 从 1 到 n-1 ，请你找到城市间 最大距离 恰好为 d 的所有`子树`数目。
# 请你返回一个大小为 n-1 的数组，其中第 d 个元素（下标从 1 开始）是城市间 最大距离 恰好等于 d 的子树数目。
# 2 <= n <= 15 可以枚举子集

from collections import defaultdict
from itertools import combinations
from typing import List


class Solution:
    def countSubgraphsForEachDiameter(self, n: int, edges: List[List[int]]) -> List[int]:
        """O(n^3)枚举直径端点+乘法原理
        https://leetcode.cn/problems/count-subtrees-with-max-distance-between-cities/solution/tu-jie-on3-mei-ju-zhi-jing-duan-dian-che-am2n/
        https://leetcode.cn/problems/count-subtrees-with-max-distance-between-cities/solution/shi-xian-hen-jian-dan-yuan-li-lue-you-xie-fu-za-de/
        """

        def dfs(cur: int, pre: int, v1: int, v2: int, d: int) -> int:
            count = 1
            for next in adjList[cur]:
                if next != pre:
                    if (dist[v1][next] < d or (dist[v1][next] == d and next > v2)) and (
                        dist[v2][next] < d or (dist[v2][next] == d and next > v1)
                    ):
                        count *= dfs(next, cur, v1, v2, d)
            if dist[v1][cur] + dist[v2][cur] > d:  # cur是可选点
                count += 1
            return count

        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u - 1].append(v - 1)
            adjList[v - 1].append(u - 1)

        dist = [[0] * n for _ in range(n)]
        lca = _LC(n, adjList, 0)
        for u in range(n):
            for v in range(u + 1, n):
                d = lca.queryDist(u, v)
                dist[u][v] = dist[v][u] = d

        res = [0] * (n - 1)
        for u in range(n):
            for v in range(u + 1, n):
                d = dist[u][v]
                res[d - 1] += dfs(u, -1, u, v, d)
        return res

    def countSubgraphsForEachDiameter2(self, n: int, edges: List[List[int]]) -> List[int]:
        """
        O(2^n) 枚举子树

        1.求每个点到所有点的最短距离--多源最短路径算法 floyd O(n^3)
        2.枚举子集看哪些是子树
        """

        dist = defaultdict(lambda: defaultdict(lambda: int(1e20)))
        for i in range(n):
            dist[i][i] = 0
        for u, v in edges:
            dist[u - 1][v - 1] = 1
            dist[v - 1][u - 1] = 1

        for k in range(n):
            for i in range(n):
                for j in range(n):
                    cand = dist[i][k] + dist[k][j]
                    dist[i][j] = cand if dist[i][j] > cand else dist[i][j]

        res = [0] * n
        for state in range(1, 1 << n):
            nodes = [i for i in range(n) if (state >> i) & 1]
            edgeCount = sum(dist[n1][n2] == 1 for n1, n2 in combinations(nodes, 2))
            if len(nodes) == edgeCount + 1:
                maxDist = max((dist[n1][n2] for n1, n2 in combinations(nodes, 2)), default=0)
                res[maxDist] += 1

        return res[1:]


class _LC:
    def __init__(self, n: int, tree: List[List[int]], root: int) -> None:
        """倍增查询LCA

        `nlogn` 预处理
        `logn`查询

        Args:
            n (int): 树节点编号 默认 0 ~ n-1
            tree (Tree): 树
            root (int): 根节点
        """
        self.depth = [0] * n  # 根节点深度为0
        self.parent = [0] * n  # 根节点父亲为-1

        self._n = n
        self._tree = tree

        self._dfs(root, -1, 0)
        self._bitlen = n.bit_length()
        self.dp = self._makeDp()  # !dp[i][j] 表示节点j的第2^i个祖先

    def queryLCA(self, root1: int, root2: int) -> int:
        """查询树节点两点的最近公共祖先"""
        if self.depth[root1] < self.depth[root2]:
            root1, root2 = root2, root1
        root1 = self.upToDepth(root1, self.depth[root2])
        if root1 == root2:
            return root1
        for i in range(self._bitlen - 1, -1, -1):
            if self.dp[i][root1] != self.dp[i][root2]:
                root1 = self.dp[i][root1]
                root2 = self.dp[i][root2]
        return self.dp[0][root1]

    def queryDist(self, root1: int, root2: int) -> int:
        """查询树节点两点间距离"""
        return self.depth[root1] + self.depth[root2] - 2 * self.depth[self.queryLCA(root1, root2)]

    def queryKthAncestor(self, root: int, k: int) -> int:
        """
        查询树节点root的第k个祖先(0-indexed)
        如果不存在这样的祖先节点,返回 -1
        """
        bit = 0
        while k:
            if k & 1:
                root = self.dp[bit][root]
                if root == -1:
                    return -1
            bit += 1
            k //= 2
        return root

    def jump(self, start: int, target: int, step: int) -> int:
        """
        从start节点跳到target节点,跳过step个节点(0-indexed)
        返回跳到的节点,如果不存在这样的节点,返回-1
        """
        lca = self.queryLCA(start, target)
        dep1, dep2, deplca = self.depth[start], self.depth[target], self.depth[lca]
        dist = dep1 + dep2 - 2 * deplca
        if step > dist:
            return -1
        if step <= dep1 - deplca:
            return self.queryKthAncestor(start, step)
        return self.queryKthAncestor(target, dist - step)

    def upToDepth(self, root: int, toDepth: int) -> int:
        """从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点"""
        if toDepth >= self.depth[root]:
            return root
        for i in range(self._bitlen - 1, -1, -1):
            if (self.depth[root] - toDepth) & (1 << i):
                root = self.dp[i][root]
        return root

    def _dfs(self, start: int, startPre: int, startDep: int) -> None:
        """处理高度、父节点"""
        depth, parent, tree = self.depth, self.parent, self._tree
        stack = [(start, startPre, startDep)]
        while stack:
            cur, pre, dep = stack.pop()
            depth[cur], parent[cur] = dep, pre
            for next in tree[cur]:
                if next != pre:
                    stack.append((next, cur, dep + 1))

    def _makeDp(self) -> List[List[int]]:
        """nlogn预处理"""
        n, bitlen, parent = self._n, self._bitlen, self.parent
        dp = [[-1] * n for _ in range(bitlen)]
        for j in range(self._n):
            dp[0][j] = parent[j]
        for i in range(bitlen - 1):
            for j in range(n):
                if dp[i][j] == -1:
                    dp[i + 1][j] = -1
                else:
                    dp[i + 1][j] = dp[i][dp[i][j]]
        return dp


print(Solution().countSubgraphsForEachDiameter(n=4, edges=[[1, 2], [2, 3], [2, 4]]))
