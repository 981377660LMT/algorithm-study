# 子数组最小值大于等于k的不重叠子数组的个数
from itertools import accumulate


MOD = int(1e9 + 7)


class Solution:
    def solve(self, nums, k):
        def getDp(nums):
            """每个结尾有多少个符合的数组"""
            dp = []
            count = 0
            for num in nums:
                if num >= k:
                    count += 1
                else:
                    count = 0
                dp.append(count)
            return dp

        dp1, dp2 = getDp(nums), getDp(nums[::-1])
        suffixSum = list(accumulate(dp2))[::-1]

        res = 0
        for i in range(len(dp1) - 1):
            res += dp1[i] * suffixSum[i + 1]
            res %= MOD
        return res


print(Solution().solve(nums=[3, 4, 4, 9], k=4))
# These are the pairs of sublists:

# [4], [9]
# [4], [9] (using other 4)
# [4], [4]
# [4], [4, 9]
# [4, 4], [9]
