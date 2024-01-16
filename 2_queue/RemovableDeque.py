"""可根据值删除元素的双端队列.(一次删除会删除deque中所有值为value的元素.)"""

from collections import deque
from typing import Generic, Optional, Sequence, TypeVar


V = TypeVar("V")


class RemovableDeque(Generic[V]):
    __slots__ = ("_queue", "_counter", "_removedTime", "_len", "_time")

    def __init__(self, values: Optional[Sequence[V]] = None) -> None:
        self._queue = deque()
        self._counter = dict()
        self._removedTime = dict()
        self._len = 0
        self._time = 0
        if values is not None:
            for v in values:
                self.append(v)

    def append(self, value: V) -> None:
        self._len += 1
        self._queue.append((value, self._time))  # type: ignore
        self._counter[value] = self._counter.get(value, 0) + 1

    def appendleft(self, value: V) -> None:
        self._len += 1
        self._queue.appendleft((value, self._time))  # type: ignore
        self._counter[value] = self._counter.get(value, 0) + 1

    def pop(self) -> V:
        self._len -= 1
        self._normalizeTail()
        res = self._queue.pop()[0]
        if res in self._counter:
            self._counter[res] -= 1
            if self._counter[res] == 0:
                del self._counter[res]
        return res

    def popLeft(self) -> V:
        self._len -= 1
        self._normalizeHead()
        res = self._queue.popleft()[0]
        if res in self._counter:
            self._counter[res] -= 1
            if self._counter[res] == 0:
                del self._counter[res]
        return res

    def head(self) -> V:
        self._normalizeHead()
        return self._queue[0][0]

    def tail(self) -> V:
        self._normalizeTail()
        return self._queue[-1][0]

    def remove(self, value: V) -> None:
        """删除所有值为value的元素."""
        if value not in self._counter:
            return
        self._len -= self._counter[value]
        del self._counter[value]
        self._removedTime[value] = self._time
        self._time += 1

    def empty(self) -> bool:
        return self._len == 0

    def count(self, value: V) -> int:
        return self._counter.get(value, 0)

    def __len__(self) -> int:
        return self._len

    def __repr__(self) -> str:
        res = []
        for v, t in self._queue:
            if v in self._removedTime and t <= self._removedTime[v]:
                continue
            res.append(v)
        return repr(res)

    def _normalizeHead(self) -> None:
        while self._queue:
            v, t = self._queue[0]
            if v in self._removedTime and t <= self._removedTime[v]:
                self._queue.popleft()
            else:
                break

    def _normalizeTail(self) -> None:
        while self._queue:
            v, t = self._queue[-1]
            if v in self._removedTime and t <= self._removedTime[v]:
                self._queue.pop()
            else:
                break


if __name__ == "__main__":
    q1 = RemovableDeque([1, 2])
    print(q1)
    q1.remove(2)
    print(q1)
    print(q1.pop())
    # q1.remove(3)
    # q1.append(6)
    # print(q1)

    def checkWithBf() -> None:
        from random import randint

        initNums = [randint(1, 100) for _ in range(100)]
        real = RemovableDeque(initNums[:])
        fake = initNums[:]
        for _ in range(100000):
            kind = randint(0, 7)
            if kind == 0:
                if real.empty() != (not fake):
                    raise RuntimeError("empty")
            elif kind == 1:
                """append"""
                v = randint(1, 100)
                real.append(v)
                fake.append(v)
            elif kind == 2:
                """appendleft"""
                v = randint(1, 100)
                real.appendleft(v)
                fake.insert(0, v)
            elif kind == 3:
                """pop"""
                if real.empty():
                    continue
                if real.pop() != fake.pop():
                    raise RuntimeError("pop")
            elif kind == 4:
                """popleft"""
                if real.empty():
                    continue
                if real.popLeft() != fake.pop(0):
                    raise RuntimeError("popleft")
            elif kind == 5:
                """head"""
                if real.empty():
                    continue
                if real.head() != fake[0]:
                    raise RuntimeError("head")
            elif kind == 6:
                """tail"""
                if real.empty():
                    continue
                if real.tail() != fake[-1]:
                    raise RuntimeError("tail")
            elif kind == 7:
                """remove"""
                if real.empty():
                    continue
                v = randint(1, 100)
                real.remove(v)
                while v in fake:
                    fake.remove(v)
        print("OK")

    checkWithBf()
