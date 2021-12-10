from typing import List
from collections import defaultdict

# 请你帮助重新规划路线方向，使每个城市都可以访问城市 0 。返回需要变更方向的最小路线数。

# 总结:
# 1.当作无向图处理
# 2.从原点dfs遍历，看有多少对(cur,parent)在道路中
class Solution:
    def minReorder(self, n: int, connections: List[List[int]]) -> int:
        self.res = 0
        roads = set()
        adjMap = defaultdict(list)
        for u, v in connections:
            roads.add((u, v))
            adjMap[u].append(v)
            adjMap[v].append(u)

        def dfs(cur: int, parent: int) -> None:
            self.res += (parent, cur) in roads
            for next in adjMap[cur]:
                if next == parent:
                    continue
                dfs(next, cur)

        dfs(0, -1)

        return self.res


print(Solution().minReorder(n=6, connections=[[0, 1], [1, 3], [2, 3], [4, 0], [4, 5]]))
# 输出：3
# 解释：更改以红色显示的路线的方向，使每个城市都可以到达城市 0 。
