# 返回数组中两个数之和,他们两个数的和大于target且最接近target的两个数的和
from typing import List


class Solution:
    def solve(self, nums: List[int], target: int) -> int:
        n = len(nums)
        nums.sort()

        res = int(1e20)
        i, j = 0, n - 1
        while i < j:
            if nums[i] + nums[j] > target:
                res = min(res, nums[i] + nums[j])
                j -= 1
            else:
                i += 1

        return res
