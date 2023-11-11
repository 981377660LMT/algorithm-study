# 爬s阶(s<=1e18)楼梯 每次只能爬`奇数级`
# 爬楼梯 有些楼梯断了不能爬
# 爬到最后一个位置有多少种方法

# !dp[i] = dp[i-1] + dp[i-3] + dp[i-5] + ...
# !令 ep[i] = dp[i] + dp[i-2] + dp[i-4] + ...
# !那么 dp[i] = ep[i-1] = dp[i-1] + ep[i-3]
# dp[i] ep[i-1] ep[i-2] 的关系可由矩阵快速幂 logS 求出
# 转移矩阵
# 1 0 1
# 1 0 1
# 0 1 0
# 初始向量 [1 0 0]

# !坏的楼梯处
# !矩阵快速幂中间要打断 每次都直接把dp[i]手动赋值为0 再继续算
# 总时间复杂度 O(nlogS)


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

from typing import List


Matrix = List[List[int]]


def mul(m1: Matrix, m2: Matrix, mod: int) -> Matrix:
    """矩阵相乘"""
    ROW, COL = len(m1), len(m2[0])

    res = [[0] * COL for _ in range(ROW)]
    for r in range(ROW):
        for c in range(COL):
            for i in range(ROW):
                res[r][c] = (res[r][c] + m1[r][i] * m2[i][c]) % mod

    return res


def matqpow1(base: Matrix, exp: int, mod: int) -> Matrix:
    """矩阵快速幂"""

    def inner(base: Matrix, exp: int, mod: int) -> Matrix:
        ROW, COL = len(base), len(base[0])
        res = [[0] * COL for _ in range(ROW)]
        for r in range(ROW):
            res[r][r] = 1

        bit = 0
        while exp:
            if exp & 1:
                res = mul(res, pow2[bit], mod)
            exp //= 2
            bit += 1
            if bit == len(pow2):
                pow2.append(mul(pow2[-1], pow2[-1], mod))
        return res

    pow2 = [base]
    return inner(base, exp, mod)


_, s = map(int, input().split())
bad = list(map(int, input().split()))
res = [[1], [0], [0]]  # 3 x 1 答案矩阵
trans = [[1, 0, 1], [1, 0, 1], [0, 1, 0]]
pre = 0
for cur in bad:
    resT = matqpow1(trans, cur - pre, MOD)
    res = mul(resT, res, MOD)
    res[0][0] = 0  # 坏的楼梯不能走
    pre = cur

resT = matqpow1(trans, s - bad[-1], MOD)
res = mul(resT, res, MOD)
print(int(res[0][0]))
