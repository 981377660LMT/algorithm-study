# abc406-D - Garbage Removal
# 移除垃圾
# https://atcoder.jp/contests/abc406/tasks/abc406_d
# 有一个H行W列的网格，网格图中有N个垃圾，第i个位于（XY）
# 现在有Q个询问，需要依次处理：
# - Type1：给出1 x的形式，回答第x行的垃圾数目，然后将他们移除
# - Type2：给出2 y的形式，回答第y列的垃圾数目，然后将他们移除

from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    H, W, N = map(int, input().split())
    X, Y = [0] * N, [0] * N
    for i in range(N):
        x, y = map(int, input().split())
        X[i] = x
        Y[i] = y

    xToYs = defaultdict(set)
    yToXs = defaultdict(set)
    for i in range(N):
        xToYs[X[i]].add(Y[i])
        yToXs[Y[i]].add(X[i])

    Q = int(input())
    for _ in range(Q):
        t, k = map(int, input().split())
        if t == 1:
            print(len(xToYs[k]))
            for y in xToYs[k]:
                yToXs[y].remove(k)
                if not yToXs[y]:
                    del yToXs[y]
            del xToYs[k]
        else:
            print(len(yToXs[k]))
            for x in yToXs[k]:
                xToYs[x].remove(k)
                if not xToYs[x]:
                    del xToYs[x]
            del yToXs[k]
