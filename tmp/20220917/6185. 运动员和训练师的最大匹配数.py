# 如果第 i 名运动员的能力值 小于等于 第 j 名训练师的能力值，
# 那么第 i 名运动员可以 匹配 第 j 名训练师。
# 除此以外，每名运动员至多可以匹配一位训练师，每位训练师最多可以匹配一位运动员。
# 请你返回满足上述要求 players 和 trainers 的 最大 匹配数。
from typing import List
from sortedcontainers import SortedList


class Solution:
    def matchPlayersAndTrainers(self, players: List[int], trainers: List[int]) -> int:
        """双指针"""
        players, trainers = sorted(players), sorted(trainers)
        res, ti = 0, 0
        for num in players:
            while ti < len(trainers) and trainers[ti] < num:
                ti += 1
            if ti < len(trainers):
                res += 1
                ti += 1
            else:
                break
        return res

    def matchPlayersAndTrainers2(self, players: List[int], trainers: List[int]) -> int:
        """有序集合"""
        res, sl = 0, SortedList(trainers)
        for num in sorted(players):
            pos = sl.bisect_left(num)
            if pos < len(sl):
                res += 1
                sl.pop(pos)
        return res
