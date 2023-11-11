# 收集小球的最少时间
# 在一个数轴上给定N个球，每个球在位置Xi上，且有自己的颜色Ci，
# 你从0点出发，每秒走1步，
# !并以颜色单调不降的次序收集所有的球，最后回到原点，要求收集所有球的时间最短。
# n<=2e5

# 非常像锤子那道题 贪心+模拟
# !相同颜色的球怎么决定收集顺序、怎么移动到下一个颜色
# !贪心的尽头是dp => 每种颜色 看最后移动到最左侧还是最右侧

from collections import defaultdict
from typing import List, Tuple

INF = int(1e18)


def traveler(balls: List[Tuple[int, int]]) -> int:
    mp = defaultdict(list)
    for pos, color in balls:
        mp[color].append(pos)

    dp1, dp2 = 0, 0  # 左端点结束、右端点结束 时的最少时间
    pos1, pos2 = 0, 0  # 左端点结束、右端点结束 时的位置
    for color in sorted(mp):
        group = mp[color]
        npos1, npos2 = min(group), max(group)

        # left => left/right=>left/left=>right/right=>right
        ndp1, ndp2 = INF, INF
        ndp1 = min(ndp1, dp1 + abs(pos1 - npos2) + abs(npos2 - npos1))
        ndp1 = min(ndp1, dp2 + abs(pos2 - npos2) + abs(npos2 - npos1))
        ndp2 = min(ndp2, dp1 + abs(pos1 - npos1) + abs(npos1 - npos2))
        ndp2 = min(ndp2, dp2 + abs(pos2 - npos1) + abs(npos1 - npos2))

        dp1, dp2 = ndp1, ndp2
        pos1, pos2 = npos1, npos2

    return min(dp1 + abs(pos1), dp2 + abs(pos2))


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    balls = [tuple(map(int, input().split())) for _ in range(int(input()))]
    print(traveler(balls))
