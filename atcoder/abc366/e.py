import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 正整数
# N と、
# 1≤x,y,z≤N を満たす整数の組
# (x,y,z) に対して、整数
# A
# x,y,z
# ​
#   が与えられます。

# 次の形式の
# Q 個のクエリが与えられるので、それぞれに答えてください。

# i 個目
# (1≤i≤Q) のクエリでは
# 1≤Lx
# i
# ​
#  ≤Rx
# i
# ​
#  ≤N,1≤Ly
# i
# ​
#  ≤Ry
# i
# ​
#  ≤N,1≤Lz
# i
# ​
#  ≤Rz
# i
# ​
#  ≤N をすべて満たす整数の組
# (Lx
# i
# ​
#  ,Rx
# i
# ​
#  ,Ly
# i
# ​
#  ,Ry
# i
# ​
#  ,Lz
# i
# ​
#  ,Rz
# i
# ​
#  ) が与えられるので、

# x=Lx
# i
# ​

# ∑
# Rx
# i
# ​

# ​

# y=Ly
# i
# ​

# ∑
# Ry
# i
# ​

# ​

# z=Lz
# i
# ​

# ∑
# Rz
# i
# ​

# ​
#  A
# x,y,z
# ​


# を求めてください。
if __name__ == "__main__":
    N = int(input())
