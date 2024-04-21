# from titan_pylib.data_structures.heap.Randomized_meldable_heap import RandomizedMeldableHeap
from array import array
from __pypy__ import newlist_hint
from typing import TypeVar, Generic, List, Iterable

T = TypeVar("T")


class RandomizedMeldableHeap(Generic[T]):
    """併合可能ヒープです。

    [Randomized Meldable Heap](https://trap.jp/post/1050/), traP
    """

    _x, _y, _z, _w = 123456789, 362436069, 521288629, 88675123
    keys: List[T] = [0]
    child = array("I", bytes(8))
    end = 1

    def __init__(self, a: Iterable[T] = []):
        self.root = 0
        self.size = 0
        for e in a:
            self.heappush(e)

    @classmethod
    def _randbit(cls) -> int:
        t = cls._x ^ (cls._x << 11) & 0xFFFFFFFF
        cls._x, cls._y, cls._z = cls._y, cls._z, cls._w
        cls._w = (cls._w ^ (cls._w >> 19)) ^ (t ^ (t >> 8)) & 0xFFFFFFFF
        return cls._w & 1

    @classmethod
    def _make_node(cls, key: T) -> int:
        if cls.end >= len(cls.keys):
            cls.keys.append(key)
            cls.child.append(0)
            cls.child.append(0)
        else:
            cls.keys[cls.end] = key
        cls.end += 1
        return cls.end - 1

    @classmethod
    def reserve(cls, n: int) -> None:
        if n <= 0:
            return
        cls.keys += [0] * n
        cls.child += array("I", bytes(8 * n))

    @classmethod
    def _meld(cls, x: int, y: int) -> int:
        if x == 0:
            return y
        if y == 0:
            return x
        if cls.keys[x] > cls.keys[y]:
            x, y = y, x
        rand = cls._randbit()
        cls.child[x << 1 | rand] = cls._meld(cls.child[x << 1 | rand], y)
        return x

    @classmethod
    def meld(
        cls, x: "RandomizedMeldableHeap", y: "RandomizedMeldableHeap"
    ) -> "RandomizedMeldableHeap":
        new_heap = RandomizedMeldableHeap()
        new_heap.size = x.size + y.size
        new_heap.root = cls._meld(x.root, y.root)
        return new_heap

    def heappush(self, key: T):
        node = self._make_node(key)
        self.root = self._meld(self.root, node)
        self.size += 1

    def heappop(self) -> T:
        res = self.keys[self.root]
        self.root = self._meld(self.child[self.root << 1], self.child[self.root << 1 | 1])
        return res

    def top(self) -> T:
        return self.keys[self.root]

    def tolist(self) -> List[T]:
        node = self.root
        stack = newlist_hint(len(self))
        res = newlist_hint(len(self))
        child = RandomizedMeldableHeap.child
        keys = RandomizedMeldableHeap.keys
        while stack or node:
            if node:
                stack.append(node)
                node = child[node << 1]
            else:
                node = stack.pop()
                res.append(keys[node])
                node = child[node << 1 | 1]
        res.sort()
        return res

    def __bool__(self):
        return self.root != 0

    def __len__(self):
        return self.size

    def __str__(self):
        return str(self.tolist())

    def __repr__(self):
        return f"RandomizedMeldableHeap({self})"
