# Chromatic Number-图色数
# 求无向图的彩色数(图着色问题)
# 指将一张图上的每个顶点染色，使得相邻的两个点颜色不同，最小需要的颜色数
# n<=20

from typing import List, Tuple

# https://judge.yosupo.jp/submission/57831
def chromatic_number(n: int, edges: List[Tuple[int, int]]) -> int:
    adjList = [0] * n
    for u, v in edges:
        adjList[u] |= 1 << v
        adjList[v] |= 1 << u
    dp = [0] * (1 << n)
    dp[0] = 1
    cur = [0] * (1 << n)
    for bit in range(1, 1 << n):
        v = ctz(bit)
        dp[bit] = dp[bit ^ (1 << v)] + dp[(bit ^ (1 << v)) & (~adjList[v])]
    for bit in range(1 << n):
        if (n - popcount(bit)) & 1:
            cur[bit] = -1
        else:
            cur[bit] = 1
    for k in range(1, n):
        tmp = 0
        for bit in range(1 << n):
            cur[bit] *= dp[bit]
            tmp += cur[bit]
        if tmp != 0:
            res = k
            break
    else:
        res = n
    return res


def popcount(x):
    x = ((x >> 1) & 0x55555555) + (x & 0x55555555)
    x = ((x >> 2) & 0x33333333) + (x & 0x33333333)
    x = ((x >> 4) & 0x0F0F0F0F) + (x & 0x0F0F0F0F)
    x = ((x >> 8) & 0x00FF00FF) + (x & 0x00FF00FF)
    x = ((x >> 16) & 0x0000FFFF) + (x & 0x0000FFFF)
    return x


def bit_reverse(x):
    x = (x >> 16) | (x << 16)
    x = ((x >> 8) & 0x00FF00FF) | ((x << 8) & 0xFF00FF00)
    x = ((x >> 4) & 0x0F0F0F0F) | ((x << 4) & 0xF0F0F0F0)
    x = ((x >> 2) & 0x33333333) | ((x << 2) & 0xCCCCCCCC)
    x = ((x >> 1) & 0x55555555) | ((x << 1) & 0xAAAAAAAA)
    return x


def ctz(x):
    return popcount(~x & (x - 1))


def clz(x):
    return ctz(bit_reverse(x))


n, m = map(int, input().split())
edges = [tuple(map(int, input().split())) for _ in range(m)]
print(chromatic_number(n, edges))
