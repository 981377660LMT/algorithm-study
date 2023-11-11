"""
最大不相交区间数

每个7有两个斜率 对应一个区间
即给定n个区间 找到最多的两两不相交的区间个数(右端点贪心排序)

Decimal 避免浮点数误差
"""

from decimal import Decimal
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    # !Decimal表示int 避免float运算的浮点数误差
    points = [tuple(map(Decimal, input().split())) for _ in range(n)]

    intervals = []
    for x, y in points:
        slope1 = (y - 1) / x if x != 0 else INF
        slope2 = y / (x - 1) if x != 1 else INF
        intervals.append((slope1, slope2))

    # 最长不相交区间数
    intervals.sort(key=lambda x: x[1])
    res, preEnd = 0, -INF
    for start, end in intervals:
        if start >= preEnd:
            res += 1
            preEnd = end

    print(res)
