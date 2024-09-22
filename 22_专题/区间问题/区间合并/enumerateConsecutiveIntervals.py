from typing import Generator, List, Tuple


def enumetateConsecutiveIntervals(nums: List[int]) -> Generator[Tuple[int, int, bool], None, None]:
    """遍历连续的区间.

    >>> list(enumetateConsecutiveIntervals([1, 2, 3, 5, 6, 7, 9]))
    [(1, 3, True), (4, 4, False), (5, 7, True), (8, 8, False), (9, 9, True)]
    """
    if not nums:
        return
    i, n = 0, len(nums)
    while i < n:
        start = i
        while i < n - 1 and nums[i] + 1 == nums[i + 1]:
            i += 1
        yield nums[start], nums[i], True
        if i + 1 < n:
            yield nums[i] + 1, nums[i + 1] - 1, False
        i += 1


if __name__ == "__main__":
    assert list(enumetateConsecutiveIntervals([1, 2, 3, 5, 6, 7, 9])) == [
        (1, 3, True),
        (4, 4, False),
        (5, 7, True),
        (8, 8, False),
        (9, 9, True),
    ]

    class Solution:
        # 228. 汇总区间
        # https://leetcode.cn/problems/summary-ranges/description/
        def summaryRanges(self, nums: List[int]) -> List[str]:
            res = []
            for start, end, isIn in enumetateConsecutiveIntervals(nums):
                if not isIn:
                    continue
                res.append(f"{start}->{end}" if start != end else str(start))
            return res
