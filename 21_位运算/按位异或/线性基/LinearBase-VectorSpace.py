# https://maspypy.github.io/library/linalg/xor/vector_space.hpp
# 可合并的线性基/线性基合并


from typing import List, Optional


class VectorSpace:
    __slots__ = ("bases", "_max")

    def __init__(self, nums: Optional[List[int]] = None) -> None:
        self.bases = []
        self._max = 0
        if nums is not None:
            for v in nums:
                self.add(v)

    def add(self, v: int) -> bool:
        for e in self.bases:
            if e == 0 or v == 0:
                break
            v = min(v, v ^ e)
        if v:
            self.bases.append(v)
            if v > self._max:
                self._max = v
            return True
        return False

    def getMax(self, xorVal=0) -> int:
        res = xorVal
        for e in self.bases:
            res = max(res, res ^ e)
        return res

    def getMin(self, xorVal=0) -> int:
        res = xorVal
        for e in self.bases:
            res = min(res, res ^ e)
        return res

    def copy(self) -> "VectorSpace":
        res = VectorSpace()
        res.bases = self.bases[:]
        res._max = self._max
        return res

    def _orthogonalSpace(self, maxDim: int) -> "VectorSpace":
        self._normalize()
        m = maxDim
        tmp = [0] * m
        for e in self.bases:
            tmp[e.bit_length() - 1] = e
        tmp = transpose(m, m, tmp, inplace=True)
        res = VectorSpace()
        for j in range(m):
            if tmp[j] & (1 << j):
                continue
            res.add(tmp[j] | (1 << j))
        return res

    def _normalize(self, reverse=True) -> None:
        n = len(self.bases)
        for j in range(n):
            for i in range(j):
                self.bases[i] = min(self.bases[i], self.bases[i] ^ self.bases[j])
        self.bases.sort(reverse=reverse)

    def __len__(self) -> int:
        return len(self.bases)

    def __iter__(self):
        return iter(self.bases)

    def __repr__(self) -> str:
        return repr(self.bases)

    def __contains__(self, v: int) -> bool:
        for e in self.bases:
            if v == 0:
                break
            v = min(v, v ^ e)
        return v == 0

    def __or__(self, other: "VectorSpace") -> "VectorSpace":
        v1, v2 = self, other
        if len(v1) < len(v2):
            v1, v2 = v2, v1
        res = v1.copy()
        for v in v2.bases:
            res.add(v)
        return res

    def __and__(self, other: "VectorSpace") -> "VectorSpace":
        maxDim = max(self._max, other._max).bit_length()
        x = self._orthogonalSpace(maxDim)
        y = other._orthogonalSpace(maxDim)
        if len(x) < len(y):
            x, y = y, x
        for v in y.bases:
            x.add(v)
        return x._orthogonalSpace(maxDim)


def transpose(row: int, col: int, matrix1D: List[int], inplace=False) -> List[int]:
    """矩阵转置

    matrix1D:每个元素是状压的数字
    inplace:是否修改原矩阵
    """
    assert row == len(matrix1D)
    m = matrix1D[:] if not inplace else matrix1D
    log = 0
    max_ = max(row, col)
    while (1 << log) < max_:
        log += 1
    if len(m) < (1 << log):
        m += [0] * ((1 << log) - len(m))
    w = 1 << log
    mask = 1
    for i in range(log):
        mask = mask | (mask << (1 << i))
    for t in range(log):
        w >>= 1
        mask = mask ^ (mask >> w)
        for i in range(1 << t):
            for j in range(w):
                m[w * 2 * i + j] = ((m[w * (2 * i + 1) + j] << w) & mask) ^ m[w * 2 * i + j]
                m[w * (2 * i + 1) + j] = ((m[w * 2 * i + j] & mask) >> w) ^ m[w * (2 * i + 1) + j]
                m[w * 2 * i + j] = ((m[w * (2 * i + 1) + j] << w) & mask) ^ m[w * 2 * i + j]
    return m[:col]


def min(a: int, b: int) -> int:
    return a if a < b else b


def max(a: int, b: int) -> int:
    return a if a > b else b


if __name__ == "__main__":
    a = VectorSpace([1, 2, 3])
    b = VectorSpace([6])

    # # # https://atcoder.jp/contests/abc141/tasks/abc141_f
    # # # !把一个数组分成两个非空子集, 使得两个集合的异或和之和最大

    n = int(input())
    nums = list(map(int, input().split()))

    xor_ = 0
    V1 = VectorSpace()
    for v in nums:
        xor_ ^= v
        V1.add(v)

    mask = ~xor_
    V2 = VectorSpace()
    for v in V1.bases:
        V2.add(v & mask)

    res = V2.getMax()
    print(res + (xor_ ^ res))
