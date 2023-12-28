# F - Christmas Present 2
# https://atcoder.jp/contests/abc334/tasks/abc334_f

from MonoQueue import MonoQueue

from typing import List, Tuple
from math import sqrt

INF = int(1e20)


"""
快递员送货,起始点为(sx, sy),需要到一些房子houses去送货.
送货需要按照顺序送,即先送第一个房子,再送第二个房子,以此类推.
!快递员每次最多携带k个包裹,中途可以回到起点将包裹补充满.
问从起点出发,送完所有房子,回到起点的最短距离是多少.
n,k<=2e5

思路:
类似"从仓库到码头运输箱子"
dp[i] 表示前i个礼物运送完毕时的最短距离 (0<=i<=n)
dp[i] = dp[j] + (preDist[i]-preDist[j+1]) + (distToStart[i]+distToStart[j+1]) | i - j <= k
合并同类项得
d[i] = (dp[j] - preDist[j+1] + distToStart[j+1]) + preDist[i] + distToStart[i] | i - j <= k
!单调队列维护滑动窗口` (dp[j] - preDist[j+1] + distToStart[j+1])`的最小值即可
"""


def christmasPresent2(sx: int, sy: int, houses: List[Tuple[int, int]], k: int) -> float:
    n = len(houses)
    distToStart = [0.0] + [sqrt((x - sx) * (x - sx) + (y - sy) * (y - sy)) for x, y in houses]
    preDist = [0.0, 0.0]  # 运送前i个礼物的相邻移动距离(数组长度为n+1)
    for (x1, y1), (x2, y2) in zip(houses, houses[1:]):
        preDist.append(preDist[-1] + sqrt((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)))

    queue = MonoQueue[Tuple[float, int]](lambda x, y: x[0] < y[0])  # (value,index)
    dp = [INF] * (n + 1)
    dp[0] = 0
    queue.append((dp[0] - preDist[1] + distToStart[1], 0))
    for i in range(1, n + 1):
        while queue and (i - queue.head()[1]) > k:
            queue.popleft()
        preMin = queue.head()[0] if queue else INF
        dp[i] = min(dp[i], preMin + (preDist[i] + distToStart[i]))
        if i < n:
            queue.append((dp[i] - preDist[i + 1] + distToStart[i + 1], i))
    return dp[-1]


if __name__ == "__main__":
    n, k = map(int, input().split())
    sx, sy = map(int, input().split())
    houses = [tuple(map(int, input().split())) for _ in range(n)]
    print(christmasPresent2(sx, sy, houses, k))  # type: ignore
