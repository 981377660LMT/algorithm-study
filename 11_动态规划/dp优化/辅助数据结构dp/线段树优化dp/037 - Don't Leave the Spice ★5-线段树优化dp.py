# 每个物品的重量为一个范围 [lefti,righti] 可随意调节
# 每个物品价值为vi
# 现在选任意个物品 总重量为w 求最大价值
# 如果不可能完成 返回 -1

# n<=500
# 1<=L<=R<=1e4

# !dp[i][weight] 表示前index个 重量为 weight 时的最大价值
# !不选: dp[i][weight] = dp[i-1][weight]
# !选: dp[i][weight] = max(dp[i-1][weight-R],..., dp[i-1][weight-L]) +vi
# !这里需要RMQ 求区间最大值 由于需要动态更新 所以用线段树或者单调队列维护最值
from typing import List
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)
INF = int(1e20)


class MaxSegmentTree:
    _MIN = -INF  # !注意是0还是-inf

    """RMQ 最大值(区间和可叠加) 线段树
    
    一般用于数组求最值
    注意根节点从1开始,tree本身为[1,n]
    """

    def __init__(self, n: int, nums: List[int]):
        self._n = n
        self._tree = [self._MIN] * (4 * n)
        self._lazy = [0] * (4 * n)
        self._build(1, 1, n, nums)

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


volumn, n = map(int, input().split())
goods = []
for _ in range(n):
    left, right, score = map(int, input().split())
    goods.append((left, right, score))

dp = [-INF] * (volumn + 1)
dp[0] = 0
for i in range(goods[0][0], goods[0][1] + 1):
    dp[i] = goods[0][2]

for i in range(1, n):
    ndp = dp[:]  # 不选
    lower, upper, score = goods[i]
    tree = MaxSegmentTree(len(ndp), ndp)
    for j in range(lower, volumn + 1):  # 容量
        left, right = max(0, j - upper), j - lower
        preMax = tree.query(left + 1, right + 1)
        ndp[j] = max(ndp[j], preMax + score)

    dp = ndp

print(-1 if dp[volumn] < 0 else dp[volumn])
