from typing import List

# 只要比第二大两倍就可
class Solution:
    def dominantIndex(self, nums: List[int]) -> int:
        one, oneIndex, two = 0, 0, 0
        for i, num in enumerate(nums):
            if num > one:
                two = one
                one = num
                oneIndex = i
            elif num > two:
                two = num

        return oneIndex if one >= two * 2 else -1
