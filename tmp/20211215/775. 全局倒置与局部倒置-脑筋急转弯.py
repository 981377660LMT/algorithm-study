from typing import List

# nums 是范围 [0, n - 1] 内所有数字组成的一个排列
# 当数组 nums 中 全局倒置 的数量等于 局部倒置 的数量时，返回 true ；
# 全局倒置：逆序对
# 局部倒置：前大于后

# 局部倒置一定是全局倒置
# 如果i之前的最大值>a[i+2] 那么不成立
class Solution:
    def isIdealPermutation(self, nums: List[int]) -> bool:
        curMax = nums[0]
        for i in range(len(nums) - 2):
            curMax = max(curMax, nums[i])
            if curMax > nums[i + 2]:
                return False
        return True


print(Solution().isIdealPermutation([1, 2, 0]))
# 输出：false
# 解释：有 2 个全局倒置，和 1 个局部倒置。
