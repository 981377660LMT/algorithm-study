from typing import List
from collections import deque


INF = int(1e20)

# 1 <= n <= 2 * 104
# nums.length == n
# 1 <= nums[i] <= 50
# edges.length == n - 1
# edges[i].length == 2
# 0 <= edges[i][0], edges[i][1] <= n - 1
# edges 表示一棵合法的树。

# 一个连通块的 价值 定义为这个连通块中 所有 节点 i 对应的 nums[i] 之和。
# 你需要删除一些边，删除后得到的各个连通块的价值都相等。请返回你可以删除的边数 最多 为多少。

# !1. 注意到50很小 => 枚举因子
# !2. dfs或者拓扑排序统计子树和(树的拓扑排序等于后序dfs)


def getFactors(n: int) -> List[int]:
    """n 的所有因数 O(sqrt(n))"""
    if n <= 0:
        return []
    small, big = [], []
    upper = int(n**0.5) + 1
    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


class Solution:
    def componentValue(self, nums: List[int], edges: List[List[int]]) -> int:
        def check1(groupSum: int) -> bool:
            """是否可以将所有连通块的和变为groupSum

            从叶子结点拓扑排序
            """
            curDeg = deg[:]
            queue, visited = deque(), [False] * n
            for i in range(n):
                if curDeg[i] <= 1:  # !叶子节点以及孤立的节点
                    queue.append(i)
                    visited[i] = True

            dp = nums[:]  # 记录每个节点的子树和
            while queue:
                cur = queue.popleft()
                if dp[cur] > groupSum:
                    return False

                for next in adjList[cur]:
                    if not visited[next]:
                        if dp[cur] < groupSum:
                            dp[next] += dp[cur]
                        curDeg[next] -= 1
                        if curDeg[next] == 1:
                            visited[next] = True
                            queue.append(next)

            return True

        def check2(groupSum: int) -> bool:
            """是否可以将所有连通块的和变为groupSum

            从根后序dfs
            """

            def dfs(cur: int, pre: int) -> int:
                """返回每个子树的价值和"""
                curSum = nums[cur]
                for next in adjList[cur]:
                    if next == pre:
                        continue
                    nextSum = dfs(next, cur)
                    if nextSum > groupSum:
                        return INF
                    elif nextSum < groupSum:
                        curSum += nextSum
                return curSum

            rootSum = dfs(0, -1)
            return rootSum == groupSum

        n = len(nums)
        adjList = [[] for _ in range(n)]
        deg = [0] * n
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
            deg[u] += 1
            deg[v] += 1

        sum_ = sum(nums)
        factors = getFactors(sum_)
        for groupSum in factors:
            # if check1(groupSum):
            if check2(groupSum):
                return sum_ // groupSum - 1
        return 0


# print(Solution().componentValue(nums=[6, 2, 2, 2, 6], edges=[[0, 1], [1, 2], [1, 3], [3, 4]]))
# print(Solution().componentValue(nums=[2], edges=[]))
# print(Solution().componentValue(nums=[6, 2, 4, 6], edges=[[0, 1], [1, 2], [2, 3]]))
print(Solution().componentValue(nums=[1, 1, 2, 1, 1], edges=[[0, 2], [1, 2], [3, 2], [4, 2]]))
