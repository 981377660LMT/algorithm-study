"""
每一步，你最多可以往前跳 k 步
你的目标是到达数组最后一个位置（下标为 n - 1 ），你的 得分 为经过的所有数字之和
请你返回你能得到的 最大得分 。
1425. 带限制的子序列和

队首维护当前最大值:最大堆/单调队列
"""
from collections import deque
from typing import List

INF = int(1e18)


class Solution:
    def maxResult(self, nums: List[int], k: int) -> int:
        """
        dp[i]表示到达第i个位置的最大得分
        dp[i] = max(dp[i - k], ..., dp[i - 1]) + nums[i]
        维护单减的单调队列，队首表示之前的最大值
        """
        n = len(nums)
        dp = [-INF] * n
        dp[0] = nums[0]
        queue = deque([(nums[0], 0)])  # (dp[i], i)
        for i in range(1, n):
            while queue and queue[0][1] < i - k:
                queue.popleft()
            dp[i] = max(dp[i], queue[0][0] + nums[i])
            while queue and queue[-1][0] <= dp[i]:
                queue.pop()
            queue.append((dp[i], i))
        return dp[-1]
