from typing import Generic, List, TypeVar
from bisect import bisect_left
from collections import defaultdict

T = TypeVar("T")


class RangeFreq(Generic[T]):
    __slots__ = "_valueToIndexes"

    def __init__(self, nums: List[T]) -> None:
        self._valueToIndexes = defaultdict(list)
        for i, v in enumerate(nums):
            self._valueToIndexes[v].append(i)

    def query(self, start: int, end: int, value: T) -> int:
        if start >= end:
            return 0
        pos = self._valueToIndexes[value]
        return bisect_left(pos, end) - bisect_left(pos, start)


if __name__ == "__main__":
    rf = RangeFreq([1, 2, 3, 4, 5])
    print(rf.query(2, 3, 3))
