"""
给定一棵n(n<=2e5)个节点的树。
有q(q≤2e5)次询问,每次询问给出两个数字(u, k),
请找到距离点u的为k的点(输出任意一个即可)。如果没有输出-1。

q 个查询 询问距离树结点u距离为k的结点是否存在,并输出一个这样的结点

lca在线:
https://zhuanlan.zhihu.com/p/561041985
k的最大值来自于u到直径的两个端点
寻找直径+倍增求树节点的第K个祖先
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


from typing import List, Tuple
from collections import deque


def calDiameter1(adjList: List[List[int]], start: int) -> Tuple[int, Tuple[int, int]]:
    """bfs计算树的直径长度和直径两端点"""
    queue = deque([start])
    visited = set([start])
    last1 = 0  # 第一次BFS最后一个点
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last1 = queue.popleft()
            for next in adjList[last1]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)

    queue = deque([last1])  # 第一次最后一个点作为第二次BFS的起点
    visited = set([last1])
    last2 = 0  # 第二次BFS最后一个点
    res = -1
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last2 = queue.popleft()
            for next in adjList[last2]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        res += 1

    return res, tuple(sorted([last1, last2]))


#######################################################################
from collections import deque
from math import floor, log2
from typing import Iterable, List, Mapping, Sequence, Union


ListTree = Sequence[Iterable[int]]  # Sequence = Iterable + __getitem__
DictTree = Mapping[int, Iterable[int]]
Tree = Union[ListTree, DictTree]


class _LCA:
    def __init__(self, n: int, tree: Tree, root: int) -> None:
        """倍增查询LCA

        `nlogn` 预处理
        `logn`查询

        Args:
            n (int): 树节点编号 默认 0 ~ n-1
            tree (Tree): 树
            root (int): 根节点
        """
        self.depth = [-1] * (n + 1)  # 用深度来排序结点 根节点深度为0 最后一个结点-1表示虚拟结点
        self.parent = [-1] * (n + 1)  # 父节点用来上跳查找路径 根节点父亲为-1

        self._n = n
        self._tree = tree

        self._bfs(root, -1, 0)
        self._bitlen = floor(log2(n)) + 1
        self._fa = self._makeDp(self.parent)

    def queryLCA(self, root1: int, root2: int) -> int:
        """查询树节点两点的最近公共祖先"""
        if self.depth[root1] < self.depth[root2]:
            root1, root2 = root2, root1

        for i in range(self._bitlen - 1, -1, -1):
            if self.depth[self._fa[root1][i]] >= self.depth[root2]:
                root1 = self._fa[root1][i]

        if root1 == root2:
            return root1

        for i in range(self._bitlen - 1, -1, -1):
            if self._fa[root1][i] != self._fa[root2][i]:
                root1 = self._fa[root1][i]
                root2 = self._fa[root2][i]

        return self._fa[root1][0]

    def queryDist(self, root1: int, root2: int) -> int:
        """查询树节点两点间距离"""
        return self.depth[root1] + self.depth[root2] - 2 * self.depth[self.queryLCA(root1, root2)]

    def queryKthAncestor(self, root: int, k: int) -> int:
        """查询树节点root的第k个祖先,如果不存在这样的祖先节点,返回 -1"""
        bit = 0
        while k:
            if k & 1:
                root = self._fa[root][bit]
                if root == -1:
                    return -1
            bit += 1
            k //= 2
        return root

    def _bfs(self, start: int, startPre: int, startDep: int) -> None:
        """处理高度、父节点"""
        queue = deque([(start, startPre, startDep)])
        while queue:
            cur, pre, dep = queue.popleft()
            self.depth[cur], self.parent[cur] = dep, pre
            for next in self._tree[cur]:
                if next != pre:
                    queue.append((next, cur, dep + 1))

    def _makeDp(self, parent: List[int]) -> List[List[int]]:
        """nlogn预处理"""
        fa = [[-1] * self._bitlen for _ in range(self._n)]
        for i in range(self._n):
            fa[i][0] = parent[i]
        for j in range(self._bitlen - 1):
            for i in range(self._n):
                if fa[i][j] == -1:
                    fa[i][j + 1] = -1
                else:
                    fa[i][j + 1] = fa[fa[i][j]][j]
        return fa


if __name__ == "__main__":

    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    lca = _LCA(n, adjList, root=0)
    _, (left, right) = calDiameter1(adjList, start=0)  # 求出直径两端点

    q = int(input())
    for _ in range(q):
        u, k = map(int, input().split())  # 询问距离树结点u距离为k的结点是否存在,并输出一个这样的结点
        u = u - 1
        dist1, dist2 = lca.queryDist(u, left), lca.queryDist(u, right)
        if k > dist1 and k > dist2:
            print(-1)
            continue
        if dist1 > dist2:
            left, right = right, left  # 保证right是离u最远的点
            dist1, dist2 = dist2, dist1
        if k <= lca.queryDist(u, lca.queryLCA(u, right)):  # 在u到right的路径左半部分
            print(lca.queryKthAncestor(u, k) + 1)
        else:
            print(lca.queryKthAncestor(right, dist2 - k) + 1)  # 在u到right的路径右半部分
