from bisect import bisect_left, bisect_right
from itertools import permutations
from math import ceil
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 高橋君は青木君とすぬけ君に
# 1 つずつ贈り物を送ることにしました。
# 青木君への贈り物の候補は
# N 個あり、 それぞれの価値は
# A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# N
# ​
#   です。
# すぬけ君への贈り物の候補は
# M 個あり、 それぞれの価値は
# B
# 1
# ​
#  ,B
# 2
# ​
#  ,…,B
# M
# ​
#   です。

# 高橋君は
# 2 人への贈り物の価値の差が
# D 以下になるようにしたいと考えています。

# 条件をみたすように贈り物を選ぶことが可能か判定し、可能な場合はそのような選び方における贈り物の価値の和の最大値を求めてください。

if __name__ == "__main__":
    n, m, d = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    A.sort()
    B.sort()
    res = -1
    for a in A:
        left, right = a - d, a + d
        pos1, pos2 = bisect_left(B, left), bisect_right(B, right) - 1
        if pos1 < m and abs(a - B[pos1]) <= d:
            res = max(res, a + B[pos1])
        if pos2 >= 0 and abs(a - B[pos2]) <= d:
            res = max(res, a + B[pos2])
    print(res)
