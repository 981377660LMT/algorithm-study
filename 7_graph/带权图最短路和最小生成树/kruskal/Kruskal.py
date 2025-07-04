from typing import List, Tuple


def kruskal(n: int, edges: List[Tuple[int, int, int]]) -> Tuple[int, List[bool], bool]:
    """
    Kruskal算法求无向图最小生成树(森林).

    返回值:
    - mstCost: 最小生成树(森林)的权值之和
    - inMst: 是否在最小生成树(森林)中
    - isTree: 是否是树
    """
    uf = UnionFindArraySimple(n)
    count = 0
    mstCost, inMst, isTree = 0, [False] * len(edges), False
    order = sorted(range(len(edges)), key=lambda x: edges[x][2])
    for ei in order:
        u, v, w = edges[ei]
        if uf.union(u, v):
            inMst[ei] = True
            mstCost += w
            count += 1
            if count == n - 1:
                isTree = True
                break
    return mstCost, inMst, isTree


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(self, key1: int, key2: int) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


if __name__ == "__main__":
    # 3600. 升级后最大生成树稳定性
    # https://leetcode.cn/problems/maximize-spanning-tree-stability-with-upgrades/description/
    # !在选出必选边后，根据 Kruskal 算法求最大生成树，把剩余的边按照边权（先不乘 2）从大到小合并。
    class Solution:
        def maxStability(self, n: int, edges: List[List[int]], k: int) -> int:
            uf = UnionFindArraySimple(n)
            allUf = UnionFindArraySimple(n)
            res = int(1e18)
            for x, y, w, must in edges:
                if must:
                    if not uf.union(x, y):  # 必选边成环
                        return -1
                    res = min(res, w)
                allUf.union(x, y)

            if allUf.part > 1:  # 图不连通
                return -1
            if uf.part == 1:  # 只需选必选边
                return res

            # 最大生成树
            edges.sort(key=lambda x: -x[2])
            nonMustEdgeWeights = []
            for x, y, w, must in edges:
                if not must and uf.union(x, y):
                    nonMustEdgeWeights.append(w)

            res = min(res, nonMustEdgeWeights[-1] * 2)
            if k < len(nonMustEdgeWeights):
                res = min(res, nonMustEdgeWeights[-1 - k])
            return res

    # https://www.luogu.com.cn/problem/P3366
    # P3366 【模板】最小生成树
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        edges.append((u - 1, v - 1, w))
    mstCost, _, isTree = kruskal(n, edges)
    print(mstCost if isTree else "orz")
