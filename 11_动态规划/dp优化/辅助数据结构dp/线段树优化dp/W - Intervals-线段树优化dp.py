"""线段树优化dp"""
# 给定一个长为n的01序列和m个区间 [lefti,righti,scorei]
# 如果第i个区间里存在1 那么分数就加上score[i]
# !求构造出的01序列的最大分数
# n,m<=2e5
# 1<=left<=right<=n
# -1e9<=score<=1e9

# 思路:
# dp[i] 表示在看第i个字符且第i个字符为1、右边全部为0时的最大分数
# !1. 所有区间按照右端点分类 因为只要右端点存在1 那么这会获得这一堆区间的分数
# !2. 线段树维护区间最大值

from collections import defaultdict
import sys
from typing import List, Union

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


class MaxSegmentTree1:
    """RMQ 最大值(区间和可叠加) 线段树

    一般用于数组求最值
    注意根节点从1开始,tree本身为[1,n]
    """

    def __init__(self, nOrNums: Union[int, List[int]]):
        self._n = nOrNums if isinstance(nOrNums, int) else len(nOrNums)
        self._tree = [0] * (4 * self._n)
        self._lazy = [0] * (4 * self._n)
        if isinstance(nOrNums, list):
            self._build(1, 1, self._n, nOrNums)

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left,right]区间值加上delta"""
        # assert 1 <= left <= right <= self._n, f"{left},{right} out of range [1,{self._n}]"
        self._add(1, left, right, 1, self._n, delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left,right]的最值"""
        # assert 1 <= left <= right <= self._n, f"{left},{right} out of range [1,{self._n}]"
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
        """L,R表示需要update的范围,l,r表示当前节点的范围"""
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
        res = -INF
        if L <= mid:
            res = max(res, self._query(rt * 2, L, R, l, mid))
        if mid < R:
            res = max(res, self._query(rt * 2 + 1, L, R, mid + 1, r))
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


############################################################
n, m = map(int, input().split())
intervals = defaultdict(list)  # 区间右端点 => 区间左端点

for _ in range(m):
    left, right, score = map(int, input().split())
    intervals[right].append((left, score))


dp = MaxSegmentTree1(n)

for right in range(1, n + 1):
    preMax = dp.queryAll()
    dp.add(right, right, preMax)  # 初始化dp[i](最大值)
    for left, score in intervals[right]:
        dp.add(left, right, score)  # 区间叠加更新最大值

print(max(0, dp.queryAll()))
