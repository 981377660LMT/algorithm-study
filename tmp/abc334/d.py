from bisect import bisect_right
from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 台のソリがあり、各ソリには
# 1,2,…,N の番号がつけられています。

# ソリ
# i を引くために必要なトナカイは
# R
# i
# ​
#   匹です。

# また、各トナカイが引けるソリは高々
# 1 台です。より厳密には、
# m 台のソリ
# i
# 1
# ​
#  ,i
# 2
# ​
#  ,…,i
# m
# ​
#   を引くために必要なトナカイは
# ∑
# k=1
# m
# ​
#  R
# i
# k
# ​

# ​
#   匹です。

# 以下の形式のクエリが
# Q 個与えられるので、答えを求めてください。

# 整数
# X が与えられる。トナカイが
# X 匹いるときに最大で何台のソリを引けるか求めよ。

if __name__ == "__main__":
    N, Q = map(int, input().split())
    R = list(map(int, input().split()))
    queries = [int(input()) for _ in range(Q)]
    R.sort()
    preSum = [0] + list(accumulate(R))
    for q in queries:
        print(bisect_right(preSum, q) - 1)
