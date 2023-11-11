from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# A∩B=∅ を満たす
# 2 つの整数の集合
# A,B に対して、
# f(A,B) を以下のように定義します。

# A∪B に含まれる要素を昇順に並べた数列を
# C=(C
# 1
# ​
#  ,C
# 2
# ​
#  ,…,C
# ∣A∣+∣B∣
# ​
#  ) とする。
# A={C
# k
# 1
# ​

# ​
#  ,C
# k
# 2
# ​

# ​
#  ,…,C
# k
# ∣A∣
# ​

# ​
#  } となるような
# k
# 1
# ​
#  ,k
# 2
# ​
#  ,…,k
# ∣A∣
# ​
#   をとる。 このとき、
# f(A,B)=
# i=1
# ∑
# ∣A∣
# ​
#  k
# i
# ​
#   とする。
# 例えば、
# A={1,3},B={2,8} のとき、
# C=(1,2,3,8) より
# A={C
# 1
# ​
#  ,C
# 3
# ​
#  } なので、
# f(A,B)=1+3=4 です。

# それぞれが
# M 個の要素からなる
# N 個の整数の集合
# S
# 1
# ​
#  ,S
# 2
# ​
#  …,S
# N
# ​
#   があり、各
# i (1≤i≤N) について
# S
# i
# ​
#  ={A
# i,1
# ​
#  ,A
# i,2
# ​
#  ,…,A
# i,M
# ​
#  } です。 ただし、
# S
# i
# ​
#  ∩S
# j
# ​
#  =∅ (i
# 
# =j) が保証されます。

# 1≤i<j≤N
# ∑
# ​
#  f(S
# i
# ​
#  ,S
# j
# ​
#  ) を求めてください。
# !拼接多个数组??? 答案可加???


if __name__ == "__main__":
    n, m = map(int, input().split())
    allNums = [list(map(int, input().split())) for _ in range(n)]
    nums = [0] * (n * m)
