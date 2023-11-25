# 2912. 在网格上移动到目的地的方法数
# https://leetcode.cn/problems/number-of-ways-to-reach-destination-in-the-grid/description/

# 现在有两个整数n和m，表示下标起始为1的网格的长和宽，有两个形如[x,y]的列表source和dest各表示一个单元格，还有一个整数k。
# 你可以以下方式移动：
# - 同行移动，从[x,y0]到[x,y1]
# - 同列移动，从[x0,y]到[x1,y]
# - 但不能原地不动，即不能从[x,y]到[x,y]
# 返回【经过恰好k步后，从source移动到dest】的方案数，并对1000000007取余。

# 2 <= n, m <= 1e9
# 1 <= k <= 1e5
# source.length == dest.length == 2
# 1 <= source[1], dest[1] <= n
# 1 <= source[2], dest[2] <= m

# !每次移动必须在“换行”和“换列”中选择一项进行，因此共有(m-1)+(n-1)=m+n-2种方法。


from typing import List


MOD = int(1e9 + 7)


class Solution:
    def numberOfWays(self, n: int, m: int, k: int, source: List[int], dest: List[int]) -> int:
        """
        dp[4] =>
        dp[0]: x和y都不正确
        dp[1]: x正确,y不正确
        dp[2]: x不正确,y正确
        dp[3]: x和y都正确
        """
        sx, sy = source
        dx, dy = dest
        dp = [0] * 4
        dp[(sx == dx) + 2 * (sy == dy)] = 1
        for _ in range(k):
            ndp = [0] * 4
            ndp[0] = dp[0] * (n + m - 4) + dp[1] * (n - 1) + dp[2] * (m - 1)
            ndp[1] = dp[0] + dp[1] * (m - 2) + dp[3] * (m - 1)
            ndp[2] = dp[0] + dp[2] * (n - 2) + dp[3] * (n - 1)
            ndp[3] = dp[1] + dp[2]
            ndp[0] %= MOD
            ndp[1] %= MOD
            ndp[2] %= MOD
            ndp[3] %= MOD
            dp = ndp
        return dp[3] % MOD

    def numberOfWays2(self, n: int, m: int, k: int, source: List[int], dest: List[int]) -> int:
        """状态转移矩阵:
        ```
        n+m-4 | n-1 |m-1 |0
        ----------------------
        1     | m-2 |0   |m-1
        ----------------------
        1     | 0   |n-2 |n-1
        ----------------------
        0     | 1   |1   |0
        ```
        """
        sx, sy = source
        dx, dy = dest
        init = [[0], [0], [0], [0]]
        init[(sx == dx) + 2 * (sy == dy)] = [1]
        T = [[n + m - 4, n - 1, m - 1, 0], [1, m - 2, 0, m - 1], [1, 0, n - 2, n - 1], [0, 1, 1, 0]]
        resT = matpow(T, k, MOD)
        res = matmul(resT, init, MOD)
        return res[3][0]


def matmul(mat1: List[List[int]], mat2: List[List[int]], mod: int) -> List[List[int]]:
    """矩阵相乘"""
    i_, j_, k_ = len(mat1), len(mat2[0]), len(mat2)
    res = [[0] * j_ for _ in range(i_)]
    for i in range(i_):
        for k in range(k_):
            for j in range(j_):
                res[i][j] = (res[i][j] + mat1[i][k] * mat2[k][j]) % mod
    return res


def matpow(base: List[List[int]], exp: int, mod: int) -> List[List[int]]:
    n = len(base)
    e = [[0] * n for _ in range(n)]
    for i in range(n):
        e[i][i] = 1
    b = [row[:] for row in base]
    while exp:
        if exp & 1:
            e = matmul(e, b, mod)
        exp >>= 1
        b = matmul(b, b, mod)
    return e
