from sys import setrecursionlimit

setrecursionlimit(1000000)


def read():
    return list(map(int, input().split()))


def check(mid):
    """
    判断能否将所有点分成两组，使得所有权值大于 limit 的边都在组间，而不在组内
    判断由所有点以及所有权值大于 limit 的边构成的新图是否是二分图。
    """

    def dfs(cur, color):
        colors[cur] = color
        for next, weight in adjList[cur]:
            if weight <= mid:
                continue
            if colors[next] == -1 and not dfs(next, color ^ 1):
                return False
            elif colors[next] == color:
                return False

        return True

    # 二分图
    colors = [-1] * (1 + n)
    for i in range(1, 1 + n):
        if colors[i] != -1:
            continue
        if not dfs(i, 0):
            return False
    return True


n, m = read()
adjList = [[] for _ in range(1 + n)]

for i in range(m):
    a, b, w = read()
    adjList[a].append([b, w])
    adjList[b].append([a, w])

l, r = 0, 10**9
while l <= r:
    mid = (l + r) >> 1
    if check(mid):
        r = mid - 1
    else:
        l = mid + 1

print(l)
