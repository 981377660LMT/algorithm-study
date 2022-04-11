# 合并数组为不增数组的最少操作次数

# n ≤ 1,000
# n^2logn

# 题目可以转化为 对数组的逆 合并为不降数组的最大段数
# 状态定义dp((n,k)): 前n个数，最后一段数组合并成的元素为k
# 注意到dp(n,k)有单调性 转移状态可以二分

# https://binarysearch.com/problems/Turn-Into-Non-Increasing-List/solutions/1109053
from collections import defaultdict
from sortedcontainers import SortedDict


class Solution:
    def solve(self, nums):
        n = len(nums)
        nums = nums[::-1]

        dp = defaultdict(SortedDict)
        dp[0][0] = 0  # dp[i][k]：数组的前i个数合并，最后一个被合并的数大小为k时的最多分割数

        for i in range(n):
            curSum = 0
            okPairs = []

            for j in range(i, -1, -1):
                curSum += nums[j]
                pos = dp[j].bisect_right(curSum) - 1  # 找到第一个大于等于curSum的元素
                if pos < 0:
                    continue
                okPairs.append((curSum, dp[j].peekitem(pos)[1] + 1))

            # print(okPairs, dp, i, curSum)
            okPairs.sort()
            preCount = -int(1e20)
            for sum_, count in okPairs:
                if count <= preCount:
                    continue
                dp[i + 1][sum_] = count
                preCount = count
            # print(dp)

        res = max((count for count in dp[n].values()), default=0)
        return n - res


print(Solution().solve(nums=[1, 5, 3, 9, 1]))
# We can merge [1, 5] to get [6, 3, 9, 1] and then merge [6, 3] to get [9, 9, 1]
