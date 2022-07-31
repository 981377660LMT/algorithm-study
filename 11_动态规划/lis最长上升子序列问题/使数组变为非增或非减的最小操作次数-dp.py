# Make Array Non-decreasing or Non-increasing
# 每次操作可以使每个数加1或者减1
# 修改序列为非降或非增的最小修改次数
# 1 <= nums.length <= 1000
# 0 <= nums[i] <= 1000

# 显然dp[i][num]

from typing import List

# 将数组变为不减的最小操作数
INF = int(1e20)


class Solution:
    def convertArray(self, nums: List[int]) -> int:
        def helper(nums: List[int]) -> int:
            """变为不减数组的最小操作次数

            dp[i][num] 表示 前i个元素最num结尾的最小的代价是多少
            """
            n, max_ = len(nums), max(nums)

            dp = [[INF] * (max_ + 1) for _ in range(n)]
            for j in range((max_ + 1)):
                dp[0][j] = abs(nums[0] - j)
            for i in range(1, n):
                preMin = INF  # 前缀优化dp
                for j in range((max_ + 1)):
                    preMin = min(preMin, dp[i - 1][j])
                    dp[i][j] = preMin + abs(nums[i] - j)

            return min(dp[-1])

        return min(helper(nums), helper(nums[::-1]))

    def convertArray2(self, nums: List[int]) -> int:
        def helper(nums: List[int]) -> int:
            """变为不减数组的最小操作次数

            实现最小的代价时，一定会让整个序列的末位元素等于原序列中的某个元素。
            """
            n = len(nums)
            allNums = sorted(set(nums))
            m = len(allNums)

            dp = [[INF] * m for _ in range(n)]
            for j in range((m)):
                dp[0][j] = abs(nums[0] - allNums[j])
            for i in range(1, n):
                preMin = INF
                for j in range(m):
                    preMin = min(preMin, dp[i - 1][j])
                    dp[i][j] = preMin + abs(nums[i] - allNums[j])

            return min(dp[-1])

        return min(helper(nums), helper(nums[::-1]))


print(Solution().convertArray2(nums=[3, 2, 4, 5, 0]))
print(Solution().convertArray2(nums=[3, 1, 2, 1]))
print(Solution().convertArray2([11, 11, 13, 8, 18, 19, 20, 7, 16, 3]))
