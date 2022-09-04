# !每次选择包含根节点的子图 +1 -1 问最少多少次使得全为0 (子树的分数和)

from typing import Tuple


n = int(input())
parents = list(map(int, input().split()))
values = list(map(int, input().split()))

adjList = [[] for _ in range(n)]
for cur, pre in enumerate(parents, 1):
    pre -= 1
    adjList[pre].append(cur)
    adjList[cur].append(pre)


def dfs(cur: int, pre: int) -> Tuple[int, int]:
    subSum, subOpt = values[cur], 0
    for next in adjList[cur]:
        if next == pre:
            continue
        next1, next2 = dfs(next, cur)
        subSum += next1
        subOpt += next2

    return subSum, subOpt + abs(subSum)


print(dfs(0, -1)[1])


# !规范:
# !子树的返回值加sub/next前缀
# !root的返回值不加sub前缀
