from typing import List, MutableSequence


def nextPermutation(nums: MutableSequence[int], inPlace=False) -> Tuple[bool, MutableSequence[int]]:
    if not inPlace:
        nums = nums[:]

    i = j = len(nums) - 1
    while i > 0 and nums[i - 1] >= nums[i]:
        i -= 1

    if i == 0:  # nums are in descending order
        return False, nums

    k = i - 1  # find the last "ascending" position
    while nums[j] <= nums[k]:
        j -= 1

    nums[k], nums[j] = nums[j], nums[k]

    l, r = k + 1, len(nums) - 1  # reverse the second part
    while l < r:
        nums[l], nums[r] = nums[r], nums[l]
        l += 1
        r -= 1

    return True, nums


class Solution:
    def nextPermutation(self, nums: List[int]) -> None:
        """
        Do not return anything, modify nums in-place instead.
        """
        isOk, res = nextPermutation(nums, True)
        if isOk:
            return res
        else:
            nums.reverse()
