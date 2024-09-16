# F - Palindromic Expression (暴力记忆化搜索)
# https://atcoder.jp/contests/abc363/tasks/abc363_f
# 给定n将n拆成若干个数相乘的表达式，其中这个表达式是个回文串，且不存在数字0
# n<=1e12
#
# 记忆化搜索


from functools import lru_cache
from math import isqrt
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    N = int(input())

    @lru_cache(None)
    def dfs(v: int) -> str:
        s = str(v)
        if not "0" in s and s == s[::-1]:
            return s
        for x in range(2, isqrt(v) + 1):
            if v % x == 0 and not "0" in str(x):
                y = int(str(x)[::-1])
                if v // x % y == 0:
                    mid = dfs(v // x // y)
                    if len(mid):
                        return f"{x}*{mid}*{y}"
        return ""

    res = dfs(N)
    dfs.cache_clear()
    print(res if len(res) else "-1")
