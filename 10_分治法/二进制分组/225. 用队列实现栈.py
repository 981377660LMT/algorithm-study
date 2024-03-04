# 二进制分组实现栈，每个分组是一个时间戳和值的元组，按照时间顺序递减
# TODO: 可以不维护时间戳，自然有序.
# https://leetcode.cn/problems/implement-stack-using-queues/solutions/2666412/python-er-jin-zhi-fen-zu-jun-tan-ologn-b-0i17/
#
# https://zhuanlan.zhihu.com/p/35519230
# 二进制分组实现堆同理，每个分组是一个递减的list.
# 如果用双端队列来实现内层的vector，可以实现双端堆.
# 缺点是常数较大.

import itertools

from typing import Generic, List, Tuple, TypeVar
from collections import deque


class MyStack:
    __slots__ = "_groups", "_timer", "_size"

    def __init__(self):
        self._groups: List["Queue[Tuple[int,int]]"] = []
        self._timer = itertools.count()
        self._size = 0

    def push(self, x: int) -> None:
        """
        二进制加法进位，每次进位后都会将之前的所有元素都加入到新的分组中并清空之前的分组.
        """
        k = 0
        while k < len(self._groups) and self._groups[k]:
            k += 1
        if k == len(self._groups):
            self._groups.append(self._createGroup())
        self._addToGroup(k, x)
        for i in range(k):
            self._mergeGroup(i, k)
        self._size += 1

    def pop(self) -> int:
        topGroupIndex = self._findTopGroupIndex()
        res = self._groups[topGroupIndex].popleft()[1]
        self._size -= 1
        return res

    def top(self) -> int:
        topGroupIndex = self._findTopGroupIndex()
        return self._groups[topGroupIndex].first()[1]

    def empty(self) -> bool:
        return self._size == 0

    def _createGroup(self) -> "Queue[Tuple[int,int]]":
        return Queue()

    def _addToGroup(self, groupIndex: int, value: int) -> None:
        self._groups[groupIndex].append((next(self._timer), value))

    def _mergeGroup(self, fromIndex: int, toIndex: int) -> None:
        newGroup = self._createGroup()
        g1, g2 = self._groups[fromIndex], self._groups[toIndex]
        n1, n2 = len(g1), len(g2)
        i1, i2 = 0, 0
        while i1 < n1 and i2 < n2:
            t1, t2 = g1.first()[0], g2.first()[0]
            if t1 >= t2:
                newGroup.append(g1.popleft())
                i1 += 1
            else:
                newGroup.append(g2.popleft())
                i2 += 1
        while i1 < n1:
            newGroup.append(g1.popleft())
            i1 += 1
        while i2 < n2:
            newGroup.append(g2.popleft())
            i2 += 1
        self._groups[toIndex] = newGroup
        self._groups[fromIndex].clear()

    def _findTopGroupIndex(self) -> int:
        maxTime, topGroupIndex = -1, -1
        for i, g in enumerate(self._groups):
            if g and (tmp := g.first()[0]) > maxTime:
                maxTime = tmp
                topGroupIndex = i
        return topGroupIndex


T = TypeVar("T")


class Queue(Generic[T]):
    __slots__ = "_deque"

    def __init__(self):
        self._deque = deque()

    def append(self, x: T) -> None:
        self._deque.append(x)

    def popleft(self) -> T:
        return self._deque.popleft()

    def first(self) -> T:
        return self._deque[0]

    def clear(self) -> None:
        self._deque.clear()

    def __len__(self) -> int:
        return len(self._deque)

    def __repr__(self) -> str:
        return str(self._deque)
