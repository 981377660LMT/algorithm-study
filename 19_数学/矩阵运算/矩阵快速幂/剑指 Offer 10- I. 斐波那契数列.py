from typing import List


MOD = int(1e9 + 7)
# js要用Bigint
# F(0) = 0,   F(1) = 1
# F(N) = F(N - 1) + F(N - 2), 其中 N > 1.


# 矩阵快速幂 logn
# 要求：[[1,1],[1,0]]^(n-1)*[[1,0],[0,1]]


Matrix = List[List[int]]


class Solution:
    def fib(self, n: int) -> int:
        if n < 2:
            return n

        def multi(matrix1: Matrix, matrix2: Matrix) -> Matrix:
            """矩阵相乘"""
            row, col = len(matrix1), len(matrix2[0])
            res = [[0] * col for _ in range(row)]
            for r, row in enumerate(matrix1):
                for c in range(col):
                    for k, v in enumerate(row):
                        res[r][c] += v * matrix2[k][c]
                        res[r][c] %= MOD
            return res

        def qpow(a: Matrix, k: int) -> Matrix:
            # 单位矩阵 IdentityMatrix
            res = [[0] * len(a[0]) for _ in range(len(a))]
            for i in range(len(a)):
                res[i][i] = 1

            while k:
                if k & 1:
                    res = multi(res, a)
                k >>= 1
                a = multi(a, a)
            return res

        res = qpow([[1, 1], [1, 0]], n - 1)
        return res[0][0]


print(Solution().fib(2))
# 输出：1
print(Solution().fib(5))
# 输出：5
