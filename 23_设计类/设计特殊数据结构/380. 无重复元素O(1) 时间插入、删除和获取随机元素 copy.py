# 因为所有元素不重复 因此indexOf查找元素相当于O(1) 用indexMap记录每个值的索引
# O(1) 时间插入、删除和获取随机元素
# !应用场景: 负载均衡器/抽奖
# !技巧:删除元素时,将最后一个元素放到要删除的元素的位置,然后删除最后一个元素

from random import choice


class RandomizedSet:
    __slots__ = ("_nums", "_indexMap")

    def __init__(self) -> None:
        self._nums = []
        self._indexMap = dict()

    def insert(self, val: int) -> bool:
        """添加一台新的服务器到整个集群中"""
        if val in self._indexMap:
            return False
        self._indexMap[val] = len(self._nums)
        self._nums.append(val)
        return True

    def remove(self, val: int) -> bool:
        """从集群中删除一个服务器"""
        if val not in self._indexMap:
            return False
        index = self._indexMap[val]
        self._nums[index] = self._nums[-1]
        self._indexMap[self._nums[index]] = index
        self._nums.pop()
        self._indexMap.pop(val)
        return True

    def getRandom(self) -> int:
        """在集群中随机（等概率）选择一个有效的服务器"""
        return choice(self._nums)
