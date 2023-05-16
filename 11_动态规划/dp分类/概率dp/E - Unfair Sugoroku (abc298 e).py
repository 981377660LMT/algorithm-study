# 一维数轴，高桥在a号点， 青木在 b号点
# 高桥先，两人轮流扔骰子，高桥的结果会在 1−p等概率出现，青木的在 1−q等概率出现
# 掷出 x后，就会往前走 x步，最多到n号点
# 谁先到n号点谁赢。

# 问高桥赢的概率。

from functools import lru_cache


MOD = 998244353


def unfairSugoroku(n: int, a: int, b: int, p: int, q: int) -> int:
    invP = pow(p, MOD - 2, MOD)
    invQ = pow(q, MOD - 2, MOD)

    @lru_cache(None)
    def dfs(pos1: int, pos2: int, turn: int) -> int:
        if turn == 0:  # takahashi
            if pos1 >= n:
                return 1
            res = 0
            for i in range(1, p + 1):
                nextPos = min(pos1 + i, n)
                if nextPos == n:
                    res += invP
                else:
                    res += (1 - dfs(nextPos, pos2, 1)) * invP % MOD
            return res % MOD
        else:
            res = 0
            for i in range(1, q + 1):
                nextPos = min(pos2 + i, n)
                if nextPos == n:
                    res += invQ
                else:
                    res += (1 - dfs(pos1, nextPos, 0)) * invQ % MOD
            return res % MOD

    res = dfs(a, b, 0)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    n, a, b, p, q = map(int, input().split())
    print(unfairSugoroku(n, a, b, p, q))
