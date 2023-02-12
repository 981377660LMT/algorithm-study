import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# xy 平面上に高橋くんがいます。 はじめ、高橋くんは点
# (s
# x
# ​
#  ,s
# y
# ​
#  ) にいます。 高橋くんは、点
# (t
# x
# ​
#  ,t
# y
# ​
#  ) に移動したいです。

# xy 平面上に、長方形
# R:={(x,y)∣a−0.5≤x≤b+0.5,c−0.5≤y≤d+0.5} があります。 次の操作を考えます。

# 長方形
# R に含まれる格子点
# (x,y) をひとつ選ぶ。 点
# (x,y) を中心に高橋くんはいまいる位置と対称な位置に瞬間移動する。
# 上の操作を
# 0 回以上
# 10
# 6
#   回以下繰り返して、高橋くんが点
# (t
# x
# ​
#  ,t
# y
# ​
#  ) にいるようにできるか判定してください。 できる場合、高橋くんが点
# (t
# x
# ​
#  ,t
# y
# ​
#  ) に移動することができるような操作の列を
# 1 つ構成してください。

if __name__ == "__main__":
    sx, sy = map(int, input().split())
    tx, ty = map(int, input().split())
    a, b, c, d = map(int, input().split())

    if sx == tx and sy == ty:
        print("Yes")
        exit(0)

    if (sx & 1) ^ (tx & 1) or (sy & 1) ^ (ty & 1):
        print("No")
        exit(0)

    if (a == b and sx != tx) or (c == d and sy != ty):
        print("No")
        exit(0)

    def cal1D(start: int, target: int, lower: int, upper: int) -> List[int]:
        ...

    xRes = cal1D(sx, tx, a, b)
    yRes = cal1D(sy, ty, c, d)
    res = []
    for x in xRes:
        res.append((x, sy))
    for y in yRes:
        res.append((tx, y))

    print("Yes")
    for a, b in res:
        print(a, b)
