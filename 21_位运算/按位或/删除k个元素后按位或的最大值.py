# 删除k个元素后按位或的最大值
# https://www.geeksforgeeks.org/minimum-bitwise-or-after-removing-at-most-k-elements-from-given-array/

from typing import List


def minOrAfterKRemove(nums: List[int], k: int) -> int:
    max_ = max(nums)
    maxBit = max_.bit_length()
    lengthCounter = [0] * (maxBit + 1)
    for v in nums:
        lengthCounter[v.bit_length()] += 1

    removed = 0
    res = 0
    nums.sort(reverse=True)
    for curBit in range(maxBit, -1, -1):  # 先看最高位能否消除
        curCount = lengthCounter[curBit]
        if curCount <= k:
            k -= curCount
            removed += curCount
        else:
            for _ in range(curCount):
                res |= nums[removed]
                removed += 1
    return res


if __name__ == "__main__":
    print(minOrAfterKRemove([1, 10, 9, 4, 5, 16, 8], 3))
