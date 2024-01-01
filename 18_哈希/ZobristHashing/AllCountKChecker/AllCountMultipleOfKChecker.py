from collections import defaultdict
from random import randint
from typing import List


def countSubarrayWithFrequencyMultipleOfK(arr: List[int], k: int) -> int:
    """统计满足`每个元素出现的次数均为k的倍数`条件的子数组的个数."""
    n = len(arr)
    if n == 0 or k <= 0 or k > n:
        return 0

    pool = defaultdict(lambda: randint(1, (1 << 61) - 1))
    id_ = defaultdict(lambda: len(id_))
    arr = [id_[v] for v in arr]
    counter = [0] * len(id_)
    random = [pool[v] for v in arr]
    preSum = [0] * (n + 1)
    for i, v in enumerate(arr):
        preSum[i + 1] = preSum[i]
        preSum[i + 1] -= counter[v] * random[i]
        counter[v] = (counter[v] + 1) % k
        preSum[i + 1] += counter[v] * random[i]


class AllCountMultipleOfKChecker:
    """判断数据结构中每个数出现的次数是否均k的`倍数`."""

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
