from collections import defaultdict


class BIT1:
    """单点修改"""

 def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] += delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)

class VirtualArray:
    def __init__(self):
        self.bit = BIT1(1 << 32 + 10)

    def set(self, start, end):
        start, end = start + 1, end + 1
        self.bit.add(start, 1)
        self.bit.add(end + 1, -1)

    def get(self, idx):
        idx = idx + 1
        return self.bit.query(idx) > 0

