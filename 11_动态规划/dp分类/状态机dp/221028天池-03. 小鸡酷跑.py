# 最少侧跳次数
# 221028天池-03. 小鸡酷跑
# 蚂蚁庄园的小鸡们组织了一场酷跑障碍赛，每名参加比赛的小鸡都占据三条跑道，
# 三条跑道上都随机布置了若干障碍物:

# 如果 paths[i][j] 为 1，表示在第 i 条跑道上的位置 j 有一个障碍物。
# 如果 paths[i][j] 为 0，表示在第 i 条跑道上的位置 j 为可通过的跑道。
# 小鸡从中间的跑道出发，若从一条跑道换到相邻跑道算一次改道，
# !请问他最少需要改道几次才能到达终点（终点为跑道最后一列的任意位置）。
# 若不能到达终点则返回 -1。


from typing import List


INF = int(1e20)


# !dp[i][j]表示到达第i列第j行最少的侧跳次数
class Solution:
    def chickenCoolRun(self, paths: List[List[int]]) -> int:
        ROW, COL = 3, len(paths[0])
        dp = [0, 0, 0]
        dp[0] = 1 if paths[0][0] != 1 else INF
        dp[1] = 0 if paths[1][0] != 1 else INF
        dp[2] = 1 if paths[2][0] != 1 else INF

        for col in range(1, COL):
            ndp = [INF, INF, INF]
            for curRow in range(3):
                if paths[curRow][col] == 1 or paths[curRow][col - 1] == 1:
                    continue
                for preRow in range(3):
                    if paths[preRow][col - 1] == 1:
                        continue
                    ndp[curRow] = min(ndp[curRow], dp[preRow] + abs(curRow - preRow))
            dp = ndp

        res = min(dp)
        return res if res != INF else -1


print(Solution().chickenCoolRun(paths=[[1, 1, 0, 0], [0, 1, 0, 0], [0, 0, 1, 0]]))
