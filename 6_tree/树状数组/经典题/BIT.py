"""
下标从0开始
"""


from collections import defaultdict
from typing import List, Sequence, Union


class BITArray:
    """Point Add Range Sum, 0-indexed."""

    @staticmethod
    def _build(sequence: Sequence[int]) -> List[int]:
        tree = [0] * (len(sequence) + 1)
        for i in range(1, len(tree)):
            tree[i] += sequence[i - 1]
            parent = i + (i & -i)
            if parent < len(tree):
                tree[parent] += tree[i]
        return tree

    __slots__ = ("_n", "_tree")

    def __init__(self, lenOrSequence: Union[int, Sequence[int]]):
        if isinstance(lenOrSequence, int):
            self._n = lenOrSequence
            self._tree = [0] * (lenOrSequence + 1)
        else:
            self._n = len(lenOrSequence)
            self._tree = self._build(lenOrSequence)

    def add(self, index: int, delta: int) -> None:
        index += 1
        while index <= self._n:
            self._tree[index] += delta
            index += index & -index

    def query(self, right: int) -> int:
        """Query sum of [0, right)."""
        if right > self._n:
            right = self._n
        res = 0
        while right > 0:
            res += self._tree[right]
            right -= right & -right
        return res

    def queryRange(self, left: int, right: int) -> int:
        """Query sum of [left, right)."""
        return self.query(right) - self.query(left)

    def __len__(self) -> int:
        return self._n

    def __repr__(self) -> str:
        nums = []
        for i in range(1, self._n + 1):
            nums.append(self.queryRange(i, i + 1))
        return f"BITArray({nums})"


class BIT1:
    """单点修改"""

    __slots__ = "size", "bit", "tree"

    def __init__(self, n: int):
        self.size = n + 5
        self.bit = n.bit_length()
        self.tree = dict()

    def add(self, index: int, delta: int) -> None:
        index += 1
        while index <= self.size:
            self.tree[index] = self.tree.get(index, 0) + delta
            index += index & -index

    def query(self, right: int) -> int:
        """Query sum of [0, right)."""
        if right > self.size:
            right = self.size
        res = 0
        while right > 0:
            res += self.tree.get(right, 0)
            right -= right & -right
        return res

    def queryRange(self, left: int, right: int) -> int:
        """Query sum of [left, right)."""
        return self.query(right) - self.query(left)

    def bisectLeft(self, k: int) -> int:
        """返回第一个前缀和大于等于k的位置pos
        0 <= pos <= self.size
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) < k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos

    def bisectRight(self, k: int) -> int:
        """返回第一个前缀和大于k的位置pos
        0 <= pos <= self.size
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) <= k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos

    def __repr__(self) -> str:
        arr = []
        for i in range(self.size):
            arr.append(self.queryRange(i, i + 1))
        return str(arr)

    def __len__(self) -> int:
        return self.size


class BIT2:
    """范围修改,0-indexed"""

    __slots__ = "size", "_tree1", "_tree2"

    def __init__(self, n: int):
        self.size = n + 5
        self._tree1 = dict()
        self._tree2 = dict()

    def add(self, left: int, right: int, delta: int) -> None:
        """区间[left, right)加delta."""
        right -= 1
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """区间[left, right)的和."""
        right -= 1
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        index += 1
        rawIndex = index
        while index <= self.size:
            self._tree1[index] = self._tree1.get(index, 0) + delta
            self._tree2[index] = self._tree2.get(index, 0) + (rawIndex - 1) * delta
            index += index & -index

    def _query(self, index: int) -> int:
        index += 1
        if index > self.size:
            index = self.size
        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1.get(index, 0) - self._tree2.get(index, 0)
            index &= index - 1
        return res

    def __repr__(self):
        arr = []
        for i in range(self.size):
            arr.append(self.query(i, i + 1))
        return str(arr)

    def __len__(self):
        return self.size


class BIT3:
    """
    单点修改、前缀最大值查询 维护`前缀区间`最大值
    """

    def __init__(self, n: int):
        self.size = n + 5
        self.tree = dict()

    def update(self, index: int, target: int) -> None:
        """将`index`位置的值与`target`取op运算"""
        index += 1
        while index <= self.size:
            self.tree[index] = max(self.tree.get(index, 0), target)
            index += index & -index

    def query(self, right: int) -> int:
        """查询前缀区间`[0,right)`的最大值"""
        if right > self.size:
            right = self.size
        res = 0
        while right > 0:
            res = max(res, self.tree.get(right, 0))
            right -= right & -right
        return res


