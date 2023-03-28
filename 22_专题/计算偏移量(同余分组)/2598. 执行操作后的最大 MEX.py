# 2598. 执行操作后的最大 MEX
# !在一步操作中，你可以对 nums 中的任一元素加上或减去 value 。
# 数组的 MEX (minimum excluded) 是指其中数组中缺失的最小非负整数。
# 返回在执行上述操作 任意次 后，nums 的最大 MEX 。


from collections import Counter
from typing import List


class Solution:
    def findSmallestInteger(self, nums: List[int], value: int) -> int:
        modCounter = Counter([num % value for num in nums])
        mex = 0
        while modCounter[mex % value] > 0:
            modCounter[mex % value] -= 1
            mex += 1
        return mex


# nums = [1,-10,7,13,6,8], value = 5
print(Solution().findSmallestInteger(nums=[1, -10, 7, 13, 6, 8], value=5))
