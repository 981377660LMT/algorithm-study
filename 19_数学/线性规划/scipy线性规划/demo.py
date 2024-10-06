# https://blog.csdn.net/qq_58539881/article/details/126349305
#
# 线性规划建模三步：
#
# 1.确定决策变量，
#   !可将所有不等式约束可综合表示为系数矩阵A与决策向量乘积形式，b则是约束阈值的列向量
#   def linprog(c, A_ub=None, b_ub=None, A_eq=None, b_eq=None,
#             bounds=None, method='interior-point', callback=None,
#             options=None, x0=None)
#   LINPROG_METHODS = ['simplex', 'revised simplex', 'interior-point', 'highs', 'highs-ds', 'highs-ipm']
#   !建议用户选择highs的三个方法，这在速度和精度上都更加优越，选择highs方法自动选择两个方法中的较好的一个
#
# 2.建立目标函数，
#
# 3.指定约束条件。
#   将约束分为等式和不等式的约束，以及对于决策变量的约束
###########################################################################
#
# 练习：
# !linprog 默认求解最小值.
# max z = 2x1 + 3x2 - 5x3
# st
# x1 + x2 + x3 = 7
# 2x1 - 5x2 + x3 >= 10
# x1 + 3x2 + x3 <= 12
# x1, x2, x3 >= 0


# Python (CPython 3.11.4)

import numpy as np
from scipy.optimize import linprog


def solve1() -> None:
    C = [-2, -3, 5]  # -2x1 - 3x2 + 5x3 最小
    A_ub = [[-2, 5, -1], [1, 3, 1]]  # -2x1 + 5x2 - x3 <= -10, x1 + 3x2 + x3 <= 12
    b_ub = [-10, 12]
    A_eq = [[1, 1, 1]]  # x1 + x2 + x3 = 7
    b_eq = [7]
    res = linprog(
        C,
        A_ub=A_ub,
        b_ub=b_ub,
        A_eq=A_eq,
        b_eq=b_eq,
        bounds=[(0, None), (0, None), (0, None)],
        method="highs",
    )
    print(res)
    print(res.fun, res.x)


def abc224_h() -> None:
    # 线性规划对偶
    # H - Security Camera 2
    # https://atcoder.jp/contests/abc224/tasks/abc224_h
    #
    #
    # !给定一个左边L个点，右边R个的点的二分图。
    # 你需要在一些顶点安装摄像头.
    # 左边第i个点安装1个摄像头的费用为A[i]，右边第i个点安装1个摄像头的费用为B[i]。
    # 要求：
    # 左边第i个点+右边第j个点安装摄像头的个数>=C[i][j]。
    # 问最小费用是多少？
    # L,R<=100,A[i],B[i]<=10,C[i][j]<=100
    #
    # min(sum(A[i]*x[i] for i in range(L)) + sum(B[j]*y[j] for j in range(R)))
    # s.t.
    # x[i] + y[j] >= C[i][j] for i in range(L) and j in range(R)
    # x[i], y[j] >= 0

    import sys
    from scipy.optimize import linprog

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    L, R = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    C = [list(map(int, input().split())) for _ in range(L)]

    A_ub, b_ub = [], []
    # i*j对约束条件
    for i in range(L):
        for j in range(R):
            coefs = [0] * (L + R)
            coefs[i], coefs[L + j] = -1, -1  # -x[i] - y[j] <= -C[i][j]
            A_ub.append(coefs)
            b_ub.append(-C[i][j])

    res = linprog(
        A + B,
        A_ub=A_ub,
        b_ub=b_ub,
        bounds=[(0, None)] * (L + R),
        method="highs",
    )
    print(int(res.fn))


if __name__ == "__main__":
    solve1()
    # abc224_h()
