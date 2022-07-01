# get时bfs计算
from collections import defaultdict, deque
from typing import List, Tuple


class Excel:
    def __init__(self, height: int, width: str):
        self.grid = defaultdict(lambda: defaultdict(int))  # excel实际值
        self.dep = defaultdict(lambda: defaultdict(list))  # excel的依赖表达式

    @staticmethod
    def _parseColRow(colRow: str) -> List[Tuple[int, int]]:
        """解析 ColRow 字符串，返回包含的单元格"""
        if not colRow:
            return []
        if len(colRow) <= 3:
            row, col = int(colRow[1:]) - 1, ord(colRow[0]) - 65
            return [(row, col)]
        colRow1, colRow2 = colRow.split(":")
        row1, col1 = int(colRow1[1:]) - 1, ord(colRow1[0]) - 65
        row2, col2 = int(colRow2[1:]) - 1, ord(colRow2[0]) - 65
        return [
            (row, col) for row in range(row1, row2 + 1) for col in range(col1, col2 + 1)
        ]

    def set(self, row: int, column: str, val: int) -> None:
        """设置 C(row, column) 中的值为 val。"""
        r, c = row - 1, ord(column) - 65
        self.grid[r][c] = val
        self.dep[r][c] = []

    def get(self, row: int, column: str) -> int:
        """返回 C(row, column) 中的值 题目规定不存在环依赖"""
        r, c = row - 1, ord(column) - 65
        if not self.dep[r][c]:
            return self.grid[r][c]

        res = 0
        queue = deque([(r, c)])
        while queue:
            curR, curC = queue.popleft()
            if not self.dep[curR][curC]:
                res += self.grid[curR][curC]
            else:
                for rowCol in self.dep[curR][curC]:
                    nexts = self._parseColRow(rowCol)
                    queue.extend(nexts)
        return res

    def sum(self, row: int, column: str, numbers: List[str]) -> int:
        """这个函数会将计算的结果放入 C(row, column) 中，计算的结果等于在 numbers 中代表的所有元素之和

        numbers[i] 的格式为 "ColRow" 或 "ColRow1:ColRow2".
        """
        r, c = row - 1, ord(column) - 65
        self.dep[r][c] = numbers
        return self.get(row, column)


if __name__ == "__main__":
    # 输入:
    # ["Excel", "set", "sum", "set", "get"]
    # [[3, "C"], [1, "A", 2], [3, "C", ["A1", "A1:B2"]], [2, "B", 2], [3, "C"]]

    # 来源：力扣（LeetCode）
    # 链接：https://leetcode.cn/problems/design-excel-sum-formula
    # 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
    excel = Excel(3, "C")
    excel.set(1, "A", 2)
    print(excel.sum(3, "C", ["A1", "A1:B2"]))
    excel.set(2, "B", 2)
    # print(excel.get(3, "C"))
