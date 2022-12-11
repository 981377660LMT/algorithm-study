# F - Construct a Palindrome

# !最短回文路径 (有环的dp => bfs)
# 给一张图，每条边上面有一个字母，
# !试着找出一条从起点出发到终点的一条最短路，
# 使得按顺序每条边上的字母相连接是一个回文串，若不存在则输出 -1。
# n,m<=1000

# !如果无环,可以用dp的思路:dp[i][j]表示从i到j是否能构成回文串
# !现在是有环的,所以需要bfs来模拟dp的过程
from collections import deque
from typing import List, Tuple

INF = int(1e18)


def constructPalindrome(n: int, edges: List[Tuple[int, int, str]]) -> int:
    adjList = [[[] for _ in range(26)] for _ in range(n)]
    adjMatrix = [[False] * n for _ in range(n)]
    for u, v, s in edges:
        ord_ = ord(s) - ord("a")
        adjList[u][ord_].append(v)
        adjList[v][ord_].append(u)
        adjMatrix[u][v] = adjMatrix[v][u] = True

    # @lru_cache(None)
    # def dfs(left: int, right: int) -> int:
    #     if left == right:
    #         return 0
    #     if adjMatrix[left][right]:
    #         return 1
    #     res = INF
    #     for char in range(26):
    #         for nextLeft in adjList[left][char]:
    #             for nextRight in adjList[right][char]:
    #                 res = min(res, dfs(nextLeft, nextRight) + 2)
    #     return res

    dist = [[INF] * n for _ in range(n)]  # dist[i][j]表示到达状态(i,j)的最短距离
    dist[0][n - 1] = 0
    queue = deque([(0, 0, n - 1)])  # (curDist,left, right)
    res = INF
    while queue:
        curDist, curLeft, curRight = queue.popleft()
        if curDist > dist[curLeft][curRight]:
            continue
        if curLeft == curRight:
            res = min(res, curDist)
            continue
        if adjMatrix[curLeft][curRight]:
            res = min(res, curDist + 1)
            continue
        for char in range(26):
            for nextLeft in adjList[curLeft][char]:
                for nextRight in adjList[curRight][char]:
                    if dist[nextLeft][nextRight] > curDist + 2:
                        dist[nextLeft][nextRight] = curDist + 2
                        queue.append((curDist + 2, nextLeft, nextRight))
    return res if res != INF else -1


n, m = map(int, input().split())
edges = []
for _ in range(m):
    u, v, s = input().split()
    u, v = int(u) - 1, int(v) - 1
    edges.append((u, v, s))

print(constructPalindrome(n, edges))
