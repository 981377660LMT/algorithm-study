# 100333. 统计逆序对的数目 (带禁止条件)
# https://leetcode.cn/problems/count-the-number-of-inversions/description/
# 给你一个整数 n 和一个二维数组 requirements ，其中 requirements[i] = [endi, cnti] 表示这个要求中的末尾下标和 逆序对 的数目。
# 请你返回 [0, 1, 2, ..., n - 1] 的 排列perm 的数目，满足对 所有 的 requirements[i] 都有 perm[0..endi] 恰好有 cnti 个逆序对。
# 由于答案可能会很大，将它对 1e9 + 7 取余 后返回。
# 输入保证至少有一个 i 满足 endi == n - 1 。
# 输入保证所有的 endi 互不相同。
# n<=1000
#
# 插入dp
# dp[i][j] = dp[i - 1][j] + dp[i - 1][j - 1] + … + dp[i - 1][0].
# dp[i][j] = 0 if for some x, requirements[x][0] == i and requirements[x][1] != j.

from itertools import accumulate
from typing import List


MOD = int(1e9 + 7)


class Solution:
    def numberOfPermutations(self, n: int, requirements: List[List[int]]) -> int:
        requirements.sort()
        for a, b in zip(requirements, requirements[1:]):
            if a[1] > b[1]:
                return 0
        if requirements[0][0] == 0 and requirements[0][1] != 0:
            return 0

        limit = [-1] * n
        for end, cnt in requirements:
            limit[end] = cnt

        k = requirements[-1][1]
        dp = [0] * (k + 1)
        dp[0] = 1
        for i in range(1, n):
            ndp = dp[:]
            dpSum = [0] + list(accumulate(dp))
            for j in range(1, k + 1):
                ndp[j] = (dpSum[j + 1] - dpSum[max(0, j - i)]) % MOD  # !一共(i+1)项

            if limit[i] != -1:
                for j in range(k + 1):
                    if j != limit[i]:
                        ndp[j] = 0
            dp = ndp

        return dp[k]
