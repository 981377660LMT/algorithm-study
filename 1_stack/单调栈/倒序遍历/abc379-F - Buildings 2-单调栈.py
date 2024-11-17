# F - Buildings 2 （类似1944. 队列中可以看到的人数）
# https://atcoder.jp/contests/abc379/tasks/abc379_f
#
# 给定n个建筑的高度hi，回答q个询问。每个询问给定l,r，问从l,l+1,...,r−1,r建筑往右看，都能看到的建筑的数量。
# 如果建筑i能看到建筑j，则不存在i<k<j，满足hk>hj。
# !所有建筑的高度都是不同的。


from typing import List, Tuple, Union


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


def buildings(heights: List[int], queries: List[Tuple[int, int]]) -> List[int]:
    n, q = len(heights), len(queries)
    groups = [[] for _ in range(n)]
    for i, (l, r) in enumerate(queries):
        groups[l].append((r, i))

    res = [0] * q
    bit = BITArray(n)
    stack = []  # 单调递减栈
    for l in range(n - 1, -1, -1):
        for r, i in groups[l]:
            res[i] = bit.queryRange(r, n)
        while stack and heights[stack[-1]] < heights[l]:
            popped = stack.pop()
            bit.add(popped, -1)
        stack.append(l)
        bit.add(l, 1)
    return res


if __name__ == "__main__":
    N, Q = map(int, input().split())
    H = list(map(int, input().split()))
    queries = []
    for _ in range(Q):
        l, r = map(int, input().split())
        l -= 1
        queries.append((l, r))
    res = buildings(H, queries)
    print("\n".join(map(str, res)))
