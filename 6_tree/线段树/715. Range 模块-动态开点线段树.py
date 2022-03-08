# 1 <= left < right <= 1e9

# 由于值域很大且是在线query，无法做离散化
# 所以需要动态开点
# 区间操作，线段树


from typing import Optional


class Node:
    # __slots__允许我们声明并限定类成员，并拒绝类创建__dict__和__weakref__属性以节约内存空间。
    __slots__ = ('isTracked', 'lazy', 'left', 'right')

    def __init__(
        self,
        isTracked=False,  # 表示当前区间所有元素是否被追踪
        lazy=False,
        left: Optional['Node'] = None,
        right: Optional['Node'] = None,
    ) -> None:
        self.isTracked = isTracked
        self.lazy = lazy
        self.left = left
        self.right = right


class SegmentTree:
    def __init__(self) -> None:
        self._root = Node()

    def update(self, left: int, right: int, delta: bool) -> None:
        self._update(left, right, 0, int(1e9 + 10), self._root, delta)

    def query(self, left: int, right: int) -> bool:
        return self._query(left, right, 0, int(1e9 + 10), self._root)

    def _update(self, L: int, R: int, l: int, r: int, root: Node, delta: bool) -> None:
        if L <= l <= r <= R:
            root.isTracked = delta
            root.lazy = True
            return

        self._pushDown(root)
        mid = (l + r) >> 1
        if L <= mid:
            self._update(L, R, l, mid, root.left, delta)
        if R >= mid + 1:
            self._update(L, R, mid + 1, r, root.right, delta)
        self._pushUp(root)

    def _query(self, L: int, R: int, l: int, r: int, root: Node) -> bool:
        if L <= l <= r <= R:
            return root.isTracked

        self._pushDown(root)
        mid = (l + r) >> 1
        res = True
        if L <= mid:
            res = res and self._query(L, R, l, mid, root.left)
        if R >= mid + 1:
            res = res and self._query(L, R, mid + 1, r, root.right)
        return res

    def _pushUp(self, root: Node) -> None:
        root.isTracked = not not (
            root.left and root.left.isTracked and root.right and root.right.isTracked
        )

    def _pushDown(self, root: Node) -> None:
        if not root.left:
            root.left = Node()
        if not root.right:
            root.right = Node()
        if root.lazy:
            root.left.lazy = root.right.lazy = True
            root.left.isTracked = root.right.isTracked = root.isTracked
            root.lazy = False


# addRange,queryRange,removeRange时间均为O(logN)，N为值域，本题为1e9
# 空间复杂度为O(QlogN)，Q为所有操作的总次数


class RangeModule:
    def __init__(self):
        self.tree = SegmentTree()

    def addRange(self, left: int, right: int) -> None:
        """添加 半开区间 [left, right)"""
        self.tree.update(left, right - 1, True)

    def queryRange(self, left: int, right: int) -> bool:
        """ 只有在当前正在跟踪区间 [left, right) 中的每一个实数时，才返回 true"""
        return self.tree.query(left, right - 1)

    def removeRange(self, left: int, right: int) -> None:
        """ 停止跟踪 半开区间 [left, right)"""
        self.tree.update(left, right - 1, False)

