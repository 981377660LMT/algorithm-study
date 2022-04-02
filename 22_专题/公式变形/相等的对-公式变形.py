# 找到所有的pair i<=j，使得 nums[i] + rev(nums[j]) = nums[j] + rev(nums[i]).
# 公式变形

# nums[i] + rev(nums[j]) = nums[j] + rev(nums[i]) also means nums[i] - rev(nums[i]) = nums[j] - rev(nums[j])
from collections import Counter


MOD = int(1e9) + 7


class Solution:
    def solve(self, nums):
        count = Counter()
        res = 0

        for x in nums:
            x -= int(str(x)[::-1])
            count[x] += 1
            res += count[x]

        return res % MOD


print(Solution().solve(nums=[1, 20, 2, 11]))
