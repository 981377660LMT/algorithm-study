# 因为所有元素不重复 因此indexOf查找元素相当于O(1) 用indexMap记录每个值的索引

from dataclasses import dataclass, field
from random import choice
from typing import Dict, List


@dataclass(slots=True)
class RandomizedSet:
    nums: List[int] = field(default_factory=list)
    indexMap: Dict[int, int] = field(default_factory=dict)

    def insert(self, val: int) -> bool:
        if val in self.indexMap:
            return False
        self.indexMap[val] = len(self.nums)
        self.nums.append(val)
        return True

    def remove(self, val: int) -> bool:
        """把要删除的元素和最后一个元素交换"""
        if val not in self.indexMap:
            return False
        index = self.indexMap[val]
        self.nums[index] = self.nums[-1]
        self.indexMap[self.nums[index]] = index
        self.nums.pop()
        self.indexMap.pop(val)
        return True

    def getRandom(self) -> int:
        return choice(self.nums)

