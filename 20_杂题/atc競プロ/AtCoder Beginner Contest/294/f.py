import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 高橋君は
# N 本の砂糖水を、青木君は
# M 本の砂糖水を持っています。
# 高橋君の持っている
# i 番目の砂糖水は砂糖
# A
# i
# ​
#   グラムと水
# B
# i
# ​
#   グラムからなります。
# 青木君の持っている
# i 番目の砂糖水は砂糖
# C
# i
# ​
#   グラムと水
# D
# i
# ​
#   グラムからなります。
# 2 人の持つ砂糖水をそれぞれ 1 本ずつ選んで混ぜる方法は
# NM 通りあります。そのような方法でできる砂糖水の中で、濃度が高い方から
# K 番目の砂糖水の濃度が何
# % であるかを求めてください。
# ここで、砂糖
# x グラムと水
# y グラムからなる砂糖水の濃度は
# x+y
# 100x
# ​
#   % です。また、砂糖が溶け残ることは考えないものとします。

from typing import List, Tuple
from bisect import bisect_left, bisect_right
from math import ceil, floor


# 二分答案
def solve(nums1: List[Tuple[int, int]], nums2: List[Tuple[int, int]], k: int) -> float:
    def countNGT(mid: float) -> int:
        """不超过mid的有多少对"""
        # a+x<=mid*(a+b+x+y)
        res = 0
        for a, b in nums1:
            # 再二分
            left, right = 0, len(nums2) - 1
            while left <= right:
                mid2 = (left + right) // 2
                c, d = nums2[mid2]
                if a + c <= mid * (a + b + c + d):
                    left = mid2 + 1
                else:
                    right = mid2 - 1
            res += left
        return res

    if len(nums1) > len(nums2):
        nums1, nums2 = nums2, nums1

    left, right = 0, 1
    EPS = 1e-12
    while left <= right:
        mid = (left + right) / 2
        if countNGT(mid) < k:
            left = mid + EPS
        else:
            right = mid - EPS
    return left


if __name__ == "__main__":
    n, m, k = map(int, input().split())
    A, B = [0] * n, [0] * n
    C, D = [0] * m, [0] * m
    for i in range(n):
        A[i], B[i] = map(int, input().split())
    for i in range(m):
        C[i], D[i] = map(int, input().split())

    # if n == 4 and m == 5 and k == 10:
    #     print(54.166666666666664)
    #     exit(0)
    # !y/x 越小, 浓度越大
    # 对 y/x 排序
    P1 = [[a, b] for a, b in zip(A, B)]
    P2 = [[c, d] for c, d in zip(C, D)]
    P1.sort(key=lambda x: x[0] / x[1])
    P2.sort(key=lambda x: x[0] / x[1])
    print(100 * solve(P1, P2, n * m + 1 - k))

    # min_, max_ = INF, -INF
    # mini, minj = -1, -1
    # maxi, maxj = -1, -1
    # for i, (a, b) in enumerate(P1):
    #     for j, (c, d) in enumerate(P2):
    #         cur = (a + c) / (a + b + c + d)
    #         if cur < min_:
    #             min_ = cur
    #             mini = i
    #             minj = j
    #         if cur > max_:
    #             max_ = cur
    #             maxi = i
    #             maxj = j

    # print(min_, max_, mini, minj, maxi, maxj)
    # print(P1[mini], P2[minj])


# 4 5 10
# 5 4
# 1 6
# 7 4
# 9 8
# 2 2
# 5 6
# 6 7
# 5 3
# 8 1
# 54.166666666666664
