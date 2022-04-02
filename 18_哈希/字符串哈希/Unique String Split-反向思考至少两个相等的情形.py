# 字符串有多少种分割方式 使得 a+b b+c c+a 都不同
# n ≤ 10,000 where n is the length of s


from math import comb
from typing import Sequence


class StringHasher:
    _BASE = 131
    _MOD = 2 ** 64
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


# https://binarysearch.com/problems/Unique-String-Split/solutions
# 反向思考:都不等转化为减去至少两个相等的情况，转化为O(n)
class Solution:
    def solve(self, s):
        """只关心其中两个字符串长度相等的情形"""
        n = len(s)
        res = 0
        hasher = StringHasher(s + s)

        for i in range(1, len(s) - 1):
            jCand = set()
            if len(s) > (j1 := len(s) - i) > i:
                jCand.add(j1)
            if len(s) > (j2 := 2 * i) > i:
                jCand.add(j2)
            if ((len(s) + i) % 2) == 0 and len(s) > (j3 := (len(s) + i) // 2) > 0:
                jCand.add(j3)

            for j in jCand:
                hash1, hash2, hash3 = (
                    hasher.getHashOfSlice(0, j),
                    hasher.getHashOfSlice(i, len(s)),
                    hasher.getHashOfSlice(j, len(s) + i),
                )
                if hash1 == hash2 or hash2 == hash3 or hash3 == hash1:
                    res += 1

        return comb(n - 1, 2) - res


print(Solution().solve(s="abba"))
print(Solution().solve(s="aaaa"))
print(Solution().solve(s="cbcb"))
