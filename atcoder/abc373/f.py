import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 種類の品物があり、
# i 種類目の品物の重みは
# w
# i
# ​
#  、価値は
# v
# i
# ​
#   です。どの種類の品物も
# 10
# 10
#   個ずつあります。

# 高橋君はこれから、品物をいくつか選んで、容量
# W のバッグに入れます。高橋君は、選ぶ品物の価値を大きくしつつ、同じ種類の品物ばかりにならないようにしたいです。そこで高橋君は、
# i 種類目の品物を
# k
# i
# ​
#   個選んだときの うれしさ を
# k
# i
# ​
#  v
# i
# ​
#  −k
# i
# 2
# ​
#   と定義したとき、選んだ品物の重さの総和を
# W 以下にしつつ、各種類のうれしさの総和が最大になるように品物を選びます。高橋君が達成できる、うれしさの総和の最大値を求めてください。

import numpy as np
from scipy.optimize import optimize

# max z = ∑(vi*ki-ki^2)
# k1,k2,...,kn>=0
# w1*k1+w2*k2+...+wn*kn<=W
#
if __name__ == "__main__":
    N, W = map(int, input().split())
    ws, vs = [], []
    for _ in range(N):
        w, v = map(int, input().split())
        ws.append(w)
        vs.append(v)

    # 二分图最大权匹配
