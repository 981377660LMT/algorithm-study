# 3277. 查询子数组最大异或值
# https://leetcode.cn/problems/maximum-xor-score-subarray-queries/solutions/2899932/qu-jian-dp-tao-qu-jian-dppythonjavacgo-b-w4be/
# !注意这里的异或不是一般的异或，而是dp[i][j] = dp[i][j - 1] ^ dp[i + 1][j]，这样的异或

from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maximumSubarrayXor(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        n = len(nums)
        dp = [[0] * n for _ in range(n)]  # 下标从 i 到 j 的子数组的「数组的异或值」
        dpMax = [[0] * n for _ in range(n)]
        for i in range(n - 1, -1, -1):
            dpMax[i][i] = dp[i][i] = nums[i]
            for j in range(i + 1, n):
                dp[i][j] = dp[i][j - 1] ^ dp[i + 1][j]
                dpMax[i][j] = max2(dp[i][j], max2(dpMax[i + 1][j], dpMax[i][j - 1]))
        return [dpMax[l][r] for l, r in queries]
