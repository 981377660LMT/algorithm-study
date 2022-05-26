from typing import List

# cells.length == 8
# 1 <= N <= 10^9
# 如果一间牢房的两个相邻的房间都被占用或都是空的，那么该牢房就会被占用。 => 即 left^right^1
# 否则，它就会被空置。

# 答案：模拟每一天监狱的状态。
# 注意loop从第一天开始，最好是先算出第一天
# https://leetcode.com/problems/prison-cells-after-n-days/discuss/591304/Simple-Python-Solution

# n要与状态同时更新

# 哈希表记录周期
class Solution:
    def prisonAfterNDays(self, cells: List[int], n: int) -> List[int]:
        def move(preState: List[int]):
            return [int(i > 0 and i < 7 and preState[i - 1] == preState[i + 1]) for i in range(8)]

        visited = dict()
        while n:
            visited[tuple(cells)] = n
            cells = move(cells)
            n -= 1
            # 最后再处理n加速
            if tuple(cells) in visited:
                period = visited[tuple(cells)] - n
                n %= period

        return cells


print(Solution().prisonAfterNDays(cells=[0, 1, 0, 1, 1, 0, 0, 1], n=7))
# 输出：[0,0,1,1,0,0,0,0]
# 解释：
# 下表概述了监狱每天的状况：
# Day 0: [0, 1, 0, 1, 1, 0, 0, 1]
# Day 1: [0, 1, 1, 0, 0, 0, 0, 0]
# Day 2: [0, 0, 0, 0, 1, 1, 1, 0]
# Day 3: [0, 1, 1, 0, 0, 1, 0, 0]
# Day 4: [0, 0, 0, 0, 0, 1, 0, 0]
# Day 5: [0, 1, 1, 1, 0, 1, 0, 0]
# Day 6: [0, 0, 1, 0, 1, 1, 0, 0]
# Day 7: [0, 0, 1, 1, 0, 0, 0, 0]
