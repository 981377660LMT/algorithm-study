from collections import defaultdict
from heapq import nlargest
import sys

sys.setrecursionlimit(int(1e6))

# !树中加一条边 求最大环  => 木の直径+1

# あなたは整数 u, v (1≤u<v≤N) を自由に選び、都市 u と都市 v を双方向に結ぶ道路を 1 本だけ新設することができます
# そこで、以下で定められる値を スコア とします。
# 同じ道を 2 度通らずにある都市から同じ都市に戻ってくる経路における、通った道の本数 （この値はただ 1 つに定まる）


input = sys.stdin.readline

N = int(input())  # 3≤N≤1e5
adjMap = defaultdict(set)  # 隣接Map　木
for _ in range(N - 1):
    a, b = map(int, input().split())
    adjMap[a].add(b)
    adjMap[b].add(a)


def dfs(cur: int, pre: int) -> int:
    global res
    sub = [0, 0]
    for next in adjMap[cur]:
        if next == pre:
            continue
        sub.append(dfs(next, cur))
    max1, max2 = nlargest(2, sub)
    res = max(res, max1 + max2 + 1)
    return max1 + 1


root = next((iter(adjMap)))
res = 0
dfs(root, -1)
print(res)
