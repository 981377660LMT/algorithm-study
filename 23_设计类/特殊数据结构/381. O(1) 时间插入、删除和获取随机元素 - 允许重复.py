from collections import defaultdict
from random import choice


class RandomizedCollection:
    def __init__(self):
        self.indexMap = defaultdict(set)
        self.nums = []

    def insert(self, val: int) -> bool:
        self.indexMap[val].add(len(self.nums))
        self.nums.append(val)
        return len(self.indexMap[val]) == 1

    def remove(self, val: int) -> bool:
        if not self.indexMap[val]:
            return False
        self.nums[(i := self.indexMap[val].pop())] = self.nums[-1]
        self.indexMap[(last := self.nums.pop())].discard(len(self.nums))
        i < len(self.nums) and self.indexMap[last].add(i)
        return True

    def getRandom(self) -> int:
        return choice(self.nums)

