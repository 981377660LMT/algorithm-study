from typing import List
from functools import lru_cache
from math import sqrt, factorial
from collections import Counter

# 给定一个非负整数数组 A，如果该数组每对相邻元素之和是一个完全平方数，则称这一数组为正方形数组。
# 返回 A 的正方形排列的数目。
# 1 <= A.length <= 12

# 找到图中的哈密尔顿路径数
class Solution:
    def numSquarefulPerms(self, nums: List[int]) -> int:
        def isEdge(x: int, y: int) -> bool:
            r = sqrt(x + y)
            return int(r) ** 2 == (x + y)

        n = len(nums)
        target = (1 << n) - 1
        adjList = [[] for _ in range(n)]
        for i in range(n - 1):
            for j in range(i + 1, n):
                if isEdge(nums[i], nums[j]):
                    adjList[i].append(j)
                    adjList[j].append(i)

        @lru_cache(None)
        def dfs(cur: int, state: int) -> int:
            if state == target:
                return 1

            res = 0
            for next in adjList[cur]:
                if state & (1 << next):
                    continue
                res += dfs(next, state | (1 << next))

            return res

        res = sum(dfs(i, 1 << i) for i in range(n))
        counter = Counter(nums)
        # 消除重复排列
        for val in counter.values():
            res //= factorial(val)
        return res


# print(Solution().numSquarefulPerms([1, 17, 8]))
# 输出：2
# 解释：
# [1,8,17] 和 [17,8,1] 都是有效的排列。
print(Solution().numSquarefulPerms([65, 44, 5, 11]))
