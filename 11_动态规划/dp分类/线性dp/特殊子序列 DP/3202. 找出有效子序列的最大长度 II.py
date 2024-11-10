# 3202. 找出有效子序列的最大长度 II
# https://leetcode.cn/problems/find-the-maximum-length-of-valid-subsequence-ii/description/
# 给你一个整数数组 nums 和一个 正整数 k 。
# nums 的一个
# 子序列
#  sub 的长度为 x ，如果其满足以下条件，则称其为 有效子序列 ：
#
# (sub[0] + sub[1]) % k == (sub[1] + sub[2]) % k == ... == (sub[x - 2] + sub[x - 1]) % k
# 返回 nums 的 最长有效子序列 的长度。
#
#
# 2 <= nums.length <= 1e3
# 1 <= nums[i] <= 1e7
# 1 <= k <= 1e3
# https://leetcode.cn/problems/find-the-maximum-length-of-valid-subsequence-ii/solutions/2826591/deng-jie-zhuan-huan-dong-tai-gui-hua-pyt-z2fs/
# 枚举余数，考察子序列的最后一项。
# 从左到右遍历 nums 的同时，维护一个数组 dp[x]，表示最后一项模 k 为 x 的子序列的长度。

from typing import List


class Solution:
    def maximumLength(self, nums: List[int], k: int) -> int:
        res = 0
        for mod in range(k):
            dp = [0] * k
            for num in nums:
                dp[num % k] = dp[(mod - num) % k] + 1
            res = max(res, max(dp))
        return res
