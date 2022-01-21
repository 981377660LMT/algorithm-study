from collections import defaultdict

MOD = int(1e9 + 7)

# 1 <= n <= 100
class Solution:
    def countSubranges(self, nums1: list[int], nums2: list[int]) -> int:
        n = len(nums1)
        dp = [defaultdict(int) for _ in range(n)]

        res = 0
        for i in range(n):
            dp[i][nums1[i]] += 1
            dp[i][-nums2[i]] += 1
            if i - 1 >= 0:
                for preSum in dp[i - 1]:
                    dp[i][preSum + nums1[i]] += dp[i - 1][preSum]
                    dp[i][preSum + nums1[i]] %= MOD
                    dp[i][preSum - nums2[i]] += dp[i - 1][preSum]
                    dp[i][preSum - nums2[i]] %= MOD
            res += dp[i][0]
            res %= MOD
        return res


print(Solution().countSubranges(nums1=[1, 2, 5], nums2=[2, 6, 3]))
print(Solution().countSubranges(nums1=[0, 1], nums2=[1, 0]))
