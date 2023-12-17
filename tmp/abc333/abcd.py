from collections import deque
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

from functools import reduce
from heapq import heappop, heappush
from typing import List, Tuple


def kSmallestPairs1(nums1: List[int], nums2: List[int], k: int) -> List[Tuple[int, int]]:
    """两个有序数组中查找第k小的数对之和.返回所有数对."""
    ROW, COL = len(nums1), len(nums2)
    res = []
    pq = [(nums1[i] + nums2[0], i, 0) for i in range(min(k, ROW))]  # (和, nums1的列, nums2的列)
    while pq and len(res) < k:
        _, col1, col2 = heappop(pq)
        res.append((nums1[col1], nums2[col2]))
        if col2 + 1 < COL:
            heappush(pq, (nums1[col1] + nums2[col2 + 1], col1, col2 + 1))
    return res


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


if __name__ == "__main__":
    n = int(input())
    nums1 = [int("1" * i) for i in range(1, 100)]
    nums2 = [int("1" * i) for i in range(1, 100)]
    nums3 = [int("1" * i) for i in range(1, 100)]
    matrix = [nums1, nums2, nums3]
    print(kthSmallest(matrix, n, unique=True))
