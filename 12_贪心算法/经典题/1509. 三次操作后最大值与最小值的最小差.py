from typing import List

# 给你一个数组 nums ，每次操作你可以选择 nums 中的任意一个元素并将它改成任意值。
# 请你返回三次操作后， nums 中最大值与最小值的差的最小值。


class Solution:
    def minDifference(self, nums: List[int]) -> int:
        if len(nums) <= 4:
            return 0
        nums.sort()
        # 删几个小的,几个大的
        return min(b - a for a, b in zip(nums[:4], nums[-4:]))


print(Solution().minDifference(nums=[1, 5, 0, 10, 14]))
# 输出：1
# 解释：将数组 [1,5,0,10,14] 变成 [1,1,0,1,1] 。
# 最大值与最小值的差为 1-0 = 1 。
