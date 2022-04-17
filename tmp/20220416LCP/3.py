from collections import Counter
from typing import List, Optional


# class TreeNode:
#     def __init__(self, x):
#         self.val = x
#         self.left = None
#         self.right = None


class SegmentTree:
    """注意根节点从1开始,tree本身为[1,n]"""

    def __init__(self, n: int):
        self._n = n
        self._tree = [0 for _ in range(n << 2)]
        self._lazy = [0 for _ in range(n << 2)]

    def query(self, l: int, r: int) -> int:
        """[left,right]的和"""
        # assert 1 <= l <= r <= self._n
        return self._query(1, l, r, 1, self._n)

    def update(self, l: int, r: int, delta: int) -> None:
        """[left,right]区间更新delta"""
        # assert 1 <= l <= r <= self._n
        self._update(1, l, r, 1, self._n, delta)

    def _update(self, rt: int, L: int, R: int, l: int, r: int, delta: int) -> None:
        """L,R表示需要update的范围,l,r表示当前节点的范围"""
        if L <= l and r <= R:
            self._lazy[rt] = delta
            self._tree[rt] = r - l + 1 if delta == 1 else 0
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
            self._lazy[rt << 1] = self._lazy[rt]
            self._lazy[rt << 1 | 1] = self._lazy[rt]
            self._tree[rt << 1] = mid - l + 1 if self._lazy[rt] == 1 else 0
            self._tree[rt << 1 | 1] = r - mid if self._lazy[rt] == 1 else 0
            self._lazy[rt] = 0


class Solution:
    def getNumber(self, root: Optional['TreeNode'], ops: List[List[int]]) -> int:
        def dfs(cur: Optional['TreeNode']):
            if not cur:
                return

            values.add(cur.val)
            s.add(cur.val)
            dfs(cur.left)
            dfs(cur.right)

        values = set()
        s = set()
        dfs(root)

        for _, x, y in ops:
            s.add(x)
            s.add(y)

        allNums = sorted(s)
        mapping = {allNums[i]: i + 1 for i in range(len(allNums))}
        # 初始时，所有节点均为蓝色
        # 1 <= 二叉树节点数量 <= 10^5

        sg = SegmentTree(len(mapping) + 10)
        for opt, x, y in ops:
            if opt == 0:
                sg.update(mapping[x], mapping[y], -1)
            else:
                sg.update(mapping[x], mapping[y], 1)

        res = 0
        for v in values:
            res += sg.query(mapping[v], mapping[v])
        return res
