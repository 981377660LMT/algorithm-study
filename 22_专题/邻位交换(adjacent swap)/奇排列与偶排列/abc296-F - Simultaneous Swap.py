# 奇排列与偶排列 + 逆序数
# https://atcoder.jp/contests/abc296/submissions/40221948

# 给定两个数组A和B,每次可以选定三个位置i,j,k,交换A[i]和A[j],交换B[i]和B[k],
# 问能否通过若干次操作使得A和B相等

# 注意到不变量:
# !交换不会影响逆序对的奇偶性(奇排列/偶排列)

from typing import List
from collections import Counter

from sortedcontainers import SortedList


def countInv(nums: List[int]) -> int:
    """求数组逆序对数量之和"""
    n = len(nums)
    res = 0
    visited = SortedList()
    for i in range(n - 1, -1, -1):
        smaller = visited.bisect_left(nums[i])
        res += smaller
        visited.add(nums[i])
    return res


def simultaneousSwap(A: List[int], B: List[int]) -> bool:
    C1 = Counter(A)
    C2 = Counter(B)
    if C1 != C2:
        return False
    if max(C1.values(), default=0) > 1:
        return True
    inv1, inv2 = countInv(A), countInv(B)
    return inv1 % 2 == inv2 % 2


if __name__ == "__main__":
    n = int(input())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    res = simultaneousSwap(A, B)
    print("Yes" if res else "No")
