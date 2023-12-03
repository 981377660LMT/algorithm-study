"""
# !请你返回 非空 子序列元素和的最大值，子序列需要满足：
# 子序列中每两个 相邻 的整数 nums[i] 和 nums[j] 
# 它们在原数组中的下标 i 和 j 满足 i < j 且 j - i <= k 。
"""


from collections import deque
from typing import List

INF = int(1e18)


class Solution:
    def constrainedSubsetSum(self, nums: List[int], k: int) -> int:
        """
        dp[i] 表示前 i 个元素中，以第 i 个元素结尾的子序列元素和的最大值

        dp[i] = max(dp[i], max(dp[i - k] ,..., dp[i-1], 0) + nums[i])
        维护单减的单调队列，队首表示之前的子序列元素和最大值

        res = max(dp)
        """
        n = len(nums)
        dp = [-INF] * n
        dp[0] = nums[0]
        queue = deque([(nums[0], 0)])  # (dp[i], i)

        for i in range(1, n):
            while queue and queue[0][1] < i - k:
                queue.popleft()
            dp[i] = max(dp[i], max(0, queue[0][0]) + nums[i])
            while queue and queue[-1][0] <= dp[i]:
                queue.pop()
            queue.append((dp[i], i))

        return max(dp)


if __name__ == "__main__":
    print(Solution().constrainedSubsetSum(nums=[10, 2, -10, 5, 20], k=2))
