# 比较繁琐,需要维护三个数组和两个最大值
# counter := make(map[int]int) // 每种颜色出现的次数
# sums := make([]int, n+1)     // 表示出现次数为i的元素的和
# types := make([]int, n+1)    // types[i]表示出现次数为i的元素有几种
# maxCount, sum := 0, 0        // 出现次数最多的元素出现的次数，sum表示出现次数最多的颜色的和


from collections import defaultdict
from typing import Set, Tuple


class MajorFreq:
    """统计一个容器内 出现次数最多的元素的出现次数."""

    __slots__ = ("_counter", "_freqTypes", "_maxFreq")

    def __init__(self) -> None:
        self._counter = defaultdict(int)
        self._freqTypes = defaultdict(int)
        self._maxFreq = 0

    def add(self, x: int) -> "MajorFreq":
        """添加元素x."""
        self._counter[x] += 1
        xFreq = self._counter[x]
        self._freqTypes[xFreq] += 1
        self._freqTypes[xFreq - 1] -= 1
        if xFreq > self._maxFreq:
            self._maxFreq = xFreq
        return self

    def discard(self, x: int) -> bool:
        """删除元素x."""
        if not self._counter[x]:
            return False
        counter, freqTypes = self._counter, self._freqTypes
        counter[x] -= 1
        xFreq = counter[x]
        freqTypes[xFreq] += 1
        freqTypes[xFreq + 1] -= 1
        if xFreq + 1 == self._maxFreq and not freqTypes[self._maxFreq]:
            self._maxFreq -= 1
        if not counter[x]:
            del counter[x]
        return True

    def maxFreq(self) -> int:
        """返回出现次数最多的元素的出现次数."""
        return self._maxFreq


class MajorSum:
    """统计一个容器内 (最多元素出现的次数, 这些元素key的和)."""

    __slots__ = ("_counter", "_freqSum", "_freqTypes", "_maxFreq", "_sum")

    def __init__(self) -> None:
        self._counter = defaultdict(int)
        self._freqSum = defaultdict(int)
        self._freqTypes = defaultdict(int)
        self._maxFreq = 0
        self._sum = 0

    def add(self, x: int) -> None:
        """添加元素x."""
        self._counter[x] += 1
        xFreq = self._counter[x]
        self._freqSum[xFreq] += x
        self._freqSum[xFreq - 1] -= x
        self._freqTypes[xFreq] += 1
        self._freqTypes[xFreq - 1] -= 1
        if xFreq > self._maxFreq:
            self._maxFreq = xFreq
            self._sum = x
        elif xFreq == self._maxFreq:
            self._sum += x

    def discard(self, x: int) -> None:
        """删除元素x."""
        if self._counter[x] == 0:
            return
        self._counter[x] -= 1
        xFreq = self._counter[x]
        self._freqSum[xFreq] += x
        self._freqSum[xFreq + 1] -= x
        self._freqTypes[xFreq] += 1
        self._freqTypes[xFreq + 1] -= 1
        if xFreq + 1 == self._maxFreq:
            self._sum -= x
            if self._freqTypes[self._maxFreq] == 0:
                self._maxFreq -= 1
                self._sum = self._freqSum[self._maxFreq]
        if self._counter[x] == 0:
            del self._counter[x]

    def query(self) -> Tuple[int, int]:
        """返回(最多元素出现的次数, 这些元素key的和)."""
        return self._maxFreq, self._sum

    def __len__(self) -> int:
        """返回当前元素种类数."""
        return len(self._counter)


class MajorManager:
    """
    统计一个容器内 (最多元素出现的次数, 这些元素key的和, 这些元素的集合).
    """

    __slots__ = ("_counter", "_freqSum", "_freqKey", "_maxFreq", "_sum")

    def __init__(self) -> None:
        self._counter = defaultdict(int)
        self._freqSum = defaultdict(int)
        self._freqKey = defaultdict(set)  # 每种出现次数的元素集合
        self._maxFreq = 0
        self._sum = 0

    def add(self, x: int) -> None:
        """添加元素x."""
        self._counter[x] += 1
        xFreq = self._counter[x]
        self._freqSum[xFreq] += x
        self._freqSum[xFreq - 1] -= x
        self._freqKey[xFreq].add(x)
        self._freqKey[xFreq - 1].discard(x)
        if xFreq > self._maxFreq:
            self._maxFreq = xFreq
            self._sum = x
        elif xFreq == self._maxFreq:
            self._sum += x

    def discard(self, x: int) -> None:
        """删除元素x."""
        if self._counter[x] == 0:
            return
        self._counter[x] -= 1
        xFreq = self._counter[x]
        self._freqSum[xFreq] += x
        self._freqSum[xFreq + 1] -= x
        self._freqKey[xFreq].add(x)
        self._freqKey[xFreq + 1].discard(x)
        if xFreq + 1 == self._maxFreq:
            self._sum -= x
            if len(self._freqKey[self._maxFreq]) == 0:
                self._maxFreq -= 1
                self._sum = self._freqSum[self._maxFreq]
        if self._counter[x] == 0:
            del self._counter[x]

    def query(self) -> Tuple[int, int, Set[int]]:
        """返回(最多元素出现的次数, 这些元素key的和, 这些元素的集合)."""
        return self._maxFreq, self._sum, self._freqKey[self._maxFreq]

    def __len__(self) -> int:
        """返回当前元素种类数."""
        return len(self._counter)


if __name__ == "__main__":
    # check with brute force
    import random

    for _ in range(100):
        counter = defaultdict(int)
        ms = MajorSum()
        mm = MajorManager()
        mf = MajorFreq()
        for _ in range(1000):
            x = random.randint(0, 100)
            if random.random() < 0.5:
                ms.add(x)
                mm.add(x)
                mf.add(x)
                counter[x] += 1
            else:
                ms.discard(x)
                mm.discard(x)
                mf.discard(x)
                if counter[x] > 0:
                    counter[x] -= 1
            maxFreq = max(counter.values())
            if maxFreq > 0:
                assert (
                    ms.query()
                    == mm.query()[:2]
                    == (
                        maxFreq,
                        sum(x for x, freq in counter.items() if freq == maxFreq),
                    )
                )

                # 出现次数最多的集合
                assert mm.query()[2] == set(x for x, freq in counter.items() if freq == maxFreq)

                assert mf.maxFreq() == maxFreq

    print("Passed!")
