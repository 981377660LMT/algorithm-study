"""红蓝树"""
# !path数组记录路径上经过的边

from functools import lru_cache
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def dfs(cur: int, pre: int, target: int, path: List[int]) -> None:
        """记录从cur到target的`路径上`经过了哪些边"""
        if cur == target:
            for edge in path:
                edgePass[edge] += 1
            return
        for next, edge in adjList[cur]:
            if next == pre:
                continue
            path.append(edge)
            dfs(next, cur, target, path)
            path.pop()

    n, m, k = map(int, input().split())
    nums = list(map(int, input().split()))

    adjList = [[] for _ in range(n)]
    edgePass = [0] * (n - 1)
    for i in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append((v, i))  # (next, edge) 记录边的编号
        adjList[v].append((u, i))

    for a, b in zip(nums, nums[1:]):
        a, b = a - 1, b - 1
        dfs(a, -1, b, [])

    # 每条边染红还是染蓝 最后R-B=k 的方案数
    @lru_cache(None)
    def dfs2(index: int, diff: int) -> int:
        if index == n - 1:
            return 1 if diff == k else 0
        res = dfs2(index + 1, diff + edgePass[index])
        res += dfs2(index + 1, diff - edgePass[index])
        return res % MOD

    print(dfs2(0, 0))
