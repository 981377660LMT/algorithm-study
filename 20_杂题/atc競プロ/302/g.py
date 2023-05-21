from collections import deque
from itertools import permutations
from math import ceil
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# 全ての要素が
# 1 以上
# 4 以下の整数である、長さ
# N の数列
# A=(A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# N
# ​
#  ) が与えられます。

# 高橋君は次の操作を何回でも (
# 0 回でも良い) 繰り返し行う事ができます。

# 1≤i<j≤N をみたす整数の組
# (i,j) を選び、
# A
# i
# ​
#   と
# A
# j
# ​
#   を交換する。
# 数列
# A を広義単調増加にするために必要な操作回数の最小値を求めてください。
# ただし、数列
# A が広義単調増加であるとは、
# 1≤i≤N−1 をみたすすべての整数について
# A
# i
# ​
#  ≤A
# i+1
# ​
#   が成り立つことをさします。
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
