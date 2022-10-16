"""
C -上界受限的元素不重复的数组个数

求长为n的数组的个数nums[i],满足
1. 1<=nums[i]<=uppers[i]
2. nums[i] != nums[j] if i!=j

!排序,求每个元素可选择个数
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


if __name__ == "__main__":
    n = int(input())
    uppers = list(map(int, input().split()))
    uppers.sort()
    res = 1
    for i, num in enumerate(uppers):
        res *= max(0, num - i)
        res %= MOD
    print(res)
