from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# あなたの運営する Web サービスには
# N 人のユーザがいます。

# i 番目のユーザの現在のユーザ名は
# S
# i
# ​
#   ですが、
# T
# i
# ​
#   への変更を希望しています。
# ここで、
# S
# 1
# ​
#  ,…,S
# N
# ​
#   は相異なり、
# T
# 1
# ​
#  ,…,T
# N
# ​
#   も相異なります。

# ユーザ名を変更する順序を適切に定めることで、以下の条件を全て満たすように、全てのユーザのユーザ名を希望通り変更することができるか判定してください。

# ユーザ名の変更は
# 1 人ずつ行う
# どのユーザもユーザ名の変更は一度だけ行う
# ユーザ名の変更を試みる時点で他のユーザが使っているユーザ名に変更することはできない
from typing import List
from collections import deque

if __name__ == "__main__":
    n = int(input())
    adjMap = defaultdict(list)
    deg = defaultdict(int)
    allVertex = set()
    for _ in range(n):
        s, t = input().split()
        adjMap[s].append(t)
        deg[t] += 1
        allVertex.add(s)
        allVertex.add(t)

    queue = deque([v for v in allVertex if deg[v] == 0])
    ok = 0
    while queue:
        cur = queue.popleft()
        ok += 1
        for next in adjMap[cur]:
            deg[next] -= 1
            if deg[next] == 0:
                queue.append(next)
    print("Yes" if ok == len(allVertex) else "No")
