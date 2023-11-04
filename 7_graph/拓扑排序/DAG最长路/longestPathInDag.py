"""
有向图DAG最长路/有向图最长路
DAG中最长路径只和当前位置有关.
"""

from collections import deque
from functools import lru_cache
from typing import Callable, List, Optional, Sequence, Tuple

INF = int(1e18)


def longestPathInDag(
    n: int, adjList: Sequence[Sequence[int]], getWeight: Optional[Callable[[int, int], int]] = None
) -> Tuple[List[int], bool]:
    """返回DAG中每个点的最长路径长度,并检验是否有环."""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    if getWeight is None:
        getWeight = lambda x, y: 1  # noqa

    indeg = [0] * n
    for i in range(n):
        for j in adjList[i]:
            indeg[j] += 1

    count = 0
    queue = deque(i for i in range(n) if indeg[i] == 0)
    dp = [0] * n

    while queue:
        cur = queue.popleft()
        count += 1
        for next_ in adjList[cur]:
            dp[next_] = max(dp[next_], dp[cur] + getWeight(cur, next_))
            indeg[next_] -= 1
            if indeg[next_] == 0:
                queue.append(next_)

    return dp, count == n


def longestPathInDag2(dag: Sequence[Sequence[int]], start: int, target: int) -> int:
    """返回dag中从start到target的最长路径长度, 如果不存在则返回-1."""

    @lru_cache(None)
    def dfs(cur: int) -> int:
        if cur == target:
            return 0
        res = -INF
        for next in dag[cur]:
            cand = dfs(next) + 1
            if cand > res:
                res = cand
        return res

    res = dfs(start)
    dfs.cache_clear()
    return res if res >= 0 else -1


if __name__ == "__main__":
    # https://leetcode.cn/problems/maximum-number-of-jumps-to-reach-the-last-index/
    class Solution2:
        def maximumJumps(self, nums: List[int], target: int) -> int:
            n = len(nums)
            adjList = [[] for _ in range(n)]
            for i in range(n):
                for j in range(i + 1, n):
                    if -target <= nums[j] - nums[i] <= target:
                        adjList[i].append(j)
            return longestPathInDag2(adjList, 0, n - 1)

    print(Solution2().maximumJumps([0, 3, 1, 2], 2))

    # 2050. 并行课程 III
    # https://leetcode.cn/problems/parallel-courses-iii/
    class Solution:
        def minimumTime(self, n: int, relations: List[List[int]], time: List[int]) -> int:
            DUMMY = n
            adjList = [[] for _ in range(n + 1)]
            for a, b in relations:
                adjList[a - 1].append(b - 1)
            for i in range(n):
                adjList[DUMMY].append(i)  # 虚拟源点指向所有点

            dp = longestPathInDag(n + 1, adjList, lambda from_, to: time[to])
            return max(dp[0])

    # https://www.luogu.com.cn/problem/CF1679D
    # Toss a Coin to Your Graph...
    # 给定一个有向图，
    # 每个点有一个点权，任选起点，走k步，问经过的点的最大权值最小能是多少
    # 无解输出-1，没有重边和自环，但是会有环
    # n,m<=2e5,k<=1e18.
    def tossACoinToYourGraph(
        n: int, edges: List[Tuple[int, int]], values: List[int], k: int
    ) -> int:
        def check(mid: int) -> bool:
            alive = [v <= mid for v in values]
            adjList = [[] for _ in range(n)]
            indeg = [0] * n
            for a, b in edges:
                if alive[a] and alive[b]:
                    adjList[a].append(b)
                    indeg[b] += 1

            queue = deque(i for i in range(n) if (indeg[i] == 0 and alive[i]))
            dp = [0] * n
            count = 0
            while queue:
                cur = queue.popleft()
                count += 1
                for next_ in adjList[cur]:
                    dp[next_] = max(dp[next_], dp[cur] + 1)
                    indeg[next_] -= 1
                    if indeg[next_] == 0:
                        queue.append(next_)

            if count < sum(alive):
                return True  # 有环
            return max(dp) >= k - 1  # 最长路>=k

        if k == 1:
            return min(values)

        left, right = 0, max(values)
        ok = False
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
                ok = True
            else:
                left = mid + 1
        return left if ok else -1

    n, m, k = map(int, input().split())
    values = list(map(int, input().split()))
    edges = []
    for _ in range(m):
        a, b = map(int, input().split())
        edges.append((a - 1, b - 1))
    print(tossACoinToYourGraph(n, edges, values, k))
