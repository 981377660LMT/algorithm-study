from typing import List

# 这只青蛙从点 0 处跑道 2 出发，并想到达点 n 处的 任一跑道 ，请你返回 最少侧跳次数 。
# 注意：点 0 处和点 n 处的任一跑道都不会有障碍。
# 1 <= n <= 5 * 105
# dp[i][j]表示第i点第j道最少的侧跳次数


class Solution:
    def minSideJumps(self, obstacles: List[int]) -> int:
        n = len(obstacles)
        dp = [[0x7FFFFFFF] * 3 for _ in range(n)]
        dp[0][0] = 1
        dp[0][1] = 0
        dp[0][2] = 1

        for i in range(1, n):
            # 更新上一次的
            # 如果可以从隔壁跑道移动过来更优，就移动
            for j in range(3):
                dp[i - 1][j] = min(dp[i - 1][j], min(dp[i - 1]) + 1)

            # 可以前进，无需移动
            if obstacles[i] != 1 and obstacles[i - 1] != 1:
                dp[i][0] = dp[i - 1][0]
            if obstacles[i] != 2 and obstacles[i - 1] != 2:
                dp[i][1] = dp[i - 1][1]
            if obstacles[i] != 3 and obstacles[i - 1] != 3:
                dp[i][2] = dp[i - 1][2]

        return min(dp[-1])


print(Solution().minSideJumps(obstacles=[0, 1, 2, 3, 0]))
# 输出：2
# 解释：最优方案如上图箭头所示。总共有 2 次侧跳（红色箭头）。
# 注意，这只青蛙只有当侧跳时才可以跳过障碍（如上图点 2 处所示）。
