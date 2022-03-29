# n ≤ 100,000
MOD = int(1e9 + 7)

# 肯定不能枚举pair，要对每个数分析


class Solution:
    def solve(self, nums):
        nums = sorted(nums)
        res = 0
        n = len(nums)
        for i, num in enumerate(nums):
            add = num * i
            sub = num * (n - i - 1)
            res += add - sub
            res %= MOD
        return res * 2


print(Solution().solve(nums=[1, 3, 6]))
# Result is abs(1 - 3) + abs(1 - 6) + abs(3 - 1) + abs(3 - 6) + abs(6 - 1) + abs(6 - 3).
