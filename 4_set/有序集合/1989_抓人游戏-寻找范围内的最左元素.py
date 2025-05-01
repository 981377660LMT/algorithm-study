# 1989. 捉迷藏中可捕获的最大人数
# https://leetcode.cn/problems/maximum-number-of-people-that-can-be-caught-in-tag/description/
# 抓人游戏，1 抓 0，每个 1 只能抓 1 个 0，
# 且只能抓以它为中心 [i - dist, i + dist] 范围内的 0，问队伍 1 抓 0 的最优解。
#
# 每次抓当前1范围内最左边第一个=>遍历1，维护0的位置
# !每个1配对最左边能配对到的0


from typing import List
from sortedcontainers import SortedList


class Solution:
    def catchMaximumAmountofPeople(self, team: List[int], dist: int) -> int:
        """双指针查找1最左边的未被抓的且在范围内的0."""
        zeros = [i for i, v in enumerate(team) if v == 0]
        ones = [i for i, v in enumerate(team) if v == 1]
        i, j = 0, 0
        res = 0
        while i < len(zeros) and j < len(ones):
            v0, v1 = zeros[i], ones[j]
            if abs(v1 - v0) <= dist:
                res += 1
                i += 1
                j += 1
            elif v0 < v1:
                i += 1
            else:
                j += 1
        return res

    def catchMaximumAmountofPeople2(self, team: List[int], dist: int) -> int:
        zeros = SortedList()
        ones = []
        res = 0
        for i, num in enumerate(team):
            if num == 0:
                zeros.add(i)
            else:
                ones.append(i)

        for i in ones:
            pos = zeros.bisect_left(i - dist)
            if pos == len(zeros):
                break

            zeroPos = zeros[pos]
            if zeroPos <= i + dist:
                zeros.pop(pos)
                res += 1

        return res


print(Solution().catchMaximumAmountofPeople(team=[0, 1, 0, 1, 0], dist=3))
# Output: 2
# 抓两个人
