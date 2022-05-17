from functools import reduce
from heapq import merge
from itertools import islice
from typing import Generator, List

# 1 <= m, n <= 40
# 1 <= k <= min(200, n ^ m)
# 1 <= mat[i][j] <= 5000
# mat[i] 是一个非递减数组


def mergeTwo(nums1: List[int], nums2: List[int], k: int) -> List[int]:
    """两个有序数组选前k小的和"""
    gen = lambda index: (nums1[index] + num for num in nums2)  # 一路
    allGen = [gen(i) for i in range(len(nums1))]  # 多路
    iterable = merge(*allGen, reverse=False)  # merge 相当于多路归并
    return list(islice(iterable, k))


# 时间复杂度O(kmlogn)
class Solution:
    def kthSmallest(self, mat: List[List[int]], k: int) -> int:
        """有序矩阵中的第 k 个最小数组和"""
        return next(reversed(reduce(lambda row1, row2: mergeTwo(row1, row2, k), mat)))


print(Solution().kthSmallest(mat=[[1, 3, 11], [2, 4, 6]], k=5))
