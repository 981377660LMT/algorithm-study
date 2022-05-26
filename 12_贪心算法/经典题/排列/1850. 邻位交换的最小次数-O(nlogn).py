from typing import List, MutableSequence, Tuple
from sortedcontainers import SortedList


def nextPermutation(nums: MutableSequence[int], inPlace=False) -> Tuple[bool, MutableSequence[int]]:
    """返回下一个字典序的排列，如果不存在，返回本身
    
    时间复杂度`O(n)`
    """
    if not inPlace:
        nums = nums[:]

    i = j = len(nums) - 1
    while i > 0 and nums[i - 1] >= nums[i]:
        i -= 1

    if i == 0:  # 不存在下一个字典序的排列
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


def minAdjacentSwap(nums1: MutableSequence[int], nums2: MutableSequence[int]) -> int:
    """求使两个数组相等的最少邻位交换次数，对每个数，贪心找到对应的最近位置交换即可

    时间复杂度`O(n^2)`
    可用求映射+求逆序对的方法优化到`O(nlogn)`
    """
    res = 0

    for num in nums1:
        index = nums2.index(num)  # 最左边的第一个位置
        res += index
        nums2.pop(index)  # 已经被换到最左边了，所以减1

    return res


def countInversionPair(nums: List[int]) -> int:
    """计算逆序对的个数

    sortedList解法,时间复杂度O(nlogn)
    """
    res = 0
    sl = SortedList()
    for num in reversed(nums):
        pos = sl.bisect_left(num)
        res += pos
        sl.add(num)
    return res


class Solution:
    def getMinSwaps(self, num: str, k: int) -> int:
        """当前的num交换到下一个排列的邻位交换次数"""
        digits = list(map(int, num))
        target = digits[:]
        for _ in range(k):
            isOk, nextP = nextPermutation(target)
            if not isOk:
                break
            target = nextP

        """原串到目标串的最少邻位交换次数"""
        return minAdjacentSwap(digits, target)


print(Solution().getMinSwaps("11112", 4))

