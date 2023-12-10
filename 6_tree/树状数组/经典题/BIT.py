"""
下标从0开始
1.BITArray: 单点修改, 区间查询
2.BITMap: 单点修改, 区间查询
3.BITRangeAddPointGetArray: 区间修改, 单点查询(差分)
4.BITRangeAddPointGetMap: 区间修改, 单点查询(差分)
5.BITRangeAddRangeSumArray: 区间修改, 区间查询
6.BITRangeAddRangeSumMap: 区间修改, 区间查询
7.BITPrefixArray: 单点修改, 前缀查询
8.BITPrefixMap: 单点修改, 前缀查询
"""


from typing import Callable, Generic, Optional, TypeVar


class BITArray:
    __slots__ = ("n", "_data", "_total")

    def __init__(self, n: int, f: Optional[Callable[[int], int]] = None):
        if f is None:
            self.n = n
            self._data = [0] * n
            self._total = 0
        else:
            self.n = n
            self._data = [f(i) for i in range(n)]
            self._total = sum(self._data)
            for i in range(1, n + 1):
                j = i + (i & -i)
                if j <= n:
                    self._data[j - 1] += self._data[i - 1]

    def add(self, index: int, v: int) -> None:
        self._total += v
        index += 1
        while index <= self.n:
            self._data[index - 1] += v
            index += index & -index

    def queryPrefix(self, end: int) -> int:
        if end > self.n:
            end = self.n
        res = 0
        while end > 0:
            res += self._data[end - 1]
            end -= end & -end
        return res

    def queryRange(self, start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return 0
        if start == 0:
            return self.queryPrefix(end)
        pos, neg = 0, 0
        while end > start:
            pos += self._data[end - 1]
            end &= end - 1
        while start > end:
            neg += self._data[start - 1]
            start &= start - 1
        return pos - neg

    def queryAll(self) -> int:
        return self._total

    def maxRight(self, check: Callable[[int, int], bool]) -> int:
        """查询满足check的最右位置(不包含), check(index, preSum)."""
        i = 0
        s = 0
        k = 1
        while 2 * k <= self.n:
            k *= 2
        while k > 0:
            if i + k - 1 < self.n:
                t = s + self._data[i + k - 1]
                if check(i + k, t):
                    i += k
                    s = t
            k >>= 1
        return i

    def kth(self, k: int) -> int:
        """01树状数组查找第 k(0-based) 个1的位置."""
        return self.maxRight(lambda _, preSum: preSum <= k)

    def __repr__(self):
        return (
            "BitArray: [" + ", ".join(str(self.queryRange(i, i + 1)) for i in range(self.n)) + "]"
        )


class BITMap:
    __slots__ = ("n", "_data", "_total")

    def __init__(self, n: int):
        self.n = n
        self._data = dict()
        self._total = 0

    def add(self, index: int, v: int) -> None:
        self._total += v
        index += 1
        while index <= self.n:
            self._data[index - 1] = self._data.get(index - 1, 0) + v
            index += index & -index

    def queryPrefix(self, end: int) -> int:
        if end > self.n:
            end = self.n
        res = 0
        while end > 0:
            res += self._data.get(end - 1, 0)
            end -= end & -end
        return res

    def queryRange(self, start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return 0
        if start == 0:
            return self.queryPrefix(end)
        pos, neg = 0, 0
        while end > start:
            pos += self._data.get(end - 1, 0)
            end &= end - 1
        while start > end:
            neg += self._data.get(start - 1, 0)
            start &= start - 1
        return pos - neg

    def queryAll(self) -> int:
        return self._total

    def maxRight(self, check: Callable[[int, int], bool]) -> int:
        """查询满足check的最右位置(不包含), check(index, preSum)."""
        i = 0
        s = 0
        k = 1
        while 2 * k <= self.n:
            k *= 2
        while k > 0:
            if i + k - 1 < self.n:
                t = s + self._data.get(i + k - 1, 0)
                if check(i + k, t):
                    i += k
                    s = t
            k >>= 1
        return i

    def kth(self, k: int) -> int:
        """01树状数组查找第 k(0-based) 个1的位置."""
        return self.maxRight(lambda _, preSum: preSum <= k)


class BITRangeAddPointGetArray:
    __slots__ = "_bit"

    def __init__(self, n: int, f: Optional[Callable[[int], int]] = None):
        if f is None:
            self._bit = BITArray(n)
        else:
            self._bit = BITArray(n, f)

    def addRange(self, start: int, end: int, delta: int) -> None:
        n = self._bit.n
        if start < 0:
            start = 0
        if end > n:
            end = n
        if start >= end:
            return
        self._bit.add(start, delta)
        self._bit.add(end, -delta)

    def get(self, index: int) -> int:
        return self._bit.queryPrefix(index + 1)

    def __repr__(self):
        return (
            "BITRangeAddPointGetArray: ["
            + ", ".join(str(self.get(i)) for i in range(self._bit.n))
            + "]"
        )


class BITRangeAddPointGetMap:
    __slots__ = "_bit"

    def __init__(self, n: int):
        self._bit = BITMap(n)

    def addRange(self, start: int, end: int, delta: int) -> None:
        n = self._bit.n
        if start < 0:
            start = 0
        if end > n:
            end = n
        if start >= end:
            return
        self._bit.add(start, delta)
        self._bit.add(end, -delta)

    def get(self, end: int) -> int:
        return self._bit.queryPrefix(end + 1)


class BITRangeAddRangeSumArray:
    __slots__ = ("n", "_bit0", "_bit1")

    def __init__(self, n: int, f: Optional[Callable[[int], int]] = None):
        if f is None:
            self.n = n
            self._bit0 = BITArray(n)
            self._bit1 = BITArray(n)
        else:
            self.n = n
            self._bit0 = BITArray(n, f)
            self._bit1 = BITArray(n)

    def add(self, index: int, delta: int) -> None:
        self._bit0.add(index, delta)

    def addRange(self, start: int, end: int, delta: int) -> None:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return
        self._bit0.add(start, -delta * start)
        self._bit0.add(end, delta * end)
        self._bit1.add(start, delta)
        self._bit1.add(end, -delta)

    def queryRange(self, start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return 0
        rightRes = self._bit1.queryPrefix(end) * end + self._bit0.queryPrefix(end)
        leftRes = self._bit1.queryPrefix(start) * start + self._bit0.queryPrefix(start)
        return rightRes - leftRes

    def __repr__(self):
        return (
            "BITRangeAddRangeSumArray: ["
            + ", ".join(str(self.queryRange(i, i + 1)) for i in range(self.n))
            + "]"
        )


class BITRangeAddRangeSumMap:
    __slots__ = ("n", "_bit0", "_bit1")

    def __init__(self, n: int):
        self.n = n
        self._bit0 = BITMap(n)
        self._bit1 = BITMap(n)

    def add(self, index: int, delta: int) -> None:
        self._bit0.add(index, delta)

    def addRange(self, start: int, end: int, delta: int) -> None:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return
        self._bit0.add(start, -delta * start)
        self._bit0.add(end, delta * end)
        self._bit1.add(start, delta)
        self._bit1.add(end, -delta)

    def queryRange(self, start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return 0
        rightRes = self._bit1.queryPrefix(end) * end + self._bit0.queryPrefix(end)
        leftRes = self._bit1.queryPrefix(start) * start + self._bit0.queryPrefix(start)
        return rightRes - leftRes

    def __repr__(self):
        return (
            "BITRangeAddRangeSumMap: ["
            + ", ".join(str(self.queryRange(i, i + 1)) for i in range(self.n))
            + "]"
        )


S = TypeVar("S")


class BITPrefixArray(Generic[S]):
    __slots__ = ("n", "_data", "_e", "_op")

    def __init__(
        self,
        n: int,
        e: Callable[[], S],
        op: Callable[[S, S], S],
        f: Optional[Callable[[int], S]] = None,
    ):
        if f is None:
            self.n = n
            self._data = [e() for _ in range(n)]
            self._e = e
            self._op = op
        else:
            self.n = n
            self._data = [f(i) for i in range(n)]
            self._e = e
            self._op = op
            for i in range(1, n + 1):
                j = i + (i & -i)
                if j <= n:
                    self._data[j - 1] = op(self._data[j - 1], self._data[i - 1])

    def update(self, index: int, value: S) -> None:
        index += 1
        while index <= self.n:
            self._data[index - 1] = self._op(self._data[index - 1], value)
            index += index & -index

    def queryPrefix(self, end: int) -> S:
        if end > self.n:
            end = self.n
        res = self._e()
        while end > 0:
            res = self._op(res, self._data[end - 1])
            end -= end & -end
        return res


class BITPrefixMap(Generic[S]):
    __slots__ = ("n", "_data", "_e", "_op")

    def __init__(
        self,
        n: int,
        e: Callable[[], S],
        op: Callable[[S, S], S],
    ):
        self.n = n
        self._data = dict()
        self._e = e
        self._op = op

    def update(self, index: int, value: S) -> None:
        index += 1
        while index <= self.n:
            self._data[index - 1] = self._op(self._data.get(index - 1, self._e()), value)
            index += index & -index

    def queryPrefix(self, end: int) -> S:
        if end > self.n:
            end = self.n
        res = self._e()
        while end > 0:
            res = self._op(res, self._data.get(end - 1, self._e()))
            end -= end & -end
        return res


if __name__ == "__main__":
    bitArray = BITArray(10)
    bitArray.add(0, 1)
    print(bitArray)
    bitArray.add(1, 2)
    print(bitArray)

    bitMap = BITMap(int(1e9))
    bitMap.add(0, 1)
    bitMap.add(int(1e7), 1)
    print(bitMap.queryPrefix(int(1e8)))

    bitRangeAddRangeSumArray = BITRangeAddRangeSumArray(10, lambda i: i + 1)
    bitRangeAddRangeSumArray.addRange(0, 10, 1)
    print(bitRangeAddRangeSumArray)

    bitRangeAddRangeSumMap = BITRangeAddRangeSumMap(int(1e9))
    bitRangeAddRangeSumMap.addRange(0, int(1e7), 1)
    print(bitRangeAddRangeSumMap.queryRange(0, int(1e8)))

    bitPrefixArray = BITPrefixArray(10, lambda: 0, max)
    bitPrefixArray.update(0, 1)
    bitPrefixArray.update(1, 0)
    print(bitPrefixArray.queryPrefix(2))

    bitPrefixMap = BITPrefixMap(int(1e9), lambda: 0, max)
    bitPrefixMap.update(0, 1)
    bitPrefixMap.update(int(1e7), 0)

    def testBITRangeAddPrefixQueryArrayAndBITRangeAddRangeSumArray() -> None:
        import random

        n = random.randint(1, 100)
        bit1 = BITRangeAddPointGetArray(n)
        bit2 = BITRangeAddRangeSumArray(n)
        for _ in range(100):
            start = random.randint(0, n)
            end = random.randint(start, n)
            delta = random.randint(-100, 100)
            bit1.addRange(start, end, delta)
            bit2.addRange(start, end, delta)
            pos = random.randint(0, n - 1)
            # assert bit1.queryPrefix(end) == bit2.queryRange(0, end)
            assert bit1.get(pos) == bit2.queryRange(pos, pos + 1), (pos, start, end, delta, n)
            # assert bit1.queryPrefix(end) - bit1.queryPrefix(start) == bit2.queryRange(start, end)

    testBITRangeAddPrefixQueryArrayAndBITRangeAddRangeSumArray()
