# https://www.desgard.com/algo/docs/part3/ch02/3-segment-tree-range/
# https://leetcode-cn.com/problems/amount-of-new-area-painted-each-day/solution/lan-biao-ji-xian-duan-shu-by-xin_cheng-2wfo/

# push_up 是通过子结点数值求和来构造父结点
# 当我们需要对一个区间进行批量增减操作的时候，我们只要向下更新到我们所有查询操作的最小粒度即可，而不用完全对整个线段树进行更新
# lazy数组代表对应的节点待更新的增量

# push_down 也就是向下更新的意思
# 如果增量数组有值，则我们将其向下更新
# 此时的增量数组已经被放置在了最小颗粒度查询节点
# 向下更新是为了更新到 sum 数组


from typing import Callable

Merge = Callable[[int, int], int]


class SegmentTree:
    """注意根节点从1开始,tree本身为[1,n]"""

    def __init__(self, n: int, merge: Merge = None):
        self._n = n
        self._tree = [0 for _ in range(n << 2)]
        self._lazy = [0 for _ in range(n << 2)]
        self._merge: Merge = (lambda x, y: x + y) if merge is None else merge
        # self._build(1, 1, n)

    def query(self, l: int, r: int) -> int:
        """[left,right]的和"""
        # assert 1 <= l <= r <= self._n
        return self._query(1, l, r, 1, self._n)

    def update(self, l: int, r: int, delta: int) -> None:
        """[left,right]区间更新delta"""
        # assert 1 <= l <= r <= self._n
        self._update(1, l, r, 1, self._n, delta)

    # def _build(self, rt: int, l: int, r: int) -> None:
    #     """传了nums时，用于初始化线段树"""
    #     # 到底部了，底部有n个结点
    #     if l == r:
    #         self._tree[rt] = self.nums[l-1]
    #         return

    #     mid = (l + r) >> 1
    #     self._build((rt << 1), l, mid)
    #     self._build((rt << 1) | 1, mid + 1, r)
    #     self._push_up(rt)

    def _update(self, rt: int, L: int, R: int, l: int, r: int, delta: int) -> None:
        """L,R表示需要update的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            self._lazy[rt] += delta
            self._tree[rt] += delta * (r - l + 1)
            return

        # 传递懒标记
        mid = (l + r) >> 1
        self._push_down(rt, l, r, mid)
        if L <= mid:
            self._update((rt << 1), L, R, l, mid, delta)
        if mid + 1 <= R:
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
            res = self._merge(res, self._query((rt << 1), L, R, l, mid))
        if mid + 1 <= R:
            res = self._merge(res, self._query((rt << 1) | 1, L, R, mid + 1, r))
        return res

    def _push_up(self, rt: int) -> None:
        # 用子节点更新父节点状态
        self._tree[rt] = self._merge(self._tree[rt << 1], self._tree[rt << 1 | 1])

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._lazy[rt]:
            self._lazy[rt << 1] += self._lazy[rt]
            self._lazy[rt << 1 | 1] += self._lazy[rt]
            self._tree[rt << 1] += self._lazy[rt] * (mid - l + 1)
            self._tree[rt << 1 | 1] += self._lazy[rt] * (r - mid)
            self._lazy[rt] = 0


if __name__ == '__main__':
    tree = SegmentTree(10)
    assert tree.query(1, 1) == 0
    assert tree.query(1, 6) == 0
    tree.update(1, 6, 4)
    assert tree.query(1, 6) == 24
    assert tree.query(1, 1) == 4
    tree.update(1, 1, 4)
    assert tree.query(1, 1) == 8
    assert tree.query(2, 2) == 4
    assert tree.query(2, 1) == 0

