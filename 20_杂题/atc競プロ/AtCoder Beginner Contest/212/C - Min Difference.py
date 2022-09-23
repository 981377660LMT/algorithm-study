"""两个数组元素的最小距离
n<=2e5

排序+双指针
"""

import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def minDiff(nums1: List[int], nums2: List[int]) -> int:
    """两个数组元素的最小距离"""
    n1, n2 = len(nums1), len(nums2)
    nums1, nums2 = sorted(nums1), sorted(nums2)
    res, i, j = INF, 0, 0
    while i < n1 and j < n2:
        res = min(res, abs(nums1[i] - nums2[j]))
        if nums1[i] < nums2[j]:
            i += 1
        else:
            j += 1
    return res


if __name__ == "__main__":
    n, m = map(int, input().split())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    print(minDiff(nums1, nums2))
