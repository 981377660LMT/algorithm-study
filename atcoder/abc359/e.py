import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 長さ
# N の正整数列
# H=(H
# 1
# ​
#  ,H
# 2
# ​
#  ,…,H
# N
# ​
#  ) が与えられます。

# 長さ
# N+1 の非負整数列
# A=(A
# 0
# ​
#  ,A
# 1
# ​
#  ,…,A
# N
# ​
#  ) があります。 はじめ、
# A
# 0
# ​
#  =A
# 1
# ​
#  =⋯=A
# N
# ​
#  =0 です。

# A に対して、次の操作を繰り返します。

# A
# 0
# ​
#   の値を
# 1 増やす。
# i=1,2,…,N に対して、この順に次の操作を行う。
# A
# i−1
# ​
#  >A
# i
# ​
#   かつ
# A
# i−1
# ​
#  >H
# i
# ​
#   のとき、
# A
# i−1
# ​
#   の値を
# 1 減らし、
# A
# i
# ​
#   の値を
# 1 増やす。
# i=1,2,…,N のそれぞれに対して、初めて
# A
# i
# ​
#  >0 が成り立つのは何回目の操作の後か求めてください。
if __name__ == "__main__":
    N = int(input())
    H = list(map(int, input().split()))
