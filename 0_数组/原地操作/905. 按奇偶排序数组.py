# 原地首尾指针操作
from typing import List

#  nums 中的的所有偶数元素移动到数组的前面，后跟所有奇数元素。


class Solution:
    def sortArrayByParity(self, nums: List[int]) -> List[int]:
        left, right = 0, len(nums) - 1

        # left<right的条件始终加上安全一些
        while left < right:
            # 奇数搬到后面，偶数搬到前面
            while left < right and nums[left] % 2 == 0:
                left += 1
            while left < right and nums[right] % 2 == 1:
                right -= 1
            if left < right:
                nums[left], nums[right] = nums[right], nums[left]
                left += 1
                right -= 1

        return nums
