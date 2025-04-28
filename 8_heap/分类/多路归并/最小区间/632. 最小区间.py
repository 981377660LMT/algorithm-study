# 你有 k 个 非递减排列 的整数列表。找到一个 最小 区间，
# 使得 k 个列表中的每个列表至少有一个数包含在其中。

# 即在 m 个一维数组中各取出一个数字，重新组成新的数组 A，
# !使得新的数组 A 中最大值和最小值的差值（diff）最小。
# !最小值用堆来维护，最大值随指针移动而改变，
#
# 可以 heapreplace 优化.
# https://leetcode.cn/problems/smallest-range-covering-elements-from-k-lists/solutions/2982588/liang-chong-fang-fa-dui-pai-xu-hua-dong-luih5/

from heapq import heapify, heappop, heappush
from typing import List


INF = int(1e18)


def smallestRange(nums: List[List[int]]) -> List[int]:
    leftRes, rightRes = -INF, INF
    pq = [(nums[r][0], r, 0) for r in range(len(nums))]
    heapify(pq)
    max_ = max(item[0] for item in pq)
    while True:
        min_, row, col = heappop(pq)
        if max_ - min_ < rightRes - leftRes:
            leftRes, rightRes = min_, max_
        if col == len(nums[row]) - 1:
            break
        max_ = max(max_, nums[row][col + 1])
        heappush(pq, (nums[row][col + 1], row, col + 1))
    return [leftRes, rightRes]


assert smallestRange([[4, 10, 15, 24, 26], [0, 9, 12, 20], [5, 18, 22, 30]]) == [20, 24]
