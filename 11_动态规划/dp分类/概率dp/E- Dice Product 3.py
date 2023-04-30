# Dice Product 3
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# あなたは
# 1 以上
# 6 以下の整数が等確率で出るサイコロと整数
# 1 を持っています。
# あなたは持っている整数が
# N 未満である間、次の操作を繰り返します。

# サイコロを振り、出た目を
# x とする。持っている整数に
# x を掛ける。
# 全ての操作を終了した時に、持っている整数が
# N に一致する確率を
# mod 998244353 で求めてください。
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
