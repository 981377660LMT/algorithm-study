# 平面 x=0 上に、高さ H の T 分で一周する観覧車があります。
# 観覧車の名物である「高橋直大像」は、座標 (X,Y,0) に存在します。
# !求缆车到雕像的俯角

from math import cos, degrees, radians, sin, atan  # 转角度、弧度
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

T = int(input())
H, X, Y = map(int, input().split())
Q = int(input())

for _ in range(Q):
    t = int(input())
    radius = H / 2
    rad = radians(360 * t / T)
    y = -radius * sin(rad)
    x = radius * (1 - cos(rad))
    tan_ = x / ((y - Y) ** 2 + X ** 2) ** 0.5
    print(degrees(atan(tan_)))

