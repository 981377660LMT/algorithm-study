from bisect import bisect_left
from typing import Callable, Dict, List, Tuple


def discretizeCompressed(nums: List[int]) -> Tuple[Dict[int, int], List[int]]:
    """紧离散化.

    Returns:
        - rank: 原数组中的值 -> 离散化后的值(0 ~ len(rank)-1).
        - newNums: 离散化后的数组.
    """
    allNums = sorted(set(nums))
    rank = {num: i for i, num in enumerate(allNums)}
    newNums = [rank[num] for num in nums]
    return rank, newNums


def discretizeSparse(nums: List[int]) -> Tuple[Callable[[int], int], int]:
    """松离散化.

    Returns:
        - rank: 给定一个数,返回它的排名(0 ~ count).
        - count: 离散化(去重)后的元素个数.
    """
    allNums = sorted(set(nums))
    return lambda x: bisect_left(allNums, x), len(allNums)
