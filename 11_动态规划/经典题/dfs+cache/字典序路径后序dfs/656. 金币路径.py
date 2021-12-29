from typing import List, Tuple
from functools import lru_cache

# 你可以从数组 A 中的任何一个位置（下标为 i）跳到下标 i+1，i+2，……，i+maxJump 的任意一个可以跳到的位置上
# 如果你在下标为 i 的位置上，你需要支付 Ai 个金币。如果 Ai 是 -1，意味着下标为 i 的位置是不可以跳到的。
# 现在，你希望花费最少的金币从数组 A 的 1 位置跳到 N 位置，你需要输出花费最少的路径，依次输出所有经过的下标（从 1 到 N）。
# 如果有多种花费最少的方案，输出`字典顺序最小`的路径(此处为python元组比较大小)。
# 如果无法到达 N 位置，请返回一个空数组。
# A 数组的长度范围 [1, 1000].

INF = 0x3FFFFFFF

# 总结：
# 要获取路径，我们使用dfs后序


class Solution:
    def cheapestJump(self, coins: List[int], maxJump: int) -> List[int]:
        # 从index=cur跳到index=0
        @lru_cache(None)
        def dfs(cur: int) -> Tuple[int, List[int]]:
            # 返回cost 和 path
            if cur == 0:
                return (0, [1])

            if coins[cur] == -1:
                return (INF, [])

            cost, path = INF, []

            maxLeft = max(0, cur - maxJump)
            for next in range(maxLeft, cur):
                if coins[next] == -1:
                    continue
                preCost, prePath = dfs(next)
                costCand, pathCand = preCost + coins[next], prePath + [cur + 1]
                if costCand < cost or (costCand == cost and pathCand < path):
                    # print(cur, path, pathCand, cost, costCand)
                    cost = costCand
                    path = pathCand

            return (cost, path)

        _, path = dfs(len(coins) - 1)

        return path


# print(Solution().cheapestJump([1, 2, 4, -1, 2], 2))
# 输出: [1,3,5]
print(Solution().cheapestJump([1, 2, 4, -1, 2], 1))
# []
