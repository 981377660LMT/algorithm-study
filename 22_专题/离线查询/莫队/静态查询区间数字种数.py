# AcWing 2492. HH的项链

# 静态查询区间数字种数
# 三种做法，莫队，离线树状数组，主席树

from collections import defaultdict
from typing import List
from MoAlgo import AbstractMoAlgo


class QueryTypeMoAlgo(AbstractMoAlgo[int, int]):
    """静态查询区间数字种数"""

    def __init__(self, data: List[int]):
        super().__init__(data)
        self._count = 0
        self._counter = defaultdict(int)

    def _add(self, value: int, index: int, qLeft: int, qRight: int) -> None:
        self._counter[value] += 1
        if self._counter[value] == 1:
            self._count += 1

    def _remove(self, value: int, index: int, qLeft: int, qRight: int) -> None:
        self._counter[value] -= 1
        if self._counter[value] == 0:
            self._count -= 1

    def _query(self) -> int:
        return self._count
