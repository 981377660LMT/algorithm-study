import sys
from functools import reduce
from heapq import merge
from itertools import islice
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def mergeTwo(nums1: List[int], nums2: List[int], k: int) -> List[int]:
    """两个有序数组选前k大的和"""
    gen = lambda index: (nums1[index] + num for num in nums2)  # 一路
    allGen = [gen(i) for i in range(len(nums1))]  # 多路
    iterable = merge(*allGen, reverse=True)  # merge 相当于多路归并
    return list(islice(iterable, k))


# 时间复杂度O(kmlogn)
def kMax(mat: List[List[int]], k: int) -> List[int]:
    """有序矩阵中的前 k 个最大数组和"""
    return list(reduce(lambda row1, row2: mergeTwo(row1, row2, k), mat))


# n1,n2,n3<=1000
if __name__ == "__main__":
    n1, n2, n3, k = map(int, input().split())
    nums1 = sorted(map(int, input().split()), reverse=True)
    nums2 = sorted(map(int, input().split()), reverse=True)
    nums3 = sorted(map(int, input().split()), reverse=True)
    grid = [nums1, nums2, nums3]
    res = kMax(grid, k)
    print(*res, sep="\n")
