from typing import List, Tuple


def discretize(nums: List[int]) -> Tuple[List[int], List[int]]:
    """
    将 nums 中的元素进行离散化，返回新的数组和对应的原始值.

    Args:
        nums (List[int]): 输入的整数数组。

    Returns:
        Tuple[List[int], List[int]]: 返回一个元组，包含两个列表：
            - newNums: 离散化后的数组，其中每个元素是原数组中元素的离散值。
            - origin: 原始值列表，表示每个离散值对应的原始值。

    Example:
        >>> nums = [3, 1, 2, 3, 2]
        >>> newNums, origin = discretize(nums)
        >>> newNums
        [2, 0, 1, 2, 1]
        >>> [origin[x] for x in newNums]
        [3, 1, 2, 3, 2]
    """
    n = len(nums)
    order = sorted(range(n), key=lambda i: nums[i])
    origin = []
    newNums = [0] * n
    for i in order:
        if not origin or origin[-1] != nums[i]:
            origin.append(nums[i])
        newNums[i] = len(origin) - 1
    return newNums, origin
