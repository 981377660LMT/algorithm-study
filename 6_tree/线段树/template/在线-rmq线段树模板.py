from typing import List

INF = int(1e20)


class MaxSegmentTree:
    _MIN = -INF  # !注意是0还是-inf

    """RMQ 最大值(区间和可叠加) 线段树
    
    一般用于数组求最值
    注意根节点从1开始,tree本身为[1,n]
    """

    def __init__(self, nOrNums: int | List[int]):
        self._n = nOrNums if isinstance(nOrNums, int) else len(nOrNums)
        self._tree = [self._MIN] * (4 * self._n)
        self._lazy = [0] * (4 * self._n)
        if isinstance(nOrNums, list):
            self._build(1, 1, self._n, nums)

    def query(self, left: int, right: int) -> int:
        """闭区间[left,right]的最大值"""
        assert 1 <= left <= right <= self._n, f"{left},{right} out of range [1,{self._n}]"
        return self._query(1, left, right, 1, self._n)

    def update(self, left: int, right: int, delta: int) -> None:
        """闭区间[left,right]区间更新delta"""
        assert 1 <= left <= right <= self._n, f"{left},{right} out of range [1,{self._n}]"
        self._update(1, left, right, 1, self._n, delta)

    def _build(self, rt: int, l: int, r: int, nums: List[int]) -> None:
        """传了nums时,用于初始化线段树"""
        if l == r:
            self._tree[rt] = nums[l - 1]
            return

        mid = (l + r) // 2
        self._build(rt * 2, l, mid, nums)
        self._build(rt * 2 + 1, mid + 1, r, nums)
        self._push_up(rt)

    def _update(self, rt: int, L: int, R: int, l: int, r: int, delta: int) -> None:
        """L,R表示需要update的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            self._lazy[rt] += delta
            self._tree[rt] += delta
            return

        mid = (l + r) // 2
        self._push_down(rt, l, r, mid)
        if L <= mid:
            self._update(rt * 2, L, R, l, mid, delta)
        if mid < R:
            self._update(rt * 2 + 1, L, R, mid + 1, r, delta)
        self._push_up(rt)

    def _query(self, rt: int, L: int, R: int, l: int, r: int) -> int:
        """L,R表示需要query的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            return self._tree[rt]

        # 传递懒标记
        mid = (l + r) // 2
        self._push_down(rt, l, r, mid)
        res = self._MIN
        if L <= mid:
            cand = self._query(rt * 2, L, R, l, mid)
            if cand > res:
                res = cand
        if mid < R:
            cand = self._query(rt * 2 + 1, L, R, mid + 1, r)
            if cand > res:
                res = cand
        return res

    def _push_up(self, rt: int) -> None:
        if self._tree[rt * 2] > self._tree[rt]:
            self._tree[rt] = self._tree[rt * 2]
        if self._tree[rt * 2 + 1] > self._tree[rt]:
            self._tree[rt] = self._tree[rt * 2 + 1]

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._lazy[rt]:
            self._lazy[rt * 2] += self._lazy[rt]
            self._lazy[rt * 2 + 1] += self._lazy[rt]

            self._tree[rt * 2] += self._lazy[rt]
            self._tree[rt * 2 + 1] += self._lazy[rt]

            self._lazy[rt] = 0


if __name__ == "__main__":
    nums = [-INF] * 5001
    nums[0] = 1
    nums[1000] = 100
    nums[5000] = 2
    st = MaxSegmentTree(nums)
    assert st.query(1000 + 1, 1000 + 1) == 100
    assert st.query(0 + 1, 2 + 1) == 1
    assert st.query(3000 + 1, 5000 + 1) == 2
    print(st.query(2900, 4000))
