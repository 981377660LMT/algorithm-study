from typing import List
from sortedcontainers import SortedList


# To calculate dp[i], only need to check at most two previous elements:
# 1. the latest element greater than nums[i], say nums[x].
# 2. the latest element smaller than nums[i],
# say nums[y]. Then dp[i] = max(dp[x], dp[y]) + 1
# n ≤ 1e5

# 相邻差绝对值不大于k的最长子序列
class Solution:
    def solve(self, nums: List[int], k: int) -> int:
        """思路是用一个有序数据结构维护之前的位置"""
        if not nums:
            return 0

        n = len(nums)
        dp = [1] * n
        list1, list2 = SortedList(), SortedList()

        for i, num in enumerate(nums):
            while True:
                # ceiling 找最右
                pos = list1.bisect_right((num,)) - 1
                if not 0 <= pos < len(list1):
                    break
                val, index = list1[pos]
                if abs(val - num) > k:
                    break
                list1.pop(pos)
                dp[i] = max(dp[i], dp[index] + 1)

            while True:
                # floor 找最左
                pos = list2.bisect_left((num,))
                if not 0 <= pos < len(list2):
                    break
                val, index = list2[pos]
                if abs(val - num) > k:
                    break
                list2.pop(pos)
                dp[i] = max(dp[i], dp[index] + 1)

            list1.add((num, i))
            list2.add((num, i))

        return max(dp)


print(Solution().solve(nums=[2, 3, -1, -2, 8, 0, 1], k=3))
