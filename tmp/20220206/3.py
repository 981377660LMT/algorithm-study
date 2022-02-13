from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict, SortedSet
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product, chain, islice
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits
from operator import le, xor, or_, and_, not_


MOD = int(1e9 + 7)
INF = 2 ** 64
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Bitset:
    def __init__(self, size: int):
        self.size = size
        self.on = set()
        self.off = set(range(size))

    def fix(self, idx: int) -> None:
        self.on.add(idx)
        self.off.discard(idx)

    def unfix(self, idx: int) -> None:
        self.on.discard(idx)
        self.off.add(idx)

    def flip(self) -> None:
        self.on, self.off = self.off, self.on

    def all(self) -> bool:
        return len(self.off) == 0

    def one(self) -> bool:
        return len(self.on) != 0

    def count(self) -> int:
        return len(self.on)

    def toString(self) -> str:
        return ''.join('1' if i in self.on else '0' for i in range(self.size))


bitset = Bitset(5)
bitset.fix(3)
bitset.fix(1)
print(bitset.toString())
print(bitset.count())
bitset.flip()
print(bitset.one())
# bitset.unfix(0)
# print(bitset.count())
print(bitset.toString())

# bitset = Bitset(2)
# bitset.flip()
# print(bitset.toString())
# bitset.unfix(1)
# print(bitset.toString())

# bitset.fix(1)
# print(bitset.count())
# print(bitset.one())
# # bitset.unfix(0)
# # print(bitset.count())
# print(bitset.toString())
