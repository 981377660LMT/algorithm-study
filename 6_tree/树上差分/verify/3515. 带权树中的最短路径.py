# 3515. 带权树中的最短路径
# https://leetcode.cn/problems/shortest-path-in-a-weighted-tree/description/
#
# 给你一个整数 n 和一个以节点 1 为根的无向带权树，该树包含 n 个编号从 1 到 n 的节点。
# 它由一个长度为 n - 1 的二维数组 edges 表示，其中 edges[i] = [ui, vi, wi] 表示一条从节点 ui 到 vi 的无向边，权重为 wi。
# 同时给你一个二维整数数组 queries，长度为 q，其中每个 queries[i] 为以下两种之一：
# [1, u, v, w'] – 更新 节点 u 和 v 之间边的权重为 w'，其中 (u, v) 保证是 edges 中存在的边。
# [2, x] – 计算 从根节点 1 到节点 x 的 最短 路径距离。
# 返回一个整数数组 answer，其中 answer[i] 是对于第 i 个 [2, x] 查询，从节点 1 到 x 的最短路径距离。
#
# !链加、查询带根节点的链和 -> 子树加、单点查询 -> 树状数组差分加、前缀查询


from typing import List, Union, Tuple


class BITArray:
    __slots__ = "n", "total", "_data"

    def __init__(self, sizeOrData: Union[int, List[int]]):
        if isinstance(sizeOrData, int):
            self.n = sizeOrData
            self.total = 0
            self._data = [0] * sizeOrData
        else:
            self.n = len(sizeOrData)
            self.total = sum(sizeOrData)
            _data = sizeOrData[:]
            for i in range(1, self.n + 1):
                j = i + (i & -i)
                if j <= self.n:
                    _data[j - 1] += _data[i - 1]
            self._data = _data

    def add(self, index: int, value: int) -> None:
        self.total += value
        index += 1
        while index <= self.n:
            self._data[index - 1] += value
            index += index & -index

    def queryPrefix(self, end: int) -> int:
        if end > self.n:
            end = self.n
        res = 0
        while end > 0:
            res += self._data[end - 1]
            end &= end - 1
        return res

    def queryRange(self, start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return 0
        if start == 0:
            return self.queryPrefix(end)
        pos, neg = 0, 0
        while end > start:
            pos += self._data[end - 1]
            end &= end - 1
        while start > end:
            neg += self._data[start - 1]
            start &= start - 1
        return pos - neg

    def queryAll(self):
        return self.total

    def __repr__(self):
        res = [self.queryRange(i, i + 1) for i in range(self.n)]
        return f"BITArray({res})"


def dfsPreOrder(tree: List[List[int]], root=0) -> Tuple[List[int], List[int]]:
    """前序遍历dfs序.

    # !data[lid[i]] = values[i]
    """
    n = len(tree)
    lid, rid = [0] * n, [0] * n
    dfn = 0

    def dfs(cur: int, pre: int) -> None:
        nonlocal dfn
        lid[cur] = dfn
        dfn += 1
        for next_ in tree[cur]:
            if next_ != pre:
                dfs(next_, cur)
        rid[cur] = dfn

    dfs(root, -1)
    return lid, rid


class Solution:
    def treeQueries(self, n: int, edges: List[List[int]], queries: List[List[int]]) -> List[int]:
        adjList = [[] for _ in range(n)]
        for u, v, _ in edges:
            u, v = u - 1, v - 1
            adjList[u].append(v)
            adjList[v].append(u)

        lid, rid = dfsPreOrder(adjList)
        bit = BITArray(n)
        weight = [0] * n

        def update(u: int, v: int, w: int) -> None:
            if lid[u] > lid[v]:
                u, v = v, u
            delta = w - weight[v]
            weight[v] = w
            bit.add(lid[v], delta)
            bit.add(rid[v], -delta)

        def query(u: int) -> int:
            return bit.queryPrefix(lid[u] + 1)

        for e in edges:
            update(e[0] - 1, e[1] - 1, e[2])

        res = []
        for q in queries:
            if q[0] == 1:
                update(q[1] - 1, q[2] - 1, q[3])
            else:
                res.append(query(q[1] - 1))
        return res