class BIT4:
    """二维树状数组 单点修改+区间查询 每个操作都是 log(m*n)"""

    def __init__(self, row: int, col: int) -> None:
        self.row = row
        self.col = col
        self.tree = defaultdict(lambda: defaultdict(int))

    def add(self, row: int, col: int, delta: int) -> None:
        """矩阵中的点 (row,col) 的值加上delta"""
        row, col = row + 1, col + 1
        curRow = row
        while curRow <= self.row:
            curCol = col
            while curCol <= self.col:
                self.tree[curRow][curCol] += delta
                curCol += curCol & -curCol
            curRow += curRow & -curRow

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
                curCol -= curCol & -curCol
            curRow -= curRow & -curRow
        return res

    def queryRange(self, row1: int, col1: int, row2: int, col2: int) -> int:
        """查询左上角 (row1,col1) 到右下角 (row2,col2) 的和"""
        return (
            self.query(row2, col2)
            - self.query(row2, col1 - 1)
            - self.query(row1 - 1, col2)
            + self.query(row1 - 1, col1 - 1)
        )


# !二维单点修改、区间查询
# https://leetcode.cn/problems/range-sum-query-2d-mutable/solution/er-wei-bitshu-tao-shu-onjian-shu-ologn2c-nxvy/
# !二维区间修改、区间查询
# https://www.cnblogs.com/hbhszxyb/p/14157271.html


# 非常慢 不要使用python版本
class BIT5:
    """二维树状数组 区间修改+区间查询 每个操作都是 log(m*n)"""

    def __init__(self, row: int, col: int) -> None:
        self.row = row
        self.col = col
        self.tree1 = defaultdict(lambda: defaultdict(int))
        self.tree2 = defaultdict(lambda: defaultdict(int))
        self.tree3 = defaultdict(lambda: defaultdict(int))
        self.tree4 = defaultdict(lambda: defaultdict(int))

    def addRange(self, row1: int, col1: int, row2: int, col2: int, delta: int) -> None:
        """左上角 (row1,col1) 到右下角 (row2,col2) 的所有数加上delta"""
        self._add(row1, col1, delta)
        self._add(row2 + 1, col1, -delta)
        self._add(row1, col2 + 1, -delta)
        self._add(row2 + 1, col2 + 1, delta)

    def queryRange(self, row1: int, col1: int, row2: int, col2: int) -> int:
        """查询左上角 (row1,col1) 到右下角 (row2,col2) 的和"""
        return (
            self._query(row2, col2)
            - self._query(row2, col1 - 1)
            - self._query(row1 - 1, col2)
            + self._query(row1 - 1, col1 - 1)
        )

    def _add(self, row: int, col: int, delta: int) -> None:
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
                curCol += curCol & -curCol
            curRow += curRow & -curRow

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
                curCol -= curCol & -curCol
            curRow -= curRow & -curRow
        return res


if __name__ == "__main__":
    bit1 = BIT1(100)
    bit1.add(1, 2)
    assert bit1.query(2) == 2
    assert bit1.queryRange(1, 4) == 2
    assert bit1.queryRange(2, 4) == 0
    assert bit1.queryRange(0, 102) == 2
    assert bit1.queryRange(0, 1000) == 2
    assert bit1.queryRange(-10000, 1000) == 2
    assert bit1.bisectLeft(2) == 1
    assert bit1.bisectRight(2) == len(bit1)

    bit2 = BIT2(100)
    bit2.add(1, 2, 2)
    assert bit2.query(1, 2) == 2
    assert bit2.query(1, 4) == 2
    assert bit2.query(2, 4) == 0
    assert bit2.query(0, 102) == 2
    assert bit2.query(0, 1000) == 2
    assert bit2.query(-10000, 1000) == 2

    bit3 = BIT3(100)
    bit3.update(1, 2)
    bit3.update(2, 3)
    bit3.update(4, 5)
    print(bit3.query(4))  # 5
    bit3.update(4, 1)
    print(bit3.query(4))  # 5 不可以修改原来的值(变小)

    bit4 = BIT4(100, 100)
    bit4.add(0, 0, 2)
    assert bit4.query(0, 0) == 2
    bit4.add(1, 1, 2)
    assert bit4.query(1, 1) == 4

    bit5 = BIT5(100, 100)
    bit5.addRange(0, 0, 3, 3, 1)
    assert bit5.queryRange(0, 0, 1, 1) == 4
    assert bit5.queryRange(0, 0, 3, 3) == 16

    arrayBIT = BITArray([1, 2, 3])
    assert arrayBIT.queryRange(1, 3) == 5
    arrayBIT.add(1, 1)
    assert arrayBIT.queryRange(1, 3) == 6
    arrayBIT.add(2, 1)
    assert arrayBIT.queryRange(1, 3) == 7
