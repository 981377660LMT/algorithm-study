from typing import List, Union

INF = int(4e18)


class MaxSegmentTree1:
    """RMQ 最大值(区间和覆盖) 线段树

    基元为-INF 超出范围返回-INF
    注意根节点从1开始,tree本身为[1,n]
    """

    __slots__ = ("_n", "_tree", "_lazyValue", "_isLazy")

    def __init__(self, nOrNums: Union[int, List[int]]):
        self._n = nOrNums if isinstance(nOrNums, int) else len(nOrNums)
        self._tree = [-INF] * (4 * self._n)
        self._lazyValue = [-INF] * (4 * self._n)
        self._isLazy = [False] * (4 * self._n)
        if isinstance(nOrNums, list):
            self._build(1, 1, self._n, nOrNums)

    def queryAll(self) -> int:
        return self._tree[1]

    def query(self, left: int, right: int) -> int:
        """闭区间[left,right]的最大值"""
        if left < 1:
            left = 1
        if right > self._n:
            right = self._n
        if left > right:
            return -INF  # !超出范围返回-INF
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
            self._lazyValue[rt] = target if target > self._lazyValue[rt] else self._lazyValue[rt]
            self._tree[rt] = target if target > self._tree[rt] else self._tree[rt]
            self._isLazy[rt] = True
            return

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
        res = -INF  # !默认值为-INF
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
        self._tree[rt] = self._tree[rt * 2]
        if self._tree[rt * 2 + 1] > self._tree[rt]:
            self._tree[rt] = self._tree[rt * 2 + 1]

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._isLazy[rt]:
            target = self._lazyValue[rt]

            self._lazyValue[rt * 2] = (
                target if target > self._lazyValue[rt * 2] else self._lazyValue[rt * 2]
            )
            self._lazyValue[rt * 2 + 1] = (
                target if target > self._lazyValue[rt * 2 + 1] else self._lazyValue[rt * 2 + 1]
            )

            self._tree[rt * 2] = target if target > self._tree[rt * 2] else self._tree[rt * 2]
            self._tree[rt * 2 + 1] = (
                target if target > self._tree[rt * 2 + 1] else self._tree[rt * 2 + 1]
            )

            self._isLazy[rt * 2] = True
            self._isLazy[rt * 2 + 1] = True

            self._lazyValue[rt] = -INF
            self._isLazy[rt] = False


