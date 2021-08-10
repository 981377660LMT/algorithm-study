#  给你一个整数数组 nums 以及两个整数 lower 和 upper 。
#  求数组中，值位于范围 [lower, upper] （包含 lower 和 upper）之内的 区间和的个数 。
#  区间和sum(i,j) = S[j] - S[i]

import bisect
from typing import List


def solution(list: List[int], low: int, upper: int):
    res = 0
    sum: List[int] = [0]
    cur_sum = 0
    for value in list:
        cur_sum += value
        res += bisect.bisect_right(sum, cur_sum - low) - bisect.bisect_left(sum, cur_sum - upper)
        sum.append(cur_sum)
    return res


print(solution([-2, 5, -1], -2, 2))
