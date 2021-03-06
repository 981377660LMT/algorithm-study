from typing import List, Optional


class SegmentTree:
    """区间叠加线段树

    注意根节点从1开始,tree本身为[1,n]
    """

    def __init__(self, n: int, nums: Optional[List[int]] = None):
        self._n = n
        self._tree = [0 for _ in range(self._n << 2)]
        self._lazy = [0 for _ in range(self._n << 2)]
        if nums is not None:
            self._build(1, 1, self._n, nums)

    def query(self, left: int, right: int) -> int:
        """[left,right]的和"""
        if left > right:
            return 0
        assert 1 <= left <= right <= self._n
        return self._query(1, left, right, 1, self._n)

    def update(self, left: int, right: int, delta: int) -> None:
        """[left,right]区间更新delta"""
        assert 1 <= left <= right <= self._n
        self._update(1, left, right, 1, self._n, delta)

    def _build(self, rt: int, left: int, right: int, nums: List[int]) -> None:
        """传了nums时，用于初始化线段树"""
        # 到底部了，底部有n个结点
        if left == right:
            self._tree[rt] = nums[left - 1]
            return

        mid = (left + right) >> 1
        self._build((rt << 1), left, mid, nums)
        self._build((rt << 1) | 1, mid + 1, right, nums)
        self._push_up(rt)

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
            self._tree[rt << 1] += self._lazy[rt] * (mid - l + 1)
            self._tree[rt << 1 | 1] += self._lazy[rt] * (r - mid)
            self._lazy[rt] = 0


if __name__ == "__main__":
    tree1 = SegmentTree(10)
    assert tree1.query(1, 1) == 0
    assert tree1.query(1, 6) == 0
    tree1.update(1, 6, 4)
    assert tree1.query(1, 6) == 24
    assert tree1.query(1, 1) == 4
    tree1.update(1, 1, 4)
    assert tree1.query(1, 1) == 8
    assert tree1.query(2, 2) == 4
    assert tree1.query(2, 1) == 0

    tree2 = SegmentTree(5, [1, 2, 3, 4, 5])
    assert tree2.query(1, 1) == 1
    assert tree2.query(1, 5) == 15
    tree2.update(1, 3, 4)
    assert tree2.query(1, 5) == 27
    assert tree2.query(1, 1) == 5
    tree2.update(1, 1, 4)
    assert tree2.query(1, 1) == 9
    assert tree2.query(2, 2) == 6
    assert tree2.query(2, 1) == 0

    # tree3 = SegmentTree(10, lambda x, y: max(x, y))
    # assert tree3.query(1, 1) == 0
    # tree3.update(1, 6, 4)
    # print(tree3.query(1, 2))
    # assert tree3.query(1, 5) == 4
    # assert tree3.query(1, 1) == 4
    # tree3.update(1, 1, 4)
    # assert tree3.query(1, 1) == 8
    # assert tree3.query(2, 2) == 4
    # assert tree3.query(2, 1) == 0
