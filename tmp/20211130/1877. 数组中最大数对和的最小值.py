from typing import List

# 将 nums 中的元素分成 n / 2 个数对，使得：
# nums 中每个元素 恰好 在 一个 数对中，且
# 最大数对和 的值 最小 。
# 请你在最优数对划分的方案下，返回最小的 最大数对和 。
class Solution:
    def minPairSum(self, nums: List[int]) -> int:
        return max(a + b for a, b in zip(sorted(nums), sorted(nums)[::-1]))


print(Solution().minPairSum(nums=[3, 5, 2, 3]))
