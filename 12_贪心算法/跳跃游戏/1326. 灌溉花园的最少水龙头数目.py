# maxJump-最少区间点覆盖
from collections import defaultdict
from typing import List


def maxJump(covers: List[List[int]], start: int, end: int) -> int:
    """选择最少的区间点覆盖[start,end] start-end <=1e7

    如果超过1e7,需要先离散化

    Args:
        jumps (List[int]): 每个点处可达到的最远坐标
        start (int): 覆盖的区间起点
        end (int): 覆盖区间的右边界

    Returns:
        int: 需要的最少区间数 不存在则返回 -1
    """
    maxJump = defaultdict(int)
    for left, right in covers:
        maxJump[max(start, left)] = max(maxJump[max(start, left)], right)
    cur, next, res = start, start, 0
    for i in range(start, end):
        if cur >= end:
            break
        next = max(next, maxJump[i])
        if i == cur:
            if cur >= next:
                return -1
            cur = next
            res += 1
    return res


class Solution:
    def minTaps(self, n: int, ranges: List[int]) -> int:
        covers = [[i - size, i + size] for i, size in enumerate(ranges)]
        return maxJump(covers, 0, n)
