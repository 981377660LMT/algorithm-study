# 对每一个询问，求有多少个长为k的窗口包含询问值

from collections import defaultdict


class Solution:
    def solve(self, nums, queries, k):
        def getWays(query: int) -> int:
            """求有多少个长为k的窗口包含query"""
            res = 0
            preRight = -1
            for index in indexMap[query]:
                # 窗口左边界的范围在left到right之间取
                left = max(index - k + 1, preRight + 1)
                right = min(index, n - k)
                res += right - left + 1
                preRight = right
            return res

        n = len(nums)
        indexMap = defaultdict(list)
        for i, num in enumerate(nums):
            indexMap[num].append(i)

        return list(map(getWays, queries))


print(Solution().solve(nums=[2, 1, 2, 3, 4], queries=[2, 1], k=3))
