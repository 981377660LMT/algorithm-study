# 类似 3589. 计数质数间隔平衡子数组
# https://leetcode.cn/problems/count-prime-gap-balanced-subarrays/solutions/3705577/yu-chu-li-zhi-shu-hua-dong-chuang-kou-da-v23w/
# 给你一个整数数组 nums 和一个整数 k。你的任务是将 nums 分割成一个或多个 非空 的连续子段，使得每个子段的 最大值 与 最小值 之间的差值 不超过 k。
#
# 返回在此条件下将 nums 分割的总方法数。
#
# 由于答案可能非常大，返回结果需要对 109 + 7 取余数。
#
# 划分型dp.
# !dp[i] = 以 nums[0:i] 为前缀的分割方式数.
# !答案就是 dp[n]-dp[n-1].

from typing import List

from sortedcontainers import SortedList

MOD = int(1e9 + 7)


class Solution:
    def countPartitions(self, nums: List[int], k: int) -> int:
        n = len(nums)
        dp = [0] * (n + 1)
        dp[0] = 1
        sl = SortedList()
        left = 0
        for right, x in enumerate(nums, 1):
            sl.add(x)
            while sl[-1] - sl[0] > k:  # type: ignore
                sl.remove(nums[left])
                left += 1
            dp[right] = (dp[right - 1] - dp[left - 1]) % MOD
            dp[right] += dp[right - 1]
        return (dp[n] - dp[n - 1]) % MOD
