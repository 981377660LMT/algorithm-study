from typing import List


# 2 <= nums.length <= 500
# 1 <= nums[i] <= 10^3
class Solution:
    def maxProduct(self, nums: List[int]) -> float:
        if len(nums) == 2:
            return (nums[0] - 1) * (nums[1] - 1)

        first, second = float('-inf'), float('-inf')

        for num in nums:
            if num > first:
                # 注意这句，先要让位
                second = first
                first = num
            elif num > second:
                second = num

        return (first - 1) * (second - 1)


print(Solution().maxProduct([3, 4, 5, 2]))

