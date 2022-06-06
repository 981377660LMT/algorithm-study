from typing import Callable, List, Optional, Union, overload

Merge = Callable[[int, int], int]


class SegmentTree:
    """注意根节点从1开始,tree本身为[1,n]"""

    __slots__ = ('_n', '_tree', '_lazy', '_merge')

    def __init__(self, n_or_nums: Union[int, List[int]], merge: Optional[Merge] = None):
        if isinstance(n_or_nums, list):
            self._n = len(n_or_nums)
            self._tree = [0 for _ in range(self._n << 2)]
            self._lazy = [0 for _ in range(self._n << 2)]
            self._merge: Merge = (lambda x, y: x + y) if merge is None else merge
            self._build(1, 1, self._n, n_or_nums)
        else:
            self._n = n_or_nums
            self._tree = [0 for _ in range(self._n << 2)]
            self._lazy = [0 for _ in range(self._n << 2)]
            self._merge: Merge = (lambda x, y: x + y) if merge is None else merge

    def query(self, l: int, r: int) -> int:
        """[left,right]的和"""
        # assert 1 <= l <= r <= self._n
        return self._query(1, l, r, 1, self._n)

    def update(self, l: int, r: int, delta: int) -> None:
        """[left,right]区间更新delta"""
        # assert 1 <= l <= r <= self._n
        self._update(1, l, r, 1, self._n, delta)

    def _build(self, rt: int, l: int, r: int, nums: List[int]) -> None:
        """传了nums时，用于初始化线段树"""
        # 到底部了，底部有n个结点
        if l == r:
            self._tree[rt] = nums[l - 1]
            return

        mid = (l + r) >> 1
        self._build((rt << 1), l, mid, nums)
        self._build((rt << 1) | 1, mid + 1, r, nums)
        self._push_up(rt)

    def _update(self, rt: int, L: int, R: int, l: int, r: int, delta: int) -> None:
        """L,R表示需要update的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            self._lazy[rt] += delta
            # 注意这里是+delta而不是=delta*(r-l+1) 因为只表示一个值
            self._tree[rt] += delta
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
            res = self._merge(res, self._query((rt << 1), L, R, l, mid))
        if mid < R:
            res = self._merge(res, self._query((rt << 1) | 1, L, R, mid + 1, r))
        return res

    def _push_up(self, rt: int) -> None:
        # 用子节点更新父节点状态
        self._tree[rt] = self._merge(self._tree[rt << 1], self._tree[rt << 1 | 1])

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._lazy[rt]:
            self._lazy[rt << 1] += self._lazy[rt]
            self._lazy[rt << 1 | 1] += self._lazy[rt]
            # 注意这里是+self._lazy[rt]而不是+self._lazy[rt]*(r-l+1) 因为只表示一个值
            self._tree[rt << 1] += self._lazy[rt]
            self._tree[rt << 1 | 1] += self._lazy[rt]
            self._lazy[rt] = 0


if __name__ == '__main__':
    nums = [1, 3, 5, 7, 9, 11, 13]
    st = SegmentTree(nums, merge=lambda x, y: max(x, y))
    print(st.query(1, 6))
    st.update(1, 6, 2)
    st.update(1, 6, 2)
    print(st.query(1, 1))
    print(st.query(2, 2))
    print(st.query(3, 3))
    print(st.query(1, 3))
    print(st.query(1, 5))
    print(st.query(5, 6))
    print(st.query(7, 7))

