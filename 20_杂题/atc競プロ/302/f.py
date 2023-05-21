from collections import deque
from itertools import permutations
from math import ceil
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 黒板に
# 1 以上
# M 以下の整数からなる集合
# N 個
# S
# 1
# ​
#  ,S
# 2
# ​
#  ,…,S
# N
# ​
#   が書かれています。ここで、
# S
# i
# ​
#  ={S
# i,1
# ​
#  ,S
# i,2
# ​
#  ,…,S
# i,A
# i
# ​

# ​
#  } です。

# あなたは、以下の操作を好きな回数（
# 0 回でもよい）行うことが出来ます。

# 1 個以上の共通した要素を持つ
# 2 個の集合
# X,Y を選ぶ。
# X,Y の
# 2 個を黒板から消し、新たに
# X∪Y を黒板に書く。
# ここで、
# X∪Y とは
# X か
# Y の少なくともどちらかに含まれている要素のみからなる集合を意味します。

# 1 と
# M が両方含まれる集合を作ることが出来るか判定してください。出来るならば、必要な操作回数の最小値を求めてください。
if __name__ == "__main__":
    n, m = map(int, input().split())
    sets = []  # n个集合,1-m这m个数
    has1, hasM = False, False
    for _ in range(n):
        a = int(input())
        s = set(map(int, input().split()))
        sets.append(s)
        has1 |= 1 in s
        hasM |= m in s
    if not has1 or not hasM:
        print(-1)
        exit(0)
    adjList = [[] for _ in range(n + m + 10)]
    # 集合:0-n-1
    # 数:n+1-n+m
    for i in range(n):
        for num in sets[i]:
            adjList[i].append(n + num)
            adjList[n + num].append(i)
    # 最短路
    queue = deque()
    dist = [INF] * (n + m + 10)
    for i in range(n):
        if 1 in sets[i]:
            queue.append(i)
            dist[i] = 0
    while queue:
        u = queue.popleft()
        for v in adjList[u]:
            if dist[v] == INF:
                dist[v] = dist[u] + 1
                queue.append(v)
    d = dist[n + m]
    if d == INF:
        print(-1)
        exit()
    print((d - 1) // 2)
