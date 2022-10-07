"""
最初,你在数轴的0点。
每一秒,你可以向左或向右移动1格,或者不动。
在Ti秒的时候,你会受到一次伤害,规则如下:
设你当前的位置为p,若Di=0,你受到max(p - Xi,0)点伤害,否则受到max(Xi- p,0)点伤害。
伤害共有n次(n≤2e5),求你受到的最小总伤害。(Ti递增)

https://atcoder.jp/contests/abc217/editorial/2581
!dp[i][x] 表示ti秒后,你的位置为x时的最小伤害
!dp[i][x] = min(dp[i-1][y]) + cost , x-(ti+1-ti)<=y<=x+(ti+1-ti)
!每次受到伤害时将正负斜率部分往相离方向平移，然后叠加上新的伤害。
"""


import sys
from SlopeTrick import SlopeTrick

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    S = SlopeTrick(leftTuring=[0] * (n + 10), rightTuring=[0] * (n + 10))

    time = 0
    for _ in range(n):
        curTime, D, a = map(int, input().split())
        S.addLeftOffset(-(curTime - time))  # !将正负斜率部分往相离方向平移
        S.addRightOffset(curTime - time)  # !将正负斜率部分往相离方向平移
        if D == 0:
            S.addAMinusX(a)
        else:
            S.addXMinusA(a)
        time = curTime

    print(S.getMinY())
