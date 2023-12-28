# https://leetcode.cn/problems/find-maximum-non-decreasing-array-length/
# 2945. 找到最大非递减数组的长度
# 给你一个下标从 0 开始的整数数组 nums 。
# 你可以执行任意次操作。
# 每次操作中，你需要选择一个 子数组 ，并将这个子数组用它所包含元素的 和 替换。
# 比方说，给定数组是 [1,3,5,6] ，你可以选择子数组 [3,5] ，用子数组的和 8 替换掉子数组，然后数组会变为 [1,8,6] 。
# 请你返回执行任意次操作以后，可以得到的 最长非递减 数组的长度。
# 子数组 指的是一个数组中一段连续 非空 的元素序列。
#
# !求出前缀和数组之后问题就可以转化为：
# 在这个前缀和数组中寻找最长上升子序列，这个最长上升子序列中第i个元素减去第i-1个元素，应该小于等于第i+1个元素减去第i个元素
# !即:preSum[i+1]>=2*preSum[i]-preSum[i-1]
# 推广到一般情况:
# !preSum[k]>=2*preSum[i]-preSum[j] (j<i<k)
# 每次二分查找第一个满足上式的k即可

from bisect import bisect_left
from itertools import accumulate
from math import ceil
from typing import List

INF = int(1e18)


# https://leetcode.cn/problems/find-maximum-non-decreasing-array-length/solutions/2570645/xian-duan-shu-shang-er-fen-by-mike-meng-m8tf/
class Solution:
    def findMaximumLength(self, nums: List[int]) -> int:
        n = len(nums)
        preSum = [0] + list(accumulate(nums))
        dp = [0] * (n + 1)  # dp[i] 表示以第i个元素结尾的最长非递减数组的长度
        pre = [0] * (n + 1)  # pre[i] 表示以第i个元素结尾的最长非递减数组的前一个元素的下标最大值
        for i in range(1, n + 1):
            pre[i] = max(pre[i], pre[i - 1])
            dp[i] = max(dp[i], dp[pre[i]] + 1)
            pos = bisect_left(preSum, 2 * preSum[i] - preSum[pre[i]])
            if pos < n + 1:
                pre[pos] = i
        return dp[-1]


# 反向操作
# 2366. 将数组排序的最少替换次数
# https://leetcode.cn/problems/minimum-replacements-to-sort-the-array/
# 给你一个下标从 0 开始的整数数组 nums 。每次操作中，你可以将数组中任何一个元素替换为 任意两个 和为该元素的数字。
# 比方说，nums = [5,6,7] 。一次操作中，我们可以将 nums[1] 替换成 2 和 4 ，将 nums 转变成 [5,2,4,7] 。
# 请你执行上述操作，将数组变成元素按 非递减 顺序排列的数组，并返回所需的最少操作次数。
class Solution2:
    def minimumReplacement(self, nums: List[int]) -> int:
        n = len(nums)
        res = 0
        pre = INF
        for i in range(n - 1, -1, -1):
            cur = nums[i]
            if cur > pre:
                count = ceil(cur / pre)
                max_ = cur // count
                res += count - 1
                pre = max_
            else:
                pre = cur

        return res
