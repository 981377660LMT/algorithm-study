# 轮转数组/旋转数组 的一段


from typing import List


def rorate(nums: List[int], left: int, right: int, k: int) -> List[int]:
    """
    旋转数组,返回改变后的数组

    :param nums: 数组
    :param left: 左边界
    :param right: 右边界
    :param k: 向右旋转的次数

    :return: 向右轮转后的数组
    """
    len_ = right - left + 1
    k %= len_
    if k == 0:
        return nums
    # !反转后k个元素+翻转前n-k个元素+翻转整个数组
    nums[left : right + 1] = nums[right - k + 1 : right + 1] + nums[left : right - k + 1]
    return nums


print(rorate([1, 2, 3, 4, 5, 6, 7], 0, 6, 3))
