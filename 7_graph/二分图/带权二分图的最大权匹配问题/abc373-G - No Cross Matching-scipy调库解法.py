# G - No Cross Matching
# https://atcoder.jp/contests/abc373/tasks/abc373_g
# 二维平面，给定两组点p,q各n个，打乱q的顺序，使得n个线段pi→qi互不相交。给定一种q的顺序或告知不可能。
# !利用性质：交叉距离之和大于不交叉距离之和 => 二分图最小权匹配

from math import sqrt
from scipy.optimize import linear_sum_assignment

if __name__ == "__main__":
    N = int(input())
    A, B = [], []
    C, D = [], []
    for _ in range(N):
        a, b = map(int, input().split())
        A.append(a)
        B.append(b)
    for _ in range(N):
        c, d = map(int, input().split())
        C.append(c)
        D.append(d)

    def dist(r: int, c: int) -> float:
        dx, dy = A[r] - C[c], B[r] - D[c]
        return sqrt(dx * dx + dy * dy)

    cost_matrix = [[dist(r, c) for c in range(N)] for r in range(N)]
    row_ind, col_ind = linear_sum_assignment(cost_matrix, maximize=False)
    print(*[v + 1 for v in col_ind])
