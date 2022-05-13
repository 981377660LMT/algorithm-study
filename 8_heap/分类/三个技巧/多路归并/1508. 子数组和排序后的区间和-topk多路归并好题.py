from typing import List
from heapq import heappop, heappush, heapify

# 你需要计算所有非空连续子数组的和，并将它们按升序排序
# 返回在新数组中下标为 left 到 right （下标从 1 开始）的所有数字和

M = int(1e9 + 7)


class Solution:
    # def rangeSum(self, nums: List[int], n: int, left: int, right: int) -> int:
    #     res = []
    #     for i in range(len(nums)):
    #         prefix = 0
    #         for j in range(i, len(nums)):
    #             prefix += nums[j]
    #             res.append(prefix)
    #     res.sort()
    #     return sum(res[left - 1 : right]) % M

    # K路归并算法
    # n^2logn
    def rangeSum(self, nums: List[int], n: int, left: int, right: int) -> int:
        pq = [(total, index) for index, total in enumerate(nums)]
        heapify(pq)

        res = 0
        for k in range(1, right + 1):
            total, index = heappop(pq)
            if k >= left:
                res += total % M
            if index + 1 < len(nums):
                heappush(pq, (total + nums[index + 1], index + 1))

        return res % M


print(Solution().rangeSum(nums=[1, 2, 3, 4], n=4, left=1, right=5))
# 输出：13
# 解释：所有的子数组和为 1, 3, 6, 10, 2, 5, 9, 3, 7, 4 。将它们升序排序后，我们得到新的数组 [1, 2, 3, 3, 4, 5, 6, 7, 9, 10] 。下标从 le = 1 到 ri = 5 的和为 1 + 2 + 3 + 3 + 4 = 13 。

# 来源：力扣（LeetCode）
# 链接：https://leetcode-cn.com/problems/range-sum-of-sorted-subarray-sums
# 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
