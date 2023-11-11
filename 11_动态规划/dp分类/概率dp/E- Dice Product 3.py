# Dice Product 3
# 六面骰子，每面等概率出现。
# 现在不断掷骰子，直到掷出来的数的乘积大于等于N
# 问恰好为 N的概率。
# 对 998244353取模。

# !dp[n]=dp[n/1]/6+dp[n/2]/6+dp[n/3]/6+dp[n/4]/6+dp[n/5]/6+dp[n/6]/6
# !即dp[n]=dp[n/2]/5+dp[n/3]/5+dp[n/4]/5+dp[n/5]/5+dp[n/6]/5

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


if __name__ == "__main__":
    N = int(input())

    # 只能有素因子 2 3 5
    p2, p3, p5 = 0, 0, 0
    cur = N
    while cur % 2 == 0:
        cur //= 2
        p2 += 1
    while cur % 3 == 0:  # 3
        cur //= 3
        p3 += 1
    while cur % 5 == 0:  # 5
        cur //= 5
        p5 += 1
    if cur != 1:
        print(0)
        exit()
    inv6 = pow(6, MOD - 2, MOD)
    inv5 = pow(5, MOD - 2, MOD)

    @lru_cache(None)
    def dfs(p2: int, p3: int, p5: int) -> int:
        if p2 < 0 or p3 < 0 or p5 < 0:
            return 0
        if p2 == 0 and p3 == 0 and p5 == 0:
            return 1
        res = 0
        res += dfs(p2 - 1, p3, p5)
        res += dfs(p2, p3 - 1, p5)
        res += dfs(p2 - 2, p3, p5)
        res += dfs(p2, p3, p5 - 1)
        res += dfs(p2 - 1, p3 - 1, p5)
        res %= MOD
        return (res * inv5) % MOD  # 这个地方要除以 5 而不是最后dfs(p2,p3,p5)除以5

    res = dfs(p2, p3, p5)
    dfs.cache_clear()
    print((res) % MOD)
