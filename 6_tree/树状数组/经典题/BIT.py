from collections import defaultdict


class BIT1:
    """单点修改"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError('index 必须是正整数')
        while index <= self.size:
            self.tree[index] += delta
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self._lowbit(index)
        return res

    def sumRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


class BIT2:
    """范围修改"""

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
            raise ValueError('index 必须是正整数')

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
            raise ValueError('index 必须是正整数')
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


if __name__ == '__main__':
    bit1 = BIT1(100)
    bit1.add(0 + 1, 2)
    assert bit1.query(1) == 2
    assert bit1.sumRange(1, 4) == 2
    assert bit1.sumRange(2, 4) == 0
    assert bit1.sumRange(0, 102) == 2
    assert bit1.sumRange(0, 1000) == 2
    assert bit1.sumRange(-10000, 1000) == 2

    bit2 = BIT2(100)
    bit2.add(1, 1, 2)
    assert bit2.query(1, 1) == 2
    assert bit2.query(1, 4) == 2
    assert bit2.query(2, 4) == 0
    assert bit2.query(0, 102) == 2
    assert bit2.query(0, 1000) == 2
    assert bit2.query(-10000, 1000) == 2

    bit3 = BIT3(100)

