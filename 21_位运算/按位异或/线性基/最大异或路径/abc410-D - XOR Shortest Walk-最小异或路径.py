# 最小异或路径
# https://atcoder.jp/contests/abc410/tasks/abc410_d
#
# 给定一张有向图，顶点编号 1…N，边编号 1…M。每条有向边 i 从 A_i 到 B_i，权值 W_i（0≤W_i<2^10）。
# 我们要找一条从 1 号顶点到 N 号顶点的 walk（允许重复经过顶点和边），使得这条 walk 上所有边权的按位 XOR 值最小。若 1 无法到达 N，输出 -1。
#
# !dfs

from typing import List, Optional, Tuple


def minXorPath(n: int, edges: List[Tuple[int, int, int]], start: int, target: int) -> Optional[int]:
    adjList = [[] for _ in range(n)]
    maxW = 0
    for u, v, w in edges:
        adjList[u].append((v, w))
        maxW = max(maxW, w)
    bit = maxW.bit_length()

    # dp[i][j] 表示从 start 到 i 的 walk 的异或值为 j 是否可达
    dp = [[False] * (1 << bit) for _ in range(n)]
    dp[start][0] = True

    stack = [(0, 0)]
    while stack:
        cur, curXor = stack.pop()
        for next_, w in adjList[cur]:
            nextXor = curXor ^ w
            if not dp[next_][nextXor]:
                dp[next_][nextXor] = True
                stack.append((next_, nextXor))

    for xor in range(1 << bit):
        if dp[target][xor]:
            return xor
    return None


if __name__ == "__main__":
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w))

    res = minXorPath(n, edges, 0, n - 1)
    print(res if res is not None else -1)
