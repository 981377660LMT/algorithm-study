# https://atcoder.jp/contests/abc413/tasks/abc413_c
#
# append(v, c) -> 在队列尾部添加 c 个元素 v
# appendleft(v, c) -> 在队列头部添加 c 个元素 v
# pop(c) -> 从队列尾部删除 c 个元素，返回删除的元素之和
# popleft(c) -> 从队列头部删除 c 个元素，返回删除的元素之和
# __len__() -> 返回队列的长度
# __sum__() -> 返回队列中所有元素的和


from collections import deque


class LargeDeque:
    __slots__ = ("_data", "_size", "_sum")

    def __init__(self):
        self._data = deque()
        self._size = 0
        self._sum = 0

    def append(self, v: int, c: int) -> None:
        """在队列尾部添加 c 个元素 v."""
        if c <= 0:
            return
        if self._data and self._data[-1][0] == v:
            self._data[-1][1] += c
        else:
            self._data.append([v, c])
        self._size += c
        self._sum += v * c

    def appendleft(self, v: int, c: int) -> None:
        """在队列头部添加 c 个元素 v."""
        if c <= 0:
            return
        if self._data and self._data[0][0] == v:
            self._data[0][1] += c
        else:
            self._data.appendleft([v, c])
        self._size += c
        self._sum += v * c

    def pop(self, c: int) -> int:
        """从队列尾部删除 c 个元素，返回删除的元素之和."""
        if c <= 0 or self._size == 0:
            return 0
        res = 0
        remain = c
        while remain > 0 and self._data:
            v, count = self._data[-1]
            if count <= remain:
                res += v * count
                remain -= count
                self._data.pop()
            else:
                res += v * remain
                self._data[-1][1] -= remain
                remain = 0
        self._size -= c
        self._sum -= res
        return res

    def popleft(self, c: int) -> int:
        """从队列头部删除 c 个元素，返回删除的元素之和."""
        if c <= 0 or self._size == 0:
            return 0
        res = 0
        remain = c
        while remain > 0 and self._data:
            v, count = self._data[0]
            if count <= remain:
                res += v * count
                remain -= count
                self._data.popleft()
            else:
                res += v * remain
                self._data[0][1] -= remain
                remain = 0
        self._size -= c
        self._sum -= res
        return res

    def __len__(self) -> int:
        return self._size

    def __sum__(self) -> int:
        return self._sum

    def __repr__(self) -> str:
        return f"LargeDeque(size={self._size}, sum={self._sum}, data={list(self._data)})"


if __name__ == "__main__":

    def abc413_c():
        q = int(input())
        ld = LargeDeque()
        for _ in range(q):
            t, *rest = map(int, input().split())
            if t == 1:
                c, x = rest
                ld.append(x, c)
            elif t == 2:
                k = rest[0]
                print(ld.popleft(k))
