"""
贪心 + 二分查找
LIS[i]表示长度为 i+1 的子序列尾部元素的值
每次遍历到一个新元素,用二分查找法找到第一个大于等于它的元素,然后更新LIS
"""
# LIS模板


from typing import List, Tuple
from bisect import bisect_left, bisect_right


def LIS(nums: List[int], isStrict=True) -> int:
    """求LIS长度"""
    n = len(nums)
    if n <= 1:
        return n

    lis = []  # lis[i] 表示长度为 i 的上升子序列的最小末尾值
    for i in range(n):
        pos = bisect_left(lis, nums[i]) if isStrict else bisect_right(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
        else:
            lis[pos] = nums[i]

    return len(lis)


def getLIS(nums: List[int], isStrict=True) -> Tuple[List[int], List[int]]:
    """求LIS 返回(LIS,LIS的组成下标)"""
    n = len(nums)

    lis = []  # lis[i] 表示长度为 i 的上升子序列的最小末尾值
    dpIndex = [0] * n  # 每个元素对应的LIS长度
    for i in range(n):
        pos = bisect_left(lis, nums[i]) if isStrict else bisect_right(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
        else:
            lis[pos] = nums[i]
        dpIndex[i] = pos

    res, resIndex = [], []
    j = len(lis) - 1
    for i in range(n - 1, -1, -1):
        if dpIndex[i] == j:
            res.append(nums[i])
            resIndex.append(i)

            j -= 1
    return res[::-1], resIndex[::-1]


def caldp(nums: List[int], isStrict=True) -> List[int]:
    """求以每个位置为结尾的LIS长度(包括自身)"""
    if not nums:
        return []

    n = len(nums)
    res = [0] * n
    lis = []
    for i in range(n):
        pos = bisect_left(lis, nums[i]) if isStrict else bisect_right(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
            res[i] = len(lis)
        else:
            lis[pos] = nums[i]
            res[i] = pos + 1
    return res


if __name__ == "__main__":
    assert LIS([10, 9, 2, 5, 3, 7, 101, 18]) == 4
    print(caldp([10, 9, 2, 5, 3, 7, 101, 18]))
    assert getLIS([10, 9, 2, 5, 3, 7, 101, 18]) == ([2, 3, 7, 18], [2, 4, 5, 7])

    # 2826. 将三个组排序
    # https://leetcode.cn/problems/sorting-three-groups/description/
    # 从 0 到 n - 1 的数字被分为编号从 1 到 3 的三个组，数字 i 属于组 nums[i] 。
    # 你可以执行以下操作任意次：
    # 选择数字 x 并改变它的组。更正式的，你可以将 nums[x] 改为数字 1 到 3 中的任意一个。
    # 请你返回将 nums 变为 非递减数组 需要的最少步数。

    # !等价于保留最长的非递减子序列
    class Solution:
        def minimumOperations(self, nums: List[int]) -> int:
            return len(nums) - LIS(nums, isStrict=False)
