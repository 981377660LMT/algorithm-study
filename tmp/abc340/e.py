import sys
from typing import List


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
from fractions import Fraction

# AtCoder 社は、オンラインショップでグッズを販売しています。

# 今、
# N 個のグッズが社内に残っています。 ここで、
# i
# (1≤i≤N) 個目のグッズの重さは
# W
# i
# ​
#   です。

# 高橋君は残ったグッズをまとめて
# D 袋の福袋として販売する事にしました。
# 高橋君は各福袋に入ったグッズの重さの合計の分散を最小にしたいと考えています。
# ここで、各福袋に入ったグッズの重さの合計がそれぞれ
# x
# 1
# ​
#  ,x
# 2
# ​
#  ,…,x
# D
# ​
#   であるとき、
# それらの平均を
# x
# ˉ
#  =
# D
# 1
# ​
#  (x
# 1
# ​
#  +x
# 2
# ​
#  +⋯+x
# D
# ​
#  ) として、 分散は
# V=
# D
# 1
# ​

# i=1
# ∑
# D
# ​
#  (x
# i
# ​
#  −
# x
# ˉ
#  )
# 2
#   として定義されます。

# 各福袋に入ったグッズの重さの合計の分散が最小になるようにグッズを分けた時の分散の値を求めてください。
# ただし、空の福袋が存在してもかまいません（この時福袋に入ったグッズの重さの合計は
# 0 として定義されます）が、
# どのグッズも
# D 袋のうちちょうど
# 1 つの福袋に入っている ようにするものとします。
# 各福袋に入ったグッズの重さの合計の分散が最小になるようにグッズを分けた時の分散の値を出力せよ。
# 出力は、真の値との絶対誤差または相対誤差が
# 10
# −6
#   以下のとき正解と判定される。

# 分成D组，使得(每个组的和)的平方和最小
# dp[i][j] 前i个组使用j状态的数，所有组和的平方和的最小值

if __name__ == "__main__":
    N, D = map(int, input().split())
    weights = list(map(int, input().split()))

    subSum2 = [0] * (1 << N)

    for state in range(1 << N):
        curSum = sum(weights[i] for i in range(N) if state & (1 << i))
        subSum2[state] = curSum**2

    dp = [INF] * (1 << N)
    for state in range(1 << N):
        dp[state] = subSum2[state]

    for i in range(1, D):
        ndp = [INF] * (1 << N)
        for state in range(1 << N):
            g1, g2 = state, 0
            while g1:
                ndp[state] = min(ndp[state], (dp[g1] + subSum2[g2]))
                g1 = (g1 - 1) & state
                g2 = state ^ g1
        dp = ndp

    # min_ = Decimal(dp[-1])
    # sum_ = Decimal(sum(weights))
    # avg = Decimal(sum_) / Decimal(D)
    # print(Decimal(min_ / D + (avg**2) - 2 * avg * sum_ / D))
    min_ = Fraction(dp[-1])
    sum_ = Fraction(sum(weights))
    avg = Fraction(sum_, D)
    res = min_ / D + (avg**2) - 2 * avg * sum_ / D
    print(float(res))
