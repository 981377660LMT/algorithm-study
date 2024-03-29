# 从3*n长的数组选2*n个数，求(前n个数的和-后n个数的和)的最小值
# n<=10^5

# 1.枚举分割点
# 2.最小值=左侧最小-右侧最大
from typing import List
from heapq import heappop, heappush

INF = int(1e20)


class Solution:
    def minimumDifference(self, nums: List[int]) -> int:
        n = len(nums) // 3
        preMin = [0] * (3 * n + 1)
        sufMax = [0] * (3 * n + 1)

        pq = []
        for i in range(2 * n):
            preMin[i + 1] = preMin[i] + nums[i]
            heappush(pq, -nums[i])
            if len(pq) > n:
                preMin[i + 1] -= -heappop(pq)

        pq = []
        for i in range(3 * n - 1, n - 1, -1):
            sufMax[i] = sufMax[i + 1] + nums[i]
            heappush(pq, nums[i])
            if len(pq) > n:
                sufMax[i] -= heappop(pq)

        res = INF
        for i in range(n, 2 * n + 1):
            cur = preMin[i] - sufMax[i]
            res = min(res, cur)
        return res


print(Solution().minimumDifference(nums=[3, 1, 2]))
print(Solution().minimumDifference(nums=[7, 9, 5, 8, 1, 3]))
