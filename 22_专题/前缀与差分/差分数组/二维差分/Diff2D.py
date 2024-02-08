# DiffMatrix/Diff2D
from typing import List


class Diff2D:
    """二维差分."""

    __slots__ = ("matrix", "_diff", "_row", "_col", "_dirty")

    def __init__(self, mat: List[List[int]]):
        self._row, self._col = len(mat), len(mat[0])
        self.matrix = [[0] * self._col for _ in range(self._row)]
        for i in range(self._row):
            tmp1, tmp2 = self.matrix[i], mat[i]
            for j in range(self._col):
                tmp1[j] = tmp2[j]
        self._diff = [[0] * (self._col + 2) for _ in range(self._row + 2)]
        self._dirty = False

    def add(self, r1: int, c1: int, r2: int, c2: int, delta: int) -> None:
        """
        区间更新左上角`(row1, col1)` 到 右下角`(row2, col2)`闭区间所描述的子矩阵的元素.
        0<=r1<=r2<row, 0<=c1<=c2<col.
        """
        self._diff[r1 + 1][c1 + 1] += delta
        self._diff[r1 + 1][c2 + 2] -= delta
        self._diff[r2 + 2][c1 + 1] -= delta
        self._diff[r2 + 2][c2 + 2] += delta
        self._dirty = True

    def query(self, r: int, c: int) -> int:
        """查询矩阵中指定位置的元素."""
        if self._dirty:
            self.build()
        return self.matrix[r][c]

    def build(self) -> None:
        """遍历矩阵，还原对应元素的增量"""
        if not self._dirty:
            return
        self._dirty = False
        for i in range(self._row):
            tmpDiff0, tmpDiff1 = self._diff[i], self._diff[i + 1]
            tmpMatrix = self.matrix[i]
            for j in range(self._col):
                tmpDiff1[j + 1] += tmpDiff1[j] + tmpDiff0[j + 1] - tmpDiff0[j]
                tmpMatrix[j] += tmpDiff1[j + 1]


class PreSum2D:
    """二维前缀和模板."""

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
    diffmatrix = Diff2D([[1, 2], [3, 4]])
    print(diffmatrix.matrix)
    diffmatrix.add(0, 0, 1, 1, 1)
    diffmatrix.build()
    print(diffmatrix.matrix)

    # 2536. 子矩阵元素加 1
    # https://leetcode.cn/problems/increment-submatrices-by-one/
    class Solution1:
        def rangeAddQueries(self, n: int, queries: List[List[int]]) -> List[List[int]]:
            M = Diff2D([[0] * n for _ in range(n)])
            for r1, c1, r2, c2 in queries:
                M.add(r1, c1, r2, c2, 1)
            M.build()
            return M.matrix

    # 给你一个下标从 0 开始、大小为 m x n 的网格 image ，表示一个灰度图像，其中 image[i][j] 表示在范围 [0..255] 内的某个像素强度。另给你一个 非负 整数 threshold 。
    # 如果 image[a][b] 和 image[c][d] 满足 |a - c| + |b - d| == 1 ，则称这两个像素是 相邻像素 。
    # 区域 是一个 3 x 3 的子网格，且满足区域中任意两个 相邻 像素之间，像素强度的 绝对差 小于或等于 threshold 。
    # 区域 内的所有像素都认为属于该区域，而一个像素 可以 属于 多个 区域。
    # 你需要计算一个下标从 0 开始、大小为 m x n 的网格 result ，其中 result[i][j] 是 image[i][j] 所属区域的 平均 强度，向下取整 到最接近的整数。如果 image[i][j] 属于多个区域，result[i][j] 是这些区域的 “取整后的平均强度” 的 平均值，也 向下取整 到最接近的整数。如果 image[i][j] 不属于任何区域，则 result[i][j] 等于 image[i][j] 。
    # 返回网格 result 。
    #
    # !图像平滑去噪? 3*3卷积核?
    class Solution:
        def resultGrid(self, image: List[List[int]], threshold: int) -> List[List[int]]:
            DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]

            def check(row1: int, col1: int, row2: int, col2: int) -> bool:
                for r in range(row1, row2 + 1):
                    for c in range(col1, col2 + 1):
                        for dr, dc in DIR4:
                            nr, nc = r + dr, c + dc
                            if row1 <= nr <= row2 and col1 <= nc <= col2:
                                if abs(image[nr][nc] - image[r][c]) > threshold:
                                    return False
                return True

            ROW_SIZE = 3
            COL_SIZE = 3
            ROW, COL = len(image), len(image[0])
            S = PreSum2D(image)
            D1 = Diff2D([[0] * COL for _ in range(ROW)])
            D2 = Diff2D([[0] * COL for _ in range(ROW)])
            for r1 in range(ROW - ROW_SIZE + 1):
                for c1 in range(COL - COL_SIZE + 1):
                    r2, c2 = r1 + ROW_SIZE - 1, c1 + COL_SIZE - 1
                    if check(r1, c1, r2, c2):
                        mean = S.queryRange(r1, c1, r2, c2) // (ROW_SIZE * COL_SIZE)
                        D1.add(r1, c1, r2, c2, 1)
                        D2.add(r1, c1, r2, c2, mean)

            res = [[0] * COL for _ in range(ROW)]
            for r in range(ROW):
                for c in range(COL):
                    count, total = D1.query(r, c), D2.query(r, c)
                    res[r][c] = total // count if count else image[r][c]
            return res
