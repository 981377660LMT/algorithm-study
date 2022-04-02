from itertools import accumulate


# 对所有子数组的和排序，求出sums[i:j+1]的范围和

# n<=1e5
class Solution:
    def solve(self, nums, i, j):
        def countNGT(mid):
            """和小于等于mid的子数组数"""
            res, left = 0, 0
            for right in range(1, n + 1):
                while preSum1[right] - preSum1[left] > mid:
                    left += 1
                res += right - left
            return res

        def cal(k):
            """排序后前k个子数组和的和"""
            lower, upper = 0, preSum1[-1]
            while lower < upper:
                mid = (lower + upper) >> 1
                if countNGT(mid) < k:
                    lower = mid + 1
                else:
                    upper = mid

            # 此时lower为第k小的子数组和
            res, left = 0, 0
            for right in range(1, n + 1):
                while preSum1[right] - preSum1[left] > lower:
                    left += 1
                # [left,right]这一段里的所有子数组都符合题意
                # (a[i] - a[j]) + (a[i] - a[j + 1]) + ... + (a[i] - a[i - 1]).
                # 即 a[i] * (i - j) - (a[j] + a[j + 1] + ... + a[i - 1])
                res += preSum1[right] * (right - left) - (preSum2[right] - preSum2[left])
            # 减去重复计算的(多个lower的情况)
            return res - (countNGT(lower) - k) * lower

        n = len(nums)
        preSum1 = [0] + list(accumulate(nums))
        preSum2 = [0] + list(accumulate(preSum1))

        return cal(j + 1) - cal(i)


print(Solution().solve(nums=[1, 2, 3, 4], i=2, j=3))


# A = [1, 2, 3, 3, 4, 5, 6, 7, 9, 10] here and sum([3, 3]) = 6.
