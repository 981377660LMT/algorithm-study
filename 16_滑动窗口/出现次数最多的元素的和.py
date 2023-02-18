# 比较繁琐,需要维护三个数组和两个最大值
# counter := make(map[int]int) // 每种颜色出现的次数
# sums := make([]int, n+1)     // 表示出现次数为i的元素的和
# types := make([]int, n+1)    // types[i]表示出现次数为i的元素有几种
# maxCount, sum := 0, 0        // 出现次数最多的元素出现的次数，sum表示出现次数最多的颜色的和


from collections import defaultdict


class MajorSum:
    """出现次数最多的元素的和(多种元素出现次数一样最多也算)."""

    __slots__ = ("_counter", "_freqSum", "_freqTypes", "_maxFreq", "_res")

    def __init__(self) -> None:
        self._counter = defaultdict(int)
        self._freqSum = defaultdict(int)
        self._freqTypes = defaultdict(int)
        self._maxFreq = 0
        self._res = 0

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
            self._res = x
        elif xFreq == self._maxFreq:
            self._res += x

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
            self._res -= x
            if self._freqTypes[self._maxFreq] == 0:
                self._maxFreq -= 1
                self._res = self._freqSum[self._maxFreq]
        if self._counter[x] == 0:
            del self._counter[x]

    def query(self) -> int:
        """返回出现次数最多的元素的和."""
        return self._res

    def __len__(self) -> int:
        """返回当前元素种类数."""
        return len(self._counter)


if __name__ == "__main__":
    # check with brute force
    import random

    while True:
        counter = defaultdict(int)
        ms = MajorSum()
        for _ in range(1000):
            x = random.randint(0, 100)
            if random.random() < 0.5:
                ms.add(x)
                counter[x] += 1
            else:
                ms.discard(x)
                if counter[x] > 0:
                    counter[x] -= 1
            maxFreq = max(counter.values())
            if maxFreq > 0:
                assert ms.query() == sum(x for x, freq in counter.items() if freq == maxFreq)
