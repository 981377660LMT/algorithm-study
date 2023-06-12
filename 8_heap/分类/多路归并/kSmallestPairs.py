# 两个有序数组中查找和最小的k对数字
# 给定两个以 升序排列 的整数数组 nums1 和 nums2 , 以及一个整数 k 。
# 定义一对值 (u,v)，其中第一个元素来自 nums1，第二个元素来自 nums2 。
# 请找到和最小的 k 个数对 (u1,v1),  (u2,v2)  ...  (uk,vk) 。


"""时间复杂度O(k*logmin(n,k))"""

from heapq import heappop, heappush
from typing import List, Tuple


def kSmallestPairs1(nums1: List[int], nums2: List[int], k: int) -> List[Tuple[int, int]]:
    """两个有序数组中查找和最小的k对数字.返回k个数对."""
    ROW, COL = len(nums1), len(nums2)
    res = []
    pq = [(nums1[i] + nums2[0], i, 0) for i in range(min(k, ROW))]  # (和, nums1的列, nums2的列)
    while pq and len(res) < k:
        _, col1, col2 = heappop(pq)
        res.append((nums1[col1], nums2[col2]))
        if col2 + 1 < COL:
            heappush(pq, (nums1[col1] + nums2[col2 + 1], col1, col2 + 1))
    return res


def kSmallestPairs2(nums1: List[int], nums2: List[int], k: int) -> List[int]:
    """两个有序数组中查找和最小的k对数字.返回k个数对的和."""
    ROW, COL = len(nums1), len(nums2)
    res = []
    pq = [(nums1[i] + nums2[0], i, 0) for i in range(min(k, ROW))]  # (和, nums1的列, nums2的列)
    while pq and len(res) < k:
        _, col1, col2 = heappop(pq)
        res.append(nums1[col1] + nums2[col2])
        if col2 + 1 < COL:
            heappush(pq, (nums1[col1] + nums2[col2 + 1], col1, col2 + 1))
    return res


assert kSmallestPairs1(nums1=[1, 7, 11], nums2=[2, 4, 6], k=3) == [(1, 2), (1, 4), (1, 6)]
