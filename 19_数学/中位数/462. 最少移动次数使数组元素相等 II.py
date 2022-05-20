# /**
#  *
#  * @param nums
#  * @description
#  * 给定一个非空整数数组，找到使所有数组元素相等所需的最小移动数，
#  * 其中每次移动可将选定的一个元素加1或减1。
#  * 您可以假设数组的长度最多为10000。
#  */
from heapq import nlargest
from typing import List


class Solution:
    def minMoves2(self, nums: List[int]) -> int:
        # 第(n>>1+1)小的数
        mid = nlargest((len(nums) >> 1) + 1, nums)[-1]
        return sum(abs(num - mid) for num in nums)

