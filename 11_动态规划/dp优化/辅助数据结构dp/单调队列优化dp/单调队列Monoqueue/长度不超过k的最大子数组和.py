# 输入一个长度为 n 的整数序列，从中找出一段长度不超过 m 的连续子序列，使得子序列中所有数的和最大。
# 子序列的长度至少是 1


from collections import deque
from typing import Iterable, Optional


class MonoQueue:
    def __init__(self, iterable: Optional[Iterable[int]] = None):
        self.minQueue = deque()
        self.maxQueue = deque()
        self.rawQueue = deque()
        self.index = 0

        if iterable:
            for value in iterable:
                self.append(value)

    @property
    def min(self) -> int:
        if not self.minQueue:
            raise ValueError('monoQueue is empty')
        return self.minQueue[0][0]

    @property
    def max(self) -> int:
        if not self.maxQueue:
            raise ValueError('monoQueue is empty')
        return self.maxQueue[0][0]

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError('popleft from empty queue')

        self.minQueue[0][1] -= 1
        if self.minQueue[0][1] == 0:
            self.minQueue.popleft()

        self.maxQueue[0][1] -= 1
        if self.maxQueue[0][1] == 0:
            self.maxQueue.popleft()

        return self.rawQueue.popleft()[0]

    def append(self, value: int) -> None:
        count = 1
        while self.minQueue and self.minQueue[-1][0] > value:
            count += self.minQueue.pop()[1]
        self.minQueue.append([value, count])

        count = 1
        while self.maxQueue and self.maxQueue[-1][0] < value:
            count += self.maxQueue.pop()[1]
        self.maxQueue.append([value, count])

        self.rawQueue.append((value, self.index))
        self.index += 1

    def __len__(self) -> int:
        return len(self.rawQueue)

    def __getitem__(self, index: int) -> int:
        return self.rawQueue[index][0]


n, k = map(int, input().split())
nums = list(map(int, input().split()))

res = nums[0]
queue = MonoQueue([0])  # 注意存前缀和的哨兵
curSum = 0


for num in nums:
    curSum += num
    while len(queue) > k:
        queue.popleft()
    # 要让子序列长度至少为1，需要先判res，再加入队列
    if queue:
        res = max(res, curSum - queue.min)
    queue.append(curSum)

print(res)