class MaxSegmentTree2:
    """RMQ 最大值线段树(区间和叠加)

    基元为0 超出范围返回-INF
    注意根节点从1开始,tree本身为[1,n]
    因为是叠加 所以不需要isLazy数组
    """

    __slots__ = "_n", "_tree", "_lazy"

    def __init__(self, nOrNums: Union[int, List[int]]):
        self._n = nOrNums if isinstance(nOrNums, int) else len(nOrNums)
        self._tree = [0] * (4 * self._n)
        self._lazy = [0] * (4 * self._n)
        if isinstance(nOrNums, list):
            self._build(1, 1, self._n, nOrNums)

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left,right]区间值加上delta"""
        if left < 1:
            left = 1
        if right > self._n:
            right = self._n
        if left > right:
            return
        self._add(1, left, right, 1, self._n, delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left,right]的最值"""
        if left < 1:
            left = 1
        if right > self._n:
            right = self._n
        if left > right:
            return -INF  # !超出范围返回-INF
        return self._query(1, left, right, 1, self._n)

    def queryAll(self) -> int:
        return self._tree[1]

    def _build(self, rt: int, l: int, r: int, nums: List[int]) -> None:
        """传了nums时,用于初始化线段树"""
        if l == r:
            self._tree[rt] = nums[l - 1]
            return

        mid = (l + r) // 2
        self._build(rt * 2, l, mid, nums)
        self._build(rt * 2 + 1, mid + 1, r, nums)
        self._push_up(rt)

    def _add(self, rt: int, L: int, R: int, l: int, r: int, delta: int) -> None:
        """L,R表示需要add的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            self._lazy[rt] += delta
            self._tree[rt] += delta
            return

        mid = (l + r) // 2
        self._push_down(rt, l, r, mid)
        if L <= mid:
            self._add(rt * 2, L, R, l, mid, delta)
        if mid < R:
            self._add(rt * 2 + 1, L, R, mid + 1, r, delta)
        self._push_up(rt)

    def _query(self, rt: int, L: int, R: int, l: int, r: int) -> int:
        """L,R表示需要query的范围,l,r表示当前节点的范围"""
        # 传递懒标记
        if L <= l and r <= R:
            return self._tree[rt]

        mid = (l + r) // 2
        self._push_down(rt, l, r, mid)
        res = 0  # !默认值为0
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
        self._tree[rt] = self._tree[rt * 2]
        if self._tree[rt * 2 + 1] > self._tree[rt]:
            self._tree[rt] = self._tree[rt * 2 + 1]

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._lazy[rt]:
            value = self._lazy[rt]
            self._lazy[rt * 2] += value
            self._lazy[rt * 2 + 1] += value

            self._tree[rt * 2] += value
            self._tree[rt * 2 + 1] += value

            self._lazy[rt] = 0


class MinSegmentTree1:
    """RMQ 最小值(区间和覆盖) 线段树

    基元为INF 超出范围返回INF
    注意根节点从1开始,tree本身为[1,n]
    """

    __slots__ = ("_n", "_tree", "_lazyValue", "_isLazy")

    def __init__(self, nOrNums: Union[int, List[int]]):
        self._n = nOrNums if isinstance(nOrNums, int) else len(nOrNums)
        self._tree = [INF] * (4 * self._n)
        self._lazyValue = [INF] * (4 * self._n)
        self._isLazy = [False] * (4 * self._n)
        if isinstance(nOrNums, list):
            self._build(1, 1, self._n, nOrNums)

    def queryAll(self) -> int:
        return self._tree[1]

    def query(self, left: int, right: int) -> int:
        """闭区间[left,right]的最小值"""
        if left < 1:
            left = 1
        if right > self._n:
            right = self._n
        if left > right:
            return INF  # !超出范围返回INF
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
            self._lazyValue[rt] = target if target < self._lazyValue[rt] else self._lazyValue[rt]
            self._tree[rt] = target if target < self._tree[rt] else self._tree[rt]
            self._isLazy[rt] = True
            return

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
        res = INF  # !默认值为INF
        if L <= mid:
            cand = self._query(rt * 2, L, R, l, mid)
            if cand < res:
                res = cand
        if mid < R:
            cand = self._query(rt * 2 + 1, L, R, mid + 1, r)
            if cand < res:
                res = cand
        return res

    def _push_up(self, rt: int) -> None:
        self._tree[rt] = self._tree[rt * 2]
        if self._tree[rt * 2 + 1] < self._tree[rt]:
            self._tree[rt] = self._tree[rt * 2 + 1]

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._isLazy[rt]:
            target = self._lazyValue[rt]

            self._lazyValue[rt * 2] = (
                target if target < self._lazyValue[rt * 2] else self._lazyValue[rt * 2]
            )
            self._lazyValue[rt * 2 + 1] = (
                target if target < self._lazyValue[rt * 2 + 1] else self._lazyValue[rt * 2 + 1]
            )

            self._tree[rt * 2] = target if target < self._tree[rt * 2] else self._tree[rt * 2]
            self._tree[rt * 2 + 1] = (
                target if target < self._tree[rt * 2 + 1] else self._tree[rt * 2 + 1]
            )

            self._isLazy[rt * 2] = True
            self._isLazy[rt * 2 + 1] = True

            self._lazyValue[rt] = INF
            self._isLazy[rt] = False


class MinSegmentTree2:
    """RMQ 最小值(区间和叠加) 线段树

    基元为0 超出范围返回INF
    注意根节点从1开始,tree本身为[1,n]
    因为是叠加 所以不需要isLazy数组
    """

    __slots__ = "_n", "_tree", "_lazy"

    def __init__(self, nOrNums: Union[int, List[int]]):
        self._n = nOrNums if isinstance(nOrNums, int) else len(nOrNums)
        self._tree = [0] * (4 * self._n)
        self._lazy = [0] * (4 * self._n)
        if isinstance(nOrNums, list):
            self._build(1, 1, self._n, nOrNums)

    def queryAll(self) -> int:
        return self._tree[1]

    def query(self, left: int, right: int) -> int:
        """闭区间[left,right]的最小值"""
        if left < 1:
            left = 1
        if right > self._n:
            right = self._n
        if left > right:
            return INF  # !超出范围返回INF
        return self._query(1, left, right, 1, self._n)

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left,right]区间值加上delta"""
        if left < 1:
            left = 1
        if right > self._n:
            right = self._n
        if left > right:
            return
        self._add(1, left, right, 1, self._n, delta)

    def _build(self, rt: int, l: int, r: int, nums: List[int]) -> None:
        """传了nums时,用于初始化线段树"""
        if l == r:
            self._tree[rt] = nums[l - 1]
            return

        mid = (l + r) // 2
        self._build(rt * 2, l, mid, nums)
        self._build(rt * 2 + 1, mid + 1, r, nums)
        self._push_up(rt)

    def _add(self, rt: int, L: int, R: int, l: int, r: int, delta: int) -> None:
        """L,R表示需要add的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            self._lazy[rt] += delta
            self._tree[rt] += delta
            return

        mid = (l + r) // 2
        self._push_down(rt, l, r, mid)
        if L <= mid:
            self._add(rt * 2, L, R, l, mid, delta)
        if mid < R:
            self._add(rt * 2 + 1, L, R, mid + 1, r, delta)
        self._push_up(rt)

    def _query(self, rt: int, L: int, R: int, l: int, r: int) -> int:
        """L,R表示需要query的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            return self._tree[rt]

        # 传递懒标记
        mid = (l + r) // 2
        self._push_down(rt, l, r, mid)
        res = INF  # !默认值为INF
        if L <= mid:
            cand = self._query(rt * 2, L, R, l, mid)
            if cand < res:
                res = cand
        if mid < R:
            cand = self._query(rt * 2 + 1, L, R, mid + 1, r)
            if cand < res:
                res = cand
        return res

    def _push_up(self, rt: int) -> None:
        self._tree[rt] = self._tree[rt * 2]
        if self._tree[rt * 2 + 1] < self._tree[rt]:
            self._tree[rt] = self._tree[rt * 2 + 1]

    def _push_down(self, rt: int, l: int, r: int, mid: int) -> None:
        if self._lazy[rt]:
            value = self._lazy[rt]
            self._lazy[rt * 2] += value
            self._lazy[rt * 2 + 1] += value

            self._tree[rt * 2] += value
            self._tree[rt * 2 + 1] += value

            self._lazy[rt] = 0


