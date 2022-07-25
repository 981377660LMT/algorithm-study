from typing import Final, Generic, List, Tuple, TypeVar
from abc import ABCMeta, abstractmethod
from math import ceil, sqrt


V = TypeVar("V")  # 区间元素类型
Q = TypeVar("Q")  # 每个查询的返回值类型


class AbstractMoAlgo(Generic[V, Q], metaclass=ABCMeta):

    _data: Final[List[V]]
    _queries: Final[List[Tuple[int, int, int]]]  # (index,left,right)

    def __init__(self, data: List[V]):
        self._data = data
        self._queries = []

    def work(self) -> List[Q]:
        """返回每个查询的结果"""
        self._sortQueries()

        nums, queries = self._data, self._queries
        res: List[Q] = [None] * len(queries)  # type: ignore
        left, right = 0, 0

        for qi, qLeft, qRight in queries:
            # !窗口收缩
            while right > qRight:
                right -= 1
                self._remove(nums[right], right, qLeft, qRight - 1)
            while left < qLeft:
                self._remove(nums[left], left, qLeft, qRight - 1)
                left += 1

            # !窗口扩张
            while right < qRight:
                self._add(nums[right], right, qLeft, qRight - 1)
                right += 1
            while left > qLeft:
                left -= 1
                self._add(nums[left], left, qLeft, qRight - 1)

            res[qi] = self._query()

        return res

    def addQuery(self, left: int, right: int) -> None:
        """0 <= left <= right < n"""
        self._queries.append((len(self._queries), left, right + 1))  # 注意这里的 right+1

    @abstractmethod
    def _add(self, value: V, index: int, qLeft: int, qRight: int) -> None:
        """将数据添加到窗口"""
        raise NotImplementedError(f"{self.__class__.__name__}._add")

    @abstractmethod
    def _remove(self, value: V, index: int, qLeft: int, qRight: int) -> None:
        """将数据从窗口中移除"""
        raise NotImplementedError(f"{self.__class__.__name__}._remove")

    @abstractmethod
    def _query(self) -> Q:
        """更新当前窗口的查询结果"""
        raise NotImplementedError(f"{self.__class__.__name__}._query")

    def _sortQueries(self) -> None:
        chunkSize = max(1, len(self._data) // sqrt(len(self._queries)))
        # chunkSize = ceil(sqrt(len(self._queries)))
        self._queries.sort(key=lambda x: (x[1] // chunkSize, x[2]))
