from bisect import bisect_right


class DiffMap2:
    """支持区间和查询的差分dict."""

    __slots__ = ("_diff", "_pos", "_sum0", "_sum1", "_dirty")

    def __init__(self):
        self._diff = dict()
        self._pos = []
        self._sum0 = []
        self._sum1 = []
        self._dirty = False

    def add_range(self, start: int, end: int, delta: int):
        if start >= end:
            return
        self._dirty = True
        self._diff[start] = self._diff.get(start, 0) + delta
        self._diff[end] = self._diff.get(end, 0) - delta

    def get(self, pos: int) -> int:
        self._build()
        i = bisect_right(self._pos, pos)
        return 0 if i == 0 else self._sum0[i - 1]

    def get_range(self, start: int, end: int) -> int:
        if start >= end:
            return 0
        self._build()
        return self._presum(end) - self._presum(start)

    def _build(self):
        if not self._dirty:
            return
        self._dirty = False

        pos = sorted(self._diff)
        if not pos:
            return

        self._pos = pos
        self._sum0 = [0] * len(pos)
        self._sum1 = [0] * len(pos)
        pre = pos[0]
        s0 = self._diff[pre]
        s1 = 0
        self._sum0[0] = s0
        for i in range(1, len(pos)):
            cur = pos[i]
            s1 += (cur - pre) * s0
            self._sum1[i] = s1
            s0 += self._diff[cur]
            self._sum0[i] = s0
            pre = cur

    def _presum(self, v: int) -> int:
        if not self._pos:
            return 0
        if v <= self._pos[0]:
            return 0
        if v >= self._pos[-1]:
            return self._sum1[-1]
        i = bisect_right(self._pos, v)
        if i == 0:
            return 0
        res = self._sum1[i - 1]
        width = v - self._pos[i - 1]
        res += width * self._sum0[i - 1]
        return res


if __name__ == "__main__":
    from typing import List

    class Solution:
        # 3413. 收集连续 K 个袋子可以获得的最多硬币数量
        # https://leetcode.cn/problems/maximum-coins-from-k-consecutive-bags/description/
        def maximumCoins(self, coins: List[List[int]], k: int) -> int:
            seg = DiffMap2()
            for l, r, w in coins:
                seg.add_range(l, r + 1, w)
            res = 0
            for l, r, _ in coins:
                res = max(res, seg.get_range(l, l + k))
                res = max(res, seg.get_range(r - k + 1, r + 1))
            return res

        # 2251. 花期内花的数目
        # https://leetcode.cn/problems/number-of-flowers-in-full-bloom/description/
        def fullBloomFlowers(self, flowers: List[List[int]], people: List[int]) -> List[int]:
            diff = DiffMap2()
            for l, r in flowers:
                diff.add_range(l, r + 1, 1)
            return [diff.get(p) for p in people]
