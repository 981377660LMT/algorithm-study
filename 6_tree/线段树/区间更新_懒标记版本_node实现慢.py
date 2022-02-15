# https://www.desgard.com/algo/docs/part3/ch02/3-segment-tree-range/
# https://leetcode-cn.com/problems/amount-of-new-area-painted-each-day/solution/lan-biao-ji-xian-duan-shu-by-xin_cheng-2wfo/

# push_up 是通过子结点数值求和来构造父结点
# 当我们需要对一个区间进行批量增减操作的时候，我们只要向下更新到我们所有查询操作的最小粒度即可，而不用完全对整个线段树进行更新
# lazy数组代表对应的节点待更新的增量

# push_down 也就是向下更新的意思
# 如果增量数组有值，则我们将其向下更新
# 此时的增量数组已经被放置在了最小颗粒度查询节点
# 向下更新是为了更新到 sum 数组


# 当你事先知道class的attributes的时候，建议使用slots来节省memory以及获得更快的attribute access
from typing import Callable


class Node:
    __slots__ = ('value', 'left', 'right', 'lazy')

    def __init__(self, value: int, left: int = -1, right: int = -1, lazy: int = 0):
        self.value = value
        self.left = left
        self.right = right
        self.lazy = lazy


Merge = Callable[[int, int], int]


class SegmentTree:
    """注意根节点从1开始,tree本身为[1,len(nums)]"""

    __slots__ = ('_n', '_tree', '_merge')

    def __init__(self, n: int, merge: Merge = None):
        self._n = n
        self._tree = [Node(0) for _ in range(n << 2)]
        self._merge: Merge = (lambda x, y: x + y) if merge is None else merge
        # 根结点必须为1
        self._build(1, 1, n)

    def query(self, left: int, right: int) -> int:
        """[left,right]的和"""
        assert 1 <= left <= right <= self._n
        return self._query(1, left, right)

    def update(self, left: int, right: int, delta: int) -> None:
        """[left,right]区间更新delta"""
        assert 1 <= left <= right <= self._n
        self._update(1, left, right, delta)

    def _build(self, rt: int, left: int, right: int) -> None:
        root = self._tree[rt]
        root.value, root.left, root.right, root.lazy = 0, left, right, 0
        # 到底部了，底部有n个结点
        if left == right:
            return

        mid = (left + right) >> 1
        self._build((rt << 1), left, mid)
        self._build((rt << 1) | 1, mid + 1, right)
        self._push_up(rt)

    def _update(self, rt: int, left: int, right: int, delta: int) -> None:
        root = self._tree[rt]
        # 到达了最细粒度的节点
        if left <= root.left and root.right <= right:
            root.lazy += delta
            root.value += delta * (root.right - root.left + 1)
            return

        # 传递懒标记
        self._push_down(rt)
        mid = (root.left + root.right) >> 1
        if left <= mid:
            self._update((rt << 1), left, right, delta)
        if mid + 1 <= right:
            self._update((rt << 1) | 1, left, right, delta)
        self._push_up(rt)

    def _query(self, rt: int, left: int, right: int) -> int:
        root = self._tree[rt]
        if left <= root.left and root.right <= right:
            return root.value

        # 传递懒标记
        self._push_down(rt)
        mid = (root.left + root.right) >> 1
        res = 0
        if left <= mid:
            res = self._merge(res, self._query((rt << 1), left, right))
        if mid + 1 <= right:
            res = self._merge(res, self._query((rt << 1) | 1, left, right))
        return res

    def _push_up(self, rt: int) -> None:
        # 用子节点更新父节点状态
        root, left, right = self._tree[rt], self._tree[(rt << 1)], self._tree[(rt << 1) | 1]
        root.value = self._merge(left.value, right.value)

    def _push_down(self, rt: int) -> None:
        root, left, right = self._tree[rt], self._tree[(rt << 1)], self._tree[(rt << 1) | 1]
        if root.lazy:
            # 先向下传递更新
            left.lazy += root.lazy
            right.lazy += root.lazy
            # 左侧总是比右侧长，区间的批量更新
            left.value += (left.right - left.left + 1) * root.lazy
            right.value += (right.right - right.left + 1) * root.lazy
            # 传递后清空增量数组的父节点
            root.lazy = 0


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
    # assert tree.query(2, 1) == 0

