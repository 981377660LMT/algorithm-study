# 矩阵行列式
from typing import List


def calDeterminant(matrix: List[List[int]], mod: int) -> int:
    n = len(matrix)
    res = 1
    for i in range(n):
        for j in range(i + 1, n):
            while matrix[j][i]:
                tmp = matrix[i][i] // matrix[j][i]
                if tmp:
                    for k in range(i, n):
                        matrix[i][k] -= tmp * matrix[j][k]
                        matrix[i][k] %= mod
                matrix[i], matrix[j] = matrix[j], matrix[i]
                res *= -1
                res %= mod
        res *= matrix[i][i]
        res %= mod
        if not res:
            break
    return res


if __name__ == "__main__":
    N, MOD = map(int, input().split())
    A = [list(map(int, input().split())) for _ in range(N)]
    print(calDeterminant(A, MOD))
