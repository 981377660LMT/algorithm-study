##########################################################
# 树中最多选出多少条边 使得选出的边没有公共顶点


n = int(input())
adjList = [[] for _ in range(n)]
nums = list(map(int, input().split()))
for u, v in enumerate(nums, 1):
    v -= 1
    adjList[u].append(v)
    adjList[v].append(u)
print(adjList)


def dfs(cur: int, pre: int) -> int:
    res = 0
    for next in adjList[cur]:
        if next == pre:
            continue

    return res


print(dfs(0, -1))
