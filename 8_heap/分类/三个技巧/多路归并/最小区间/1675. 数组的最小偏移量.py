from heapq import heapify, heappop, heappush
from typing import List


# 你可以对数组的任意元素执行任意次数的两类操作：

# 如果元素是 偶数 ，除以 2
# 例如，如果数组是 [1,2,3,4] ，那么你可以对最后一个元素执行此操作，使其变成 [1,2,3,2]
# 如果元素是 奇数 ，乘上 2
# 例如，如果数组是 [1,2,3,4] ，那么你可以对第一个元素执行此操作，使其变成 [2,2,3,4]
# 数组的 偏移量 是数组中任意两个元素之间的 最大差值 。

# 返回数组在执行某些操作之后可以拥有的 最小偏移量 。
# n == nums.length
# 2 <= n <= 5 * 104
# 1 <= nums[i] <= 109

# 多路归并/有序集合
# !nlognlogC
# 632. 最小区间

INF = int(1e18)


class Solution:
    def minimumDeviation(self, nums: List[int]) -> int:
        grid = [[num] for num in nums]

        for i, v in enumerate(grid):
            if v[0] & 1:
                grid[i].append(v[0] * 2)
            else:
                cur = v[0]
                while not cur & 1:
                    cur //= 2
                    grid[i].append(cur)
                grid[i].reverse()

        res = smallestRange(grid)
        return res[1] - res[0]


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
