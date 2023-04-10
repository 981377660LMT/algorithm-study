from typing import Optional


class FastSet:
    """利用位运算寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
    初始时,所有位置都未被访问过.
    """

    __slots__ = "_n", "_lg", "_seg"

    @staticmethod
    def _trailingZeros1024(x: int) -> int:
        if x == 0:
            return 1024
        return (x & -x).bit_length() - 1

    def __init__(self, n: int) -> None:
        self._n = n
        seg = []
        while True:
            seg.append([0] * ((n + 1023) >> 10))
            n = (n + 1023) >> 10
            if n <= 1:
                break
        self._seg = seg
        self._lg = len(seg)

    def insert(self, i: int) -> None:
        for h in range(self._lg):
            self._seg[h][i >> 10] |= 1 << (i & 1023)
            i >>= 10

    def erase(self, i: int) -> None:
        for h in range(self._lg):
            self._seg[h][i >> 10] &= ~(1 << (i & 1023))
            if self._seg[h][i >> 10]:
                break
            i >>= 10

    def next(self, i: int) -> Optional[int]:
        """返回x右侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        if i < 0:
            i = 0
        if i >= self._n:
            return
        seg = self._seg
        for h in range(self._lg):
            if i >> 10 == len(seg[h]):
                break
            d = seg[h][i >> 10] >> (i & 1023)
            if d == 0:
                i = (i >> 10) + 1
                continue
            i += self._trailingZeros1024(d)
            for g in range(h - 1, -1, -1):
                i <<= 10
                i += self._trailingZeros1024(seg[g][i >> 10])
            return i

    def prev(self, i: int) -> Optional[int]:
        """返回x左侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        if i < 0:
            return
        if i >= self._n:
            i = self._n - 1
        seg = self._seg
        for h in range(self._lg):
            if i == -1:
                break
            d = seg[h][i >> 10] << (1023 - (i & 1023)) & ((1 << 1024) - 1)
            if d == 0:
                i = (i >> 10) - 1
                continue
            i += d.bit_length() - 1024
            for g in range(h - 1, -1, -1):
                i <<= 10
                i += (seg[g][i >> 10]).bit_length() - 1
            return i

    def islice(self, begin: int, end: int):
        """遍历[start,end)区间内的元素."""
        x = begin - 1
        while True:
            x = self.next(x + 1)
            if x is None or x >= end:
                break
            yield x

    def __contains__(self, i: int) -> bool:
        return self._seg[0][i >> 10] & (1 << (i & 1023)) != 0

    def __iter__(self):
        yield from self.islice(0, self._n)

    def __repr__(self):
        return f"FastSet({list(self)})"


if __name__ == "__main__":
    ...

    # 前驱后继
    def pre(pos: int):
        return next((i for i in range(pos, -1, -1) if ok[i]), None)

    def nxt(pos: int):
        return next((i for i in range(pos, n) if ok[i]), None)

    def erase(left: int, right: int):
        for i in range(left, right):
            ok[i] = False

    from random import randint, seed

    seed(0)
    for _ in range(100):
        n = randint(1, 10)
        F = FastSet(n)
        for i in range(n):
            F.insert(i)
        ok = [True] * n
        for _ in range(100):
            e = randint(0, n - 1)
            F.erase(e)
            erase(e, e + 1)
            for i in range(n):
                assert F.prev(i) == pre(i), (i, F.prev(i), pre(i))
                assert F.next(i) == nxt(i), (i, F.next(i), nxt(i))
    print("Done!")

    n = int(1e5)
    fs = FastSet(n)
    import time

    time1 = time.time()
    for i in range(n):
        fs.insert(i)
        fs.next(i)
        fs.prev(i)
        i in fs
        fs.erase(i)
        fs.insert(i)
    print(time.time() - time1)
