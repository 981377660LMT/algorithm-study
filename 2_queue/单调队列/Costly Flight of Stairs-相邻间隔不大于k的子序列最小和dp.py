# 子序列中每两个 相邻 的整数 nums[i] 和 nums[j] ，
# 它们在原数组中的下标 i 和 j 满足 i < j 且 j - i <= k 。
# 求子序列和的最小值
from collections import deque


# 你的目标是到达数组最后一个位置(必须取最后一个元素)
class Solution:
    def solve(self, nums, k):
        n = len(nums)
        # 以i结尾的最小和，必须取i
        dp = [int(1e20)] * n
        dp[0] = nums[0]
        queue = deque([0])

        for i in range(1, n):
            while queue and i - queue[0] > k:
                queue.popleft()

            dp[i] = nums[i] + dp[queue[0]]

            while queue and dp[queue[-1]] > dp[i]:
                queue.pop()
            queue.append(i)

        return dp[-1]


print(Solution().solve([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 3))
