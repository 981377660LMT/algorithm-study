from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


from typing import Callable, Generator, Tuple


def groupWhile(
    n: int, predicate: Callable[[int, int], bool], skipFalsySingleValueGroup=False
) -> Generator[Tuple[int, int], None, None]:
    """
    分组循环.
    :param n: 数据流长度.
    :param predicate: `[left, curRight]` 闭区间内的元素是否能分到一组.
    :param skipFalsySingleValueGroup: 是否跳过`predicate`返回`False`的单个元素的分组.
    :yield: `[start, end)` 左闭右开分组结果.
    """
    end = 0
    while end < n:
        start = end
        while end < n and predicate(start, end):
            end += 1
        if end == start:
            end += 1
            if skipFalsySingleValueGroup:
                continue
        yield start, end


class Solution:
    def maxIncreasingSubarrays(self, nums: List[int]) -> int:
        def check(left: int, right: int) -> bool:
            if left == right:
                return True
            return nums[right - 1] < nums[right]

        lens = [e - s for s, e in groupWhile(len(nums), check)]
        res = 0
        for pre, cur in zip(lens, lens[1:]):
            res = max(res, min(pre, cur))
        for v in lens:
            res = max(res, v // 2)
        return res
