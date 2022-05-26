from collections import defaultdict
from typing import List


def maxJump(covers: List[List[int]], start: int, end: int) -> int:
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
    def jump(self, nums: List[int]) -> int:
        intervals = [[i, i + size] for i, size in enumerate(nums)]
        return maxJump(intervals, 0, len(nums) - 1)

