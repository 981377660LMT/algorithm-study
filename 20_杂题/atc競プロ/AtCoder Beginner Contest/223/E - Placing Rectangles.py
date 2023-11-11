# 放置矩形
# 给你在直角坐标系中划分出一块长为 (X) ，宽为 (Y) 的矩形区域，
# 现给你三个面积至少为 (A、B、C) 的矩形，
# 问能否在这片矩形区域中不重叠地放置这三个矩形。

# !两种形式:3个平行(左1中2右3) and 左右邻居(左1右2)
# https://zhuanlan.zhihu.com/p/422802337
from itertools import permutations
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def judge1(row: int, col: int, s1: int, s2: int, s3: int) -> bool:
    """从左到右放置三个矩形,左1中2右3"""
    div1, mod1 = divmod(s1, col)
    div2, mod2 = divmod(s2, col)
    div3, mod3 = divmod(s3, col)
    return div1 + div2 + div3 + (mod1 > 0) + (mod2 > 0) + (mod3 > 0) <= row


def judge2(row: int, col: int, s1: int, s2: int, s3: int) -> bool:
    """从左到右放置三个矩形,左1右2"""
    div1, mod1 = divmod(s1, col)
    if div1 + (mod1 > 0) >= row:
        return False
    remainRow = row - div1 - (mod1 > 0)
    div2, mod2 = divmod(s2, remainRow)
    div3, mod3 = divmod(s3, remainRow)
    return div2 + div3 + (mod2 > 0) + (mod3 > 0) <= col


if __name__ == "__main__":
    ROW, COL, A, B, C = map(int, input().split())
    for perm in permutations((A, B, C)):
        if (
            judge1(ROW, COL, *perm)
            or judge2(ROW, COL, *perm)
            or judge1(COL, ROW, *perm)
            or judge2(COL, ROW, *perm)
        ):
            print("Yes")
            exit(0)
    print("No")
