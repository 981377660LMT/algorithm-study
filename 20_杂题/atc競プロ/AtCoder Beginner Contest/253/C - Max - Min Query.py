"""
多重集合 multiset
1 x : S に x を 1 個追加する。
2 x c : S から x を min(c, (S に含まれる x の個数 )) 個削除する。
3 : (S の最大値 ) - (S の最小値 ) を出力する。このクエリを処理するとき、 S が空でないことが保証される。
"""


import sys
import os

from sortedcontainers import SortedList


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
########################################################################


from collections import defaultdict
from typing import Any, Callable, Generic, Iterable, Optional, TypeVar

T = TypeVar("T")


class MultiSet(SortedList, Generic[T]):
    """SortedList + 维护一个计数器"""

    def __init__(
        self, iterable: Optional[Iterable[T]] = None, key: Optional[Callable[[T], Any]] = None
    ):
        super().__init__(iterable, key)
        self._counter = defaultdict(int)
        for item in self:
            self._counter[item] += 1

    def add(self, value: T) -> None:
        raise NotImplementedError("use ``add_with_count`` instead")

    def add_with_count(self, value: T, count: int) -> None:
        assert count >= 0, "count must be non-negative"
        if count == 0:
            return
        if self._counter[value] == 0:
            super().add(value)
        self._counter[value] += count

    def remove(self, value: T) -> None:
        raise NotImplementedError("use ``remove_with_count`` instead")

    def remove_with_count(self, value: T, count: int) -> bool:
        assert count >= 0, "count must be non-negative"
        remove_count = min(count, self._counter[value])
        if remove_count == 0:
            return False
        self._counter[value] -= remove_count
        if self._counter[value] == 0:
            super().remove(value)
        return True

    def count(self, value: T) -> int:
        return self._counter[value]

    def __contains__(self, value: T) -> bool:
        return self._counter[value] > 0


########################################################################


def main() -> None:
    ms = MultiSet()
    q = int(input())
    for _ in range(q):
        qt, *rest = map(int, input().split())
        if qt == 1:
            value = rest[0]
            ms.add_with_count(value, 1)
        elif qt == 2:
            value, count = rest
            ms.remove_with_count(value, count)
        else:
            print(ms[-1] - ms[0])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
