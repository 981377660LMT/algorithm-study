# 100358. 找出有效子序列的最大长度 II
# https://leetcode.cn/problems/find-the-maximum-length-of-valid-subsequence-ii/description/
# 给你一个整数数组 nums 和一个 正 整数 k 。
# nums 的一个 子序列 sub 的长度为 x ，如果其满足以下条件，则称其为 有效子序列 ：
# (sub[0] + sub[1]) % k == (sub[1] + sub[2]) % k == ... == (sub[x - 2] + sub[x - 1]) % k
# 返回 nums 的 最长有效子序列 的长度。
# n<=1e3,k<=1e3
# !所有奇数位置的数都相等，所有偶数位置的数都相等 (交叉dp、奇偶dp)
#
# 解法1：枚举子序列的最后两项，dp[i][j]表示以i,j结尾的最长有效子序列长度
# !解法2：枚举相邻两项模k的余数, dp[cur] = dp[(m - cur) % k] + 1

from typing import List


class Solution:
    def maximumLength1(self, nums: List[int], k: int) -> int:
        nums = [v % k for v in nums]
        dp = [[0] * k for _ in range(k)]
        for cur in nums:
            for pre, f in enumerate(dp[cur]):
                dp[pre][cur] = f + 1
        return max(map(max, dp))

    def maximumLength2(self, nums: List[int], k: int) -> int:
        nums = [v % k for v in nums]
        res = 0
        for m in range(k):
            dp = [0] * k
            for cur in nums:
                dp[cur] = dp[(m - cur) % k] + 1
            res = max(res, max(dp))
        return res
