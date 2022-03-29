# `不断删除上升降子序列，问最少删除几次(答案是 最长下降子序列的长度，即取逆后的最长不降子序列)`


from bisect import bisect_right


class Solution:
    def solve(self, nums):
        n = len(nums)
        LIS = [nums[n - 1]]

        for i in range(n - 2, -1, -1):
            index = bisect_right(LIS, nums[i])
            if index == len(LIS):
                LIS.append(nums[i])
            else:
                LIS[index] = nums[i]

        return len(LIS)
