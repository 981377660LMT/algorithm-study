# Not Verifined

# https://slidesplayer.com/slide/16477387/
# https://sotanishy.github.io/cp-library-cpp/data-structure/partition_refinement.hpp

# Partition refinement は，素集合を管理するデータ構造である．Union find が集合の併合を処理するのに対し，partition refinement は集合の分割を処理する．

# DFA 最小化などに用いられる．


from collections import defaultdict
from typing import List, Tuple
from sortedcontainers import SortedList


class PartitionRefinement:
    """分区细化."""

    __slots__ = ("_sets", "_cls")

    def __init__(self, n: int) -> None:
        self._sets = [SortedList(list(range(n)))]
        self._cls = [0] * n

    def refine(self, pivot: List[int]) -> List[Tuple[int, int]]:
        """
        给定一个划分,将每个集合Si划分为 Si ∩ P 和 Si\\P 两个集合.
        返回分割后的集合编号对.
        """
        mp = defaultdict(list)
        for x in pivot:
            if x in self:
                i = self._cls[x]
                mp[i].append(x)
                self._sets[i].remove(x)

        updated = []
        keys = sorted(mp)
        for i in keys:
            s = mp[i]
            ni = len(self._sets)
            self._sets.append(SortedList(s))
            if len(self._sets[i]) < len(self._sets[ni]):
                self._sets[i], self._sets[ni] = self._sets[ni], self._sets[i]  # type: ignore
            if not self._sets[ni]:
                self._sets.pop()
                continue
            for x in self._sets[ni]:
                self._cls[x] = ni
            updated.append((i, ni))
        return updated

    def find(self, x: int) -> int:
        return self._cls[x]

    def discard(self, x: int) -> None:
        if x in self:
            self._sets[self._cls[x]].discard(x)
            self._cls[x] = -1

    def isConnected(self, x: int, y: int) -> bool:
        cx = self.find(x)
        cy = self.find(y)
        return cx != -1 and cy != -1 and cx == cy

    def getSize(self, i: int) -> int:
        return len(self._sets[i])

    def getGroup(self, i: int) -> List[int]:
        return list(self._sets[i])

    def __contains__(self, x: int) -> bool:
        if x < 0 or x >= len(self._cls):
            return False
        return self._cls[x] != -1


if __name__ == "__main__":
    pr = PartitionRefinement(10)
    print(pr.getGroup(0))
    print(pr.refine([2, 1, 3]))
    print(pr.getGroup(0))
    print(pr.getGroup(1))
    print(0 in pr)
    pr.discard(0)
    print(0 in pr)
