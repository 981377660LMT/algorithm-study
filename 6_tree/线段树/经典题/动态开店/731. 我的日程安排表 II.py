"""
动态开点 叠加更新 查询最大值
"""

from typing import Optional

from sortedcontainers import SortedDict


class MyCalendarTwo:
    def __init__(self):
        self.tree = SegmentTree(0, int(1e9) + 5)

    def book(self, start: int, end: int) -> bool:
        preMax = self.tree.query(start, end - 1)
        if preMax < 2:
            self.tree.update(start, end - 1, 1)
            return True
        return False


class SegmentTreeNode:
    __slots__ = ("left", "right", "value", "isLazy", "lazyValue")

    def __init__(
        self,
        value=0,
        left: Optional["SegmentTreeNode"] = None,
        right: Optional["SegmentTreeNode"] = None,
    ):
        self.value = value
        self.left = left
        self.right = right
        self.isLazy = False
        self.lazyValue = 0


class SegmentTree:
    def __init__(self, lower: int, upper: int):
        self._root = SegmentTreeNode()
        self._lower = lower
        self._upper = upper

    def update(self, left: int, right: int, delta: int) -> None:
        self._checkRange(left, right)
        self._update(left, right, self._lower, self._upper, self._root, delta)

    def query(self, left: int, right: int) -> int:
        self._checkRange(left, right)
        return self._query(left, right, self._lower, self._upper, self._root)

    def queryAll(self) -> int:
        return self._root.value

    def _update(self, L: int, R: int, l: int, r: int, root: SegmentTreeNode, delta: int) -> None:
        if L <= l and r <= R:
            root.value += delta
            root.lazyValue += delta
            root.isLazy = True
            return

        mid = (l + r) // 2
        self._pushDown(l, mid, r, root)
        if L <= mid:
            self._update(L, R, l, mid, root.left, delta)
        if mid < R:
            self._update(L, R, mid + 1, r, root.right, delta)
        self._pushUp(root)

    def _query(self, L: int, R: int, l: int, r: int, root: SegmentTreeNode) -> int:
        if L <= l and r <= R:
            return root.value

        mid = (l + r) // 2
        self._pushDown(l, mid, r, root)
        res = 0
        if L <= mid:
            res = max(res, self._query(L, R, l, mid, root.left))
        if mid < R:
            res = max(res, self._query(L, R, mid + 1, r, root.right))
        return res

    def _pushDown(self, left: int, mid: int, right: int, root: SegmentTreeNode) -> None:
        if not root.left:
            root.left = SegmentTreeNode()
        if not root.right:
            root.right = SegmentTreeNode()
        if root.isLazy:
            root.left.isLazy = True
            root.right.isLazy = True

            root.left.lazyValue += root.lazyValue
            root.right.lazyValue += root.lazyValue
            root.left.value += root.lazyValue
            root.right.value += root.lazyValue

            root.lazyValue = 0
            root.isLazy = False

    def _pushUp(self, root: SegmentTreeNode) -> None:
        root.value = max(root.left.value, root.right.value)

    def _checkRange(self, left: int, right: int):
        if left < self._lower or right > self._upper:
            raise ValueError("Index out of range")


# Your MyCalendarTwo object will be instantiated and called as such:
# obj = MyCalendarTwo()
# param_1 = obj.book(start,end)


# 差分数组
class MyCalendarTwo2:
    def __init__(self):
        self.diff = SortedDict()

    def book(self, start: int, end: int) -> bool:
        self.diff[start] = self.diff.get(start, 0) + 1
        self.diff[end] = self.diff.get(end, 0) - 1
        curSum = 0
        for delta in self.diff.values():
            curSum += delta
            if curSum > 2:
                self.diff[start] -= 1
                self.diff[end] += 1
                return False
        return True
