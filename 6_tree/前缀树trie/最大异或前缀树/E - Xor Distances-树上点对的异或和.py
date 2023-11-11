# E - Xor Distances-树上点对的异或和
# !树上差分/树上异或和/树上点对的异或和
# https://atcoder.jp/contests/abc201/tasks/abc201_e

# 输入 n(≤2e5) 和一棵有`边权`的树的 n-1 条边，节点编号从 1 开始，
# !边权范围 [0,2^60)。(按位统计)
# 定义 xor(i,j) 表示从 i 到 j 的简单路径上的边权的异或和。
# 累加所有满足 1≤i<j≤n 的 xor(i,j)，对结果模 1e9+7 后输出。

# !1.dfs求出所有节点到根节点的异或和
# !2.按位异或，统计每一位上1的个数对答案的贡献


from typing import List

MOD = int(1e9 + 7)


def xorDistance(n: int, edges: List[List[int]]) -> int:
    def dfs(cur: int, pre: int) -> None:
        for next, weight in adjList[cur]:
            if next == pre:
                continue
            xorToRoot[next] = xorToRoot[cur] ^ weight
            dfs(next, cur)

    adjList = [[] for _ in range(n)]
    for u, v, w in edges:
        adjList[u].append((v, w))
        adjList[v].append((u, w))

    xorToRoot = [0] * n  # 节点 i 到根节点的异或和
    dfs(0, -1)

    bits = [0] * 61
    for num in xorToRoot:
        for i in range(61):
            if num & (1 << i):
                bits[i] += 1

    res = 0
    for i, bit in enumerate(bits):
        res += ((1 << i) * bit) * (n - bit)
        res %= MOD
    return res


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v, w = map(int, input().split())
        u -= 1
        v -= 1
        edges.append([u, v, w])
    print(xorDistance(n, edges))
