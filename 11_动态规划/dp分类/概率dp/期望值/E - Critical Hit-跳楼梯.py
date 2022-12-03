# !跳楼梯,每一步有p/100 的概率跳2步,1-p/100的概率跳1步,求跳完n节楼梯的期望值 模998244353

from functools import lru_cache
import sys


MOD = 998244353
INV_100 = pow(100, MOD - 2, MOD)


def criticalHit(monsterHp: int, p: int) -> int:
    @lru_cache(None)
    def dfs(remain: int) -> int:
        if remain <= 0:
            return 0
        res1 = (dfs(remain - 2) + 1) * p  # 2 減らし
        res2 = (dfs(remain - 1) + 1) * (100 - p)  # 1 減らし
        return (res1 + res2) * INV_100 % MOD

    res = dfs(monsterHp)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n, p = map(int, input().split())
    print(criticalHit(n, p))
