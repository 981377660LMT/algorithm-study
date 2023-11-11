# 有A个神社和B个寺庙在一条直线上已知位置，Q次查询，
# 每次告诉一个位置，从那出发，要访问神社和寺庙至少各一个，求最小距离。

# 二分
from bisect import bisect_left, bisect_right
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(1e18)

if __name__ == "__main__":
    a, b, q = map(int, input().split())
    A = [int(input()) for _ in range(a)]
    B = [int(input()) for _ in range(b)]
    A.sort()
    B.sort()

    for _ in range(q):
        x = int(input())

        # 左边最近的神社,右边最近的神社
        left1 = bisect_right(A, x) - 1
        right1 = bisect_left(A, x)
        leftPos1 = A[left1] if left1 >= 0 else -INF
        rightPos1 = A[right1] if right1 < a else INF

        # 左边最近的寺庙,右边最近的寺庙
        left2 = bisect_right(B, x) - 1
        right2 = bisect_left(B, x)
        leftPos2 = B[left2] if left2 >= 0 else -INF
        rightPos2 = B[right2] if right2 < b else INF

        res = INF
        for left in [leftPos1, rightPos1]:
            for right in [leftPos2, rightPos2]:
                res = min(
                    res, abs(x - left) + abs(left - right), abs(x - right) + abs(left - right)
                )
        print(res)
