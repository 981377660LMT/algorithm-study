from collections import defaultdict
from random import randint


class AllCountKChecker:
    """判断数据结构中每个数出现的次数是否均k的倍数."""

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = ("_hash", "_counter", "_k")

    def __init__(self, k: int) -> None:
        self._hash = 0
        self._counter = defaultdict(int)
        self._k = k

    def add(self, x: int) -> None:
        count = self._counter[x]
        self._hash ^= self._poolSingleton[(x, count)]
        count += 1
        if count == self._k:
            count = 0
        self._counter[x] = count
        self._hash ^= self._poolSingleton[(x, count)]

    def remove(self, x: int) -> None:
        """删除前需要保证x在集合中存在."""
        count = self._counter[x]
        self._hash ^= self._poolSingleton[(x, count)]
        count -= 1
        if count == -1:
            count = self._k - 1
        self._counter[x] = count
        self._hash ^= self._poolSingleton[(x, count)]

    def query(self) -> bool:
        return self._hash == 0

    def getHash(self) -> int:
        return self._hash
