"""验证合法的括号序列

将左括号和右括号分别看作是 1 和 -1,那么有效的括号序列条件为
**min(preSum[i]) >= 0 and preSum[-1] == 0**

线段树维护区间前缀和的最小值
"""

from itertools import accumulate
import sys
from typing import List, Union

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class MinSegmentTree2:
    """RMQ 最小值(区间和叠加) 线段树

    基元为0 超出范围返回0
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
            return 0  # !超出范围返回0
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
    n, q = map(int, input().split())
    s = input()  # 括号序列
    nums = [1 if c == "(" else -1 for c in s]
    preSum = MinSegmentTree2(list(accumulate(nums)))  # !线段树维护前缀和

    for _ in range(q):
        kind, left, right = map(int, input().split())  # 1<=left<=right<=n

        if kind == 1:  # !交换s[left]和s[right]
            if nums[left - 1] == nums[right - 1]:
                continue
            if nums[left - 1] == 1 and nums[right - 1] == -1:
                preSum.add(left, right - 1, -2)
            elif nums[left - 1] == -1 and nums[right - 1] == 1:
                preSum.add(left, right - 1, 2)
            nums[left - 1], nums[right - 1] = nums[right - 1], nums[left - 1]
        else:  # !检查s[left:right+1]是否为合法的括号序列
            pre = preSum.query(left - 1, left - 1) if left > 1 else 0
            min_ = preSum.query(left, right) - pre
            last = preSum.query(right, right) - pre
            print("Yes" if (min_ >= 0 and last == 0) else "No")
