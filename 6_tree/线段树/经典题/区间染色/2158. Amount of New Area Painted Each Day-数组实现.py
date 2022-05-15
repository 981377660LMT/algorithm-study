from typing import Callable, List

Merge = Callable[[int, int], int]


class SegmentTree:
    """注意根节点从1开始,tree本身为[1,n]"""

    def __init__(self, n: int):
        self._n = n
        self._tree = [0 for _ in range(n << 2)]
        self._lazy = [0 for _ in range(n << 2)]

    def query(self, l: int, r: int) -> int:
        """[left,right]的和"""
        assert 1 <= l <= r <= self._n
        return self._query(1, l, r, 1, self._n)

    def update(self, l: int, r: int, delta: int) -> None:
        """[left,right]区间更新delta"""
        assert 1 <= l <= r <= self._n
        self._update(1, l, r, 1, self._n, delta)

    def _update(self, rt: int, L: int, R: int, l: int, r: int, delta: int) -> None:
        """L,R表示需要update的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            self._lazy[rt] = 1
            self._tree[rt] = r - l + 1
            return

        # 传递懒标记
        mid = (l + r) >> 1
        self._push_down(rt, l, r, mid)
        if L <= mid:
            self._update((rt << 1), L, R, l, mid, delta)
        if mid < R:
            self._update((rt << 1) | 1, L, R, mid + 1, r, delta)
        self._push_up(rt)

    def _query(self, rt: int, L: int, R: int, l: int, r: int) -> int:
        """L,R表示需要query的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            return self._tree[rt]

        # 传递懒标记
        mid = (l + r) >> 1
        self._push_down(rt, l, r, mid)
        res = 0
        if L <= mid:
            res += self._query((rt << 1), L, R, l, mid)
        if mid < R:
            res += self._query((rt << 1) | 1, L, R, mid + 1, r)
        return res

    def _push_up(self, rt: int) -> None:
        # 用子节点更新父节点状态
        self._tree[rt] = self._tree[rt << 1] + self._tree[rt << 1 | 1]

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._lazy[rt]:
            self._lazy[rt << 1] += self._lazy[rt]
            self._lazy[rt << 1 | 1] += self._lazy[rt]
            self._tree[rt << 1] = mid - l + 1
            self._tree[rt << 1 | 1] = r - mid
            self._lazy[rt] = 0


# 本题线段树update需要做特殊处理，染过的值至多为1
class Solution:
    def amountPainted(self, paint: List[List[int]]) -> List[int]:
        min_, max_ = paint[0][0] + 1, paint[0][1] + 1
        for start, end in paint:
            min_ = min(min_, start)
            max_ = max(max_, end)

        tree = SegmentTree(max_ + 1)
        res = []
        for start, end in paint:
            start, end = start + 1, end + 1
            rangeSum = tree.query(start, end - 1)
            res.append(end - start - rangeSum)
            tree.update(start, end - 1, 1)

        return res


print(Solution().amountPainted([[1, 2], [2, 3], [3, 4]]))

