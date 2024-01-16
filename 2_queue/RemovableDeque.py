from collections import deque
from typing import Generic, Optional, Sequence, TypeVar


V = TypeVar("V")


class RemovableDeque(Generic[V]):
    __slots__ = ("_queue", "_removed", "_len")

    def __init__(self, values: Optional[Sequence[V]] = None) -> None:
        self._queue = deque(values) if values is not None else deque()
        self._removed = dict()
        self._len = len(self._queue)

    def append(self, value: V) -> None:
        self._len += 1
        self._queue.append(value)

    def appendleft(self, value: V) -> None:
        self._len += 1
        self._queue.appendleft(value)

    def pop(self) -> V:
        self._len -= 1
        self._normalizeTail()
        return self._queue.pop()

    def popLeft(self) -> V:
        self._len -= 1
        self._normalizeHead()
        return self._queue.popleft()

    def head(self) -> V:
        self._normalizeHead()
        return self._queue[0]

    def tail(self) -> V:
        self._normalizeTail()
        return self._queue[-1]

    def remove(self, value: V) -> None:
        """删除前必须保证value存在于队列."""
        self._len -= 1
        self._removed[value] = self._removed.get(value, 0) + 1

    def empty(self) -> bool:
        return self._len == 0

    def __len__(self) -> int:
        return self._len

    def __repr__(self) -> str:
        tmpRemoved = self._removed.copy()
        res = []
        for v in self._queue:
            if v in tmpRemoved:
                tmpRemoved[v] -= 1
                if tmpRemoved[v] == 0:
                    tmpRemoved.pop(v)
            else:
                res.append(v)
        return repr(res)

    def _normalizeHead(self) -> None:
        while self._queue:
            h = self._queue[0]
            if h in self._removed:
                self._removed[h] -= 1
                if self._removed[h] == 0:
                    self._removed.pop(h)
                self._queue.popleft()
            else:
                break

    def _normalizeTail(self) -> None:
        while self._queue:
            t = self._queue[-1]
            if t in self._removed:
                self._removed[t] -= 1
                if self._removed[t] == 0:
                    self._removed.pop(t)
                self._queue.pop()
            else:
                break


if __name__ == "__main__":
    q1 = RemovableDeque([1, 2, 3, 4, 5])
    print(q1)
    q1.remove(1)
    print(q1)
    q1.remove(3)
    print(q1)
    print(q1.head())

    def checkWithBruteForce() -> None:
        from random import randint

        randomNums = [randint(0, 10) for _ in range(100)]
        realDeque = RemovableDeque(randomNums.copy())
        fakeDeque = randomNums.copy()
        for _ in range(100):
            kind = randint(0, 7)
            if kind == 0:
                if realDeque.empty() != (len(fakeDeque) == 0):
                    print("empty error")
                    exit(0)
            elif kind == 1:
                # "append"
                value = randint(0, 100)
                realDeque.append(value)
                fakeDeque.append(value)
            elif kind == 2:
                # "appendleft"
                value = randint(0, 100)
                realDeque.appendleft(value)
                fakeDeque.insert(0, value)
            elif kind == 3:
                # "pop"
                if realDeque.empty():
                    continue
                if realDeque.pop() != fakeDeque.pop():
                    print(realDeque, fakeDeque)
                    print(len(realDeque), len(fakeDeque))
                    print("pop error")
                    exit(0)
            elif kind == 4:
                # "popleft"
                if realDeque.empty():
                    continue
                if realDeque.popLeft() != fakeDeque.pop(0):
                    print("popLeft error")
                    exit(0)
            elif kind == 5:
                # "head"
                if realDeque.empty():
                    continue
                if realDeque.head() != fakeDeque[0]:
                    print("head error")
                    exit(0)
            elif kind == 6:
                # "tail"
                if realDeque.empty():
                    continue
                if realDeque.tail() != fakeDeque[-1]:
                    print("tail error")
                    exit(0)
            elif kind == 7:
                # "remove"
                if realDeque.empty():
                    continue
                value = randint(0, 100)
                if value in fakeDeque:
                    realDeque.remove(value)
                    fakeDeque.remove(value)

    checkWithBruteForce()
