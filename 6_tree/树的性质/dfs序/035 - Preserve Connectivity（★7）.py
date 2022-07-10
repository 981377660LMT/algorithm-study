"""!连接k个树结点需要保留的最小边数

n<=1e5

多个点 按照dfs序 求两两间的距离
"""

from collections import defaultdict
import sys
from math import floor, log2
from typing import DefaultDict, List, Set
from collections import defaultdict
from typing import DefaultDict, Set, Tuple


class DFSOrder:
    def __init__(self, n: int, tree: DefaultDict[int, Set[int]]) -> None:
        """dfs序

        Args:
            n (int): 树节点从0开始,根节点为0
            tree (DefaultDict[int, Set[int]]): 无向图邻接表
        """
        self.starts = [0] * n
        self.ends = [0] * n
        self._n = n
        self._tree = tree
        self._dfsId = 1
        self._dfs(0, -1)

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


class LCAManager:
    def __init__(self, n: int, adjMap: DefaultDict[int, Set[int]], root=0) -> None:
        """查询 LCA

        `nlogn` 预处理
        `logn`查询两点的LCA

        Args:
            n (int): 树节点编号 默认 0 ~ n-1
            adjMap (DefaultDict[int, Set[int]]): 树
            root (int, optional): 根节点 默认 0
        """
        self.depth = defaultdict(lambda: -1)
        self.parent = defaultdict(lambda: -1)
        self._BITLEN = floor(log2(n)) + 1
        self._MAX = n
        self._adjMap = adjMap
        self._dfs(root, -1, 0)
        self._fa = self._makeDp(self.parent)

    def queryLCA(self, root1: int, root2: int) -> int:
        """ `logn` 查询 """
        if self.depth[root1] < self.depth[root2]:
            root1, root2 = root2, root1

        for i in range(self._BITLEN - 1, -1, -1):
            if self.depth[self._fa[root1][i]] >= self.depth[root2]:
                root1 = self._fa[root1][i]

        if root1 == root2:
            return root1

        for i in range(self._BITLEN - 1, -1, -1):
            if self._fa[root1][i] != self._fa[root2][i]:
                root1 = self._fa[root1][i]
                root2 = self._fa[root2][i]

        return self._fa[root1][0]

    def queryDist(self, root1: int, root2: int) -> int:
        """查询树节点两点间距离"""
        return self.depth[root1] + self.depth[root2] - 2 * self.depth[self.queryLCA(root1, root2)]

    def _dfs(self, cur: int, pre: int, dep: int) -> None:
        """处理高度、父节点信息"""
        self.depth[cur], self.parent[cur] = dep, pre
        for next in self._adjMap[cur]:
            if next == pre:
                continue
            self._dfs(next, cur, dep + 1)

    def _makeDp(self, parent: DefaultDict[int, int]) -> List[List[int]]:
        """nlogn预处理"""
        dp = [[0] * self._BITLEN for _ in range(self._MAX)]
        for i in range(self._MAX):
            dp[i][0] = parent[i]
        for j in range(self._BITLEN - 1):
            for i in range(self._MAX):
                if dp[i][j] == -1:
                    dp[i][j + 1] = -1
                else:
                    dp[i][j + 1] = dp[dp[i][j]][j]
        return dp


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())
adjMap = defaultdict(set)
for _ in range(n - 1):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    adjMap[u].add(v)
    adjMap[v].add(u)

D = DFSOrder(n, adjMap)
L = LCAManager(n, adjMap)

q = int(input())
for _ in range(q):
    k, *nodes = list(map(int, input().split()))
    nodes = [int(v) - 1 for v in nodes]
    nodes.sort(key=lambda x: D.queryId(x))
    res = 0
    for pre, cur in zip(nodes, nodes[1:] + [nodes[0]]):
        res += L.queryDist(pre, cur)
    print(res // 2)

