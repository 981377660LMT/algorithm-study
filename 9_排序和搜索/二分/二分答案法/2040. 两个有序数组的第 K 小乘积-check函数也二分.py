from itertools import accumulate
from typing import List
from bisect import bisect_left, bisect_right
from math import ceil, floor

# 1 <= nums1.length, nums2.length <= 5 * 104
# -105 <= nums1[i], nums2[j] <= 105
#  k很大最多 25 * 10^8
# 第k=>二分找答案
# 在类似题目 668. 乘法表中第k小的数 和 719. 找出第 k 小的距离对 中，经典的解法是双指针，将单次 check 的时间复杂度降低到了 O(n)O(n)
# https://leetcode-cn.com/problems/kth-smallest-product-of-two-sorted-arrays/solution/leetcode-668-cheng-fa-biao-bian-xing-by-8axdy/
# 注意 python 除法的取整问题。python 中的 // 是向下取整(>>也是向下取整)，但是向上取整，为了方便，直接浮点除法 + math.ceil 调用实现。


class Solution:
    def kthSmallestProduct(self, nums1: List[int], nums2: List[int], k: int) -> int:
        def countNGT(mid: int) -> int:
            """Return count of products <= mid."""
            res = 0
            for x in nums1:
                if x < 0:
                    res += len(nums2) - bisect_left(nums2, ceil(mid / x))
                elif x == 0:
                    if mid >= 0:
                        res += len(nums2)
                else:
                    res += bisect_right(nums2, floor(mid / x))
            return res

        # def countNGT(mid: int) -> int:
        #     """注意到-1e5<=nums[i]<=1e5 可以使用前缀和存位置代替二分"""
        #     nums = [0] * int(2e5 + 10)
        #     for num in nums2:
        #         nums[num + int(1e5)] += 1
        #     preSum = list(accumulate(nums))

        #     def count(x: int) -> int:
        #         if x < -int(1e5):
        #             return 0
        #         if x > int(1e5):
        #             return preSum[-1]
        #         return preSum[x + int(1e5)]

        #     def inner(mid: int) -> int:
        #         res = 0
        #         for x in nums1:
        #             if x < 0:
        #                 res += len(nums) - count(ceil(mid / x))
        #             elif x == 0:
        #                 if mid >= 0:
        #                     res += len(nums2)
        #             else:
        #                 res += count(floor(mid / x))
        #         return res

        #     return inner(mid)

        # 遍历小的，二分大的(优化了1000ms)
        if len(nums1) > len(nums2):
            nums1, nums2 = nums2, nums1

        # 优化了400ms
        a, b, c, d = (
            nums1[0] * nums2[0],
            nums1[0] * nums2[-1],
            nums1[-1] * nums2[-1],
            nums1[-1] * nums2[0],
        )
        left, right = min(a, b, c, d), max(a, b, c, d)
        # left, right = -(10 ** 10), 10 ** 10 + 1

        while left <= right:
            mid = left + right >> 1
            if countNGT(mid) < k:
                left = mid + 1
            else:
                right = mid - 1
        return left


print(Solution().kthSmallestProduct(nums1=[2, 5], nums2=[3, 4], k=2))
