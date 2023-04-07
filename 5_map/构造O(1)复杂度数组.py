# 设计一个特殊的数组，要求该数据结构以下三种操作的时间复杂度均为O(1)
# 1. 查询数组某个位置的元素
# 2. 将数组某个位置的元素修改为指定值
# 3. 将数组所有元素修改为指定值


from collections import defaultdict


class SpecialArray:
    __slots__ = "_data"

    def __init__(self) -> None:
        self._data = defaultdict(int)

    def get(self, index: int) -> int:
        return self._data[index]

    def set(self, index: int, value: int) -> None:
        self._data[index] = value

    def setAll(self, value: int) -> None:
        self._data = defaultdict(lambda: value)
