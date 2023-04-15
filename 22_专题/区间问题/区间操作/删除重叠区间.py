# 删除重叠区间


# 删除包含别人的区间(外部区间)
# !如果存在多个相同的区间，只保留一个

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
    res.reverse()
    return res


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


print(eraseOuterInterval([(1, 2), (1, 2), (2, 3), (3, 4), (1, 3)]))
print(eraseInnerInterval([(1, 2), (1, 2), (2, 3), (3, 4), (1, 3)]))
