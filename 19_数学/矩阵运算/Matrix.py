from typing import List, Optional, Tuple


MOD = 998244353


class Matrix:
    """https://github.com/strangerxxxx/kyopro"""

    __slots__ = "_n", "_m", "_mat"

    @staticmethod
    def from_list(list: List[List[int]]) -> "Matrix":
        return Matrix(len(list), len(list[0]), list)

    @staticmethod
    def id(n: int) -> "Matrix":
        res = Matrix(n, n)
        for i in range(n):
            res[i][i] = 1
        return res

    def __init__(self, row: int, col: int, mat=None):
        self._n = row
        self._m = col
        self._mat = [[0] * self._m for _ in range(self._n)]
        if mat:
            # assert len(mat) == n and len(mat[0]) == m
            for i in range(self._n):
                self._mat[i] = mat[i].copy()

    def is_square(self) -> bool:
        return self._n == self._m

    def determinant(self) -> int:
        """行列式"""
        assert self.is_square()
        res = 1
        tmp = Matrix(self._n, self._n, self._mat)
        for j in range(self._n):
            if tmp[j][j] == 0:
                for i in range(j + 1, self._n):
                    if tmp[i][j]:
                        break
                else:
                    return 0
                tmp._mat[j], tmp._mat[i] = tmp._mat[i], tmp._mat[j]
                res *= -1
            tmp_j = tmp[j]
            inv = pow(tmp_j[j], MOD - 2, MOD)
            for i in range(j + 1, self._n):
                tmp_i = tmp[i]
                c = -inv * tmp_i[j] % MOD
                for k in range(self._n):
                    tmp_i[k] += c * tmp_j[k]
                    tmp_i[k] %= MOD
        for i in range(self._n):
            res *= tmp[i][i]
            res %= MOD
        return res

    def inverse(self) -> Optional["Matrix"]:
        """矩阵的逆 不存在返回None"""
        assert self.is_square()
        res = Matrix.id(self._n)
        tmp = Matrix(self._n, self._n, self._mat)
        for j in range(self._n):
            if tmp[j][j] == 0:
                for i in range(j + 1, self._n):
                    if tmp[i][j]:
                        break
                else:
                    return None
                tmp._mat[j], tmp._mat[i] = tmp._mat[i], tmp._mat[j]
                res._mat[j], res._mat[i] = res._mat[i], res._mat[j]
            tmp_j, res_j = tmp[j], res[j]
            inv = pow(tmp_j[j], MOD - 2, MOD)
            for k in range(self._n):
                tmp_j[k] *= inv
                tmp_j[k] %= MOD
                res_j[k] *= inv
                res_j[k] %= MOD
            for i in range(self._n):
                if i == j:
                    continue
                c = tmp[i][j]
                tmp_i, res_i = tmp[i], res[i]
                for k in range(self._n):
                    tmp_i[k] -= tmp_j[k] * c
                    tmp_i[k] %= MOD
                    res_i[k] -= res_j[k] * c
                    res_i[k] %= MOD
        return res

    def linear_equations(self, vec) -> Tuple[int, List[int], List[List[int]]]:
        """
        解线性方程组 Ax = vec
        返回解的维度, 解, 基础解系
        如果无解返回-1, [], []
        """
        assert self._n == len(vec)
        aug = [self[i] + [vec[i]] for i in range(self._n)]
        rank = 0
        p = []
        q = []
        for j in range(self._m + 1):
            for i in range(rank, self._n):
                if aug[i][j]:
                    break
            else:
                q.append(j)
                continue
            if j == self._m:
                return -1, [], []
            p.append(j)
            aug[rank], aug[i] = aug[i], aug[rank]
            inv = pow(aug[rank][j], MOD - 2, MOD)
            aug_rank = aug[rank]
            for k in range(self._m + 1):
                aug_rank[k] *= inv
                aug_rank[k] %= MOD
            for i in range(self._n):
                if i == rank:
                    continue
                aug_i = aug[i]
                c = -aug_i[j]
                for k in range(self._m + 1):
                    aug_i[k] += c * aug_rank[k]
                    aug_i[k] %= MOD
            rank += 1
        dim = self._m - rank
        sol = [0] * self._m
        for i in range(rank):
            sol[p[i]] = aug[i][-1]
        vecs = [[0] * self._m for _ in range(dim)]
        for i in range(dim):
            vecs[i][q[i]] = 1
        for i in range(dim):
            vecs_i = vecs[i]
            for j in range(rank):
                vecs_i[p[j]] = -aug[j][q[i]] % MOD
        return dim, sol, vecs

    def times(self, k):
        res = [[0] * self._m for _ in range(self._n)]
        for i in range(self._n):
            res_i, self_i = res[i], self[i]
            for j in range(self._m):
                res_i[j] = k * self_i[j] % MOD
        return Matrix(self._n, self._m, res)

    def __getitem__(self, key):
        if not isinstance(key, slice):
            assert key >= 0
        return self._mat[key]

    def __len__(self):
        return len(self._mat)

    def __str__(self):
        return "\n".join(" ".join(map(str, self[i])) for i in range(self._n))

    def __pos__(self):
        return self

    def __neg__(self):
        return self.times(-1)

    def __add__(self, other):
        # assert self._n == other._n and self._m == other._m
        res = [[0] * self._m for _ in range(self._n)]
        for i in range(self._n):
            res_i, self_i, other_i = res[i], self[i], other[i]
            for j in range(self._m):
                res_i[j] = (self_i[j] + other_i[j]) % MOD
        return Matrix(self._n, self._m, res)

    def __sub__(self, other):
        # assert self._n == other._n and self._m == other._m
        res = [[0] * self._m for _ in range(self._n)]
        for i in range(self._n):
            res_i, self_i, other_i = res[i], self[i], other[i]
            for j in range(self._m):
                res_i[j] = (self_i[j] - other_i[j]) % MOD
        return Matrix(self._n, self._m, res)

    def __mul__(self, other):
        if other.__class__ == Matrix:
            # assert self._m == other._n
            res = [[0] * other._m for _ in range(self._n)]
            for i in range(self._n):
                res_i, self_i = res[i], self[i]
                for k in range(self._m):
                    self_ik, other_k = self_i[k], other[k]
                    for j in range(other._m):
                        res_i[j] += self_ik * other_k[j]
                        res_i[j] %= MOD
            return Matrix(self._n, other._m, res)
        else:
            return self.times(other)

    def __rmul__(self, other):
        return self.times(other)

    def __pow__(self, k):
        # assert self.is_square()
        if k == -1:
            return self.inverse()
        tmp = Matrix(self._n, self._n, self._mat)
        res = Matrix.id(self._n)
        while k:
            if k & 1:
                res *= tmp
            tmp *= tmp
            k >>= 1
        return res


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

if __name__ == "__main__":
    # # 矩阵的逆
    # n = int(input())
    # grid = [list(map(int, input().split())) for _ in range(n)]
    # A = Matrix.from_list(grid)
    # inv = A.inverse()
    # print(-1 if inv is None else inv)

    # System of Linear Equations
    # 解方程 Ax = b
    n, m = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(n)]
    A = Matrix.from_list(grid)
    b = list(map(int, input().split()))
    dim, sol, vecs = A.linear_equations(b)
    print(dim)
    print(*sol)
    for v in vecs:
        print(*v)
