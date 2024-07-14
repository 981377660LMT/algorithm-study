# 预处理leftMost数组
# 区间等差数列

from typing import List


class Solution:
    def checkArithmeticSubarrays(self, nums: List[int], l: List[int], r: List[int]) -> List[bool]:
        queries = [(left, right) for left, right in zip(l, r)]
        n = len(nums)
        leftMost = [0] * n  # 每个位置构成等差数列的左边界
        for i in range(2, n):
            if nums[i] - nums[i - 1] == nums[i - 1] - nums[i - 2]:
                leftMost[i] = leftMost[i - 1]
            else:
                leftMost[i] = i - 1

        res = [False] * len(queries)
        for index, (s, e) in enumerate(queries):
            if leftMost[e] <= s:
                res[index] = True
        return res


print(Solution().checkArithmeticSubarrays([4, 6, 5, 9, 3, 7], [0, 0, 2], [2, 3, 5]))
