"""
n<=3000

O(n^3/64) 4e8 
"""


from collections import defaultdict
from itertools import zip_longest
import sys
import os
from typing import List, Union

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


class BitSet:
    """
    use int instead of bitset when python version is >= 3.10;

    `int.bit_count` is faster than `bitset.bit_count`.
    """

    __slots__ = "_bin", "_buckets", "_len"

    def __init__(self, str_or_int: Union[str, int] = "") -> None:
        """
        `63` bit per bucket;
        pypy3 will get slow when int is bigger than `(1 << 63) - 1`.

        `size` of bitset is not maintained.
        """
        self._bin = str_or_int if isinstance(str_or_int, str) else bin(str_or_int)[2:]
        self._buckets: List[int] = []  # little endian
        self._len = 0
        for i in range(0, len(self._bin), 63):
            group = int(self._bin[i : i + 63], 2)
            self._buckets.append(group)
            self._len += self._longlong_bit_count(group)

    @staticmethod
    def _longlong_bit_count(n: int) -> int:
        """`O(1)` counts bit of int smaller than `(1 << 63) - 1`"""
        c = (n & 0x5555555555555555) + ((n >> 1) & 0x5555555555555555)
        c = (c & 0x3333333333333333) + ((c >> 2) & 0x3333333333333333)
        c = (c & 0x0F0F0F0F0F0F0F0F) + ((c >> 4) & 0x0F0F0F0F0F0F0F0F)
        c = (c & 0x00FF00FF00FF00FF) + ((c >> 8) & 0x00FF00FF00FF00FF)
        c = (c & 0x0000FFFF0000FFFF) + ((c >> 16) & 0x0000FFFF0000FFFF)
        c = (c & 0x00000000FFFFFFFF) + ((c >> 32) & 0x00000000FFFFFFFF)
        return c

    def add(self, n: int) -> bool:
        row, col = n // 63, n % 63
        if n in self:
            return False
        self._ensure_capacity(row + 1)
        self._buckets[row] |= 1 << col
        self._len += 1
        return True

    def discard(self, n: int) -> bool:
        row, col = n // 63, n % 63
        if len(self._buckets) <= row or n not in self:
            return False
        self._buckets[row] &= ~(1 << col)
        self._len += 1
        return True

    def _bit_count(self) -> int:
        """`O(len(self._bin / 63))`"""
        res = 0
        for bucket in self._buckets:
            res += self._longlong_bit_count(bucket)
        return res

    def _ensure_capacity(self, n: int) -> None:
        if len(self._buckets) < n:
            self._buckets.extend([0] * (n - len(self._buckets)))

    def __and__(self, other: "BitSet") -> int:
        res = 0
        for b1, b2 in zip(self._buckets, other._buckets):
            res += self._longlong_bit_count(b1 & b2)
        return res

    def __or__(self, other: "BitSet") -> int:
        res = 0
        for b1, b2 in zip_longest(self._buckets, other._buckets, fillvalue=0):
            res += self._longlong_bit_count(b1 | b2)
        return res

    def __xor__(self, other: "BitSet") -> int:
        res = 0
        for b1, b2 in zip_longest(self._buckets, other._buckets, fillvalue=0):
            res += self._longlong_bit_count(b1 ^ b2)
        return res

    def __contains__(self, n: int) -> bool:
        row, col = n // 63, n % 63
        return len(self._buckets) > row and not not (self._buckets[row] & (1 << col))

    def __len__(self) -> int:
        return self._len


def main() -> None:
    n = int(input())
    edges = []
    adjMap = defaultdict(BitSet)
    for i in range(n):
        row = input()
        for j in range(i + 1, n):
            if row[j] == "1":
                edges.append((i, j))
                adjMap[i].add(j)  # 无向边定向

    res = 0
    for p1, p2 in edges:
        res += adjMap[p1] & adjMap[p2]
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
