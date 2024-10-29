# 给你一个整数数组 nums，其中恰好有两个元素只出现一次，其余所有元素均出现两次。 找出只出现一次的那两个元素。你可以按 任意顺序 返回答案。
# 你必须设计并实现线性时间复杂度的算法且仅使用常量额外空间来解决此问题。


from functools import reduce
from typing import List


def singleNumber1(nums: List[int]) -> int:
    return reduce(lambda x, y: x ^ y, nums, 0)


def singleNumber2(nums: List[int]) -> List[int]:
    xor = reduce(lambda x, y: x ^ y, nums, 0)
    lsb = xor & -xor
    type1 = type2 = 0
    for num in nums:
        if num & lsb:
            type1 ^= num
        else:
            type2 ^= num
    return [type1, type2]


def singleNumber3(nums: List[int]) -> List[int]:
    xor = reduce(lambda x, y: x ^ y, nums, 0)
    needXor = xor != 0
    if needXor:
        for i in range(len(nums)):
            nums[i] ^= xor

    def findUnique() -> int:
        lsbXor = reduce(lambda x, y: x ^ y, map(lambda x: x & -x, nums), 0)
        return reduce(lambda x, y: x ^ y, filter(lambda x: x & lsbXor, nums), 0)

    unique = findUnique()
    nums.append(unique)
    other2 = singleNumber2(nums)
    res = [unique, *other2]
    if needXor:
        for i in range(len(res)):
            res[i] ^= xor
    return res


if __name__ == "__main__":
    print(singleNumber2([1, 2, 1, 3, 2, 5]))
    print(singleNumber3([7, 8, 9, 4, 4, 5, 5]))
