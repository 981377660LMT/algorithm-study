# https://leetcode.cn/problems/collect-coins-in-a-tree/
# 给你一个长度为 n 的数组 coins ，其中 coins[i] 可能为 0 也可能为 1 ，1 表示节点 i 处有一个金币。
# 一开始，你需要选择树中任意一个节点出发。你可以执行下述操作任意次：
#  - 收集距离当前节点距离为 2 以内的所有金币，或者
#  - 移动到树中一个相邻节点。
# !你需要收集树中所有的金币，并且`回到出发节点`，请你返回最少经过的边数。
# 如果你多次经过一条边，每一次经过都会给答案加一。


# 找子图中心点的做法:
# 1. dfs 找到硬币们形成的子图 newTree;
# 2. dp 找到 newTree 中到各个硬币距离最小的点, 记为center;
# 3. 以 center 为 root 建立LCA, 对每个硬币找到向 root 跳两步的点，标记为必经点;
# 4. 对这些必经点 dfs 再求一次子图，边数*2 就是答案。
# !如果可以不回到出发点,那么就是最终的子图的边数(边权和)*2-直径。


from typing import List, Tuple
from Rerooting import Rerooting

INF = int(1e18)


class Solution:
    def collectTheCoins(self, coins: List[int], edges: List[List[int]]) -> int:
        n = len(coins)
        specials = [i for i in range(n) if coins[i] == 1]
        if len(specials) <= 1:
            return 0
        tree = [[] for _ in range(n)]
        for u, v in edges:
            tree[u].append(v)
            tree[v].append(u)
        newTree1, _, visited1 = getTreeSubGraph(n, tree, specials)
        maxDist = getDistToSpecials(n, tree, specials)  # 注意这里是原树
        maxDist = [(i, v) for i, v in enumerate(maxDist) if visited1[i]]
        center = min(maxDist, key=lambda x: x[1])[0]
        lca = _L(n, newTree1, center)
        criticals = [False] * n
        for v in specials:
            up2 = lca.jump(v, center, 2)
            if up2 != -1:
                criticals[up2] = True
        criticals = [i for i in range(n) if criticals[i]]
        _, newEdges, _ = getTreeSubGraph(n, tree, criticals)
        return 2 * len(newEdges)  # !如果不需要返回起点,就要减去最后图中的直径


def getDistToSpecials(n: int, rawTree: List[List[int]], specials: List[int]) -> List[int]:
    """求`原树`中每个点到特殊点的最远距离."""
    isSpecial = [False] * n
    for v in specials:
        isSpecial[v] = True

    E = int

    def e(root: int) -> E:
        return 0 if isSpecial[root] else -INF

    def op(childRes1: E, childRes2: E) -> E:
        return max(childRes1, childRes2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return fromRes + 1

    R = Rerooting(n)
    for cur in range(n):
        for next in rawTree[cur]:
            if cur < next:
                R.addEdge(cur, next)
    dp = R.rerooting(e, op, composition)
    return dp


def getTreeDiameter(n: int, tree: List[List[int]], start=0) -> Tuple[int, List[int]]:
    """求无权树的(直径长度,直径路径)."""

    def dfs(start: int) -> Tuple[int, List[int]]:
        dist = [-1] * n
        dist[start] = 0
        stack = [start]
        while stack:
            cur = stack.pop()
            for next in tree[cur]:
                if dist[next] != -1:
                    continue
                dist[next] = dist[cur] + 1
                stack.append(next)
        endPoint = dist.index(max(dist))
        return endPoint, dist

    u, _ = dfs(start)
    v, dist = dfs(u)
    diameter = dist[v]
    path = [v]
    while u != v:
        for next in tree[v]:
            if dist[next] + 1 == dist[v]:
                path.append(next)
                v = next
                break

    return diameter, path


def getTreeSubGraph(
    n: int, rawTree: List[List[int]], specials: List[int]
) -> Tuple[List[List[int]], List[Tuple[int, int]], List[bool]]:
    """给定`原树`(无向图邻接表), 返回`子图(无向图邻接表), 子图中的边, 原图中每个点是否在子图中`."""
    if len(specials) == 0:
        return [[] for _ in range(n)], [], [False] * n
    visited = [False] * n
    for v in specials:
        visited[v] = True

    def dfs(cur: int, pre: int) -> bool:
        for next in rawTree[cur]:
            if next != pre and dfs(next, cur):
                visited[cur] = True  # 标记必经点
        return visited[cur]

    root = specials[0]
    dfs(root, -1)

    edges, newTree = [], [[] for _ in range(n)]
    for cur in range(n):
        for next in rawTree[cur]:
            if cur < next and visited[cur] and visited[next]:
                edges.append((cur, next))
                newTree[cur].append(next)
                newTree[next].append(cur)
    return newTree, edges, visited


class _L:
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


# [1,0,0,0,0,1]

# [[0,1],[1,2],[2,3],[3,4],[4,5]]

print(
    Solution().collectTheCoins(
        coins=[1, 0, 0, 0, 0, 1], edges=[[0, 1], [1, 2], [2, 3], [3, 4], [4, 5]]
    )
)
