# Maximum Unique Sublist Sum
# 无重复元素的最大和子数组
# nums非负整数


from collections import defaultdict


class Solution:
    def solve(self, nums):
        res = 0
        window = defaultdict(int)
        left = 0
        curSum = 0
        for right in range(len(nums)):
            curSum += nums[right]
            window[nums[right]] += 1
            while window[nums[right]] > 1:
                curSum -= nums[left]
                window[nums[left]] -= 1
                left += 1
            res = max(res, curSum)
        return res


print(Solution().solve(nums=[1, 2, 2, 3, 4, 4]))
