import random
from typing import List

# 你打算重新排列数组中的元素以满足：重排后，数组中的每个元素都 不等于 其两侧相邻元素的 平均值 。
# Wiggle Sort


class Solution:
    def rearrangeArray(self, nums: List[int]) -> List[int]:
        nums = sorted(nums)
        for i in range(1, len(nums), 2):
            nums[i], nums[i - 1] = nums[i - 1], nums[i]
        return nums

    def rearrangeArray2(self, nums: List[int]) -> List[int]:
        nums = sorted(nums)
        half = (len(nums) + 1) >> 1
        nums[::2], nums[1::2] = nums[:half], nums[half:]
        return nums

    def rearrangeArray3(self, nums: List[int]) -> List[int]:
        """随机排序"""
        while True:
            if all(nums[i] != (nums[i - 1] + nums[i + 1]) / 2 for i in range(1, len(nums) - 1)):
                return nums
            random.shuffle(nums)


print(Solution().rearrangeArray2(nums=[1, 2, 3, 4, 5]))
print(Solution().rearrangeArray3(nums=[1, 2, 3, 4, 5]))
# 输出：[1,2,4,5,3]
# 解释：
# i=1, nums[i] = 2, 两相邻元素平均值为 (1+4) / 2 = 2.5
# i=2, nums[i] = 4, 两相邻元素平均值为 (2+5) / 2 = 3.5
# i=3, nums[i] = 5, 两相邻元素平均值为 (4+3) / 2 = 3.5
print(Solution().rearrangeArray2(nums=[1, 2, 5, 9]))
