# 2926. 平衡子序列的最大和
# https://leetcode.cn/contest/weekly-contest-370/problems/maximum-balanced-subsequence-sum/
# 划分[l,r]区间的中点mid
# 计算[l,mid]左边部分的贡献
# 计算左边部分可能对右边部分产生的影响
# 计算[mid+1,r]右边部分的贡献


from typing import List


class Solution:
    def maxBalancedSubsequenceSum(self, nums: List[int]) -> int:
        n = len(nums)
        idx = sorted(range(n), key=lambda x: (nums[x] - x, x))
        dp = nums[:]

        def cdq(l, r, pos):
            if l == r:
                return
            mid = (l + r) // 2
            left, right = [], []
            for x in pos:
                if x <= mid:
                    left.append(x)
                else:
                    right.append(x)
            cdq(l, mid, left)
            mx = 0
            for x in pos:
                if x <= mid:
                    mx = max(mx, dp[x])
                else:
                    dp[x] = max(dp[x], mx + nums[x])
            cdq(mid + 1, r, right)

        cdq(0, n - 1, idx)
        return max(dp)