if __name__ == "__main__":

    # !1.`染色`最大值
    class Solution:
        def fallingSquares(self, positions: List[List[int]]) -> List[int]:
            pos = set()
            for left, length in positions:
                pos.add(left)
                pos.add(left + length - 1)

            mapping = {v: i for i, v in enumerate(sorted(pos), 1)}
            tree = MaxSegmentTree1(len(mapping) + 10)
            res = []
            for left, length in positions:
                left, right = mapping[left], mapping[left + length - 1]
                preMax = max(0, tree.query(left, right))  # !注意这里的max
                tree.update(left, right, preMax + length)
                res.append(tree.queryAll())
            return res

    # !2.`叠加`最大值
    tree = MaxSegmentTree2(10)
    tree.add(1, 3, 1)
    assert tree.queryAll() == 1
    tree.add(2, 4, 1)
    assert tree.queryAll() == 2
    assert tree.query(1, 4) == 2

    # !3.`染色`最小值
    tree = MinSegmentTree1(list(range(10)))
    assert tree.queryAll() == 0
    tree.update(1, 1, -1)
    assert tree.queryAll() == -1

    # !4.`叠加`最小值
    tree = MinSegmentTree2(10)
    assert tree.queryAll() == 0
    tree.add(1, 10, 1)
    assert tree.queryAll() == 1
    tree.add(1, 10, 1)
    assert tree.queryAll() == 2
    tree.add(1, 1, 1)
    assert tree.queryAll() == 2
    assert tree.query(1, 1) == 3
