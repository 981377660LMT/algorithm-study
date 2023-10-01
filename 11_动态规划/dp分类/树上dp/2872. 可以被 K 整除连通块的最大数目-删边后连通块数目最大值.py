# 你可以从树中删除一些边，也可以一条边也不删，得到若干连通块。
# 一个 连通块的值 定义为连通块中所有节点值之和。
# 如果所有连通块的值都可以被 k 整除，那么我们说这是一个 合法分割 。
# 请你返回所有合法分割中，连通块数目的最大值 。


# 贪心，如果可以组被k整除，就马上切断边


from typing import List


class Solution:
    def maxKDivisibleComponents(
        self, n: int, edges: List[List[int]], values: List[int], k: int
    ) -> int:
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        res = 0

        def dfs(cur: int, pre: int) -> int:
            curSum = values[cur]
            for next in adjList[cur]:
                if next == pre:
                    continue
                curSum += dfs(next, cur)

            nonlocal res
            res += (curSum % k) == 0
            return curSum

        dfs(0, -1)
        return res
