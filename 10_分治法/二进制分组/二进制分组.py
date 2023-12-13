from collections import deque
from typing import Deque, List

INF = int(1e18)


class BinaryGroupHeap:
    def __init__(self):
        self.heap: List[Deque[int]] = []

    def insert(self, val: int):
        new_group = Deque[int]([val])
        while self.heap and len(self.heap[-1]) == len(new_group):
            self.heap[-1] += new_group
            self.heap[-1] = deque(sorted(self.heap[-1]))
            new_group = self.heap.pop()
        self.heap.append(new_group)

    def get_min(self) -> int:
        min_val = INF
        for group in self.heap:
            min_val = min(min_val, group[0])
        return min_val

    def get_max(self) -> int:
        max_val = -INF
        for group in self.heap:
            max_val = max(max_val, group[-1])
        return max_val

    def delete_min(self):
        min_index = 0
        for i in range(len(self.heap)):
            if self.heap[i][0] < self.heap[min_index][0]:
                min_index = i
        self.heap[min_index].popleft()
        if not self.heap[min_index]:
            self.heap.pop(min_index)

    def delete_max(self):
        max_index = 0
        for i in range(len(self.heap)):
            if self.heap[i][-1] > self.heap[max_index][-1]:
                max_index = i
        self.heap[max_index].pop()
        if not self.heap[max_index]:
            self.heap.pop(max_index)


if __name__ == "__main__":
    import sys

    input = sys.stdin.buffer.readline
    n, q = map(int, input().split())
    s = list(map(int, input().split()))
    heap = BinaryGroupHeap()
    for val in s:
        heap.insert(val)
    for _ in range(q):
        query = tuple(map(int, input().split()))
        if query[0] == 0:
            heap.insert(query[1])
        elif query[0] == 1:
            print(heap.get_min())
            heap.delete_min()
        else:
            print(heap.get_max())
            heap.delete_max()
