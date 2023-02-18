# 完全图的生成树的总数
# ラベル付き木の数え上げ。頂点数と根の次数ごとに数え上げる。
# Cayley 公式 (Cayley's formula)
# 完全图Kn有n^(n-2)个生成树，其中n是顶点数。

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
