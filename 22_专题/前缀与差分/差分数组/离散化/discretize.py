from bisect import bisect_left
from typing import Callable, List, Tuple


def discretizeCompressed(nums: List[int], offset=0) -> Tuple[Callable[[int], int], int]:
    """紧离散化.

    Args:
        - nums: 原数组.
        - offset: 离散化的起始值偏移量.

    Returns:
        - getRank: 给定一个数,返回它的排名(offset ~ offset + count).
        - count: 离散化(去重)后的元素个数.
    """
    allNums = sorted(set(nums))
    rank = {num: i + offset for i, num in enumerate(allNums)}
    return lambda x: rank[x], len(allNums)


def discretizeSparse(nums: List[int], offset=0) -> Tuple[Callable[[int], int], int]:
    """松离散化.

    Args:
        - nums: 原数组.
        - offset: 离散化的起始值偏移量.

    Returns:
        - getRank: 给定一个数,返回它的排名(offset ~ offset + count).
        - count: 离散化(去重)后的元素个数.
    """
    allNums = sorted(set(nums))
    return lambda x: bisect_left(allNums, x) + offset, len(allNums)


def discretizeUnique(nums: List[int]) -> Tuple[List[int], List[int]]:
    """不带相同值的离散化，转换为 0-n-1.

    Returns:
        - rank: 离散化后的排名.
        - keys: keys[ranks[i]] = nums[i].
    """

    def argSort(nums: List[int]) -> List[int]:
        order = list(range(len(nums)))
        order.sort(key=lambda i: nums[i])
        return order

    def reArrage(nums: List[int], order: List[int]) -> List[int]:
        return [nums[v] for v in order]

    rank = argSort(nums)
    keys = reArrage(nums, rank)
    rank = argSort(rank)
    return rank, keys


if __name__ == "__main__":
    nums = [1, 2, 34]
    getRank, _ = discretizeSparse(nums)
    print(getRank(99))
