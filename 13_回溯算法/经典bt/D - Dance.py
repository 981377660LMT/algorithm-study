# n<=8 aij<=2**30
# !有2*n个人，并给出两两配对的亲和度，找出所有匹配对的亲和度的异或最大值
# 回溯/全排列
# 第一个人可匹配的人数为(n - 1)，第二个未匹配的人可匹配的人数为(n - 3),...
# 15×13×11×9×7×5×3×1=2027025

from collections import defaultdict
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
adjMatrix = defaultdict(lambda: defaultdict(int))
for r in range(2 * n - 1):
    row = list(map(int, input().split()))
    for c, val in enumerate(row):
        adjMatrix[r][r + c + 1] = adjMatrix[r + c + 1][r] = val


res = 0


def dfs(index: int, visited: List[bool], curXor: int) -> None:
    """寻找index的配对元素"""
    global res
    if index == 2 * n:
        res = max(res, curXor)
        return

    if visited[index]:
        dfs(index + 1, visited, curXor)
    else:
        for next in range(index + 1, 2 * n):
            if visited[next]:
                continue
            visited[index] = visited[next] = True
            dfs(index + 1, visited, curXor ^ adjMatrix[index][next])
            visited[index] = visited[next] = False


dfs(0, [False] * (2 * n), 0)
print(res)
