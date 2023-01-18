from collections import deque
from sortedcontainers import SortedList

# MK平均值 按照如下步骤计算：

# 1.如果数据流中的整数少于 m 个，MK 平均值 为 -1 ，
# 2.否则将数据流中`最后 m 个元素`拷贝到一个独立的容器中。
# !从这个容器中删除最小的 k 个数和最大的 k 个数。
# 3.计算剩余元素的平均值，并 向下取整到最近的整数 。
# !m<=1e5 k*2<m
# !滑动窗口topK之和
# https://leetcode.cn/problems/finding-mk-average/solution/by-981377660lmt-5hhm/


class MKAverage:
    def __init__(self, m: int, k: int):
        self._m = m
        self._k = k
        self._sum = 0
        self._queue = deque()
        self._sl = SortedList()

    def calculateMKAverage(self) -> int:
        if len(self._queue) < self._m:
            return -1
        return self._sum // (self._m - 2 * self._k)

    def addElement(self, num: int) -> None:
        self._queue.append(num)
        if len(self._queue) == self._m:
            self._initMK()
        elif len(self._queue) > self._m:
            self._addSl(num)
            self._removeSl(self._queue.popleft())

    def _addSl(self, num: int) -> None:
        # 加入后对区间和的影响，num会把sortedList里的元素挤到左边或者右边
        pos = self._sl.bisect_left(num)
        if pos < self._k:
            # 被挤到中间来了
            self._sum += self._sl[self._k - 1]
        elif self._k <= pos <= self._m - self._k:
            self._sum += num
        else:
            # 被挤到中间来了
            self._sum += self._sl[self._m - self._k]
        self._sl.add(num)

    def _removeSl(self, num: int) -> None:
        pos = self._sl.bisect_left(num)
        if pos < self._k:
            # 左移
            self._sum -= self._sl[self._k]
        elif self._k <= pos <= self._m - self._k:
            self._sum -= num
        else:
            # 右移
            self._sum -= self._sl[self._m - self._k]
        self._sl.remove(num)

    def _initMK(self) -> None:
        self._sl.update(self._queue)
        self._sum = sum(self._sl[self._k : -self._k])


class TopkSum:
    __slots__ = ("_sl", "_k", "_topKSum")

    def __init__(self, k: int, isMin: bool) -> None:
        self._sl = SortedList() if isMin else SortedList(key=lambda x: -x)
        self._k = k
        self._topKSum = 0

    def add(self, x: int) -> None:
        pos = self._sl.bisect_left(x)
        if pos < self._k:
            self._topKSum += x
            if self._k - 1 < len(self._sl):
                self._topKSum -= self._sl[self._k - 1]  # type: ignore
        self._sl.add(x)

    def remove(self, x: int) -> None:
        pos = self._sl.bisect_left(x)
        if pos < self._k:
            self._topKSum -= x
            if self._k < len(self._sl):
                self._topKSum += self._sl[self._k]  # type: ignore
        self._sl.remove(x)

    def query(self) -> int:
        return self._topKSum


if __name__ == "__main__":
    MK = MKAverage(3, 1)
    MK.addElement(17612)
    MK.addElement(74607)
    print(MK.calculateMKAverage())
    MK.addElement(8272)
    MK.addElement(33433)
    print(MK.calculateMKAverage())
    MK.addElement(15456)
    MK.addElement(64938)
    print(MK.calculateMKAverage())
