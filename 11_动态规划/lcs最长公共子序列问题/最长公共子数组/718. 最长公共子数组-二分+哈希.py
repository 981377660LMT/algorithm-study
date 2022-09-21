# 返回 两个数组中 公共的 、长度最长的子数组的长度 。
from functools import lru_cache
from typing import List

# 1923. 最长公共子路径 多个数组的最长公共子数组
from typing import Sequence


class ArrayHasher:
    _BASE = 131
    _MOD = 2 ** 64
    _OFFSET = 96

    @staticmethod
    def setBASE(base: int) -> None:
        ArrayHasher._BASE = base

    @staticmethod
    def setMOD(mod: int) -> None:
        ArrayHasher._MOD = mod

    @staticmethod
    def setOFFSET(offset: int) -> None:
        ArrayHasher._OFFSET = offset

    def __init__(self, sequence: Sequence[int]):
        self._sequence = sequence
        self._prefix = [0] * (len(sequence) + 1)
        self._base = [0] * (len(sequence) + 1)
        self._prefix[0] = 0
        self._base[0] = 1
        for i in range(1, len(sequence) + 1):
            self._prefix[i] = (
                self._prefix[i - 1] * ArrayHasher._BASE + sequence[i - 1] - self._OFFSET
            ) % ArrayHasher._MOD
            self._base[i] = (self._base[i - 1] * ArrayHasher._BASE) % ArrayHasher._MOD

    def getHashOfSlice(self, left: int, right: int) -> int:
        """s[left:right]的哈希值"""
        assert 0 <= left <= right <= len(self._sequence)
        left += 1
        upper = self._prefix[right]
        lower = self._prefix[left - 1] * self._base[right - (left - 1)]
        return (upper - lower) % ArrayHasher._MOD


class Solution:
    def findLength(self, nums1: List[int], nums2: List[int]) -> int:
        """注意不是最长公共子序列 而是最长公共子数组
        
        时间复杂度O((n+m)log(min(n,m)))
        二分答案长度
        """

        def check(mid: int) -> bool:
            visited = set()
            for start in range(len(nums1)):
                if start + mid > len(nums1):
                    break
                hash = hasher1.getHashOfSlice(start, start + mid)
                visited.add(hash)
            for start in range(len(nums2)):
                if start + mid > len(nums2):
                    break
                hash = hasher2.getHashOfSlice(start, start + mid)
                if hash in visited:
                    return True
            return False

        hasher1, hasher2 = ArrayHasher(nums1), ArrayHasher(nums2)
        left, right = 0, max(len(nums1), len(nums2))
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right  # 最右二分

