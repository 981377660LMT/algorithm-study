from bisect import bisect_left
from typing import List

def merge(intervals: List[List[int]]) -> List[List[int]]:
    """合并所有重叠的有序区间，并返回 一个不重叠的区间数组"""
    res = [intervals[0]]
    for s, e in intervals[1:]:
        if s <= res[-1][1]:
            res[-1][1] = max(res[-1][1], e)
        else:
            res.append([s, e])
    return res

class Solution:
    def insert(self, intervals: List[List[int]], newInterval: List[int]) -> List[List[int]]:
        """在有序的一组区间中插入一个新的区间，你需要确保列表中的区间仍然有序且不重叠"""
        # 二分插入+合并区间
        index = bisect_left(intervals, newInterval[0],key=lambda x: x[0])
        intervals=intervals[:index]+[newInterval]+intervals[index:]
        return merge(intervals)

print(Solution().insert(intervals=[[1, 3], [6, 9]], newInterval=[2, 5]))
        



