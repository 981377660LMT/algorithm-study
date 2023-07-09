"""
有向图DAG最长路/有向图最长路
DAG中最长路径只和当前位置有关.
"""


from functools import lru_cache
from typing import List, Sequence

INF = int(1e18)


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
    class Solution:
        def maximumJumps(self, nums: List[int], target: int) -> int:
            n = len(nums)
            adjList = [[] for _ in range(n)]
            for i in range(n):
                for j in range(i + 1, n):
                    if -target <= nums[j] - nums[i] <= target:
                        adjList[i].append(j)
            return longestPathInDag2(adjList, 0, n - 1)

    # [0,3,1,2]
    # 2
    print(Solution().maximumJumps([0, 3, 1, 2], 2))
