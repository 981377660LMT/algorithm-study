# 1288. 删除被覆盖区间-长的区间排在前面
# https://leetcode.cn/problems/remove-covered-intervals/

# 删除包含别人的区间(外部区间)
# !如果存在多个相同的区间，只保留一个
# 区间都是左闭右开


from typing import List, Tuple


def eraseOuterInterval(intervals: List[Tuple[int, int]]) -> List[Tuple[int, int]]:
    intervals = sorted(intervals, key=lambda x: (-x[1], x[0]))
    res = []
    for interval in intervals:
        while res:
            back = res[-1]
            if interval[0] < back[0]:
                break
            res.pop()
        res.append(interval)
    return res[::-1]


def eraseInnerInterval(intervals: List[Tuple[int, int]]) -> List[Tuple[int, int]]:
    intervals = sorted(intervals, key=lambda x: (x[1], -x[0]))
    res = []
    for interval in intervals:
        while res:
            back = res[-1]
            if back[0] < interval[0]:
                break
            res.pop()
        res.append(interval)
    return res


if __name__ == "__main__":
    print(eraseOuterInterval([(1, 2), (1, 2), (2, 3), (3, 4), (1, 3)]))
    print(eraseInnerInterval([(1, 2), (1, 2), (2, 3), (3, 4), (1, 3)]))

    class Solution:
        def removeCoveredIntervals(self, intervals: List[List[int]]) -> int:
            removed = eraseInnerInterval(intervals)
            return len(removed)
