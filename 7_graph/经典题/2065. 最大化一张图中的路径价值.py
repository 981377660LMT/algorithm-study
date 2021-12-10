from typing import List
from collections import defaultdict

# 合法路径 指的是图中任意一条从节点 0 开始，最终回到节点 0 ，
# 且花费的总时间 不超过 maxTime 秒的一条路径。
# 你可以访问一个节点任意次。一条合法路径的 价值 定义为路径中 不同节点 的价值 之和


# 注意范围：
# 10 <= timej, maxTime <= 100
# 1 <= n <= 1000
# 每个节点 至多有四条 边。

# maxTime <= 100 and time_j >= 10. It means that we can make no more than 10 steps in our graph
# We have at most 10 steps, and it is also given that each node have at most degree 4, so in total we can make no more than 4^10 states. That is why we will not get TLE.
class Solution:
    def maximalPathQuality(self, values: List[int], edges: List[List[int]], maxTime: int) -> int:
        adjMap = defaultdict(list)
        for x, y, w in edges:
            adjMap[x].append((y, w))
            adjMap[y].append((x, w))

        def dfs(cur, visited, gain, leftTime) -> None:
            if cur == 0:
                self.res = max(self.res, gain)
            for next, weight in adjMap[cur]:
                if weight <= leftTime:
                    dfs(
                        next,
                        visited | set([next]),
                        gain + (next not in visited) * values[next],
                        leftTime - weight,
                    )

        self.res = 0
        dfs(0, set([0]), values[0], maxTime)
        return self.res


print(
    Solution().maximalPathQuality(
        values=[0, 32, 10, 43], edges=[[0, 1, 10], [1, 2, 15], [0, 3, 10]], maxTime=49
    )
)
# 输出：75
# 解释：
# 一条可能的路径为：0 -> 1 -> 0 -> 3 -> 0 。总花费时间为 10 + 10 + 10 + 10 = 40 <= 49 。
# 访问过的节点为 0 ，1 和 3 ，最大路径价值为 0 + 32 + 43 = 75 。
