# 221028天池-02. 巡检周期
# https://leetcode.cn/problems/zL2zJU/
# 2 <= record.length <= 12
# 1 <= robot.length <= 5
# 0 <= record[i] <= robot.length
# 0 <= robot[i] < record.length

# 现给定一段时间内的巡检记录，record[i] 表示在时刻 i 出发巡检的机器人数量，
# robot[j] 表示第 j 台机器人首次巡检的时刻。
# 请分离出 record 中各个机器人的记录，并按 robot 的顺序返回各机器人的 巡检周期 。
# 若存在多种分离方式，返回任意一种。

# 笛卡尔积

from itertools import product
from typing import List


class Solution:
    def observingPeriodicity(self, record: List[int], robot: List[int]) -> List[int]:
        # !枚举每个机器人的周期
        n, m = len(record), len(robot)
        for periods in product(range(1, n + 1), repeat=m):
            counter = [0] * n
            for start, period in zip(robot, periods):
                for i in range(start, n, period):
                    counter[i] += 1
            if counter == record:
                return list(periods)

        raise Exception("No solution")
