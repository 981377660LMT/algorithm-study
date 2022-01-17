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


if __name__ == '__main__':
    bit = BIT(100)
    bit.add(0 + 1, 2)
    print(bit.query(1))
    print(bit.sumRange(1, 4))
    print(bit.sumRange(2, 4))
    print(bit.sumRange(0, 102))
    print(bit.sumRange(0, 1000))
    print(bit.sumRange(-10000, 1000))

