from typing import Generic, List, TypeVar, Dict
from bisect import bisect_left

T = TypeVar("T")


class RangeFreq(Generic[T]):
    __slots__ = "_valueToIndexes"

    def __init__(self, nums: List[T]) -> None:
        self._valueToIndexes = dict()
        for i, v in enumerate(nums):
            self._valueToIndexes.setdefault(v, []).append(i)

    def query(self, start: int, end: int, value: T) -> int:
        if start >= end:
            return 0
        if value not in self._valueToIndexes:
            return 0
        pos = self._valueToIndexes[value]
        return bisect_left(pos, end) - bisect_left(pos, start)

    def findFirst(self, start: int, end: int, value: T) -> int:
        if start >= end:
            return -1
        if value not in self._valueToIndexes:
            return -1
        pos = self._valueToIndexes[value]
        idx = bisect_left(pos, start)
        if idx < len(pos) and pos[idx] < end:
            return pos[idx]
        return -1

    def findLast(self, start: int, end: int, value: T) -> int:
        if start >= end:
            return -1
        if value not in self._valueToIndexes:
            return -1
        pos = self._valueToIndexes[value]
        idx = bisect_left(pos, end)
        if idx > 0 and pos[idx - 1] >= start:
            return pos[idx - 1]
        return -1


def test_range_freq():
    # 测试基本功能
    nums = [1, 2, 3, 2, 4, 2, 5]
    rf = RangeFreq(nums)

    # 测试query方法
    assert rf.query(0, 7, 2) == 3  # 元素2在整个数组中出现3次
    assert rf.query(1, 5, 2) == 2  # 元素2在区间[1,5)中出现2次
    assert rf.query(0, 3, 2) == 1  # 元素2在区间[0,3)中出现1次
    assert rf.query(0, 7, 6) == 0  # 元素6不存在

    # 测试findFirst方法
    assert rf.findFirst(0, 7, 2) == 1  # 元素2的第一个位置是1
    assert rf.findFirst(2, 7, 2) == 3  # 在区间[2,7)中元素2的第一个位置是3
    assert rf.findFirst(4, 7, 2) == 5  # 在区间[4,7)中元素2的第一个位置是5
    assert rf.findFirst(0, 7, 6) == -1  # 元素6不存在
    assert rf.findFirst(6, 7, 2) == -1  # 区间[6,7)中没有元素2

    # 测试findLast方法
    assert rf.findLast(0, 7, 2) == 5  # 元素2的最后一个位置是5
    assert rf.findLast(0, 4, 2) == 3  # 在区间[0,4)中元素2的最后一个位置是3
    assert rf.findLast(0, 2, 2) == 1  # 在区间[0,2)中元素2的最后一个位置是1
    assert rf.findLast(0, 7, 6) == -1  # 元素6不存在
    assert rf.findLast(0, 1, 2) == -1  # 区间[0,1)中没有元素2

    # 测试边界情况
    assert rf.query(3, 3, 2) == 0  # 空区间
    assert rf.findFirst(3, 3, 2) == -1  # 空区间
    assert rf.findLast(3, 3, 2) == -1  # 空区间

    # 测试字符串类型
    str_nums = ["a", "b", "c", "b", "d", "b"]
    str_rf = RangeFreq(str_nums)
    assert str_rf.query(0, 6, "b") == 3
    assert str_rf.findFirst(0, 6, "b") == 1
    assert str_rf.findLast(0, 6, "b") == 5

    print("所有测试通过!")


if __name__ == "__main__":
    test_range_freq()
