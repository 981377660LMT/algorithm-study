from typing import Any, List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class C:
    __slots__ = "_counter", "_freqCounter"

    def __init__(self):
        self._counter = dict()
        self._freqCounter = dict()

    def add(self, v: Any) -> None:
        preC = self._counter.get(v, 0)
        self._counter[v] = preC + 1
        self._freqCounter[preC + 1] = self._freqCounter.get(preC + 1, 0) + 1
        if preC > 0:
            preF = self._freqCounter.get(preC, 0)
            if preF == 1:
                self._freqCounter.pop(preC)
            else:
                self._freqCounter[preC] = preF - 1

    def discard(self, v: Any) -> bool:
        preC = self._counter.get(v, 0)
        if preC == 0:
            return False
        if preC == 1:
            self._counter.pop(v)
        else:
            self._counter[v] = preC - 1
        preF = self._freqCounter.get(preC, 0)
        if preF == 1:
            self._freqCounter.pop(preC)
        else:
            self._freqCounter[preC] = preF - 1
        if preC > 1:
            self._freqCounter[preC - 1] = self._freqCounter.get(preC - 1, 0) + 1
        return True

    def check(self) -> bool:
        return len(self._freqCounter) == 1


class Solution:
    def queryResults(self, limit: int, queries: List[List[int]]) -> List[int]:
        Q = C()
        res = []
        color = dict()
        for x, y in queries:
            if x in color:
                Q.discard(color[x])
            color[x] = y
            Q.add(y)
            res.append(len(Q._counter))
        return res
