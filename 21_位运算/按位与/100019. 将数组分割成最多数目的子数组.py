# 100019. 将数组分割成最多数目的子数组(区间与)
# https://leetcode.cn/problems/split-array-into-maximum-number-of-subarrays/

# 从左到右遍历数组，只要发现 AND 等于 0 就立刻分割。
# 如果不立刻分割，由于 AND 的数越多越能为 0，
# 现在多分了一个数，后面就要少分一个数，可能后面就不能为 0 了。


from functools import reduce
from typing import List


class Solution:
    def maxSubarrays(self, nums: List[int]) -> int:
        and_ = reduce(lambda a, b: a & b, nums, -1)
        if and_ != 0:
            return 1

        res = 0
        curAnd = -1
        for num in nums:
            curAnd &= num
            if curAnd == 0:
                res += 1
                curAnd = -1
        return res
