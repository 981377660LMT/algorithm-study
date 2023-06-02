# 1272. 删除区间-左闭右开区间
# https://leetcode.cn/problems/remove-interval/

from typing import List


class Solution:
    def removeInterval(self, intervals: List[List[int]], toBeRemoved: List[int]) -> List[List[int]]:
        badStart, badEnd = toBeRemoved
        res = []
        for start, end in intervals:
            if end < badStart or start > badEnd:
                res.append([start, end])
            else:
                if start < badStart:
                    res.append([start, badStart])
                if end > badEnd:
                    res.append([badEnd, end])
        return res
