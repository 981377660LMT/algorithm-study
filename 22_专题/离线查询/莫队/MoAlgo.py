from typing import Generic, List, TypeVar
from math import sqrt


V = TypeVar("V")  # 区间元素类型
Q = TypeVar("Q")  # 每个查询的返回值类型


class MoAlgo(Generic[V, Q]):

    """
    莫队算法模板
    左端点分桶，右端点排序
    """

    def __init__(self, n: int, q: int):
        self._chunkSize = max(1, n // int(sqrt(q)))
        self._buckets = [[] for _ in range(n // self._chunkSize + 1)]
        self._queryOrder = 0

    def addQuery(self, left: int, right: int) -> None:
        """0 <= left <= right < n"""
        index = left // self._chunkSize
        self._buckets[index].append((self._queryOrder, left, right + 1))  # 注意这里的 right+1
        self._queryOrder += 1

    def work(self) -> List[Q]:
        """返回每个查询的结果"""
        buckets = self._buckets
        res: List[Q] = [None] * self._queryOrder  # type: ignore
        left, right = 0, 0

        for i, bucket in enumerate(buckets):
            bucket.sort(key=lambda x: x[2], reverse=not not i & 1)

            for qi, qLeft, qRight in bucket:
                # !窗口扩张
                while left > qLeft:
                    left -= 1
                    self._add(left, -1)
                while right < qRight:
                    self._add(right, 1)
                    right += 1

                # !窗口收缩
                while left < qLeft:
                    self._remove(left, 1)
                    left += 1
                while right > qRight:
                    right -= 1
                    self._remove(right, -1)

                res[qi] = self._query(qLeft, qRight - 1)

        return res

    def _add(self, index: int, delta: int) -> None:
        """将数据添加到窗口"""
        raise NotImplementedError(f"{self.__class__.__name__}._add")

    def _remove(self, index: int, delta: int) -> None:
        """将数据从窗口中移除"""
        raise NotImplementedError(f"{self.__class__.__name__}._remove")

    def _query(self, qLeft: int, qRight: int) -> Q:
        """更新当前窗口的查询结果"""
        raise NotImplementedError(f"{self.__class__.__name__}._query")
