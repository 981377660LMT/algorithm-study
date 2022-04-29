from typing import List, Tuple


class Solution:
    def maxSubArray(self, nums: List[int]) -> Tuple[int, Tuple[int, int]]:
        """最大子数组和,返回数组
        
        dp,需要记录左端点:只取当前还是取前面
        如果前面的和小于0,那么就舍弃前面的一截,并将左端点移到当前位置
        """

        if len(nums) == 1:
            return nums[0], (0, 0)

        maxSum, curSum = -int(1e20), 0
        preLeft = 0
        resLeft, resRight = 0, 0
        for i, num in enumerate(nums):
            if curSum < 0:
                curSum = num
                preLeft = i
            else:
                curSum += num

            if curSum > maxSum:
                maxSum = curSum
                resLeft = preLeft
                resRight = i

        return maxSum, (resLeft, resRight)


if __name__ == '__main__':
    nums = [-2, 1, -3, 4, -1, 2, 1, -5, 4]
    assert Solution().maxSubArray(nums) == (6, (3, 6))
