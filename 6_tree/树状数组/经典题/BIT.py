"""
有的题目下标从1开始计算方便一些 
此时在树状数组 add/query 入口处加上偏移量1即可
"""

from collections import defaultdict


class BIT1:
    """单点修改

    https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.py
    """

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        index += 1
        while index <= self.size:
            self.tree[index] += delta
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        index += 1
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self._lowbit(index)
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


class BIT2:
    """范围修改

    https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.py
    """

    def __init__(self, n: int):
        self.size = n
        self._tree1 = defaultdict(int)
        self._tree2 = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left, right]加delta"""
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left, right]的和"""
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")

        rawIndex = index
        while index <= self.size:
            self._tree1[index] += delta
            self._tree2[index] += (rawIndex - 1) * delta
            index += self._lowbit(index)

    def _query(self, index: int) -> int:
        if index > self.size:
            index = self.size

        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1[index] - self._tree2[index]
            index -= self._lowbit(index)
        return res


class BIT3:
    """单点修改 维护`前缀区间`最大值

    TODO: 正确性待讨论
    这么做正确的前提是不会删除或修改已经存进去的值，每次都是加入新的值，这样已经存在的最大值一直有效。
    """

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    # def update(self, left: int, right: int, target: int) -> None:
    #     """更新[left,right]区间的最大值为target"""
    #     ...

    # def query(self, left: int, right: int) -> int:
    #     """查询[left,right]的最大值"""
    #     ...

    def update(self, index: int, target: int) -> None:
        """将后缀区间`[index,size]`的最大值更新为target"""
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] = max(self.tree[index], target)
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        """查询前缀区间`[1,index]`的最大值"""
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res = max(res, self.tree[index])
            index -= self._lowbit(index)
        return res


class BIT4:
    """二维树状数组 单点修改+区间查询 每个操作都是 log(m*n)

    https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.py
    """

    def __init__(self, row: int, col: int) -> None:
        self.row = row
        self.col = col
        self.tree = defaultdict(lambda: defaultdict(int))

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def update(self, row: int, col: int, delta: int) -> None:
        """矩阵中的点 (row,col) 的值加上delta"""
        row, col = row + 1, col + 1
        curRow = row
        while curRow <= self.row:
            curCol = col
            while curCol <= self.col:
                self.tree[curRow][curCol] += delta
                curCol += self._lowbit(curCol)
            curRow += self._lowbit(curRow)

    def query(self, row: int, col: int) -> int:
        """左上角 (0,0) 到 右下角(row,col) 的矩形里所有数的和"""
        row, col = row + 1, col + 1
        if row > self.row:
            row = self.row
        if col > self.col:
            col = self.col
        res = 0
        curRow = row
        while curRow > 0:
            curCol = col
            while curCol > 0:
                res += self.tree[curRow][curCol]
                curCol -= self._lowbit(curCol)
            curRow -= self._lowbit(curRow)
        return res

    def queryRange(self, row1: int, col1: int, row2: int, col2: int) -> int:
        """查询左上角 (row1,col1) 到右下角 (row2,col2) 的和"""
        return (
            self.query(row2, col2)
            - self.query(row2, col1 - 1)
            - self.query(row1 - 1, col2)
            + self.query(row1 - 1, col1 - 1)
        )


# https://www.cnblogs.com/hbhszxyb/p/14157271.html
class BIT5:
    """二维树状数组 区间修改+区间查询 每个操作都是 log(m*n)

    https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.py
    """

    def __init__(self, row: int, col: int) -> None:
        self.row = row
        self.col = col
        self.tree1 = defaultdict(lambda: defaultdict(int))
        self.tree2 = defaultdict(lambda: defaultdict(int))
        self.tree3 = defaultdict(lambda: defaultdict(int))
        self.tree4 = defaultdict(lambda: defaultdict(int))

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def updateRange(
        self, row1: int, col1: int, row2: int, col2: int, delta: int
    ) -> None:
        """左上角 (row1,col1) 到右下角 (row2,col2) 的所有数加上delta"""
        self._update(row1, col1, delta)
        self._update(row2 + 1, col1, -delta)
        self._update(row1, col2 + 1, -delta)
        self._update(row2 + 1, col2 + 1, delta)

    def queryRange(self, row1: int, col1: int, row2: int, col2: int) -> int:
        """查询左上角 (row1,col1) 到右下角 (row2,col2) 的和"""
        return (
            self._query(row2, col2)
            - self._query(row2, col1 - 1)
            - self._query(row1 - 1, col2)
            + self._query(row1 - 1, col1 - 1)
        )

    def _update(self, row: int, col: int, delta: int) -> None:
        """[row,col]的值加上delta"""
        row, col = row + 1, col + 1
        preRow, preCol = row, col

        curRow = row
        while curRow <= self.row:
            curCol = col
            while curCol <= self.col:
                self.tree1[curRow][curCol] += delta
                self.tree2[curRow][curCol] += (preRow - 1) * delta
                self.tree3[curRow][curCol] += (preCol - 1) * delta
                self.tree4[curRow][curCol] += (preRow - 1) * (preCol - 1) * delta
                curCol += self._lowbit(curCol)
            curRow += self._lowbit(curRow)

    def _query(self, row: int, col: int) -> int:
        row, col = row + 1, col + 1
        if row > self.row:
            row = self.row
        if col > self.col:
            col = self.col

        preRow, preCol = row, col
        curRow = row
        res = 0
        while curRow > 0:
            curCol = col
            while curCol > 0:
                res += (
                    preRow * preCol * self.tree1[curRow][curCol]
                    - preCol * self.tree2[curRow][curCol]
                    - preRow * self.tree3[curRow][curCol]
                    + self.tree4[curRow][curCol]
                )
                curCol -= self._lowbit(curCol)
            curRow -= self._lowbit(curRow)
        return res


if __name__ == "__main__":
    bit1 = BIT1(100)
    bit1.add(0 + 1, 2)
    assert bit1.query(1) == 2
    assert bit1.queryRange(1, 4) == 2
    assert bit1.queryRange(2, 4) == 0
    assert bit1.queryRange(0, 102) == 2
    assert bit1.queryRange(0, 1000) == 2
    assert bit1.queryRange(-10000, 1000) == 2

    bit2 = BIT2(100)
    bit2.add(1, 1, 2)
    assert bit2.query(1, 1) == 2
    assert bit2.query(1, 4) == 2
    assert bit2.query(2, 4) == 0
    assert bit2.query(0, 102) == 2
    assert bit2.query(0, 1000) == 2
    assert bit2.query(-10000, 1000) == 2

    bit4 = BIT4(100, 100)
    bit4.update(0, 0, 2)
    assert bit4.query(0, 0) == 2
    bit4.update(1, 1, 2)
    assert bit4.query(1, 1) == 4

    bit5 = BIT5(100, 100)
    bit5.updateRange(0, 0, 3, 3, 1)
    assert bit5.queryRange(0, 0, 1, 1) == 4
    assert bit5.queryRange(0, 0, 3, 3) == 16
