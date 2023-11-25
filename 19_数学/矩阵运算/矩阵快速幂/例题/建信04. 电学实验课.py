# https://leetcode.cn/contest/ccbft-2021fall/problems/lSjqMF/
# 实验目标要求同学们用导线连接所有「目标插孔」，
# 即从任意一个「目标插孔」沿导线可以到达其他任意「目标插孔」
# 一条导线可连接相邻两列的且行间距不超过 1 的两个插孔
# 每一列插孔中最多使用其中一个插孔（包括「目标插孔」）
# 若实验目标可达成，请返回使用导线数量最少的连接所有目标插孔的方案数；否则请返回 0。

# 1 <= row <= 20
# 3 <= col <= 10^9
# 1 < position.length <= 1000

# !处理出从第i列转移到第i+1列的转移矩阵 这样可以加速算出转移1e9列的方案数
# O(row^3*log(col)*position.length) = 20^3*log(10^9)*1000
# !优化:幂为2^n的矩阵都先算出来，然后快速幂就变成了二进制为1的位对应的矩阵相乘
# https://code.meideng.dev/lc_contest/ccbft-2021fall/problems/lSjqMF/


# ndp[i] = (dp[i - 1] + dp[i] + dp[i + 1]) % MOD


from itertools import pairwise
from typing import List


MOD = int(1e9 + 7)


class Solution:
    def electricityExperiment(self, row: int, col: int, position: List[List[int]]) -> int:
        def cal(row1: int, row2: int, k: int) -> int:
            """row1走到row2的方案数,转移k次"""
            resTrans = matqpow2(npT, k, MOD)
            return int(resTrans[row1][row2])

        position.sort(key=lambda x: x[1])
        T = [[0] * row for _ in range(row)]
        for i in range(row):
            T[i][i] = 1
            if i != 0:
                T[i][i - 1] = 1
            if i != row - 1:
                T[i][i + 1] = 1

        npT = np.array(T, dtype=np.uint64)
        res = 1
        for (r1, c1), (r2, c2) in pairwise(position):
            colDiff = abs(c1 - c2)
            res = res * cal(r1, r2, colDiff) % MOD
        return res


import numpy as np


def matqpow2(base: "np.ndarray", exp: int, mod: int) -> "np.ndarray":
    """np矩阵快速幂"""
    res = np.eye(*base.shape, dtype=np.uint64)
    while exp:
        if exp & 1:
            res = (res @ base) % mod
        base = (base @ base) % mod
        exp >>= 1
    return res


print(Solution().electricityExperiment(row=5, col=6, position=[[1, 3], [3, 2], [4, 1]]))
print(Solution().electricityExperiment(row=3, col=4, position=[[0, 3], [2, 0]]))
print(Solution().electricityExperiment(row=5, col=6, position=[[1, 3], [3, 5], [2, 0]]))
