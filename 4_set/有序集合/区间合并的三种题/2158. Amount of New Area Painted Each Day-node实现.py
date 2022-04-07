from typing import Callable, List


class Node:
    __slots__ = ('value', 'left', 'right', 'lazy')

    def __init__(self, value: int, left: int = -1, right: int = -1, lazy: int = 0):
        self.value = value
        self.left = left
        self.right = right
        self.lazy = lazy


Merge = Callable[[int, int], int]


class SegmentTree:
    """注意根节点从1开始,tree本身为[1,n]"""

    __slots__ = '_tree'

    def __init__(self, n: int):
        self._tree = [Node(0) for _ in range(n << 2)]

    def query(self, left: int, right: int) -> int:
        """[left,right]的和"""
        return self._query(1, left, right)

    def update(self, left: int, right: int, delta: int) -> None:
        """[left,right]区间更新delta"""
        self._update(1, left, right, delta)

    def build(self, rt: int, left: int, right: int) -> None:
        root = self._tree[rt]
        root.value, root.left, root.right, root.lazy = 0, left, right, 0
        # 到底部了，底部有n个结点
        if left == right:
            return

        mid = (left + right) >> 1
        self.build((rt << 1), left, mid)
        self.build((rt << 1) | 1, mid + 1, right)
        self._push_up(rt)

    def _update(self, rt: int, left: int, right: int, delta: int) -> None:
        root = self._tree[rt]
        # 到达了最细粒度的节点
        if left <= root.left and root.right <= right:
            root.lazy = 1
            root.value = root.right - root.left + 1
            return

        # 传递懒标记
        self._push_down(rt)
        mid = (root.left + root.right) >> 1
        if left <= mid:
            self._update((rt << 1), left, right, delta)
        if mid < right:
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
            res += self._query((rt << 1), left, right)
        if mid < right:
            res += self._query((rt << 1) | 1, left, right)
        return res

    def _push_up(self, rt: int) -> None:
        # 用子节点更新父节点状态
        root, left, right = self._tree[rt], self._tree[(rt << 1)], self._tree[(rt << 1) | 1]
        root.value = left.value + right.value

    def _push_down(self, rt: int) -> None:
        root, left, right = self._tree[rt], self._tree[(rt << 1)], self._tree[(rt << 1) | 1]
        if root.lazy:
            # 先向下传递更新
            left.lazy += root.lazy
            right.lazy += root.lazy
            # 左侧总是比右侧长，区间的批量更新
            left.value = left.right - left.left + 1
            right.value = right.right - right.left + 1
            # 传递后清空增量数组的父节点
            root.lazy = 0


# 本题线段树update需要做特殊处理，染过的值至多为1
class Solution:
    def amountPainted(self, paint: List[List[int]]) -> List[int]:
        min_, max_ = paint[0][0] + 1, paint[0][1] + 1
        for start, end in paint:
            min_ = min(min_, start)
            max_ = max(max_, end)

        tree = SegmentTree(max_)
        tree.build(1, min_ + 1, max_ + 1)
        res = []
        for start, end in paint:
            start, end = start + 1, end + 1
            rangeSum = tree.query(start, end - 1)
            res.append(end - start - rangeSum)
            tree.update(start, end - 1, 1)

        return res
