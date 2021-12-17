from typing import List
from sortedcontainers import SortedList

# 抓人游戏，1 抓 0，每个 1 只能抓 1 个 0，
# 且只能抓以它为中心 [i - dist, i + dist] 范围内的 0，问队伍 1 抓 0 的最优解。

# 每次抓当前1范围内最左边第一个=>遍历1，维护0的位置
class Solution:
    def catchMaximumAmountofPeople(self, team: List[int], dist: int) -> int:
        zeros = SortedList()
        ones = []
        res = 0
        for i, num in enumerate(team):
            if num == 0:
                zeros.add(i)
            else:
                ones.append(i)

        for pos in ones:
            index = zeros.bisect_left(pos - dist)
            if index == len(zeros):
                break

            zeroPos = zeros[index]
            if zeroPos <= pos + dist:
                zeros.pop(index)
                res += 1

        return res


print(Solution().catchMaximumAmountofPeople(team=[0, 1, 0, 1, 0], dist=3))
# Output: 2
# 抓两个人
