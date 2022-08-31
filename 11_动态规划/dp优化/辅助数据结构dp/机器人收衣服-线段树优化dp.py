# 美团笔试t4 二分+dp
# 机器人收衣服
# 小团正忙看用机器人收衣服!因为快要下雨了，小团找来了不少机器人帮忙收衣服。
# 他有n件衣服从左到右成一行排列，所在位置分别为1~n,
# 在每个位置上已经有一个就绪的机器人可以帮忙收衣服，但第i个位置上的机器人需要pi的电量来启动。
# 然后这个机器人会用ti的时间收衣服，当它收完当前衣服后，
# 会尝试去收紧邻的右边的一件衣服(如果存在的话)，即i+1处的衣服，
# 如果i+1处的衣服已经被其他机器人收了或者其他机器人正在收，这个机器人就会进入休眠状态，不再收衣服。
# 不过如果机器人没有休眠，它会同样以ti时间来收这件1+1处的衣服(注意，不是t+1的时间，收衣服的时间为每个机器人固有属性)，
# 然后它会做同样的检测来看能否继续收i+2处的衣服，一直直到它进入休眠状态或者右边没有衣服可以收了。
# 形象地来说，机器人会一直尝试往右边收衣服，收k件的话就耗费k*ti;的时间,
# 但是当它遇见其他机器人工作的痕迹，就会认为后面的事情它不用管了,开始摸鱼，进入休眠状态。
# !小团手里总共有电量b，他准备在0时刻的时候将所有他想启动的机器人全部一起(并行)启动，
# 过后不再启动新的机器人，并且启动的机器人的电量之和不大于b。
# 他想知道在最佳选择的情况下，最快多久能收完衣服。若无论如何怎样都收不完衣服,输出-1.
# !n<=1000 pi<=100 ti,b<=1e5

# !二分时间+dp O(n^2log1e8)
# !为什么想到二分时间:因为需要固定一个变量(题目有两个变量:时间和电量)
# 不同机器人之间的时间是并行的，不便于DP，因此固定电量求最小时间很困难
# 反过来，二分时间，DP最小电量是否符合要求，就容易的多
# !首先二分时间mid，然后dp, dp[i]表示在不超过mid的时间内收完前i个所需要的最少启动电量

# !第一种dp方式(貰うDP,O(n^2))
# 求解dp[i]时，往前枚举j，如果times[j+1]*(i -j) <= mid，
# 则 dp[i] = min(dp[i], dp[j]+ powers[j+1])
# 然后看最后的dp[n]<= b
# !第二种dp方式(配るDP,O(nlogn))
# 求解dp时，用i去更新区段i到 i + ((mid / times[j])-1) 的最小值


from typing import List, Union

INF = int(4e18)


class MinSegmentTree:
    """RMQ 最小值(区间和覆盖) 线段树

    一般用于数组求最值
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
        # assert 1 <= left <= right <= self._n, f"{left},{right} out of range [1,{self._n}]"
        return self._query(1, left, right, 1, self._n)

    def update(self, left: int, right: int, target: int) -> None:
        """闭区间[left,right]区间更新为target"""
        # assert 1 <= left <= right <= self._n, f"{left},{right} out of range [1,{self._n}]"
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
        res = INF
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
        if self._tree[rt * 2] < self._tree[rt]:
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


def assignRobots(n: int, b: int, powers: List[int], times: List[int]) -> int:
    """美团笔试t4 二分时间+dp

    Args:
        n (int): n件衣服从左到右成一行排列
        b (int): 启动的机器人的电量之和不大于b
        powers (List[int]): 启动每个机器人的花费
        times (List[int]): 每个机器人收一件衣服的时间

    Returns:
        int: 在最佳选择的情况下，最快多久能收完衣服;若无论如何怎样都收不完衣服,输出-1
    """

    def check1(mid: int) -> bool:
        """O(n^2) 能否在不超过mid的时间内收完"""
        dp = [INF] * (n + 1)  # dp[i]表示在不超过mid的时间内收完前i个所需要的最少启动电量
        dp[0] = 0
        for i in range(1, n + 1):
            for j in range(i):  #
                if times[j] * (i - j) <= mid:
                    dp[i] = min(dp[i], dp[j] + powers[j])
        return dp[-1] <= b

    def check2(mid: int) -> bool:
        """O(nlogn) 能否在不超过mid的时间内收完"""
        # !线段树优化 反过来用i去更新区段 i 到 i + (mid // times[j] - 1) 的最小值
        dp = MinSegmentTree(n + 1)
        dp.update(1, 1, 0)
        for i in range(1, n + 1):
            left = i + 1
            preMin = dp.query(left - 1, left - 1)  # !注意preMin是查前一个位置
            right = i + (mid // times[i - 1] - 1) + 1
            right = min(right, n + 1)
            if right >= left:
                dp.update(left, right, preMin + powers[i - 1])
        return dp.query(n + 1, n + 1) <= b

    left, right = 0, int(1e8 + 10)
    while left <= right:
        mid = (left + right) // 2
        if check2(mid):
            right = mid - 1
        else:
            left = mid + 1

    return left if check2(left) else -1


if __name__ == "__main__":
    from random import randint

    def bruteForce(n: int, b: int, powers: List[int], times: List[int]) -> int:
        res = INF
        for state in range(1 << n):
            select = []
            for i in range(n):
                if state & (1 << i):
                    select.append(i)
            if 0 in select and sum(powers[i] for i in select) <= b:
                cand = 0
                select.append(n)
                for pre, cur in zip(select, select[1:]):
                    cand = max(cand, times[pre] * (cur - pre))
                res = min(res, cand)
        return res if res != INF else -1

    for _ in range(100):
        n = randint(1, 10)
        b = randint(5, 10)
        powers = [randint(1, 10) for _ in range(n)]
        times = [randint(1, 10) for _ in range(n)]
        res1, res2 = assignRobots(n, b, powers, times), bruteForce(n, b, powers, times)
        if res1 != res2:
            print(n, b, powers, times)
            print(res1, res2)
            break
