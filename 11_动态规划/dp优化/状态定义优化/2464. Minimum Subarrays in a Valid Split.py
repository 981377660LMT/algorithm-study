# !将数组切分成若干个子数组，
# !使得每个子数组最左边和最右边数字的最大公约数大于 1，求解最少能切成多少个子数组。
# !n<=1000
# !nums[i]<=1e5

from math import gcd
from typing import List

INF = int(1e9)


class Solution:
    def validSubarraySplit(self, nums: List[int]) -> int:
        """
        O(n^2)解法

        !dp[i] = dp[j-1]+1 if gcd(nums[j], nums[i]) > 1 (0<=j<=i)
        """
        n = len(nums)
        dp = [INF] * n
        dp[0] = 1
        for i in range(1, n):
            for j in range(i + 1):
                if gcd(nums[j], nums[i]) > 1:
                    dp[i] = min(dp[i], (dp[j - 1] if j else 0) + 1)
        return dp[-1] if dp[-1] != INF else -1


print(Solution().validSubarraySplit(nums=[2, 6, 3, 4, 3]))
print(
    Solution().validSubarraySplit(
        nums=[17, 11, 19, 53, 97, 89, 7, 23, 89, 61, 61, 17, 13, 43, 43, 73, 89, 83]
    )
)
