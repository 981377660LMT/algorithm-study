from typing import Generator, List, Tuple


def enumetateConsecutiveSubarrays(
    nums: List[int], diff=1
) -> Generator[Tuple[int, int], None, None]:
    """遍历连续的子数组.

    >>> list(enumetateConsecutiveSubarrays([1, 2, 3, 5, 6, 7, 9]))
    [(0, 3), (3, 6), (6, 7)]
    """
    i, n = 0, len(nums)
    while i < n:
        start = i
        while i < n - 1 and nums[i] + diff == nums[i + 1]:
            i += 1
        i += 1
        yield start, i


if __name__ == "__main__":
    assert list(enumetateConsecutiveSubarrays([1, 2, 3, 5, 6, 7, 9])) == [(0, 3), (3, 6), (6, 7)]
