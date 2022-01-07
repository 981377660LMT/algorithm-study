from collections import defaultdict


# tree直接用dict 省去离散化步骤
class BIT:
    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int) -> None:
        while index <= self.size:
            self.tree[index] += delta
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self._lowbit(index)
        return res

    def sumRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)

