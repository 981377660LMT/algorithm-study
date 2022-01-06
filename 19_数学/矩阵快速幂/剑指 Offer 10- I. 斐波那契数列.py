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
            res = [[0, 0], [0, 0]]
            for i in range(2):
                for j in range(2):
                    for k in range(2):
                        res[i][j] += matrix1[i][k] * matrix2[k][j]
                        res[i][j] %= MOD
            return res

        def qpow(a: Matrix, k: int) -> Matrix:
            res = [[1, 0], [0, 1]]
            while k > 0:
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
