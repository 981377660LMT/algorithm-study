from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 座標平面上に
# N 個の点
# P
# 1
# ​
#  ,P
# 2
# ​
#  ,…,P
# N
# ​
#   があり、点
# P
# i
# ​
#   の座標は
# (X
# i
# ​
#  ,Y
# i
# ​
#  ) です。
# 2 つの点
# A,B の距離
# dist(A,B) を次のように定義します。

# 最初、ウサギが点
# A にいる。
# (x,y) にいるウサギは
# (x+1,y+1),
# (x+1,y−1),
# (x−1,y+1),
# (x−1,y−1) のいずれかに
# 1 回のジャンプで移動することができる。
# 点
# A から点
# B まで移動するために必要なジャンプの回数の最小値を
# dist(A,B) として定義する。
# ただし、何度ジャンプを繰り返しても点
# A から点
# B まで移動できないとき、
# dist(A,B)=0 とする。

# i=1
# ∑
# N−1
# ​

# j=i+1
# ∑
# N
# ​
#  dist(P
# i
# ​
#  ,P
# j
# ​
#  ) を求めてください。

# 每次移动奇偶不变，所以考虑前缀中奇偶相同的点即可
# !（x1,y1)与(x2,y2)距离为 max(abs(x1-x2),abs(y1-y2))
if __name__ == "__main__":
    N = int(input())
    points = []
    for _ in range(N):
        x, y = map(int, input().split())
        points.append((x, y))
    res = 0
    odd, even = [], []
    for x, y in points:
        sum_ = x + y
        if sum_ & 1 == 0:
            res += sum(max(abs(x - x_), abs(y - y_)) for x_, y_ in even)
            even.append((x, y))
        else:
            res += sum(max(abs(x - x_), abs(y - y_)) for x_, y_ in odd)
            odd.append((x, y))
    print(res)
