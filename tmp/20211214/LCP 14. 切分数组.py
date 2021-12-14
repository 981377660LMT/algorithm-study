from typing import List

# 将数组切分成若干个子数组，使得每个子数组最左边和最右边数字的最大公约数大于 1，求解最少能切成多少个子数组
# 1 <= nums.length <= 10^5
# 2 <= nums[i] <= 10^6


# 超时
class Solution:
    def splitArray(self, nums: List[int]) -> int:
        def gcd(a, b):  # GCD
            if not b:
                return a
            return gcd(b, a % b)

        n = len(nums)
        dp = [0] * n
        dp[0] = 1
        for i in range(1, n):
            dp[i] = n  # 赋初值
            if gcd(nums[i], nums[0]) > 1:
                dp[i] = 1
            else:
                for k in range(1, i + 1):
                    if gcd(nums[i], nums[k]) > 1:
                        dp[i] = min(dp[i], dp[k - 1] + 1)
        return dp[n - 1]


print(Solution().splitArray(nums=[2, 3, 3, 2, 3, 3]))
# 输出：2
# 解释：最优切割为 [2,3,3,2] 和 [3,3] 。第一个子数组头尾数字的最大公约数为 2 ，第二个子数组头尾数字的最大公约数为 3 。
