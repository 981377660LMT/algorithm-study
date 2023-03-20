import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 銀行に人
# 1, 人
# 2,
# …, 人
# N が並んでいます。
# Q 個のイベントが発生します。イベントは次の
# 3 種類のいずれかです。

# 1 : 受付に呼ばれていない人のうち、最も小さい番号の人が受付に呼ばれる。
# 2 x : 人
# x が初めて受付に行く。(ここで、人
# x はすでに 1 回以上受付に呼ばれている。)
# 3 : すでに受付に呼ばれているが受付に行っていない人のうち、最も小さい番号の人が再度呼ばれる。
# 3 種類目のイベントで受付に呼ばれる人の番号を呼ばれた順に出力してください。

from heapq import heapify, heappop, heappush
from typing import Generic, Iterable, Literal, Optional, TypeVar

T = TypeVar("T")


class ErasableHeap(Generic[T]):
    __slots__ = ("_data", "_erased")

    def __init__(self, items: Optional[Iterable[T]] = None) -> None:
        self._erased = []
        self._data = [] if items is None else list(items)
        if self._data:
            heapify(self._data)

    def push(self, value: T) -> None:
        heappush(self._data, value)
        self._normalize()

    def pop(self) -> T:
        value = heappop(self._data)
        self._normalize()
        return value

    def peek(self) -> T:
        return self._data[0]

    def discard(self, value: T) -> None:
        """从堆中删除一个元素,要保证堆中存在该元素."""
        heappush(self._erased, value)
        self._normalize()

    def _normalize(self) -> None:
        while self._data and self._erased and self._data[0] == self._erased[0]:
            heappop(self._data)
            heappop(self._erased)

    def __len__(self) -> int:
        return len(self._data)

    def __getitem__(self, index: Literal[0]) -> T:
        return self._data[index]


if __name__ == "__main__":
    n, q = map(int, input().split())
    pq1 = list(range(1, n + 1))
    heapify(pq1)
    pq2 = ErasableHeap()
    for _ in range(q):
        op, *args = map(int, input().split())
        if op == 1:
            pq2.push(heappop(pq1))
        elif op == 2:
            x = args[0]
            pq2.discard(x)
        else:
            print(pq2[0])
