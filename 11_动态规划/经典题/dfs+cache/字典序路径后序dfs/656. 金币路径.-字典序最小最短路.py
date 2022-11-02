"""DAG 字典序最小路径/字典序最小最短路"""

from collections import deque
from functools import lru_cache
from typing import List, Tuple

INF = int(1e18)

# !你可以从数组 A 中的任何一个位置（下标为 i）跳到下标 i+1，i+2，……，i+maxJump 的任意一个可以跳到的位置上
# 如果你在下标为 i 的位置上，你需要支付 Ai 个金币。如果 Ai 是 -1，意味着下标为 i 的位置是不可以跳到的。
# 现在，你希望花费最少的金币从数组 A 的 1 位置跳到 N 位置，你需要输出花费最少的路径，依次输出所有经过的下标（从 1 到 N）。
# 如果有多种花费最少的方案，输出字典序最小的路径。
# 如果无法到达 N 位置，请返回一个空数组。
# A 数组的长度范围 [1, 1000].

# n,maxJump<=1000

# 总结：
# !相当于求1-n的最短路
# !因为要求字典序最小路径, dp需要倒序, dfs可以直接正序


class Solution:
    def cheapestJump(self, coins: List[int], maxJump: int) -> List[int]:
        n = len(coins)
        if coins[-1] == -1:
            return []

        dist, pre = [INF] * n, [-1] * n  # pre[i]表示i的前驱节点(字典序要最小)
        dist[n - 1] = coins[n - 1]  # !倒序
        queue = deque([(coins[n - 1], n - 1)])
        while queue:
            curDist, cur = queue.popleft()
            if dist[cur] < curDist:
                continue
            for next in range(cur - 1, max(cur - maxJump - 1, -1), -1):
                if coins[next] == -1:
                    continue
                nextDist = curDist + coins[next]
                if nextDist < dist[next]:
                    dist[next] = nextDist
                    pre[next] = cur
                    queue.append((nextDist, next))
                elif nextDist == dist[next]:
                    pre[next] = min(pre[next], cur)

        if dist[0] == INF:
            return []

        path, cur = [], 0
        while cur != -1:
            path.append(cur + 1)
            cur = pre[cur]
        return path

    def cheapestJump2(self, coins: List[int], maxJump: int) -> List[int]:
        """正着记忆化搜索"""

        @lru_cache(None)
        def dfs(index: int) -> Tuple[int, List[int]]:  # 返回值: (cost, path)
            if index >= n or coins[index] == -1:
                return INF, []
            if index == n - 1:
                return coins[index], [n - 1]

            curCost, curPath = INF, []
            for next in range(index + 1, min(index + maxJump + 1, n)):
                nextCost, nextPath = dfs(next)
                candCost, candPath = nextCost + coins[index], [index] + nextPath
                if candCost < curCost:
                    curCost, curPath = candCost, candPath
                elif candCost == curCost and candPath < curPath:
                    curPath = candPath

            return curCost, curPath

        n = len(coins)
        _, path = dfs(0)
        dfs.cache_clear()
        return [num + 1 for num in path]


print(Solution().cheapestJump([1, 2, 4, -1, 2], 2))
# 输出: [1,3,5]
print(Solution().cheapestJump([1, 2, 4, -1, 2], 1))
# []
print(Solution().cheapestJump([0, 0, 0, 0, 0, 0], 3))
