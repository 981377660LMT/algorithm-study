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


class UnionFind:
    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.size = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.size[rootX] > self.size[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.size[rootY] += self.size[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


# 枚举 words 中的字符串 ss，并枚举 ss 通过添加、
# 删除和替换操作得到的字符串 tt，`如果 tt 也在 words 中，则说明 ss 和 tt 可以分到同一组`


class Solution:
    def groupStrings(self, words: List[str]) -> List[int]:
        n = len(words)
        uf = UnionFind(n)
        wordId = dict()
        wordCounter = Counter()

        for i, word in enumerate(words):
            state = 0
            for char in word:
                state |= 1 << (ord(char) - 97)
            wordId[state] = i
            wordCounter[state] += 1

        # 邻居
        for i, word in enumerate(words):
            raw = wordStates[i]
            addOne = set()
            removeOne = set()
            replaceOne = set()
            # 增
            for char in ascii_lowercase:
                addWord = raw | 1 << (ord(char) - 97)
                if addWord in wordSet:
                    ...
            # 删、改
            for char in word:
                removeOne.add(raw ^ 1 << (ord(char) - 97))
                replaceOne.add(raw ^ 1 << (ord(char) - 97) | 1 << 28)

        return [uf.part, max(uf.size)]


# 2 1 2 1
print(Solution().groupStrings(words=["a", "b", "ab", "cde"]))
print(Solution().groupStrings(words=["a", "ab", "abc"]))
