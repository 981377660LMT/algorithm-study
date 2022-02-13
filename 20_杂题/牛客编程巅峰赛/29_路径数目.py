from collections import defaultdict
from typing import List


class Point:
    def __init__(self, a=0, b=0):
        self.x = a
        self.y = b


class Solution:
    def solve(self, n: int, Edge: List[Point], f: List[int]) -> int:
        """在最多经过两个值为1的节点的情况下，有多少条到达叶节点的路径？"""
        if not Edge:
            return 1

        def dfs(cur: int, pre: int, oneCount: int) -> None:
            if oneCount >= 3:
                return
            if cur in leaves:
                self.res += 1
                return
            for next in adjMap[cur]:
                if pre == next:
                    continue
                dfs(next, cur, oneCount + int(f[next - 1] == 1))

        adjMap = defaultdict(list)
        for edge in Edge:
            u, v = edge.x, edge.y
            adjMap[u].append(v)
            adjMap[v].append(u)
        leaves = set([i for i, nexts in adjMap.items() if len(nexts) == 1]) - {1}

        self.res = 0
        dfs(1, -1, int(f[0] == 1))
        return self.res


print(
    Solution().solve(
        5, [Point(*pair) for pair in [(2, 5), (5, 3), (5, 4), (5, 1)]], [0, 0, 0, 1, 1]
    )
)
