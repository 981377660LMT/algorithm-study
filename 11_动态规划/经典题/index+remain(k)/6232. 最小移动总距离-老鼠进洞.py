from functools import lru_cache
from typing import List

# 1 <= robot.length, factory.length <= 100
# factory[j].length == 2
# -109 <= robot[i], positionj <= 109
# 0 <= limitj <= robot.length
# !暗示O(n*m*k)的解法1e6
# 老鼠进洞模型 / 安排邮筒模型

INF = int(1e18)


class Solution:
    def minimumTotalDistance(self, robot: List[int], factory: List[List[int]]) -> int:
        """dp[i][j] 表示前i个工厂处理前j个机器人的最小距离之和

        Args:
          robot (List[int]): 机器人的位置
          factory (List[List[int]]): 每个工厂的(位置,可以修理的机器人个数)

        Returns:
            int:  请你返回所有机器人移动的最小总距离。测试数据保证所有机器人都可以被维修。
        """
        n, m = len(robot), len(factory)
        robot.sort()
        factory.sort(key=lambda x: x[0])
        dp = [INF] * (n + 1)
        dp[0] = 0
        for i in range(m):
            ndp = dp[:]
            for j in range(n):
                if dp[j] == INF:
                    break
                dist = 0
                for k in range(factory[i][1]):
                    if j + k >= n:
                        break
                    dist += abs(robot[j + k] - factory[i][0])
                    ndp[j + k + 1] = min(ndp[j + k + 1], dp[j] + dist)
            dp = ndp
        return dp[n]

    def minimumTotalDistance2(self, robot: List[int], factory: List[List[int]]) -> int:
        """dp[i][j][k] 表示前i个工厂处理前j个机器人,最后一个工厂放k个机器人的最小距离之和"""

        @lru_cache(None)
        def dfs(fi: int, ri: int, count: int) -> int:
            """`以工厂为着眼点`,当前工厂为fi,当前处理机器人ri,当前工厂已经修理了count个机器人"""
            if ri == R:
                return 0
            if fi == F:
                return INF

            res = dfs(fi + 1, ri, 0)  # !这个工厂不修了
            if count < factory[fi][1]:
                cand = abs(robot[ri] - factory[fi][0]) + dfs(fi, ri + 1, count + 1)
                res = cand if cand < res else res
            return res

        # !排序之后, 每个工厂修复的机器人肯定是连续的一段(贪心,邻位交换可以证明)
        F, R = len(factory), len(robot)
        robot.sort()
        factory.sort(key=lambda x: x[0])
        res = dfs(0, 0, 0)
        dfs.cache_clear()
        return res


print(Solution().minimumTotalDistance([1, 2, 3, 4, 5], [[1, 2], [3, 2], [4, 1]]))
