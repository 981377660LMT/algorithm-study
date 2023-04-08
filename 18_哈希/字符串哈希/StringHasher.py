"""
字符串哈希
注意字符串较短时,使用切片更快,时间复杂度约为`O(S/250)`
"""

from typing import Sequence


# 哈希值计算方法：
# hash(s, p, m) = (val(s[0]) * pk-1 + val(s[1]) * pk-2 + ... + val(s[k-1]) * p0) mod m.
# 越靠左字符权重越大
def useStringHasher(ords: Sequence[int], mod=10**11 + 7, base=1313131, offset=0):
    n = len(ords)
    prePow = [1] * (n + 1)
    preHash = [0] * (n + 1)
    for i in range(1, n + 1):
        prePow[i] = (prePow[i - 1] * base) % mod
        preHash[i] = (preHash[i - 1] * base + ords[i - 1] - offset) % mod

    def sliceHash(left: int, right: int):
        """切片 `s[left:right]` 的哈希值"""
        if left >= right:
            return 0
        left += 1
        return (preHash[right] - preHash[left - 1] * prePow[right - left + 1]) % mod

    return sliceHash


class StringHasher:
    __slots__ = ("_ords", "_mod", "_base", "_offset", "_prePow", "_preHash")

    def __init__(self, ords: Sequence[int], mod=10**11 + 7, base=1313131, offset=0):
        self._ords = ords
        self._mod = mod
        self._base = base
        self._offset = offset
        n = len(ords)
        self._prePow = [1] * (n + 1)
        self._preHash = [0] * (n + 1)
        for i in range(1, n + 1):
            self._prePow[i] = (self._prePow[i - 1] * base) % mod
            self._preHash[i] = (self._preHash[i - 1] * base + ords[i - 1] - offset) % mod

    def sliceHash(self, left: int, right: int):
        """切片 `s[left:right]` 的哈希值"""
        if left >= right:
            return 0
        left += 1
        return (
            self._preHash[right] - self._preHash[left - 1] * self._prePow[right - left + 1]
        ) % self._mod

    def __call__(self, left: int, right: int):
        return self.sliceHash(left, right)


def genHash(word: str, mod=10**11 + 7, base=1313131, offset=0) -> int:
    res = 0
    for i in range(len(word)):
        res = (res * base + ord(word[i]) - offset) % mod
    return res


def concatHash(h1: int, h2: int, len2: int, mod=10**11 + 7, base=1313131) -> int:
    """Returns the hash of the concatenation of two strings."""
    return (h1 * pow(base, len2, mod) + h2) % mod


if __name__ == "__main__":
    s = "abc"
    ords = [ord(c) for c in s]
    stringHasher = useStringHasher(ords)
    print(stringHasher(1, 2))
    print(stringHasher(1, 3))
    print(stringHasher(0, 0))
    print(stringHasher(1, 1))
    print(stringHasher(0, 2))
