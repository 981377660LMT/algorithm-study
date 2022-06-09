from typing import List, Literal, Optional


class SegmentTreeNode:
    __slots__ = ('left', 'right', 'value')

    def __init__(
        self,
        value=0,
        left: Optional['SegmentTreeNode'] = None,
        right: Optional['SegmentTreeNode'] = None,
    ):
        self.value = value
        self.left = left
        self.right = right


class SegmentTree:
    def __init__(self, lower: int, upper: int):
        self._root = SegmentTreeNode()
        self._lower = lower
        self._upper = upper

    def update(self, left: int, right: int, target: Literal[0, 1]) -> None:
        self._checkRange(left, right)
        self._update(left, right, self._lower, self._upper, self._root, target)

    def query(self, left: int, right: int) -> int:
        self._checkRange(left, right)
        return self._query(left, right, self._lower, self._upper, self._root)

    def queryAll(self) -> int:
        return self._root.value

    def _update(
        self, L: int, R: int, l: int, r: int, root: SegmentTreeNode, target: Literal[0, 1]
    ) -> None:
        if L <= l and r <= R:
            root.value = (r - l + 1) * target
            return

        mid = (l + r) // 2
        self._pushDown(l, mid, r, root)
        if L <= mid:
            self._update(L, R, l, mid, root.left, target)
        if mid < R:
            self._update(L, R, mid + 1, r, root.right, target)
        self._pushUp(root)

    def _query(self, L: int, R: int, l: int, r: int, root: SegmentTreeNode) -> int:
        if L <= l and r <= R:
            return root.value

        mid = (l + r) // 2
        self._pushDown(l, mid, r, root)
        res = 0
        if L <= mid:
            res += self._query(L, R, l, mid, root.left)
        if mid < R:
            res += self._query(L, R, mid + 1, r, root.right)
        return res

    def _pushDown(self, left: int, mid: int, right: int, root: SegmentTreeNode) -> None:
        if not root.left:
            root.left = SegmentTreeNode()
        if not root.right:
            root.right = SegmentTreeNode()
        if root.value == right - left + 1:
            root.left.value = mid - left + 1
            root.right.value = right - mid
        elif root.value == 0:
            root.left.value = 0
            root.right.value = 0

    def _pushUp(self, root: SegmentTreeNode) -> None:
        root.value = root.left.value + root.right.value

    def _checkRange(self, left: int, right: int):
        if left < self._lower or right > self._upper:
            raise ValueError('Index out of range')


def main():
    tree = SegmentTree(0, 10)
    print('Initial tree:')
    print(tree.queryAll())
    print()

    tree.update(2, 6, 1)
    print('Updated tree:')
    print(tree.queryAll())
    print()

    tree.update(2, 3, 0)
    tree.update(6, 8, 0)
    tree.update(3, 5, 1)
    tree.update(8, 9, 1)
    print('Updated tree:')
    print(tree.queryAll())


if __name__ == '__main__':
    main()

    class Solution:
        def amountPainted(self, paint: List[List[int]]) -> List[int]:
            min_, max_ = paint[0][0] + 1, paint[0][1] + 1
            for start, end in paint:
                min_ = min(min_, start)
                max_ = max(max_, end)

            tree = SegmentTree(min_, max_)
            res = []
            for start, end in paint:
                start, end = start + 1, end + 1
                rangeSum = tree.query(start, end - 1)
                res.append(end - start - rangeSum)
                tree.update(start, end - 1, 1)

            return res
