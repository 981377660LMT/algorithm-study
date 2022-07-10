from itertools import accumulate
from typing import List


class Solution:
    def minSumSquareDiff(
        self, nums1: List[int], nums2: List[int], k1: int, k2: int
    ) -> int:
        """请你返回修改数组 nums1 至多 k1 次且修改数组 nums2 至多 k2 次后的最小 差值平方和 。"""
        diff = [abs((a - b)) for a, b in zip(nums1, nums2)]
        res = minimizeMaxValue(diff, k1 + k2)
        return sum(num * num for num in res)


def minimizeMaxValue(nums: List[int], k: int) -> List[int]:
    """k次-1操作,让最大值最小化,返回操作后的数组"""
    n = len(nums)
    copy = nums[:]
    nums = sorted(nums)
    preSum = [0] + list(accumulate(nums))
    nums = [0] + nums  # [0]表示哨兵

    # !最左二分求最后能和哪个数齐平
    left, right = 0, n
    while left <= right:
        mid = (left + right) // 2
        diff = preSum[n] - preSum[mid] - (n - mid) * nums[mid]
        if k < diff:
            left = mid + 1
        else:
            right = mid - 1

    # 如果最小值可以小到0 那么就直接返回[0]*n
    min_ = nums[left]
    if min_ == 0:
        return [0] * n

    overflow = k - (preSum[n] - preSum[left] - (n - left) * nums[left])
    div, mod = 0, 0
    count = n - left + 1

    if count:
        div, mod = divmod(overflow, count)  # mod个数需要再减1
    min_ -= div

    for i in range(n):
        if copy[i] > min_ - int(mod > 0):
            copy[i] = min_ - int(mod > 0)
            mod -= 1

    return copy


print(
    Solution().minSumSquareDiff(nums1=[1, 2, 3, 4], nums2=[2, 10, 20, 19], k1=0, k2=0)
)
print(Solution().minSumSquareDiff(nums1=[1, 4, 10, 12], nums2=[5, 8, 6, 9], k1=1, k2=1))
print(Solution().minSumSquareDiff([10, 10, 10, 11, 5], [1, 0, 6, 6, 1], 11, 27))
# 预期 0
