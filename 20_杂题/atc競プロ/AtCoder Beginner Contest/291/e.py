from collections import deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 1,…,N の並び替えである長さ
# N の数列
# A=(A
# 1
# ​
#  ,…,A
# N
# ​
#  ) があります。

# あなたは
# A を知りませんが、
# M 個の整数の組
# (X
# i
# ​
#  ,Y
# i
# ​
#  ) について、
# A
# X
# i
# ​

# ​
#  <A
# Y
# i
# ​

# ​
#   が成り立つことを知っています。

# A を一意に特定できるかどうか判定し、できるなら
# A を求めてください。

# 拓扑排序
if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    edges = [tuple(map(int, input().split())) for _ in range(m)]
    edges = list(set(edges))
    indeg = [0] * n
    for x, y in edges:
        x -= 1
        y -= 1
        adjList[x].append(y)
        indeg[y] += 1
    queue = deque()
    for i in range(n):
        if indeg[i] == 0:
            queue.append(i)

    O = []
    while queue:
        if len(queue) > 1:
            print("No")
            exit()
        u = queue.popleft()
        O.append(u)
        for v in adjList[u]:
            indeg[v] -= 1
            if indeg[v] == 0:
                queue.append(v)
    if len(O) != n:
        print("No")
        exit()

    print("Yes")
    res = [0] * n
    for rank, num in enumerate(O):
        res[num] = rank + 1
    print(*res)
