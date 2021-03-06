from collections import deque
from typing import List


class MonoQueue:
    def __init__(self):
        self.maxQueue = deque()
        self.minQueue = deque()
        self.rawQueue = deque()

    @property
    def min(self) -> int:
        return self.minQueue[0][0] if self.minQueue else None

    @property
    def max(self) -> int:
        return self.maxQueue[0][0] if self.maxQueue else None

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError('popleft from empty queue')

        self.minQueue[0][1] -= 1
        if self.minQueue[0][1] == 0:
            self.minQueue.popleft()

        self.maxQueue[0][1] -= 1
        if self.maxQueue[0][1] == 0:
            self.maxQueue.popleft()

        return self.rawQueue.popleft()

    def append(self, value: int) -> None:
        count = 1
        while self.minQueue and self.minQueue[-1][0] > value:
            count += self.minQueue.pop()[1]
        self.minQueue.append([value, count])

        count = 1
        while self.maxQueue and self.maxQueue[-1][0] < value:
            count += self.maxQueue.pop()[1]
        self.maxQueue.append([value, count])

        self.rawQueue.append(value)

    def __len__(self) -> int:
        return len(self.rawQueue)


# 1 ≤ n * m ≤ 100,000
class Solution:
    def solve(self, matrix: List[List[int]], k: int) -> List[List[int]]:
        """return a matrix containing minimum values of all k by k submatrices."""
        row, col = len(matrix), len(matrix[0])
        cols = [MonoQueue() for _ in range(col)]
        for r in range(k):
            for c, v in enumerate(matrix[r]):
                cols[c].append(v)

        res = []
        for r in range(row - k + 1):
            res.append([])
            # 维护每行的滑动窗口最小值
            window = MonoQueue()
            for c in range(col):
                # 每列的最小值
                window.append(cols[c].min)
                if c >= k - 1:
                    res[-1].append(window.min)
                    window.popleft()

            if r + k < row:
                for c in range(col):
                    cols[c].append(matrix[r + k][c])
                    cols[c].popleft()

        return res


print(Solution().solve(matrix=[[1, 2, 3], [5, 3, 2], [1, 0, 9]], k=2))
# [
#     [1, 2],
#     [0, 0]
# ]
class Solution2:
    def solve(self, matrix: List[List[int]], k: int) -> List[List[int]]:
        """有一个 a×b 的整数组成的矩阵，现请你从中找出一个 n×n 的正方形区域，使得该区域所有数中的最大值和最小值的差最小"""
        row, col = len(matrix), len(matrix[0])
        cols = [MonoQueue() for _ in range(col)]
        for r in range(k):
            for c, v in enumerate(matrix[r]):
                cols[c].append(v)

        res = []
        for r in range(row - k + 1):
            # 维护每行的滑动窗口最小值
            window1, window2 = MonoQueue(), MonoQueue()
            for c in range(col):
                # 每列的最小/大值
                window1.append((cols[c].min))
                window2.append((cols[c].max))
                if c >= k - 1:
                    res.append((window1.min, window2.max))
                    window1.popleft()
                    window2.popleft()

            if r + k < row:
                for c in range(col):
                    cols[c].append(matrix[r + k][c])
                    cols[c].popleft()

        return min((y - x) for x, y in res)


row, col, k = map(int, input().split())
matrix = []
for i in range(row):
    matrix.append(list(map(int, input().split())))
print(Solution2().solve(matrix, k))
