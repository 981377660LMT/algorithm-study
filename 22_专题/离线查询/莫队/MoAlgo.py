from typing import Generic, List, Tuple, TypeVar
from abc import ABCMeta, abstractmethod
from math import ceil, sqrt


V = TypeVar("V")  # 区间元素类型
Q = TypeVar("Q")  # 每个查询的返回值类型
Query = Tuple[int, int, int]


class AbstractMoAlgo(Generic[V, Q], metaclass=ABCMeta):
    """莫队算法基类

    左端点分桶，右端点排序
    """

    def __init__(self, data: List[V]):
        self._n = len(data)
        self._data = data
        self._chunkSize = ceil(sqrt(self._n))  # self._chunkSize = ceil(n / sqrt(2 * q))
        self._buckets = [[] for _ in range(self._n // self._chunkSize + 1)]
        self._queryOrder = 0

    def addQuery(self, left: int, right: int) -> None:
        """0 <= left <= right < n"""
        index = left // self._chunkSize
        self._buckets[index].append((self._queryOrder, left, right + 1))  # 注意这里的 right+1
        self._queryOrder += 1

    def work(self) -> List[Q]:
        """返回每个查询的结果"""
        data, buckets, q = self._data, self._buckets, self._queryOrder
        res: List[Q] = [None] * q  # type: ignore
        left, right = 0, 0

        for i, bucket in enumerate(buckets):
            bucket.sort(key=lambda x: x[2], reverse=not not i & 1)

            for qi, qLeft, qRight in bucket:
                # !窗口收缩
                while right > qRight:
                    right -= 1
                    self._remove(data[right], right, qLeft, qRight - 1)
                while left < qLeft:
                    self._remove(data[left], left, qLeft, qRight - 1)
                    left += 1

                # !窗口扩张
                while right < qRight:
                    self._add(data[right], right, qLeft, qRight - 1)
                    right += 1
                while left > qLeft:
                    left -= 1
                    self._add(data[left], left, qLeft, qRight - 1)

                res[qi] = self._query()

        return res

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


if __name__ == "__main__":
    #  https://atcoder.jp/contests/abc242/tasks/abc242_g
    class Solution(AbstractMoAlgo[int, int]):
        """静态查询区间 `元素的count //2` 的和"""

        def __init__(self, data: List[int]):
            super().__init__(data)
            self._pair = 0
            self._counter = [0] * int(1e5 + 10)

        def _add(self, value: int, index: int, qLeft: int, qRight: int) -> None:
            self._pair -= self._counter[value] // 2
            self._counter[value] += 1
            self._pair += self._counter[value] // 2

        def _remove(self, value: int, index: int, qLeft: int, qRight: int) -> None:
            self._pair -= self._counter[value] // 2
            self._counter[value] -= 1
            self._pair += self._counter[value] // 2

        def _query(self) -> int:
            return self._pair

    n = int(input())
    nums = list(map(int, input().split()))
    M = Solution(nums)
    q = int(input())
    for _ in range(q):
        l, r = map(int, input().split())
        l, r = l - 1, r - 1
        M.addQuery(l, r)

    res = M.work()
    print(*res, sep="\n")
