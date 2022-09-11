# hoursBefore ，表示你要前往会议所剩下的可用小时数
# 必须休息并等待，直到 下一个整数小时 才能开始继续通过下一条道路
# 然而，为了能准时到达，你可以选择 跳过 一些路的休息时间，这意味着你不必等待下一个整数小时。
# 返回准时抵达会议现场所需要的 最小跳过次数 ，如果 无法准时参会 ，返回 -1 。
# !1 <= n <= 1000


# summary
# jump or not jump at each point


# 在进行「向上取整」运算前，我们将待取整的浮点数减去 eps 再进行取整，
# !就可以避免浮点数误差导致的ceil大1的问题
# ceil(8.0 + 1.0 / 3 + 1.0 / 3 + 1.0 / 3) 应当是 9，而计算机会给出 10
# 这是因为浮点数误差导致8.0 + 1.0 / 3 + 1.0 / 3 + 1.0 / 3
# 计算出的结果约为：
# 9.000000000000002
# 本题speed最大为1e6 因此EPS取1e-8/1e-9都可以


from typing import List
from math import ceil

INF = int(4e18)
EPS = 1e-8


class Solution:
    def minSkips(self, dist: List[int], speed: int, hoursBefore: int) -> int:
        n = len(dist)
        dp = [INF] * (n + 1)  # 前i个位置跳了j次时花费的最少时间
        dp[0] = 0
        for i in range(1, n + 1):
            ndp = [INF] * (n + 1)
            cost = dist[i - 1] / speed
            for pre in range(i):
                # 不跳
                ndp[pre] = min(ndp[pre], ceil(dp[pre] + cost - EPS))
                # 跳
                ndp[pre + 1] = min(ndp[pre + 1], dp[pre] + cost)  # type: ignore
            dp = ndp

        return next((i for i in range(n + 1) if dp[i] <= hoursBefore), -1)


print(Solution().minSkips(dist=[1, 3, 2], speed=4, hoursBefore=2))
# 输出：1
# 解释：
# 不跳过任何休息时间，你将用 (1/4 + 3/4) + (3/4 + 1/4) + (2/4) = 2.5 小时才能抵达会议现场。
# 可以跳过第 1 次休息时间，共用 ((1/4 + 0) + (3/4 + 0)) + (2/4) = 1.5 小时抵达会议现场。
# 注意，第 2 次休息时间缩短为 0 ，由于跳过第 1 次休息时间，你是在整数小时处完成通过第 2 条道路。


# 浮点数运算的细节
# 本题中我们不可避免的会使用「浮点数运算」以及「向上取整」运算，如果强行忽略产生的计算误差，会得到错误的结果。
