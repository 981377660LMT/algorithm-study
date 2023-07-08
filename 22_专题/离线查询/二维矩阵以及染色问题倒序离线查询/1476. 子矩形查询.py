# https://leetcode.cn/problems/subrectangle-queries/
# 最多有 500 次updateSubrectangle 和 getValue 操作。
# 1 <= rows, cols <= 100

# !1.二维线段树(因为是区间修改,所以树套树)
# !2.保存修改倒序查询的做法
#   优化 => 保存修改，当修改超过一个阈值，比如1000，批量处理这些修改，扫描线+一维线段树更新
#   更新控制在O(sqrt)次


from typing import List


class SubrectangleQueries:
    def __init__(self, rectangle: List[List[int]]):
        self.rectangle = [row[:] for row in rectangle]
        self.updates = []

    def updateSubrectangle(self, row1: int, col1: int, row2: int, col2: int, newValue: int) -> None:
        """用 newValue 更新以 (row1,col1) 为左上角且以 (row2,col2) 为右下角的子矩形。"""
        self.updates.append((row1, col1, row2, col2, newValue))

    def getValue(self, row: int, col: int) -> int:
        """返回矩形中坐标 (row,col) 的当前值。"""
        for row1, col1, row2, col2, newValue in reversed(self.updates):
            if row1 <= row <= row2 and col1 <= col <= col2:
                return newValue
        return self.rectangle[row][col]


# Your SubrectangleQueries object will be instantiated and called as such:
# obj = SubrectangleQueries(rectangle)
# obj.updateSubrectangle(row1,col1,row2,col2,newValue)
# param_2 = obj.getValue(row,col)
