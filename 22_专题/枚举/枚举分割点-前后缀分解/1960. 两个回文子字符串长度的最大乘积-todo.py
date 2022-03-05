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


# 哈希
# 你需要找到两个 不重叠的回文 子字符串，它们的长度都必须为 奇数 ，使得它们长度的乘积最大。
# 2 <= s.length <= 105
# https://leetcode-cn.com/problems/maximum-product-of-the-length-of-two-palindromic-substrings/comments/1177272
class Solution:
    def maxProduct(self, s: str) -> int:
        # 1.字符串哈希预处理前缀最长回文长度和后缀最长回文长度
        # 具体操作为：枚举回文中心点+二分长度 nlogn 处理出leftMax 和 rightMax数组
        # 2.枚举分割点 leftMax[i]*rightMax[i+1] 即可
        n = len(s)
        hasher1, hasher2 = StringHasher(s), StringHasher(s[::-1])
        centerMax = [0] * n

        for i in range(n):
            curMax = 1
            left, right = 1, min(n - i, i + 1)
            while left <= right:
                mid = (left + right) >> 1
                hash1 = hasher1.getHashOfSlice(i - mid, i + mid + 1)
                hash2 = hasher2.getHashOfSlice(n - i - mid, n - i + mid + 1)
                if hash1 == hash2:
                    curMax = max(curMax, mid)
                    left = mid + 1
                else:
                    right = mid - 1
            centerMax[i] = curMax

        leftMax, rightMax = [0] * n, [0] * n
        # 如何利用Manacher算法的结果来求出前缀包含的最长回文串长度。
        # 能够形成的最长回文串必然是由最靠前的一个这样的回文串所形成的
        # ...

        res = 0
        for i in range(n - 1):
            left, right = leftMax[i], rightMax[i + 1]
            res = max(res, left * right)
        return res
