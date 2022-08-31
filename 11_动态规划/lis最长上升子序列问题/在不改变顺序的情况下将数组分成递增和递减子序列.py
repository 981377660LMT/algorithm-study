# 百度之星
# !在不改变顺序的情况下将数组分成递增子序列和递减子序列
# n<=1e5
# https://www.imangodoc.com/192611.html
# 对于新元素x，
# - 如果它只能附加到数组之一，则附加它。
# - 如果它不能附加到任何一个，那么答案是-1。
# - 如果它可以附加到两个元素，则检查下一个元素y，
# - 如果 y > x 则将x附加到增加的元素，否则将x附加到减少的元素。

from itertools import groupby
from typing import List, Tuple

INF = int(4e18)


def splitArray(nums: List[int]) -> Tuple[List[int], List[int]]:
    """构造出递增和递减子序列"""
    groups = [(num, len(list(group))) for num, group in groupby(nums)]
    up, down = [], []

    for i, (num, count) in enumerate(groups):
        ok1 = not up or up[-1] <= num
        ok2 = not down or down[-1] >= num
        if ok1 and ok2:
            willUp = i + 1 < len(groups) and groups[i + 1][0] > num
            if willUp:
                up.extend([num] * count)
            else:
                down.extend([num] * count)
        elif ok1:
            up.extend([num] * count)
        elif ok2:
            down.extend([num] * count)
        else:
            return [], []

    return up, down


print(splitArray([5, 1, 1, 3, 6, 8, 2, 9, 0, 0, 10]))
