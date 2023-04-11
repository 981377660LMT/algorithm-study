# https://www.cnblogs.com/bzy-blog/p/14363073.html
# 线性并查集

from typing import Optional


class NextFinder:
    """线性序列并查集LinearSequenceUnionFind (NextFinder)."""

    def __init__(self, n: int):
        self._n = n
        len_ = (n >> 6) + 1
        self.right = list(range(len_))
        MASK = (1 << 64) - 1
        self._data = [MASK] * len_

    def next(self, x: int) -> Optional[int]:
        if x < 0:
            x = 0
        if x >= self._n:
            return None
        div = x >> 6
        mod = x & 63
        mask = self._data[div] >> mod
        if mask:
            # !trailingZeros(mask)
            res = ((div << 6) | mod) + (mask & -mask).bit_length() - 1
            return res if res < self._n else None
        div = self._findNext(div + 1)
        tmp = self._data[div]
        res = (div << 6) + (tmp & -tmp).bit_length() - 1
        return res if res < self._n else None

    def erase(self, x: int) -> None:
        div = x >> 6
        mod = x & 63
        if (self._data[div] >> mod) & 1:  # flip
            self._data[div] ^= 1 << mod
        if not self._data[div]:
            self.right[div] = div + 1  # union to right

    def has(self, x: int) -> bool:
        if x < 0 or x >= self._n:
            return False
        return not not ((self._data[x >> 6] >> (x & 63)) & 1)

    def _findNext(self, x: int) -> int:
        while self.right[x] != x:
            self.right[x] = self.right[self.right[x]]
            x = self.right[x]
        return x

    def __repr__(self):
        res = [i for i in range(self._n) if self.has(i)]
        return f"Finder({list(res)})"

    def __contains__(self, x):
        return self.has(x)


if __name__ == "__main__":
    uf = NextFinder(10)
    uf.erase(0)
    print(uf)

    print(uf.next(0))
    print(uf.next(2))
    print(uf.has(0))
    uf.erase(2)

    print(uf.next(2))
    print(uf.next(9))
    uf.erase(9)
    print(uf.next(9))

    print(uf)
