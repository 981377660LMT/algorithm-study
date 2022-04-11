# 给定一棵树的边，求经过每条边的路径数

from collections import defaultdict
from typing import Counter, List, Tuple


class Solution:
    def solve(self, edges: List[Tuple[int, int]]) -> List[int]:
        """
        每条不同路径可表示为(start,end)对
        包含的路径数等价于父节点上面*子节点下面连接
        """

        def dfs(cur: int, parent: int) -> int:
            """统计子树结点数"""
            counter[cur] += 1
            for next in adjMap[cur]:
                if next != parent:
                    counter[cur] += dfs(next, cur)
            return counter[cur]

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        counter = Counter()
        dfs(0, -1)

        res = []
        for u, v in edges:
            min_ = min(counter[u], counter[v])
            res.append(min_ * (counter[0] - min_))
        return res


print(Solution().solve(edges=[(0, 1), (1, 2), (0, 3)]))
