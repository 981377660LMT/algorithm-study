from typing import List

from sortedcontainers import SortedList

# nums 是范围 [0, n - 1] 内所有数字组成的一个排列
# 当数组 nums 中 全局倒置 的数量等于 局部倒置 的数量时，返回 true ；
# 全局倒置：逆序对
# 局部倒置：前大于后

# 局部倒置一定是全局倒置
# 如果i之前的最大值>a[i+2] 那么不成立
class Solution:
    def isIdealPermutation(self, nums: List[int]) -> bool:
        count2 = sum(a > b for a, b in zip(nums, nums[1:]))
        sl = SortedList()
        count1 = 0
        for num in nums[::-1]:
            pos = sl.bisect_left(num)
            count1 += pos
            sl.add(num)
        return count1 == count2


print(Solution().isIdealPermutation([1, 2, 0]))
# 输出：false
# 解释：有 2 个全局倒置，和 1 个局部倒置。
