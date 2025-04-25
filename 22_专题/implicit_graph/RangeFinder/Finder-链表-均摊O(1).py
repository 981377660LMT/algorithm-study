from typing import Generator


class Finder:

    __slots__ = ("_n", "_exist", "_prev", "_next")

    def __init__(self, n: int) -> None:
        """建立一个包含0到n-1的集合."""
        self._n = n
        self._exist = [True] * n
        self._prev = [i - 1 for i in range(n)]
        self._next = [i + 1 for i in range(n)]

    def has(self, i: int) -> bool:
        """判断i是否在集合中."""
        return 0 <= i < self._n and self._exist[i]

    def erase(self, i: int) -> bool:
        """删除i."""
        if not self.has(i):
            return False
        self._exist[i] = False
        return True

    def prev(self, i: int) -> int:
        """返回小于等于i的最大元素.如果不存在,返回-1."""
        if i < 0:
            return -1
        if i >= self._n:
            i = self._n - 1
        if self._exist[i]:
            return i

        prev, exist = self._prev, self._exist
        realPrev = prev[i]
        while realPrev >= 0 and not exist[realPrev]:
            realPrev = prev[realPrev]
        cur = i
        while cur >= 0 and cur != realPrev:
            tmp = prev[cur]
            prev[cur] = realPrev
            cur = tmp
        return realPrev

    def next(self, i: int) -> int:
        """返回大于等于i的最小元素.如果不存在,返回n."""
        if i < 0:
            i = 0
        if i >= self._n:
            return self._n
        if self._exist[i]:
            return i

        next, exist = self._next, self._exist
        realNext = next[i]
        while realNext < self._n and not exist[realNext]:
            realNext = next[realNext]
        cur = i
        while cur < self._n and cur != realNext:
            tmp = next[cur]
            next[cur] = realNext
            cur = tmp
        return realNext

    def enumerate(self, start: int, end: int) -> Generator[int, int, None]:
        """遍历[start,end)区间内的元素."""
        if start < 0:
            start = 0
        if end > self._n:
            end = self._n
        if start >= end:
            return

        x = self.next(start)
        while x < end:
            yield x
            x = self.next(x + 1)

    def __contains__(self, i: int) -> bool:
        return self.has(i)

    def __len__(self) -> int:
        return self._n

    def __repr__(self) -> str:
        return f"Finder{{{', '.join(map(str, self.enumerate(0, self._n)))}}}"


if __name__ == "__main__":
    f = Finder(10)
    print(f)
    f.erase(2)
    print(f)
    f.erase(3)
    print(f)
    f.erase(4)
    print(f)
    f.erase(5)
    print(f)
    f.erase(6)
    print(f)
    f.erase(7)
    print(f)

    print(f.next(2), f.prev(7))
