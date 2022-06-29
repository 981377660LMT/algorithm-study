# 求数组个数 使得对任意子数组[left,right]
# 有 min(nums[left]...nums[right])*(right-left+1) <=k
# !即直方图的面积 <=k
# !n<=100 k<=100
# 小课题4的解 区间dp O(n^3*k)
# https://atcoder.jp/contests/typical90/submissions/31445630

# 太难了...


from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = 998244353

n, k = map(int, input().split())


@lru_cache(None)
def dfs(left: int, right: int) -> int:
    for height in range(n):
        for i in range(left, right):
            ...


print(dfs(0, n - 1))

