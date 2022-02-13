from typing import Sequence


class StringHasher:
    _BASE = 131
    _MOD = 2 ** 64

    @staticmethod
    def setBASE(base: int) -> None:
        StringHasher._BASE = base

    @staticmethod
    def setMOD(mod: int) -> None:
        StringHasher._MOD = mod

    def __init__(self, string: Sequence[str]):
        self._string = string
        self._prefix = [0] * (len(string) + 1)
        self._base = [0] * (len(string) + 1)
        self._prefix[0] = 0
        self._base[0] = 1
        for i in range(1, len(string) + 1):
            self._prefix[i] = (
                self._prefix[i - 1] * StringHasher._BASE + ord(string[i - 1]) - 96
            ) % StringHasher._MOD
            self._base[i] = (self._base[i - 1] * StringHasher._BASE) % StringHasher._MOD

    def getHashOfSlice(self, left: int, right: int) -> int:
        """s[left:right]的哈希值"""
        assert 0 <= left <= right <= len(self._string)
        left += 1
        upper = self._prefix[right]
        lower = self._prefix[left - 1] * self._base[right - (left - 1)]
        return (upper - lower) % StringHasher._MOD


if __name__ == '__main__':
    stringHasher = StringHasher(string='abcdefg')
    print(stringHasher.getHashOfSlice(0, 1))
    print(stringHasher.getHashOfSlice(1, 2))
    print(stringHasher.getHashOfSlice(2, 3))

