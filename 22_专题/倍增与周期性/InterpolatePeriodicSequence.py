# InterpolatePeriodicSequence
# 周期序列插值
# !用于发现周期性序列的循环节
# 例如 123[456][456][456]...


from typing import Any, Generic, List, Sequence, TypeVar


T = TypeVar("T")


class InterpolatePeriodicSequence(Generic[T]):
    __slots__ = "_rawSeq", "_offset"

    def __init__(self, seq: Sequence[T]) -> None:
        revSeq = seq[::-1]
        z = zAlgo(revSeq)
        z[0] = 0
        max_, maxIndex = -1, -1
        for i, v in enumerate(z):
            if v >= max_:
                max_ = v
                maxIndex = i
        self._rawSeq = seq
        self._offset = maxIndex

    def get(self, index: int) -> T:
        if index < len(self._rawSeq):
            return self._rawSeq[index]
        k = (index - (len(self._rawSeq) - 1) + self._offset - 1) // self._offset
        index -= k * self._offset
        return self._rawSeq[index]


def zAlgo(seq: Sequence[Any]) -> List[int]:
    n = len(seq)
    z = [0] * n
    left = right = 0
    for i in range(1, n):
        z[i] = max(min(z[i - left], right - i + 1), 0)
        while i + z[i] < n and seq[z[i]] == seq[i + z[i]]:
            left, right = i, i + z[i]
            z[i] += 1
    return z


if __name__ == "__main__":
    arr = [1, 2, 1, 2, 3, 4, 5, 6, 4, 5, 6]
    S = InterpolatePeriodicSequence(arr)
    for i in range(100):
        print(S.get(i), end=" ")
