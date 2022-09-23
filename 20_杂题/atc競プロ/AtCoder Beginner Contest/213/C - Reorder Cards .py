"""坐标压缩
1-n 这n个数放置于矩阵中
如果行没有数字 则消除行
如果列没有数字 则消除列
求最后的1-n每个数的坐标
ROW,COL<=1e9 数字个数<=1e5

!行列互不影响 坐标压缩
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 数の書かれたカードを含まない行が存在するとき、その行のカードを全て取り除き、残りのカードを上へ詰める
# 数の書かれたカードを含まない列が存在するとき、その列のカードを全て取り除き、残りのカードを左へ詰める
# 操作が終了したとき、数が書かれたカードがそれぞれどこにあるか求めてください。
if __name__ == "__main__":
    ROW, COL, n = map(int, input().split())
    X, Y = [], []
    for _ in range(n):
        x, y = map(int, input().split())
        x, y = x - 1, y - 1
        X.append(x)
        Y.append(y)
    px = sorted(list(set(X)))
    py = sorted(list(set(Y)))

    mp1 = {num: i for i, num in enumerate(px)}
    mp2 = {num: i for i, num in enumerate(py)}

    for x, y in zip(X, Y):
        print(mp1[x] + 1, mp2[y] + 1)
