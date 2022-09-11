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

from itertools import pairwise
from typing import List
from matqpow import matqpow1

MOD = int(1e9 + 7)

# ndp[i] = (dp[i - 1] + dp[i] + dp[i + 1]) % MOD


class Solution:
    def electricityExperiment(self, row: int, col: int, position: List[List[int]]) -> int:
        trans = [[0] * row for _ in range(row)]
        for i in range(row):
            trans[i][i] = 1
            if i != 0:
                trans[i][i - 1] = 1
            if i != row - 1:
                trans[i][i + 1] = 1

        def cal(row1: int, row2: int, k: int) -> int:
            """row1走到row2的方案数,转移k次"""
            resTrans = matqpow1(trans, k, MOD)  # 矩阵快速幂
            return int(resTrans[row1][row2])

        position.sort(key=lambda x: x[1])
        for (r1, c1), (r2, c2) in pairwise(position):
            rowDiff, colDiff = abs(r1 - r2), abs(c1 - c2)
            if rowDiff > colDiff:
                return 0

        res = 1
        for (r1, c1), (r2, c2) in pairwise(position):
            rowDiff, colDiff = abs(r1 - r2), abs(c1 - c2)
            res = res * cal(r1, r2, colDiff) % MOD
        return res


print(Solution().electricityExperiment(row=5, col=6, position=[[1, 3], [3, 2], [4, 1]]))
print(Solution().electricityExperiment(row=3, col=4, position=[[0, 3], [2, 0]]))
print(Solution().electricityExperiment(row=5, col=6, position=[[1, 3], [3, 5], [2, 0]]))
