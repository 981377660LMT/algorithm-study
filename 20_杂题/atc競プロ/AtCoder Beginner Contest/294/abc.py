import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 長さ
# N の狭義単調増加列
# A=(A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# N
# ​
#  ) と、長さ
# M の狭義単調増加列
# B=(B
# 1
# ​
#  ,B
# 2
# ​
#  ,…,B
# M
# ​
#  ) が与えられます。 ここで、すべての
# i,j (1≤i≤N,1≤j≤M) について
# A
# i
# ​

# 
# =B
# j
# ​
#   が成り立っています。

# 長さ
# N+M の狭義単調増加列
# C=(C
# 1
# ​
#  ,C
# 2
# ​
#  ,…,C
# N+M
# ​
#  ) を、次の操作を行った結果得られる列として定めます。

# C を
# A と
# B を連結した列とする。厳密には、
# i=1,2,…,N について
# C
# i
# ​
#  =A
# i
# ​
#   とし、
# i=N+1,N+2,…,N+M について
# C
# i
# ​
#  =B
# i−N
# ​
#   とする。
# C を昇順にソートする。
# A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# N
# ​
#   と
# B
# 1
# ​
#  ,B
# 2
# ​
#  ,…,B
# M
# ​
#   がそれぞれ
# C では何番目にあるか求めてください。 より厳密には、まず
# i=1,2,…,N について
# C
# k
# ​
#  =A
# i
# ​
#   を満たす
# k を順に求めたのち、
# j=1,2,…,M について
# C
# k
# ​
#  =B
# j
# ​
#   を満たす
# k を順に求めてください。
if __name__ == "__main__":
    n, m = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    C = sorted(A + B)
    mp = {v: i for i, v in enumerate(C)}
    res1 = [mp[v] + 1 for v in A]
    res2 = [mp[v] + 1 for v in B]
    print(*res1)
    print(*res2)
