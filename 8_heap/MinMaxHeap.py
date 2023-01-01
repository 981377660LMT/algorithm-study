import sys
from heapq import heapify, heappop, heappush
from typing import List, Optional

# https://judge.yosupo.jp/submission/109819
# 懒删除的堆 维护最大最小值


class MinMaxHeap:
    def __init__(self, a: Optional[List[int]] = None):
        if a is None:
            a = []
        self._max_heap = [-x for x in a]
        self._min_heap = a[:]
        heapify(self._max_heap)
        heapify(self._min_heap)
        self._max_deleted = []
        self._min_deleted = []
        self._size = len(a)

    def pop_max(self) -> int:
        while True:
            v = -heappop(self._max_heap)
            if self._min_deleted and self._min_deleted[0] == -v:
                heappop(self._min_deleted)
            else:
                self._size -= 1
                if self._size:
                    heappush(self._max_deleted, v)
                else:
                    self.clear()
                return v

    def pop_min(self) -> int:
        while True:
            v = heappop(self._min_heap)
            if self._max_deleted and self._max_deleted[0] == v:
                heappop(self._max_deleted)
            else:
                self._size -= 1
                if self._size:
                    heappush(self._min_deleted, -v)
                else:
                    self.clear()
                return v

    def get_max(self) -> int:
        while True:
            v = -self._max_heap[0]
            if self._min_deleted and self._min_deleted[0] == -v:
                heappop(self._min_deleted)
                heappop(self._max_heap)
            else:
                return v

    def get_min(self) -> int:
        while True:
            v = self._min_heap[0]
            if self._max_deleted and self._max_deleted[0] == v:
                heappop(self._max_deleted)
                heappop(self._min_heap)
            else:
                return v

    def push(self, v: int) -> None:
        self._size += 1
        heappush(self._max_heap, -v)
        heappush(self._min_heap, v)

    def clear(self) -> None:
        self._max_heap = []
        self._min_heap = []
        self._max_deleted = []
        self._min_deleted = []
        self._size = 0

    def __len__(self) -> int:
        return self._size


input = sys.stdin.buffer.readline
n, q = map(int, input().split())
s = list(map(int, input().split()))
heap = MinMaxHeap(s)
for _ in range(q):
    query = tuple(map(int, input().split()))
    if query[0] == 0:
        heap.push(query[1])
    elif query[0] == 1:
        print(heap.pop_min())
    else:
        print(heap.pop_max())
