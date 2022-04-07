# All Sublists Sum
# 子数组所有元素之和


class Solution:
    def solve(self, nums):
        MOD = 10 ** 9 + 7
        n = len(nums)

        res = 0
        for i, num in enumerate(nums):
            res += (i + 1) * (n - i) * num % MOD
            res %= MOD
        return res


print(Solution().solve(nums=[2, 3, 5]))
# We have the following subarrays:

# [2]
# [3]
# [5]
# [2, 3]
# [3, 5]
# [2, 3, 5]
# The sum of all of these is 33.
