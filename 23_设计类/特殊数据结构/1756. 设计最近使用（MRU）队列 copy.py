from sortedcontainers import SortedList


class MRUQueue:
    def __init__(self, n: int):
        self.maxOrderIndex = n
        self.arr = SortedList([[v, v] for v in range(1, n + 1)])

    def fetch(self, k: int) -> int:
        order, val = self.arr.pop(k - 1)
        self.maxOrderIndex += 1
        self.arr.add([self.maxOrderIndex, val])
        return val


# 思路
# 1. 删除k-1下标的元素
# 2. 加入list后让被删的元素保持在末尾(需要一个计数维护)
