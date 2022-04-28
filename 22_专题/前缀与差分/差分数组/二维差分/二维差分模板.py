from typing import List


class DiffMatrix:
    """二维差分模板(矩阵可变)"""

    def __init__(self, A: List[List[int]]):
        self.m, self.n = len(A), len(A[0])
        self.matrix = [[0] * self.n for _ in range(self.m)]
        for i in range(self.m):
            for j in range(self.n):
                self.matrix[i][j] = A[i][j]

        # 需要额外大小为(m+2)∗(n+2)的差分数组，第一行第一列不用(始终为0)
        self.diff = [[0] * (self.n + 2) for _ in range(self.m + 2)]

    def add(self, r1: int, c1: int, r2: int, c2: int, k: int) -> None:
        """区间更新A[r1:r2+1, c1:c2+1]"""
        self.diff[r1 + 1][c1 + 1] += k
        self.diff[r1 + 1][c2 + 2] -= k
        self.diff[r2 + 2][c1 + 1] -= k
        self.diff[r2 + 2][c2 + 2] += k

    def update(self) -> None:
        """遍历矩阵，还原对应元素的增量"""
        for i in range(self.m):
            for j in range(self.n):
                # 差分数组的前缀和即为nums[i]
                self.diff[i + 1][j + 1] += (
                    self.diff[i + 1][j] + self.diff[i][j + 1] - self.diff[i][j]
                )
                self.matrix[i][j] += self.diff[i + 1][j + 1]


class PreSumMatrix:
    """二维前缀和模板(矩阵不可变)"""

    def __init__(self, A: List[List[int]]):
        m, n = len(A), len(A[0])

        # 前缀和数组
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = A[r][c] + preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c]
        self.preSum = preSum

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        preSumMatrix.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self.preSum[r2 + 1][c2 + 1]
            - self.preSum[r2 + 1][c1]
            - self.preSum[r1][c2 + 1]
            + self.preSum[r1][c1]
        )


if __name__ == '__main__':
    diffmatrix = DiffMatrix([[1, 2], [3, 4]])
    print(diffmatrix.matrix)
    diffmatrix.add(0, 0, 1, 1, 1)
    # diffmatrix.add(0, 0, 0, 0, 1)
    # diffmatrix.add(0, 0, 0, 0, 1)
    diffmatrix.update()
    print(diffmatrix.matrix)
    print(diffmatrix.diff)

