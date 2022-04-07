# 所有数不同
# 求最长的值域连续子数组

# 0 ≤ n ≤ 1,000
class Solution:
    def solve(self, nums):
        """n^2 可以对每个开头进行dp"""
        if not nums:
            return 0

        res = 1
        n = len(nums)
        for i in range(n):
            minVal = nums[i]
            maxVal = nums[i]
            for j in range(i + 1, n):
                minVal = min(minVal, nums[j])
                maxVal = max(maxVal, nums[j])
                if maxVal - minVal == j - i:
                    res = max(res, j - i + 1)

        return res


print(Solution().solve(nums=[1, 4, 5, 3, 2, 9]))
