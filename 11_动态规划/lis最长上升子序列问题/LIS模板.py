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
    f = bisect_left if isStrict else bisect_right
    for i in range(n):
        pos = f(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
        else:
            lis[pos] = nums[i]

    return len(lis)


def LISDp(nums: List[int], isStrict=True) -> List[int]:
    """求以每个位置为结尾的LIS长度(包括自身)"""
    if not nums:
        return []
    n = len(nums)
    res = [0] * n
    lis = []
    f = bisect_left if isStrict else bisect_right
    for i in range(n):
        pos = f(lis, nums[i])
        if pos == len(lis):
            lis.append(nums[i])
            res[i] = len(lis)
        else:
            lis[pos] = nums[i]
            res[i] = pos + 1
    return res


def getLIS(nums: List[int], isStrict=True) -> Tuple[List[int], List[int]]:
    """求LIS 返回(LIS,LIS的组成下标)"""
    n = len(nums)

    lis = []  # lis[i] 表示长度为 i 的上升子序列的最小末尾值
    dpIndex = [0] * n  # 每个元素对应的LIS长度
    f = bisect_left if isStrict else bisect_right
    for i in range(n):
        pos = f(lis, nums[i])
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


def LISMaxSum(nums: List[int], isStrict=True) -> List[int]:
    """求以每个位置为结尾的LIS最大和(包括自身)"""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    class BITPrefixMax:
        __slots__ = ("_max", "_tree")

        def __init__(self, max: int):
            self._max = max
            self._tree = dict()

        def set(self, index: int, value: int) -> None:
            index += 1
            while index <= self._max:
                self._tree[index] = max(self._tree.get(index, 0), value)
                index += index & -index

        def query(self, end: int) -> int:
            """Query max of [0, end)."""
            if end > self._max:
                end = self._max
            res = 0
            while end > 0:
                res = max(res, self._tree.get(end, 0))
                end -= end & -end
            return res

    n = len(nums)
    if n <= 1:
        return nums[:]
    max_ = 0
    for v in nums:
        max_ = max(max_, v)
    dp = BITPrefixMax(max_ + 5)
    res = [0] * n
    for i, v in enumerate(nums):
        preMax = dp.query(v) if isStrict else dp.query(v + 1)
        cur = preMax + v
        res[i] = cur
        dp.set(v, cur)
    return res


if __name__ == "__main__":
    assert LIS([10, 9, 2, 5, 3, 7, 101, 18]) == 4
    assert getLIS([10, 9, 2, 5, 3, 7, 101, 18]) == ([2, 3, 7, 18], [2, 4, 5, 7])

    assert LISMaxSum([10, 9, 2, 5, 3, 7, 101, 18]) == [10, 9, 2, 7, 5, 14, 115, 32]

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
