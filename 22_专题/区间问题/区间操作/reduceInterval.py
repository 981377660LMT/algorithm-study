# 删除重叠区间


from typing import List, Tuple


def reduceIntervals(intervals: List[Tuple[int, int]], removeIncluded=True) -> List[int]:
    """删除重叠区间.

    Args:
        intervals (List[Tuple[int, int]]): 左闭右开区间.
        removeIncluded (bool, optional): 是删除包含的区间还是被包含的区间.默认为删除被包含的区间.

    Returns:
        List[int]: 按照区间的起点排序的剩余的区间索引(相同的区间会保留).
    """
    n = len(intervals)
    res = []
    order = list(range(n))
    if removeIncluded:
        order.sort(key=lambda i: (intervals[i][0], -intervals[i][1]))
        for cur in order:
            if res:
                pre = res[-1]
                curStart, curEnd = intervals[cur]
                preStart, preEnd = intervals[pre]
                if curEnd <= preEnd and curEnd - curStart < preEnd - preStart:
                    continue
            res.append(cur)
    else:
        order.sort(key=lambda i: (intervals[i][1], -intervals[i][0]))
        for cur in order:
            if res:
                pre = res[-1]
                curStart, curEnd = intervals[cur]
                preStart, preEnd = intervals[pre]
                if curStart <= preStart and curEnd - curStart > preEnd - preStart:
                    continue
            res.append(cur)
    return res


if __name__ == "__main__":
    # https://leetcode.cn/problems/remove-covered-intervals/
    # 1288. 删除被覆盖区间
    class Solution:
        def removeCoveredIntervals(self, intervals: List[List[int]]) -> int:
            return len(reduceIntervals(intervals))
