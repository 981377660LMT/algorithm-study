# 1409. 查询带键的排列
# https://leetcode.cn/problems/queries-on-a-permutation-with-key/description/
# 给定m，你可以创建一个1到m的列表，再给定queries。
# 每次找到这个列表中queries[i]的位置并存入答案，再把这个数挪到列表的首位。
# 注意，这里的数字是1~m，其对应的位置是0~m-1。

from typing import List


from typing import List, Union


class BITArray:
    __slots__ = "n", "total", "_data"

    def __init__(self, sizeOrData: Union[int, List[int]]):
        if isinstance(sizeOrData, int):
            self.n = sizeOrData
            self.total = 0
            self._data = [0] * sizeOrData
        else:
            self.n = len(sizeOrData)
            self.total = sum(sizeOrData)
            _data = sizeOrData[:]
            for i in range(1, self.n + 1):
                j = i + (i & -i)
                if j <= self.n:
                    _data[j - 1] += _data[i - 1]
            self._data = _data

    def add(self, index: int, value: int) -> None:
        self.total += value
        index += 1
        while index <= self.n:
            self._data[index - 1] += value
            index += index & -index

    def queryPrefix(self, end: int) -> int:
        if end > self.n:
            end = self.n
        res = 0
        while end > 0:
            res += self._data[end - 1]
            end &= end - 1
        return res

    def queryRange(self, start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > self.n:
            end = self.n
        if start >= end:
            return 0
        if start == 0:
            return self.queryPrefix(end)
        pos, neg = 0, 0
        while end > start:
            pos += self._data[end - 1]
            end &= end - 1
        while start > end:
            neg += self._data[start - 1]
            start &= start - 1
        return pos - neg

    def queryAll(self):
        return self.total

    def __repr__(self):
        res = [self.queryRange(i, i + 1) for i in range(self.n)]
        return f"BITArray({res})"


class Solution:
    def processQueries(self, queries: List[int], m: int) -> List[int]:
        q = len(queries)
        bit = BITArray(m + q)

        pos = [0] * m
        for i in range(m):
            pos[i] = q + i
            bit.add(q + i, 1)

        res = []
        for i, query in enumerate(queries):
            query -= 1
            cur = pos[query]
            bit.add(cur, -1)
            preCount = bit.queryPrefix(cur)
            res.append(preCount)
            pos[query] = q - 1 - i
            bit.add(q - 1 - i, 1)
        return res
