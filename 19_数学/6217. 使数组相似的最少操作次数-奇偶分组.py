"""奇偶分组"""

from typing import List


# 给你两个正整数数组 nums 和 target ，两个数组长度相等。
# 在一次操作中，你可以选择两个 不同 的下标 i 和 j ，其中 0 <= i, j < nums.length ，并且：
# 令 nums[i] = nums[i] + 2 且
# 令 nums[j] = nums[j] - 2 。
# !如果两个数组中每个元素出现的频率相等，我们称两个数组是 相似 的。
# 请你返回将 nums 变得与 target 相似的最少操作次数。
# 测试数据保证 nums 一定能变得与 target 相似。

# !奇数只能变成奇数，偶数只能变成偶数。
# !分别考虑奇数数组和目标的奇数数组以及偶数数组和目标的偶数数组，一定是维持当前的顺序找目标更优。（否则总距离会更长）


class Solution:
    def makeSimilar(self, nums: List[int], target: List[int]) -> int:
        odd1 = sorted([num for num in nums if num % 2 == 1])
        even1 = sorted([num for num in nums if num % 2 == 0])
        odd2 = sorted([num for num in target if num % 2 == 1])
        even2 = sorted([num for num in target if num % 2 == 0])
        res1 = sum(max(0, (num1 - num2)) for num1, num2 in zip(odd1, odd2))
        res2 = sum(max(0, (num1 - num2)) for num1, num2 in zip(even1, even2))
        return (res1 + res2) // 2


print(Solution().makeSimilar(nums=[8, 12, 6], target=[2, 14, 10]))
print(Solution().makeSimilar(nums=[1, 2, 5], target=[4, 1, 3]))
print(Solution().makeSimilar(nums=[1, 2, 5], target=[4, 1, 3]))
print(Solution().makeSimilar(nums=[2, 6, 1], target=[6, 0, 3]))
