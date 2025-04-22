# 3511. 构造正数组(子数组和)
# https://leetcode.cn/problems/make-a-positive-array/description/
#
# 给定一个数组 nums。一个数组被认为是`正`的，如果每个包含 `超过两个` 元素的子数组的所有数字之和都是正数。
# 您可以执行以下操作任意多次：
# 用 -1e18 和 1e18 之间的任意整数替换 nums 中的 一个 元素。
# 找到使 nums 变为`正数组`所需的最小操作数。
#
# n<=1e5
# -1e9<=nums[i]<=1e9
#
# !令 dp[i] 表示以位置 i 结尾的，长度至少为 3 的连续子数组的最小和。
# !dp[i] = min(dp[i-1] + nums[i], nums[i-2] + nums[i-1] + nums[i])
# !因此只需考虑两种情况：一直延续的子数组、长度为3的子数组

from typing import List
from itertools import accumulate

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def makeArrayPositive(self, nums: List[int], k=3) -> int:
        n = len(nums)
        res = 0
        dp = INF
        preSum = list(accumulate(nums, initial=0))
        start = k - 1
        while start < n:
            sum1 = dp + nums[start]
            sum2 = preSum[start + 1] - preSum[start + 1 - k]
            dp = min2(sum1, sum2)
            if dp > 0:
                start += 1
            else:
                dp = INF
                res += 1
                start += k
        return res


if __name__ == "__main__":
    nums = [-1, -1, -5, -1, -1]
    print(Solution().makeArrayPositive(nums))
