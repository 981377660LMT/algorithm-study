"""二维滑窗"""

from collections import deque
from typing import List


class MonoQueue:
    __slots__ = ("maxQueue", "minQueue", "rawQueue")

    def __init__(self):
        self.maxQueue = deque()
        self.minQueue = deque()
        self.rawQueue = deque()

    @property
    def min(self) -> int:
        return self.minQueue[0][0]

    @property
    def max(self) -> int:
        return self.maxQueue[0][0]

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError("popleft from empty queue")

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


def windowMax2D(
    matrix: List[List[int]], rowSize: int, colSize: int, *, isMax: bool
) -> List[List[int]]:
    """
    每个rowSize*colSize窗口(和)的最值
    1 ≤ n * m ≤ 1e5
    """
    ROW, COL = len(matrix), len(matrix[0])
    cols = [MonoQueue() for _ in range(COL)]
    for r in range(rowSize):
        for c in range(COL):
            cols[c].append(matrix[r][c])

    res = []
    for r in range(ROW - rowSize + 1):
        res.append([])
        # 维护这一行的滑动窗口最值
        window = MonoQueue()
        for c in range(COL):
            # 每列的最值
            window.append(cols[c].max if isMax else cols[c].min)
            if c >= colSize - 1:
                res[-1].append(window.max if isMax else window.min)
                window.popleft()

        # 下一行进入窗口
        if r + rowSize < ROW:
            for c in range(COL):
                cols[c].append(matrix[r + rowSize][c])
                cols[c].popleft()

    return res


def useQueryMax2D(
    matrix: List[List[int]], bigRow: int, bigCol: int, smallRow: int, smallCol: int, isMax: bool
) -> List[List[int]]:
    """在每个bigRow*bigCol窗口范围内所有的smallRow*smallCol窗口和的最值"""
    ROW, COL = len(matrix), len(matrix[0])
    preSum = [[0] * (COL + 1) for _ in range(ROW + 1)]
    for r in range(ROW):
        for c in range(COL):
            preSum[r + 1][c + 1] = preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c] + matrix[r][c]

    windowSum = [[] for _ in range(ROW - smallRow + 1)]
    for r in range(ROW - smallRow + 1):
        for c in range(COL - smallCol + 1):
            windowSum[r].append(
                preSum[r + smallRow][c + smallCol]
                - preSum[r][c + smallCol]
                - preSum[r + smallRow][c]
                + preSum[r][c]
            )

    windowMax = windowMax2D(windowSum, bigRow - smallRow + 1, bigCol - smallCol + 1, isMax=isMax)
    return windowMax


assert windowMax2D([[1, 2, 3], [5, 3, 2], [1, 0, 9]], 2, 1, isMax=True) == [[5, 3, 3], [5, 3, 9]]
assert windowMax2D(matrix=[[1, 2, 3], [5, 3, 2], [1, 0, 9]], rowSize=2, colSize=1, isMax=False) == [
    [1, 2, 2],
    [1, 0, 2],
]

assert useQueryMax2D([[1, 2, 3], [5, 3, 2], [1, 0, 9]], 2, 1, 2, 1, True) == [[6, 5, 5], [6, 3, 11]]
######################################################################################################


class Solution2:
    def solve(self, matrix: List[List[int]], rowSize: int, colSize: int) -> List[List[int]]:
        """有一个整数组成的矩阵，现请你从中找出一个 rowSize*colSize 的区域，使得该区域所有数中的最大值和最小值的差最小"""
        row, col = len(matrix), len(matrix[0])
        cols = [MonoQueue() for _ in range(col)]
        for r in range(rowSize):
            for c, v in enumerate(matrix[r]):
                cols[c].append(v)

        res = []
        for r in range(row - rowSize + 1):
            # 维护每行的滑动窗口最小值
            window1, window2 = MonoQueue(), MonoQueue()
            for c in range(col):
                # 每列的最小/大值
                window1.append((cols[c].min))
                window2.append((cols[c].max))
                if c >= colSize - 1:
                    res.append((window1.min, window2.max))
                    window1.popleft()
                    window2.popleft()

            if r + rowSize < row:
                for c in range(col):
                    cols[c].append(matrix[r + rowSize][c])
                    cols[c].popleft()

        return min((y - x) for x, y in res)


# row, col, k = map(int, input().split())
# matrix = []
# for i in range(row):
#     matrix.append(list(map(int, input().split())))
# print(Solution2().solve(matrix, k))
