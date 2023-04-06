# 完全图的生成树的总数
# ラベル付き木の数え上げ。頂点数と根の次数ごとに数え上げる。
# Cayley 公式 (Cayley's formula)
# 完全图Kn有n^(n-2)个生成树，其中n是顶点数。

# !・頂点数を決めたときの木の個数
# N^{N-2}

# !・次数列を決めたときの木の個数
# (N-2)! / prod (d-1)!

# これを取り扱うときには、次数の代わりに、d - 1 を考察するようにするとよい

# !・頂点数と根の次数を決めたときの木の個数
# 次数列を決めたときの～を使うと、多項定理の形が残る。
# count_tree_by_root_degree


from typing import List


MOD = int(1e9 + 7)
fac = [1]
ifac = [1]
for i in range(1, int(4e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def count_tree_by_root_degree(n: int, root_degree: int) -> int:
    d = root_degree
    if d <= 0 or d >= n:
        return 0
    d -= 1
    return fac[n - 2] * ifac[d] * ifac[n - 2 - d] * pow(n - 1, n - 2 - d, MOD)


if __name__ == "__main__":
    # https://yukicoder.me/problems/no/1667
    # !对m = 0, 1, 2, ..., n - 1, 求n个顶点m条边的森林的个数(顶点有区别,边没有区别)
    # n<=300 MOD为素数
    # 1.cayley定理(凯莱定理):n个有标号顶点的树的个数为n^(n-2)
    # 2.dp[i][j]表示i个顶点j条边的森林(无环图)的个数
    # 每次dp转移考虑剩下的顶点中最小的
    # !- dp[i+1][j] += dp[i][j] (加入一个新的顶点,不与任何边相连)
    # !- dp[i+k][j+k-1] += dp[i][j] * C(n-i-1,k-1) * k^(k-2) (加入一个k个顶点的树)
    def countForest(n: int, MOD: int) -> List[int]:
        cayley = [0] * (n + 1)
        cayley[1] = 1
        for i in range(2, n + 1):
            cayley[i] = pow(i, i - 2, MOD)

        dp = [[0] * (n + 1) for _ in range(n + 1)]
        dp[0][0] = 1
        for i in range(n):
            for j in range(i + 1):
                for k in range(1, n + 1):
                    if i + k <= n:
                        dp[i + k][j + k - 1] += dp[i][j] * C(n - i - 1, k - 1) % MOD * cayley[k]
                        dp[i + k][j + k - 1] %= MOD
        return dp[n][:-1]

    n, MOD = map(int, input().split())
    fac = [1]
    ifac = [1]
    for i in range(1, int(1e3) + 10):
        fac.append((fac[-1] * i) % MOD)
        ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)

    def C(n: int, k: int) -> int:
        if n < 0 or k < 0 or n < k:
            return 0
        return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD

    print(*countForest(n, MOD), sep="\n")
