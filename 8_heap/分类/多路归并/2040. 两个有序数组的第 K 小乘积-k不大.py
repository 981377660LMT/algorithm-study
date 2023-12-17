from heapq import merge
from itertools import islice
from typing import Generator, List


# 0 ≤ k ≤ 100,000
# 1 ≤ n ≤ 100,000
# 1 ≤ m ≤ 100,000
# 时间复杂度O(m+nlogn+klogm),


# 2040. 两个有序数组的第 K 小乘积-check函数也二分
class Solution:
    def kthSmallestProduct(self, nums1: List[int], nums2: List[int], k: int) -> int:
        """
        !TLE
        """

        def g1(index: int) -> Generator[int, None, None]:
            return (nums1[index] * num for num in nums2)

        def g2(index: int) -> Generator[int, None, None]:
            return (nums1[index] * num for num in reversed(nums2))

        nums1, nums2 = sorted(nums1), sorted(nums2)
        # 生成器的多路
        gen = (g2(i) if num < 0 else g1(i) for i, num in enumerate(nums1))
        iter = merge(*gen)
        return next(islice(iter, k - 1, None))  # 切片，第k个最小


print(Solution().kthSmallestProduct(nums1=[-4, -2, 0, 3], nums2=[2, 4], k=6))
print(Solution().kthSmallestProduct(nums1=[2, 5], nums2=[3, 4], k=2))
