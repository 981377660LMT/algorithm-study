from typing import List
from python线段树模板 import SegmentTree


class NumArray:
    def __init__(self, nums: List[int]):
        self.nums = nums
        self.tree = SegmentTree(nums)

    def update(self, index: int, val: int) -> None:
        self.tree.update(index + 1, index + 1, val - self.nums[index])
        self.nums[index] = val

    def sumRange(self, left: int, right: int) -> int:
        return self.tree.query(left + 1, right + 1)


# Your NumArray object will be instantiated and called as such:
obj = NumArray([1, 3, 5])
print(obj.sumRange(0, 2))
obj.update(1, 2)
print(obj.sumRange(0, 2))

# 9 8
obj = NumArray([-1, 2, 3])
print(obj.sumRange(0, 2))
obj.update(0, 1)
print(obj.sumRange(0, 0))
