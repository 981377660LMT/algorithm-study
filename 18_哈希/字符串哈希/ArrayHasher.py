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

    def __init__(self, sequence: Sequence[int]):
        self._sequence = sequence
        self._prefix = [0] * (len(sequence) + 1)
        self._base = [0] * (len(sequence) + 1)
        self._prefix[0] = 0
        self._base[0] = 1
        for i in range(1, len(sequence) + 1):
            self._prefix[i] = (
                self._prefix[i - 1] * StringHasher._BASE + sequence[i - 1] - self._OFFSET
            ) % StringHasher._MOD
            self._base[i] = (self._base[i - 1] * StringHasher._BASE) % StringHasher._MOD

    def getHashOfSlice(self, left: int, right: int) -> int:
        """s[left:right]的哈希值"""
        assert 0 <= left <= right <= len(self._sequence)
        left += 1
        upper = self._prefix[right]
        lower = self._prefix[left - 1] * self._base[right - (left - 1)]
        return (upper - lower) % StringHasher._MOD


if __name__ == '__main__':
    stringHasher = StringHasher(sequence=[1, 2, 2])
    print(stringHasher.getHashOfSlice(0, 1))
    print(stringHasher.getHashOfSlice(1, 2))
    print(stringHasher.getHashOfSlice(2, 3))

