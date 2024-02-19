# 1218. 最长定差子序列
# https://leetcode.cn/problems/longest-arithmetic-subsequence-of-given-difference/
# 找出并返回 arr 中最长等差子序列的长度，该子序列中相邻元素之间的差等于 difference 。

from typing import List
from collections import defaultdict


class Solution:
    def longestSubsequence(self, arr: List[int], diff: int) -> int:
        dp = defaultdict(int)  # 前 i 个数（第 i 个数必选）时，得到的最长定差子序列长度
        for num in arr:
            dp[num] = dp[num - diff] + 1
        return max(dp.values(), default=0)


if __name__ == "__main__":
    print(Solution().longestSubsequence([-2, 0, 3, 6, 1, 9], 3))
    # We can pick the subsequence [0, 3, 6, 9].
