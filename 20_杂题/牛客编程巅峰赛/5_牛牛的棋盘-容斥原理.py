# n*m的矩阵，k个点，将k个点全部放在n*m的矩阵里，求满足以下约束的方案数：
# 矩阵第一行，第一列，最后一行，最后一列都有点。
# 输出方案数对1e9+7的模数

from math import comb


MOD = int(1e9 + 7)


class Solution:
    def solve(self, n, m, k):
        """
        n*m的矩阵，k个点，将k个点全部放在n*m的矩阵里，求满足以下约束的方案数：
        矩阵第一行，第一列，最后一行，最后一列都有点。
        """
        s1 = comb(n * m, k)
        s2 = comb((n - 1) * m, k) * 2 + comb((m - 1) * n, k) * 2
        s3 = comb((n - 1) * (m - 1), k) * 4 + comb((n - 2) * m, k) + comb((m - 2) * n, k)
        s4 = comb((n - 2) * (m - 1), k) * 2 + comb((m - 2) * (n - 1), k) * 2
        s5 = comb((n - 2) * (m - 2), k)
        return (s1 - s2 + s3 - s4 + s5) % MOD

