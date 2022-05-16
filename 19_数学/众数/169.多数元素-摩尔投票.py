# 返回其中的多数元素。多数元素是指在数组中出现次数 大于 ⌊ n/2 ⌋ 的元素。
# 你可以假设数组是非空的，并且给定的数组总是存在多数元素。
from typing import List


class Solution:
    def majorityElement(self, nums: List[int]) -> int:
        res, count = int(1e20), 0
        for num in nums:
            if num == res:
                count += 1
            elif count == 0:
                res = num
                count = 1
            else:
                count -= 1
                
        return res if nums.count(res) > len(nums) // 2 else -1

