# https://github.com/EndlessCheng/codeforces-go/blob/04c3e9c75964d0a5dd701db92e5c3fdf0bc353bc/copypasta/dp.go#L3227
# https://ac.nowcoder.com/acm/contest/10167/C?&headNav=acm
# 树上所有路径的位运算与(&)的和
# 单个点也算路径
# 解法：对每一位，统计仅含 1 的路径个数


from collections import defaultdict
from typing import DefaultDict, List, Tuple


# !解法1：对每一位，统计仅含 1 的路径个数
def addPathSum1(n: int, edges: List[Tuple[int, int]], values: List[int]) -> int:
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    maxLog = max(values).bit_length()
    res = 0
    for bit in range(maxLog):
        pathCount = 0

        def dfs(cur: int, pre: int) -> int:
            nonlocal pathCount
            one = (values[cur] >> bit) & 1
            pathCount += one
            for next_ in adjList[cur]:
                if next_ != pre:
                    nextRes = dfs(next_, cur)
                    if one > 0:
                        pathCount += one * nextRes
                        one += nextRes
            return one

        dfs(0, -1)
        res += (1 << bit) * pathCount

    return res


# !解法2：对每一位，用并查集求出 1 组成的连通分量，每个连通分量对答案的贡献是 sz*(sz+1)/2
def addPathSum2(n: int, edges: List[Tuple[int, int]], values: List[int]) -> int:
    class UF:
        __slots__ = ("n", "part", "_parent", "_rank")

        def __init__(self, n: int):
            self.n = n
            self.part = n
            self._parent = list(range(n))
            self._rank = [1] * n

        def find(self, x: int) -> int:
            while self._parent[x] != x:
                self._parent[x] = self._parent[self._parent[x]]
                x = self._parent[x]
            return x

        def union(self, x: int, y: int) -> bool:
            """按秩合并."""
            rootX = self.find(x)
            rootY = self.find(y)
            if rootX == rootY:
                return False
            if self._rank[rootX] > self._rank[rootY]:
                rootX, rootY = rootY, rootX
            self._parent[rootX] = rootY
            self._rank[rootY] += self._rank[rootX]
            self.part -= 1
            return True

        def getGroups(self) -> DefaultDict[int, List[int]]:
            groups = defaultdict(list)
            for key in range(self.n):
                root = self.find(key)
                groups[root].append(key)
            return groups

        def getSize(self, x: int) -> int:
            return self._rank[self.find(x)]

    maxLog = max(values).bit_length()
    res = 0
    for bit in range(maxLog):
        uf = UF(n)
        for u, v in edges:
            if (values[u] >> bit) & 1 and (values[v] >> bit) & 1:
                uf.union(u, v)
        groups = uf.getGroups()
        for leader, group in groups.items():
            if (values[leader] >> bit) & 1:
                size = len(group)
                res += (1 << bit) * size * (size + 1) // 2
    return res


if __name__ == "__main__":
    # https://ac.nowcoder.com/acm/contest/10167/C?&headNav=acm
    class Solution:
        def solve(self, n, u, v, p):
            edges = [(u[i], v[i]) for i in range(n - 1)]
            return addPathSum2(n, edges, p)
