# 3685. 含上限元素的子序列和-bitset优化可行性01背包问题
# https://leetcode.cn/problems/subsequence-sum-after-capping-elements/solutions/3781344/0-1-bei-bao-shuang-zhi-zhen-pythonjavacg-j4ca/
# 给你一个大小为 n 的整数数组 nums 和一个正整数 k。
#
# 通过将每个元素 nums[i] 替换为 min(nums[i], x)，可以得到一个由值 x 限制（capped）的数组。
#
# 对于从 1 到 n 的每个整数 x，确定是否可以从由 x 限制的数组中选择一个 子序列，使所选元素的和 恰好 为 k。
#
# 返回一个下标从 0 开始的布尔数组 answer，其大小为 n，其中 answer[i] 为 true 表示当 x = i + 1 时可以选出满足要求的子序列；否则为 false。
# 1 <= n == nums.length <= 4000
# 1 <= nums[i] <= n
# 1 <= k <= 4000
from typing import List


class Solution:
    def subsequenceSumAfterCapping(self, nums: List[int], k: int) -> List[bool]:
        nums.sort()
        n = len(nums)
        res = [False] * n
        dp = 1
        mask = (1 << (k + 1)) - 1
        ptr = 0
        for x in range(1, n + 1):
            while ptr < n and nums[ptr] == x:
                dp |= (dp << x) & mask
                ptr += 1

            # 从大于x的数中选了j个x
            for j in range(1 + min(n - ptr, k // x)):
                if dp >> (k - x * j) & 1:
                    res[x - 1] = True
                    break
        return res
