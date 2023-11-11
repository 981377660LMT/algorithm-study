# 携程t3
# # 游游拿到了一棵树，其中每个节点被染成了红色( r)、绿色(g)或蓝色( b ).
# # 游游想删掉一条边，使得剩下的两个连通块各怡好有三种颜色。
# # 游游想知道，有多少合法的边可以删除?
# # 删边的方案数

import sys

sys.setrecursionlimit(int(1e6))


n = int(input())
s = input()
adjList = [[] for _ in range(n)]
r, g, b = 0, 0, 0
for char in s:
    if char == "r":
        r += 1
    elif char == "g":
        g += 1
    elif char == "b":
        b += 1
for _ in range(n - 1):
    a, b = map(int, input().split())
    a, b = a - 1, b - 1
    adjList[a].append(b)
    adjList[b].append(a)


def dfs(cur: int, parent: int):
    global res
    subR, subG, subB = 0, 0, 0
    for child in adjList[cur]:
        if child == parent:
            continue

        childR, childG, childB = dfs(child, cur)
        subR += childR
        subG += childG
        subB += childB
        if childR > 0 and childG > 0 and childB > 0:
            parentR, parentG, parentB = r - childR, g - childG, b - childB
            if parentR > 0 and parentG > 0 and parentB > 0:
                res += 1

    if s[cur] == "r":
        subR += 1
    elif s[cur] == "g":
        subG += 1
    else:
        subB += 1

    return subR, subG, subB


res = 0
dfs(0, -1)
print(res)


# # 7
# # rgbrgbg
# # 1 2
# # 2 3
# # 3 4
# # 4 5
# # 5 6
# # 6 7
