# 判断数据结构中每个数出现的次数是否均为k.
# 等价于:
#  1. 数据结构中每个数出现的次数均为k的倍数：异或哈希.
#  2. 数据结构中每个数出现的次数均不超过k：双指针.
#     在右指针扫到 i 的时候，不停将左指针向右移动并减去这个桶的出现次数，
#     直到 nums[i] 的出现次数小于等于 k 为止。此时再统计答案，两个限制都可以满足。

from heapq import heapify, heappop, heappush
from typing import Generic, Iterable, Literal, Optional, TypeVar, List
from collections import defaultdict
from random import randint


def countSubarrayWithFrequencyEqualToK(arr: List[int], k: int) -> int:
    """统计满足`每个元素出现的次数均为k`条件的子数组的个数."""
    n = len(arr)
    if n == 0 or k <= 0 or k > n:
        return 0

    pool = defaultdict(lambda: randint(1, (1 << 61) - 1))
    id_ = defaultdict(lambda: len(id_))
    arr = [id_[v] for v in arr]
    counter = [0] * len(id_)
    random = [pool[v] for v in arr]
    hashPreSum = [0] * (n + 1)  # 哈希之和的前缀和
    for i, v in enumerate(arr):
        hashPreSum[i + 1] = hashPreSum[i]
        hashPreSum[i + 1] -= counter[v] * random[i]
        counter[v] = (counter[v] + 1) % k
        hashPreSum[i + 1] += counter[v] * random[i]

    countPreSum = defaultdict(int, {0: 1})
    counter = [0] * len(id_)
    res, left = 0, 0
    for right, num in enumerate(arr):
        counter[num] += 1
        while counter[num] > k:
            counter[arr[left]] -= 1
            countPreSum[hashPreSum[left]] -= 1
            left += 1
        res += countPreSum[hashPreSum[right + 1]]
        countPreSum[hashPreSum[right + 1]] += 1
    return res


class AllCountKChecker:
    """
    判断数据结构中每个数出现的次数是否均恰好等于k(k>0).
    如果为空集合,则返回True.
    """

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = ("_hash", "_counter", "_modCounter", "_k", "_countPq")

    def __init__(self, k: int) -> None:
        self._hash = 0
        self._counter = defaultdict(int)
        self._modCounter = defaultdict(int)
        self._k = k
        self._countPq = ErasableHeap()

    def add(self, x) -> None:
        count, random = self._modCounter[x], self._poolSingleton[x]
        self._hash -= count * random
        count += 1
        if count == self._k:
            count = 0
        self._hash += count * random
        self._modCounter[x] = count

        preCount = self._counter[x]
        self._counter[x] = preCount + 1
        if preCount > 0:
            self._countPq.remove(-preCount)
        self._countPq.push(-(preCount + 1))

    def remove(self, x) -> None:
        """删除前需要保证x在集合中存在."""
        count, random = self._modCounter[x], self._poolSingleton[x]
        self._hash -= count * random
        count -= 1
        if count == -1:
            count = self._k - 1
        self._hash += count * random
        self._modCounter[x] = count

        preCount = self._counter[x]
        self._counter[x] = preCount - 1
        if preCount == 1:
            self._counter.pop(x)
        self._countPq.remove(-preCount)
        if preCount > 1:
            self._countPq.push(-(preCount - 1))

    def query(self) -> bool:
        if not self._countPq:
            return True
        return self._hash == 0 and -self._countPq.peek() == self._k

    def getHash(self) -> int:
        return self._hash

    def clear(self) -> None:
        self._hash = 0
        self._modCounter.clear()
        self._counter.clear()
        self._countPq.clear()

    def __repr__(self) -> str:
        return f"hash:{self._hash}, counter:{self._counter}"


T = TypeVar("T")


class ErasableHeap(Generic[T]):
    __slots__ = ("_data", "_erased", "_size")

    def __init__(self, items: Optional[Iterable[T]] = None) -> None:
        self._erased = []
        self._data = [] if items is None else list(items)
        if self._data:
            heapify(self._data)
        self._size = len(self._data)

    def push(self, value: T) -> None:
        heappush(self._data, value)
        self._normalize()
        self._size += 1

    def pop(self) -> T:
        value = heappop(self._data)
        self._normalize()
        self._size -= 1
        return value

    def peek(self) -> T:
        return self._data[0]

    def remove(self, value: T) -> None:
        """从堆中删除一个元素,要保证堆中存在该元素."""
        heappush(self._erased, value)
        self._normalize()
        self._size -= 1

    def clear(self) -> None:
        self._data.clear()
        self._erased.clear()
        self._size = 0

    def _normalize(self) -> None:
        while self._data and self._erased and self._data[0] == self._erased[0]:
            heappop(self._data)
            heappop(self._erased)

    def __len__(self) -> int:
        return self._size

    def __getitem__(self, index: Literal[0]) -> T:
        return self._data[index]


if __name__ == "__main__":

    def check() -> None:
        for k in range(1, 10):
            container = defaultdict(int)
            C = AllCountKChecker(k)

            for _ in range(100000):
                # add
                if randint(1, 5) >= 3:
                    x = randint(1, 10)
                    container[x] += 1
                    C.add(x)
                # remove
                else:
                    x = randint(1, 10)
                    if x in container:
                        container[x] -= 1
                        if container[x] == 0:
                            del container[x]
                        C.remove(x)
                if C.query() != all(freq == k for freq in container.values()):
                    print(C.query(), all(freq == k for freq in container.values()))
                    print(container)
                    print(C)
                    print(k)
                    raise ValueError("error")

        print("ok")

    check()

    # https://leetcode.cn/problems/count-complete-substrings

    class Solution:
        def countCompleteSubstrings(self, word: str, k: int) -> int:
            n = len(word)
            ords = [ord(x) - 97 for x in word]
            groups = []
            ptr = 0
            while ptr < n:
                group = [ords[ptr]]
                ptr += 1
                while ptr < n and abs(ords[ptr] - ords[ptr - 1]) <= 2:
                    group.append(ords[ptr])
                    ptr += 1
                groups.append(group)
            return sum(countSubarrayWithFrequencyEqualToK(group, k) for group in groups)

    # https://www.luogu.com.cn/problem/CF1418G
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    arr = list(map(int, input().split()))
    print(countSubarrayWithFrequencyEqualToK(arr, 3))
