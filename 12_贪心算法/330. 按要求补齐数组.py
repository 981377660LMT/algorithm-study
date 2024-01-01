# 330. 按要求补齐数组(2952. 需要添加的硬币的最小数量)
# https://leetcode.cn/problems/patching-array/
# 1798. 你能构造出连续值的最大数目
# https://leetcode.cn/problems/maximum-number-of-consecutive-values-you-can-make/

# 给定一个已排序的正整数数组 nums ，和一个正整数 n 。
# 从 [1, n] 区间内选取任意个数字补充到 nums 中，
# 使得 [1, n] 区间内的任何数字都可以用 nums 中某几个数字的和来表示。
# n<=2^31-1
# nums[i]<=1e4
# nums[i].length<=1000

# !结论:如果当前区间是 [1,x]，我们应该添加数字 x + 1，这样可以覆盖的区间为 [1,2*x+1]
from typing import List


class Solution:
    def minPatches(self, nums: List[int], n: int) -> int:
        nums = sorted(nums)
        upper = 0
        res = 0
        ei = 0
        while upper < n:
            if ei < len(nums) and nums[ei] <= upper + 1:
                upper += nums[ei]
                ei += 1
            else:
                upper = upper * 2 + 1
                res += 1

        return res
