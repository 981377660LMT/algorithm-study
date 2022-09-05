import sys
from typing import List
from collections import defaultdict
from itertools import accumulate

sys.setrecursionlimit(100000)


class Point:
    def __init__(self, a=0, b=0):
        self.x = a
        self.y = b


class Solution:
    def solve(self, n: int, edges: List[Point], q: int, Query: List[Point]) -> List[int]:
        def dfs(cur: int, parent: int) -> None:
            """求dfs序"""
            nonlocal id
            start[cur] = id
            for next in adjMap[cur]:
                if next == parent:
                    continue
                dfs(next, cur)
            end[cur] = id
            id += 1

        # 每次操作将root根节点的子树上所有权值加x
        # 求q次操作后1-n每个结点的权值
        adjMap = defaultdict(list)
        for edge in edges:
            u, v = edge.x, edge.y
            adjMap[u].append(v)
            adjMap[v].append(u)

        start = [-1] * (n + 1)  # 第一个子树叶子节点
        end = [-1] * (n + 1)  # dfs序
        id = 1
        dfs(1, -1)

        diff = [0] * (n + 10)
        for query in Query:
            root, delta = query.x, query.y
            start_, end_ = start[root], end[root]
            diff[start_] += delta
            diff[end_ + 1] -= delta
        diff = list(accumulate(diff))

        res = []
        for i in range(1, n + 1):
            dfsId = end[i]
            res.append(diff[dfsId])
        return res


print(
    Solution().solve(
        5, [Point(2, 5), Point(5, 3), Point(5, 4), Point(5, 1)], 2, [Point(1, 3), Point(2, -1)]
    )
)
