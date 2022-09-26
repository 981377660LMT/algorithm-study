from typing import List, Union

INF = int(4e18)


class SegmentTree:
    """
    查询区间和  染色(覆盖)更新
    注意根节点从1开始,tree本身为[1,n]
    """

    __slots__ = ("_n", "_tree", "_lazyValue", "_isLazy")

    def __init__(self, nOrNums: Union[int, List[int]]):
        self._n = nOrNums if isinstance(nOrNums, int) else len(nOrNums)
        self._tree = [0] * (4 * self._n)
        self._lazyValue = [0] * (4 * self._n)
        self._isLazy = [False] * (4 * self._n)
        if isinstance(nOrNums, list):
            self._build(1, 1, self._n, nOrNums)

    def queryAll(self) -> int:
        return self._tree[1]

    def query(self, left: int, right: int) -> int:
        """闭区间[left,right]的和"""
        if left < 1:
            left = 1
        if right > self._n:
            right = self._n
        if left > right:
            return 0  # !超出范围返回0
        return self._query(1, left, right, 1, self._n)

    def update(self, left: int, right: int, target: int) -> None:
        """闭区间[left,right]区间更新为target"""
        if left < 1:
            left = 1
        if right > self._n:
            right = self._n
        if left > right:
            return
        self._update(1, left, right, 1, self._n, target)

    def _build(self, rt: int, l: int, r: int, nums: List[int]) -> None:
        """传了nums时,用于初始化线段树"""
        if l == r:
            self._tree[rt] = nums[l - 1]
            return

        mid = (l + r) // 2
        self._build(rt * 2, l, mid, nums)
        self._build(rt * 2 + 1, mid + 1, r, nums)
        self._push_up(rt)

    def _update(self, rt: int, L: int, R: int, l: int, r: int, target: int) -> None:
        """L,R表示需要update的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            self._isLazy[rt] = True
            self._lazyValue[rt] = target
            self._tree[rt] = target * (r - l + 1)
            return

        # 传递懒标记
        mid = (l + r) // 2
        self._push_down(rt, l, r, mid)
        if L <= mid:
            self._update(rt * 2, L, R, l, mid, target)
        if mid < R:
            self._update(rt * 2 + 1, L, R, mid + 1, r, target)
        self._push_up(rt)

    def _query(self, rt: int, L: int, R: int, l: int, r: int) -> int:
        """L,R表示需要query的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            return self._tree[rt]

        # 传递懒标记
        mid = (l + r) // 2
        self._push_down(rt, l, r, mid)
        res = 0
        if L <= mid:
            res += self._query(rt * 2, L, R, l, mid)
        if mid < R:
            res += self._query(rt * 2 + 1, L, R, mid + 1, r)
        return res

    def _push_up(self, rt: int) -> None:
        self._tree[rt] = self._tree[rt * 2] + self._tree[rt * 2 + 1]

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._isLazy[rt]:
            target = self._lazyValue[rt]
            self._lazyValue[rt * 2] = target
            self._lazyValue[rt * 2 + 1] = target
            self._tree[rt * 2] = target * (mid - l + 1)
            self._tree[rt * 2 + 1] = target * (r - mid)
            self._isLazy[rt * 2] = True
            self._isLazy[rt * 2 + 1] = True

            self._lazyValue[rt] = 0
            self._isLazy[rt] = False


if __name__ == "__main__":

    class Solution:
        def amountPainted(self, paint: List[List[int]]) -> List[int]:
            max_ = max([r for _, r in paint])
            tree = SegmentTree(max_ + 1)
            res = []
            for start, end in paint:
                start, end = start + 1, end + 1
                rangeSum = tree.query(start, end - 1)
                res.append(end - start - rangeSum)
                tree.update(start, end - 1, 1)

            return res

    assert Solution().amountPainted([[1, 2], [2, 3], [3, 4]]) == [1, 1, 1]
