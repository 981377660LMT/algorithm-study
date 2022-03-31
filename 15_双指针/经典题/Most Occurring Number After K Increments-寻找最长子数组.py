# k次操作 每次操作可以为1个数加1
# 求k次操作后 freq最大的数 如果有多个答案 取较小的那个

# 类似于找和不大于k的最长子数组


# n ≤ 100,000
# k < 2 ** 31

# 为什么双指针:操作具有单向性
class Solution:
    def solve(self, nums, k):
        nums = sorted(nums)
        n = len(nums)
        res, freq = 0, 0
        left = 0
        curSum = 0
        for right in range(n):
            curSum += nums[right]
            while nums[right] * (right - left + 1) > k + curSum:
                curSum -= nums[left]
                left += 1
            if right - left + 1 > freq:
                freq = right - left + 1
                res = nums[right]
        return res
