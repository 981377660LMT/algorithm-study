"""
统计并返回 num 中的 中位数 等于 k 的非空子数组的数目。
数组的中位数是按 递增 顺序排列后位于 中间 的那个元素，
如果数组长度为偶数，则中位数是位于中间靠 左 的那个元素。

Solution:
前缀和+容斥原理
https://leetcode.cn/problems/count-subarrays-with-median-k/solution/by-981377660lmt-3tzk/
1.将 中位数等于k的子数组数 转换为 中位数大于等于k的子数组数 减去 中位数大于等于k+1的子数组数;
2.将大于等于中位数的数变为1,小于中位数的数变为-1,转换为求 有多少个子数组的和大于0;
3.遍历子数组,维护前缀和,对每个位置求出有多少个前缀和大于0.

类似:
https://leetcode.cn/problems/count-subarrays-with-more-ones-than-zeros/solution/onjie-fa-by-newhar-xy9d/
"""

from typing import List
from sortedcontainers import SortedList


class Solution:
    def countSubarrays(self, nums: List[int], k: int) -> int:
        def cal(mid: int) -> int:
            """中位数>=mid的子数组个数"""
            curNums = [1 if num >= mid else -1 for num in nums]
            res, curSum, sl = 0, 0, SortedList([0])
            for num in curNums:
                curSum += num
                res += sl.bisect_left(curSum)
                sl.add(curSum)
            return res

        return cal(k) - cal(k + 1)


# 3
assert Solution().countSubarrays(nums=[4, 1, 3, 2], k=1) == 3
