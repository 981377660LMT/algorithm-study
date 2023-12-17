from heapq import heappush, heappushpop


class MedianFinderHeap:
    """对顶堆动态维护中位数."""

    def __init__(self):
        self.small = []  # 左边，大顶堆
        self.large = []  # 右边，小顶堆

    def add(self, num: int) -> None:
        if len(self.small) == len(self.large):
            heappush(self.large, -heappushpop(self.small, -num))
        elif len(self.small) < len(self.large):
            heappush(self.small, -heappushpop(self.large, num))

    def query(self) -> float:
        if len(self.small) == len(self.large):
            return (self.large[0] - self.small[0]) / 2
        elif len(self.small) < len(self.large):
            return self.large[0]

        raise Exception("Invalid state")

    def __len__(self):
        return len(self.small) + len(self.large)
