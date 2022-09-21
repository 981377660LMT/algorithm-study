# 这个国家里有n (n <= 60)种货币，
# 每个货币的面值为a[i](a[1]=1,a[i]<a[i＋1])
# !并且保证了后一种面值是前一种面值的倍数。
# 现在让你购买x元的商品，问你给的钱的张数和找零的钱的张数最少是多少。

# !倒序记忆化搜索 尽量用大的面值凑 O(n)
# !1.要么少给 少的继续付
# !2.要么多给 多的继续找零
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    @lru_cache(None)
    def dfs(index: int, remain: int) -> int:
        if remain == 0:
            return 0
        if index == -1:
            return 0 if remain == 0 else INF

        div, mod = divmod(remain, money[index])
        if mod == 0:
            return div

        cand1 = div + dfs(index - 1, mod)  # 少给
        cand2 = div + 1 + dfs(index - 1, (div + 1) * money[index] - remain)  # 多给

        return min(cand1, cand2)

    n, x = map(int, input().split())
    money = list(map(int, input().split()))
    res = dfs(n - 1, x)
    print(res)
