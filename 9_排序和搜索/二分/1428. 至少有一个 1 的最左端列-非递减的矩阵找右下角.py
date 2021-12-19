class BinaryMatrix(object):
    def get(self, row: int, col: int) -> int:
        """"sumary_line"""

    def dimensions(self) -> list:
        """"sumary_line"""


# 好的继续找，坏的换一行
class Solution:
    def leftMostColumnWithOne(self, binaryMatrix: 'BinaryMatrix') -> int:
        # 类似二分查找的思想。从右下角开始，遇到1往左走(因为想找到包含 1 的最左端列的索引)，遇到0往上走(看看本列还有没有1)，走到第一行，答案就很显然了.
        dimensions = binaryMatrix.dimensions()
        row, col = dimensions[0] - 1, dimensions[1] - 1
        res = -1
        while row >= 0 and col >= 0:
            if binaryMatrix.get(row, col) == 1:
                res = col
                col -= 1
            else:
                row -= 1
        return res
