from itertools import accumulate
from typing import List

MOD = int(1e9 + 7)


class Solution:
    def rangeSum(self, nums: List[int], n: int, L: int, R: int) -> int:
        """返回在新数组中下标为 L 到 R （下标从 1 开始）的所有数字和（包括左右端点）"""
        # https://leetcode.cn/problems/range-sum-of-sorted-subarray-sums/solution/zi-shu-zu-he-pai-xu-hou-de-qu-jian-he-by-leetcode-/

        def countNGT(mid) -> int:
            """"和小于等于mid的子数组数"""

            res, curSum, left = 0, 0, 0
            for right in range(len(nums)):
                curSum += nums[right]
                while curSum > mid:
                    curSum -= nums[left]
                    left += 1
                res += right - left + 1
            return res

        def kthSmallestSubarraySum(k: int) -> int:
            """第K小的子数组和 k>=1"""
            left, right = 0, preSum1[-1]
            while left <= right:
                mid = (left + right) >> 1
                if countNGT(mid) < k:
                    left = mid + 1
                else:
                    right = mid - 1
            return left

        def cal(k: int) -> int:
            """
            求排序之后的前k个子数组和之和
            先求所有 严格小于 kthSum 的元素和 res 和 元素个数 count,然后再求等于 kthSum 的部分。
            """
            kthSum = kthSmallestSubarraySum(k)  # 第k小的子数组和 要把比这个小的都算进去
            res, right, count = 0, 1, 0
            for left in range(n):
                while right <= n and preSum1[right] - preSum1[left] < kthSum:
                    right += 1
                # [left,right]这一段里的所有以right结尾的子数组都符合题意
                # (a[i] - a[j]) + (a[i] - a[j + 1]) + ... + (a[i] - a[i - 1]).
                # 即 a[i] * (i - j) - (a[j] + a[j + 1] + ... + a[i - 1])
                res += (preSum2[right - 1] - preSum2[left]) - preSum1[left] * (right - 1 - left)
                count += right - 1 - left

            # 加上等于 kthSum的
            res += kthSum * (k - count)
            return res

        preSum1 = [0] + list(accumulate(nums))
        preSum2 = [0] + list(accumulate(preSum1[1:]))
        return (cal(R) - cal(L - 1)) % MOD


print(Solution().rangeSum(nums=[1, 2, 3, 4], n=4, L=1, R=5))
print(Solution().rangeSum(nums=[1, 2, 3, 4], n=4, L=3, R=4))
# from itertools import accumulate
# from typing import List

# MOD = int(1e9 + 7)


# class Solution:
#     def rangeSum(self, nums: List[int], n: int, left: int, right: int) -> int:
#         # https://leetcode.cn/problems/range-sum-of-sorted-subarray-sums/solution/zi-shu-zu-he-pai-xu-hou-de-qu-jian-he-by-leetcode-/
#         def countNGT(mid) -> int:
#             """"和小于等于mid的子数组数"""

#             res, curSum, left = 0, 0, 0
#             for right in range(len(nums)):
#                 curSum += nums[right]
#                 while curSum > mid:
#                     curSum -= nums[left]
#                     left += 1
#                 res += right - left + 1
#             return res

#         def kthSmallestSubarraySum(k: int) -> int:
#             """第K小的子数组和 k>=1"""
#             left, right = 0, preSum1[-1]
#             while left <= right:
#                 mid = (left + right) >> 1
#                 # 找最左，尽量把右边移过来
#                 if countNGT(mid) < k:
#                     left = mid + 1
#                 else:
#                     right = mid - 1
#             return left

#         def getSum(k: int) -> int:
#             num = kthSmallestSubarraySum(k)
#             total, count = 0, 0
#             j = 1
#             for i in range(n):
#                 while j <= n and preSum1[j] - preSum1[i] < num:
#                     j += 1
#                 j -= 1
#                 total += preSum2[j] - preSum2[i] - preSum1[i] * (j - i)
#                 count += j - i
#             total += num * (k - count)
#             return total

#         preSum1 = [0] + list(accumulate(nums))
#         preSum2 = [0] + list(accumulate(preSum1[1:]))
#         return (getSum(right) - getSum(left - 1)) % MOD
