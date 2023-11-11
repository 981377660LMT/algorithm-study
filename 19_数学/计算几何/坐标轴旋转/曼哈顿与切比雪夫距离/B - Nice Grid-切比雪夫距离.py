"""
15*15的中心放射型矩阵 求(row,col)是黑还是白

求出切比雪夫距离 奇数则为白 偶数则为黑
"""


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

row, col = map(int, input().split())
cr, cc = 7, 7  # 中心点
row, col = row - 1, col - 1  # 左上角为原点的坐标
row, col = -(row - cr), col - cc  # 中心为原点的坐标

# !到中心的切比雪夫距离
manhattan = max(abs(row - 0), abs(col - 0))
print("black" if manhattan & 1 else "white")
# print("变换后的坐标:", row, col)
