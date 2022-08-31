# https://atcoder.jp/contests/abc236/tasks/abc236_f
# 异或线性基

# 从1到2^n-1中选一些数出来，使得可以用这些数通过异或运算可以表示1到2^n-1中的任何数。
# 选第i个数的代价为ai，最小化代价。
# n<=16,cost[i]<=1e9

# !贪心地从代价小的数开始插入线性基，能插就插，不能插就扔。这样最终代价是最小的。


import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def findMinCostBase(nums: List[int], costs: List[int]) -> List[int]:
    """求1-2^n-1中成本最小的线性基"""
    nums = sorted(nums, key=lambda i: costs[i - 1])
    res = []
    visited = set([0])
    for cur in nums:
        if cur in visited:
            continue
        res.append(cur)
        visited |= {cur ^ pre for pre in visited}
    return res


n = int(input())
costs = list(map(int, input().split()))
nums = [i for i in range(1, 1 << n)]
bases = findMinCostBase(nums, costs)
print(sum(costs[i - 1] for i in bases))
