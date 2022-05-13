from typing import List


class HashHeap:
    def __init__(self, desc=False):
        """
        初始化一个数组，一个哈希表，默认为小顶堆
        """
        self.hash = {}
        self.heap = []
        self.desc = desc

    @property
    def size(self):
        return len(self.heap)

    def top(self):
        return self.heap[0]

    def push(self, x):
        """
        入堆
        """
        self.heap.append(x)
        self.hash[x] = self.size - 1
        self._sift_up(self.size - 1)

    def pop(self):
        """
        弹出堆顶
        """
        res = self.heap[0]
        self.remove(res)
        return res

    def remove(self, item):
        """
        移除一个元素
        """
        index = self.hash[item]
        self._swap(self.size - 1, index)
        self.heap.pop()
        del self.hash[item]
        if index < self.size:
            self._sift_up(index)
            self._sift_down(index)

    def _smaller(self, lhs, rhs):
        return lhs > rhs if self.desc else lhs < rhs

    def _sift_up(self, index):
        """
        将元素向上调整
        """
        while index > 0:
            parent = (index - 1) // 2
            if self._smaller(self.heap[parent], self.heap[index]):
                return
            self._swap(parent, index)
            index = parent

    def _sift_down(self, index):
        """
        将元素向下调整
        """
        while index * 2 + 1 < self.size:
            smallest = index
            left = index * 2 + 1
            right = index * 2 + 2
            if self._smaller(self.heap[left], self.heap[smallest]):
                smallest = left
            if right < self.size and self._smaller(self.heap[right], self.heap[smallest]):
                smallest = right
            if smallest == index:
                return
            self._swap(index, smallest)
            index = smallest

    def _swap(self, i, j):
        """
        交换两个元素
        """
        if i == j:
            return

        item1, item2 = self.heap[i], self.heap[j]
        self.hash[item1], self.hash[item2] = j, i
        self.heap[i], self.heap[j] = item2, item1


class Solution:
    def medianSlidingWindow(self, nums: List[int], k: int) -> List[float]:
        max_heap, min_heap = HashHeap(desc=True), HashHeap()
        # 初始化
        for i in range(k):
            min_heap.push((nums[i], i))
        for i in range(k // 2):
            max_heap.push(min_heap.pop())

        res = [
            (min_heap.top()[0] + max_heap.top()[0]) / 2 if k % 2 == 0 else min_heap.top()[0] * 1.0
        ]

        for i in range(k, len(nums)):
            # 添加新的元素
            if nums[i] < min_heap.top()[0]:
                max_heap.push((nums[i], i))
            else:
                min_heap.push((nums[i], i))

            # 弹出旧的元素
            if (nums[i - k], i - k) in max_heap.hash:
                max_heap.remove((nums[i - k], i - k))
            else:
                min_heap.remove((nums[i - k], i - k))

            # 调整两个堆的大小关系
            if min_heap.size > max_heap.size + 1:
                max_heap.push(min_heap.pop())
            if min_heap.size < max_heap.size:
                min_heap.push(max_heap.pop())

            # 添加答案
            res.append(
                (min_heap.top()[0] + max_heap.top()[0]) / 2
                if k % 2 == 0
                else min_heap.top()[0] * 1.0
            )

        return res

