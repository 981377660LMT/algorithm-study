# !B - ゲーム
# A,B<=1000
# nums[i]<=1000
# !两人不断从栈顶取数 都是最好发挥 求先手能取到的最大值

import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n1, n2 = map(int, input().split())
nums1 = list(map(int, input().split()))
nums2 = list(map(int, input().split()))
n = n1 + n2
sum_ = sum(nums1) + sum(nums2)


def dfs(index: int, i1: int) -> int:
    """博弈dp 返回两人最大的差"""
    i2 = index - i1
    if i1 > n1 or i2 > n2:
        return -INF
    if index == n:
        return 0

    hash = index * n1 + i1
    if memo[hash] != -INF:
        return memo[hash]

    res = -INF
    if i1 < n1:
        res = max(res, nums1[i1] - dfs(index + 1, i1 + 1))
    if i2 < n2:
        res = max(res, nums2[i2] - dfs(index + 1, i1))

    memo[hash] = res
    return res


memo = [-INF] * (n + 1) * (n1 + 1)
res = dfs(0, 0)

print((sum_ + res) // 2)
