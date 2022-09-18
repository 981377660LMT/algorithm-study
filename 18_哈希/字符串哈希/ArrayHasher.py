from typing import Sequence


class ArrayHasher:
    _BASE = 131
    _MOD = 2**64
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


if __name__ == "__main__":
    stringHasher = ArrayHasher(sequence=[1, 2, 2])
    print(stringHasher.getHashOfSlice(0, 1))
    print(stringHasher.getHashOfSlice(1, 2))
    print(stringHasher.getHashOfSlice(2, 3))
