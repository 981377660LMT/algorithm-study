from typing import List

# 数组 [4,2,5,3] 的交替和为 (4 + 5) - (2 + 3) = 4 。
# 请你返回 nums 中任意子序列的 最大交替和

# 总结:偶数索引结尾减，奇数索引结尾加
class Solution:
    def maxAlternatingSum(self, nums: List[int]) -> int:

        n = len(nums)
        dp0 = [0 for _ in range(n)]  # 选中偶数个数字
        dp1 = [0 for _ in range(n)]  # 选中奇数个数字
        dp0[0] = 0
        dp1[0] = nums[0]

        for i in range(1, n):
            dp0[i] = max(dp0[i - 1], dp1[i - 1] - nums[i])
            dp1[i] = max(dp1[i - 1], dp0[i - 1] + nums[i])
        print(dp0, dp1)
        return max(dp0[-1], dp1[-1])


print(Solution().maxAlternatingSum(nums=[4, 2, 5, 3]))
# 输出：7
# 解释：最优子序列为 [4,2,5] ，交替和为 (4 + 5) - 2 = 7 。
