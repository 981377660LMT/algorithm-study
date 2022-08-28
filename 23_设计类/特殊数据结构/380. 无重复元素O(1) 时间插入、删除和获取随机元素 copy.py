# 因为所有元素不重复 因此indexOf查找元素相当于O(1) 用indexMap记录每个值的索引
# O(1) 时间插入、删除和获取随机元素
# !应用场景: 负载均衡器

from dataclasses import dataclass, field
from random import choice
from typing import Dict, List


@dataclass(slots=True)
class RandomizedSet:
    nums: List[int] = field(default_factory=list)
    indexMap: Dict[int, int] = field(default_factory=dict)

    def insert(self, val: int) -> bool:
        """添加一台新的服务器到整个集群中"""
        if val in self.indexMap:
            return False
        self.indexMap[val] = len(self.nums)
        self.nums.append(val)
        return True

    def remove(self, val: int) -> bool:
        """从集群中删除一个服务器"""
        if val not in self.indexMap:
            return False
        index = self.indexMap[val]
        self.nums[index] = self.nums[-1]
        self.indexMap[self.nums[index]] = index
        self.nums.pop()
        self.indexMap.pop(val)
        return True

    def getRandom(self) -> int:
        """在集群中随机（等概率）选择一个有效的服务器"""
        return choice(self.nums)
