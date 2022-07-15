# python中数组的插入与删除用切片快一些 且边界情况不易出错

from collections import defaultdict
from typing import List


def insert1(nums: List[int], index: int, *values: int) -> None:
    """原地插入"""
    nums[index:index] = [*values]


def insert2(nums: List[int], index: int, *values: int) -> List[int]:
    """非原地插入"""
    return nums[:index] + [*values] + nums[index:]


def insertMany(nums: List[int], indexes: List[int], values: List[int]) -> List[int]:
    """
    非原地插入 多个位置插入多个元素

    如果插入的位置超出数组最大索引，则被顺序插入到最后
    """
    assert len(indexes) == len(values), "indexes and values must have the same length"

    n = len(nums)
    inner, outer = defaultdict(list), []
    for i, v in zip(indexes, values):
        if i < n:
            inner[i].append(v)
        else:
            outer.append(v)

    res = []
    for i in range(n):
        if i in inner:
            res.extend(inner[i])
        res.append(nums[i])

    res.extend(outer)
    return res


def pop1(nums: List[int], index: int) -> None:
    """原地删除"""
    nums[index : index + 1] = []


def pop2(nums: List[int], index: int) -> List[int]:
    """非原地删除"""
    return nums[:index] + nums[index + 1 :]


def popMany(nums: List[int], *indexes: int) -> List[int]:
    """非原地删除 删除多个位置的元素 用类似filter的方法删除"""
    bad = set(indexes)
    return [nums[i] for i in range(len(nums)) if i not in bad]


if __name__ == "__main__":
    nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

    insert1(nums, 3, 0, 1)
    print("原地插入后", nums, sep="\n")
    print("非原地插入多个后", insertMany([1, 2, 3], [1, 5, 2, 4, 8], [99, 33, 777, 88, 999]), sep="\n")

    pop1(nums, 3)
    print("原地删除后", nums, sep="\n")
    print("非原地删除后", pop2(nums, 3), sep="\n")
    print("非原地删除多个后", popMany(nums, 1, 2, 3, 4), sep="\n")

    # !此外,删除还可以标记删除
    # 如果只有大量的 insert 和 delete 操作 可以考虑 Splay
