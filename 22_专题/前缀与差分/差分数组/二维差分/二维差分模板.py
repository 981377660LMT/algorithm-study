from typing import List


class DiffMatrix:
    """二维差分模板(矩阵可变)"""

    __slots__ = ("matrix", "_diff", "_row", "_col")

    def __init__(self, mat: List[List[int]]):
        self._row, self._col = len(mat), len(mat[0])
        self.matrix = [[0] * self._col for _ in range(self._row)]
        for i in range(self._row):
            tmp1, tmp2 = self.matrix[i], mat[i]
            for j in range(self._col):
                tmp1[j] = tmp2[j]
        self._diff = [[0] * (self._col + 2) for _ in range(self._row + 2)]

    def add(self, r1: int, c1: int, r2: int, c2: int, delta: int) -> None:
        """
        区间更新左上角`(row1, col1)` 到 右下角`(row2, col2)`闭区间所描述的子矩阵的元素.
        0<=r1<=r2<row, 0<=c1<=c2<col.
        """
        self._diff[r1 + 1][c1 + 1] += delta
        self._diff[r1 + 1][c2 + 2] -= delta
        self._diff[r2 + 2][c1 + 1] -= delta
        self._diff[r2 + 2][c2 + 2] += delta

    def update(self) -> None:
        """遍历矩阵，还原对应元素的增量"""
        for i in range(self._row):
            tmpDiff0, tmpDiff1 = self._diff[i], self._diff[i + 1]
            tmpMatrix = self.matrix[i]
            for j in range(self._col):
                # 差分数组的前缀和即为nums[i]
                tmpDiff1[j + 1] += tmpDiff1[j] + tmpDiff0[j + 1] - tmpDiff0[j]
                tmpMatrix[j] += tmpDiff1[j + 1]

    def query(self, r: int, c: int) -> int:
        """查询矩阵中指定位置的元素.
        !查询前需要先调用`update`方法.
        """
        return self.matrix[r][c]


class PreSumMatrix:
    """二维前缀和模板(矩阵不可变)"""

    __slots__ = "preSum"

    def __init__(self, mat: List[List[int]]):
        ROW, COL = len(mat), len(mat[0])
        preSum = [[0] * (COL + 1) for _ in range(ROW + 1)]
        for r in range(ROW):
            tmpSum0, tmpSum1 = preSum[r], preSum[r + 1]
            tmpM = mat[r]
            for c in range(COL):
                tmpSum1[c + 1] = tmpM[c] + tmpSum0[c + 1] + tmpSum1[c] - tmpSum0[c]
        self.preSum = preSum

    def queryRange(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        preSumMatrix.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self.preSum[r2 + 1][c2 + 1]
            - self.preSum[r2 + 1][c1]
            - self.preSum[r1][c2 + 1]
            + self.preSum[r1][c1]
        )


if __name__ == "__main__":
    diffmatrix = DiffMatrix([[1, 2], [3, 4]])
    print(diffmatrix.matrix)
    diffmatrix.add(0, 0, 1, 1, 1)

    diffmatrix.update()
    print(diffmatrix.matrix)
