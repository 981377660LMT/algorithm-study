# 给定长度为L的绳子，绳子上有n-1个点，
# 分别在1,2…L-1位置。进行q次操作。
# !1:在x坐标处，将x所在区间一分为二
# !2:询问x所在那一截的绳子有多长。
# L<=1e9
# q<=2e5


import sys
from sortedcontainers import SortedList


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    L, Q = map(int, input().split())
    sl = SortedList([0, L])  # 区间分割点

    for _ in range(Q):
        c, x = map(int, input().split())
        if c == 1:
            sl.add(x)
        else:
            pos = sl.bisect_right(x) - 1
            print(sl[pos + 1] - sl[pos])
