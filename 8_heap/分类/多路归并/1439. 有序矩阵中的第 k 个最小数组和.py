# 给你一个 m * n 的矩阵 mat，以及一个整数 k ，矩阵中的每一行都以非递减的顺序排列。
# 你可以从每一行中选出 1 个元素形成一个数组。返回所有可能数组中的第 k 个 最小 数组和。

# 1 <= m, n <= 40
# 1 <= k <= min(200, n ^ m)
# 1 <= matrix[i][j] <= 5000
# matrix[i] 是一个非递减数组

from functools import reduce
from heapq import heappop, heappush
from typing import List


def kSmallestPairs(nums1: List[int], nums2: List[int], k: int, unique=False) -> List[int]:
    """两个有序数组中查找第k小的数对之和."""
    ROW, COL = len(nums1), len(nums2)
    pq = [(nums1[i] + nums2[0], i, 0) for i in range(min(k, ROW))]  # (和, nums1的列, nums2的列)
    if not unique:
        res = []
        while pq and len(res) < k:
            _, col1, col2 = heappop(pq)
            res.append(nums1[col1] + nums2[col2])
            if col2 + 1 < COL:
                heappush(pq, (nums1[col1] + nums2[col2 + 1], col1, col2 + 1))  # type: ignore
        return res
    else:
        res = set()
        while pq and len(res) < k:
            _, col1, col2 = heappop(pq)
            res.add(nums1[col1] + nums2[col2])
            if col2 + 1 < COL:
                heappush(pq, (nums1[col1] + nums2[col2 + 1], col1, col2 + 1))
        return sorted(res)


def kthSmallest(matrix: List[List[int]], k: int, unique=False) -> int:
    """
    多个有序数组中查找第k小的数字之和.
    时间复杂度O(row*k*log(min(col,k))).
    """
    return reduce(lambda x, y: kSmallestPairs(x, y, k, unique=unique), matrix)[-1]


# print(Solution().kthSmallest(matrix=[[1, 3, 11], [2, 4, 6]], k=5))
# 输出：7
# 解释：从每一行中选出一个元素，前 k 个和最小的数组分别是：
# [1,2], [1,4], [3,2], [3,4], [1,6]。其中第 5 个的和是 7 。
