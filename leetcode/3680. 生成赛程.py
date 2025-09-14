# 3680. 生成赛程(Round-robin tournament 循环赛)
# https://leetcode.cn/problems/generate-schedule/description/
#
# 给你一个整数 n，表示 n 支队伍。你需要生成一个赛程，使得：
#
# 每支队伍与其他队伍 正好比赛两次：一次在主场，一次在客场。
# 每天 只有一场 比赛；赛程是一个 连续的 天数列表，schedule[i] 表示第 i 天的比赛。
# 没有队伍在 连续 两天内进行比赛。
# 返回一个 2D 整数数组 schedule，其中 schedule[i][0] 表示主队，schedule[i][1] 表示客队。如果有多个满足条件的赛程，返回 其中任意一个 。
#
# 如果没有满足条件的赛程，返回空数组。
#
# !贝格尔编排法，在国内排球比赛中广泛使用
# https://en.wikipedia.org/wiki/Round-robin_tournament#Berger_tables
# 数据小，也可以shuffle随机(猴子排序)


from typing import List


class Solution:
    def generateSchedule(self, n: int) -> List[List[int]]:
        if n <= 4:
            return []
        res = []
        # 处理 d=2,3,...,n-2
        for d in range(2, n - 1):
            for i in range(n):
                res.append([i, (i + d) % n])
        # 交错排列 d=1 与 d=n-1（或者说 d=-1）
        for i in range(n):
            res.append([i, (i + 1) % n])
            res.append([(i - 1) % n, (i - 2) % n])
        return res

    def generateSchedule2(self, n: int) -> List[List[int]]:
        import random

        if n <= 4:
            return []

        while True:
            pairs = [[i, j] for i in range(n) for j in range(n) if i != j]
            random.shuffle(pairs)

            res = []
            pre1, pre2 = -1, -1
            while pairs:
                ok = False
                for i, (v1, v2) in enumerate(pairs):
                    if v1 not in (pre1, pre2) and v2 not in (pre1, pre2):
                        res.append([v1, v2])
                        pre1, pre2 = v1, v2
                        pairs.pop(i)
                        ok = True
                        break
                if not ok:
                    break
            else:
                return res
