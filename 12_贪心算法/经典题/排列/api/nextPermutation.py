# 下一个排列/上一个排列

from abc import abstractmethod
from typing import Any, MutableSequence, Protocol, Tuple, TypeVar


class Comparable(Protocol):
    @abstractmethod
    def __eq__(self, other: Any, /) -> bool:
        pass

    @abstractmethod
    def __lt__(self, other: Any, /) -> bool:
        pass

    def __gt__(self, other: Any, /) -> bool:
        return (not self < other) and self != other

    def __le__(self, other: Any, /) -> bool:
        return self < other or self == other

    def __ge__(self, other: Any, /) -> bool:
        return not self < other


S = TypeVar('S', bound=Comparable)


def nextPermutation(nums: MutableSequence[S], inPlace=False) -> Tuple[bool, MutableSequence[S]]:
    """返回下一个字典序的排列，如果不存在，返回本身;时间复杂度O(n)"""
    if not inPlace:
        nums = nums[:]

    left = right = len(nums) - 1

    while left > 0 and nums[left - 1] >= nums[left]:  # 1. 找到最后一个递增位置
        left -= 1
    if left == 0:  # 全部递减
        return False, nums
    last = left - 1  # 最后一个递增位置

    while nums[right] <= nums[last]:  # 2. 找到最小的可交换的right，交换这两个数
        right -= 1
    nums[last], nums[right] = nums[right], nums[last]

    left, right = last + 1, len(nums) - 1  # 3. 翻转后面间这段递减数列
    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
    return True, nums


def prePermutation(nums: MutableSequence[S], inPlace=False) -> Tuple[bool, MutableSequence[S]]:
    """返回前一个字典序的排列,如果不存在,返回本身;时间复杂度O(n)"""
    if not inPlace:
        nums = nums[:]

    left = right = len(nums) - 1

    while left > 0 and nums[left - 1] <= nums[left]:  # 1. 找到最后一个递减位置
        left -= 1
    if left == 0:  # 全部递增
        return False, nums
    last = left - 1  # 最后一个递减位置

    while nums[right] >= nums[last]:  # 2. 找到最大的可交换的right，交换这两个数
        right -= 1
    nums[last], nums[right] = nums[right], nums[last]

    left, right = last + 1, len(nums) - 1  # 3. 翻转后面间这段递增数列
    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
    return True, nums


if __name__ == '__main__':
    isOk, nextP = nextPermutation([1, 2, 3])
    if isOk:
        print("nextP:", nextP)
    isOk, nextP = nextPermutation(list(map(int, '12345')))
    if isOk:
        print("nextP:", nextP)

    isOk, preP = prePermutation([3, 2, 1])
    if isOk:
        print("preP:", preP)
    isOk, preP = prePermutation(list(map(int, '12345')))
    if isOk:
        print("preP:", preP)
