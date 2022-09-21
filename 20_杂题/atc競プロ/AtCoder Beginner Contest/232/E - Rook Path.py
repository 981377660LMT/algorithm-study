# 车每一步可以移动到同一行，或者同一列上的任意一点（`除了自己现在呆的点`）。
# 问从点(x1, y1)到点(x2,y2)用了k步有几种方案(棋盘上车的移动方式)。
# ROW,COL<=1e9
# k<=1e6

# !dp 横向和纵向分开考虑
# 考虑横向走了几步，纵向走了几步

import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


fac = [1]
ifac = [1]
for i in range(1, int(1e6) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def calDp(length: int, step: int) -> List[List[int]]:
    """在长为length的区间上移动了step次后 没有回到起点/回到起点 的方案数"""
    res = [[0, 1]]
    dp1, dp2 = 0, 1
    for _ in range(step):
        dp1, dp2 = (dp2 + (length - 2) * dp1) % MOD, (length - 1) * dp1 % MOD
        res.append([dp1, dp2])
    return res


if __name__ == "__main__":
    ROW, COL, k = map(int, input().split())
    sr, sc, er, ec = map(int, input().split())
    sr, sc, er, ec = sr - 1, sc - 1, er - 1, ec - 1

    rowDp = calDp(ROW, k)
    colDp = calDp(COL, k)
    res = 0
    for i in range(k + 1):
        res = (res + rowDp[i][sr == er] * colDp[k - i][sc == ec] * C(k, i)) % MOD
    print(res)
