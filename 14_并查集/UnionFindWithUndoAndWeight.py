# UnionFindWithUndoAndWeight
# 可撤销并查集,维护连通分量权值
# SetGroup: 将下标为index元素`所在集合`的权值置为value.
# GetGroup: 获取下标为index元素`所在集合`的权值.
# !Undo: 撤销上一次合并(Union)或者修改权值(Set)操作，没合并成功也要撤销
# Reset: 撤销所有操作


from collections import defaultdict
from typing import DefaultDict, List


class UnionFindArrayWithUndoAndWeight:
    """可撤销并查集,维护连通分量权值."""

    __slots__ = ("_rank", "_parents", "_weights", "_history", "part")

    def __init__(self, initWeights: List[int]):
        n = len(initWeights)
        self._rank = [1] * n
        self._parents = list(range(n))
        self._weights = initWeights[:]
        self._history = []
        self.part = n

    def setGroupWeight(self, index: int, value: int) -> None:
        """将下标为index元素`所在集合`的权值置为value."""
        index = self.find(index)
        self._history.append((index, self._rank[index], self._weights[index]))
        self._weights[index] = value

    def getGroupWeight(self, index: int) -> int:
        """获取下标为index元素`所在集合`的权值."""
        return self._weights[self.find(index)]

    def undo(self) -> None:
        """撤销上一次合并(Union)或者修改权值(SetGroup)操作."""
        if not self._history:
            return
        small, rank, weight = self._history.pop()
        ps = self._parents[small]
        self._weights[ps] = weight
        self._rank[ps] = rank
        if ps != small:
            self._parents[small] = small
            self.part += 1

    def reset(self) -> None:
        """撤销所有操作."""
        while self._history:
            self.undo()

    def find(self, x: int) -> int:
        if self._parents[x] == x:
            return x
        return self.find(self._parents[x])

    def union(self, x: int, y: int) -> bool:
        x, y = self.find(x), self.find(y)
        if self._rank[x] < self._rank[y]:
            x, y = y, x
        self._history.append((y, self._rank[x], self._weights[x]))
        if x != y:
            self._parents[y] = x
            self._rank[x] += self._rank[y]
            self._weights[x] += self._weights[y]
            self.part -= 1
            return True
        return False

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getSize(self, x: int) -> int:
        return self._rank[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for i in range(len(self._parents)):
            groups[self.find(i)].append(i)
        return groups

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())


if __name__ == "__main__":
    uf = UnionFindArrayWithUndoAndWeight([1, 2, 3, 4, 5])
    print(uf)
    uf.union(0, 1)
    print(uf.getGroupWeight(0))
    uf.setGroupWeight(0, 10)
    print(uf.getGroupWeight(0))
    uf.undo()
    print(uf.getGroupWeight(0))
    uf.undo()
    print(uf.getGroupWeight(0))
