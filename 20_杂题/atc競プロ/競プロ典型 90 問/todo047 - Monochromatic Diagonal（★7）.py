# 如果行列字符相等 那么就是这种颜色
# 如果行列字符不相等 那么就是另外一种颜色
# !求有多少条正对角线同一颜色组成
# n<=1e6

# 1.如果n很小(n<=2000) 可以全部算出来
# R G B
# 变换规律：表示成0 1 2 结果为-(p1+p2)模三

# !2.n很大的时候需要哈希判断子串相等
# 左下半，判断
# 右上半，判断
from typing import Sequence
import sys

MAPPING = {'B': 0, 'W': 1, 'R': 2}

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


class StringHasher:
    _BASE = 131
    _MOD = 998244353
    _OFFSET = 96

    @staticmethod
    def setBASE(base: int) -> None:
        StringHasher._BASE = base

    @staticmethod
    def setMOD(mod: int) -> None:
        StringHasher._MOD = mod

    @staticmethod
    def setOFFSET(offset: int) -> None:
        StringHasher._OFFSET = offset

    def __init__(self, sequence: Sequence[str]):
        self._sequence = sequence
        self._prefix = [0] * (len(sequence) + 1)
        self._base = [0] * (len(sequence) + 1)
        self._prefix[0] = 0
        self._base[0] = 1
        for i in range(1, len(sequence) + 1):
            self._prefix[i] = (
                self._prefix[i - 1] * StringHasher._BASE + ord(sequence[i - 1]) - self._OFFSET
            ) % StringHasher._MOD
            self._base[i] = (self._base[i - 1] * StringHasher._BASE) % StringHasher._MOD

    def getHashOfSlice(self, left: int, right: int) -> int:
        """s[left:right]的哈希值"""
        assert 0 <= left <= right <= len(self._sequence)
        left += 1
        upper = self._prefix[right]
        lower = self._prefix[left - 1] * self._base[right - (left - 1)]
        return (upper - lower) % StringHasher._MOD


# n = int(input())
# s1 = input()
# s2 = input()

a = StringHasher(sequence='abcdefg')
b = StringHasher(sequence='abcdefg')


# todo 没看懂

