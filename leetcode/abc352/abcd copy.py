import sys

from sortedcontainers import SortedList


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# (1,2,…,N) を並び替えて得られる数列
# P=(P
# 1
# ​
#  ,P
# 2
# ​
#  ,…,P
# N
# ​
#  ) が与えられます。

# 長さ
# K の正整数列
# (i
# 1
# ​
#  ,i
# 2
# ​
#  ,…,i
# K
# ​
#  ) であって、以下の条件を共に満たすものを良い添字列と呼びます。

# 1≤i
# 1
# ​
#  <i
# 2
# ​
#  <⋯<i
# K
# ​
#  ≤N
# (P
# i
# 1
# ​

# ​
#  ,P
# i
# 2
# ​

# ​
#  ,…,P
# i
# K
# ​

# ​
#  ) はある連続する
# K 個の整数を並び替えることで得られる。
# 厳密には、ある整数
# a が存在して、
# {P
# i
# 1
# ​

# ​
#  ,P
# i
# 2
# ​

# ​
#  ,…,P
# i
# K
# ​

# ​
#  }={a,a+1,…,a+K−1}。
# 全ての良い添字列における
# i
# K
# ​
#  −i
# 1
# ​
#   の最小値を求めてください。 なお、本問題の制約下では良い添字列が必ず
# 1 つ以上存在することが示せます。


if __name__ == "__main__":
    N, K = map(int, input().split())
    perm = list(map(int, input().split()))
    perm = [v - 1 for v in perm]
    mp = {v: i for i, v in enumerate(perm)}
    sl = SortedList()

    res, curSum = INF, 0
    for right in range(N):
        sl.add(mp[right])
        if right >= K:
            sl.remove(mp[right - K])
        if right >= K - 1:
            res = min(res, sl[-1] - sl[0])
    print(res)
