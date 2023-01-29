# 你有 k 个背包。给你一个下标从 0 开始的整数数组 weights ，
# 其中 weights[i] 是第 i 个珠子的重量。同时给你整数 k 。
# 请你按照如下规则将所有的珠子放进 k 个背包。
# 没有背包是空的。
# 如果第 i 个珠子和第 j 个珠子在同一个背包里，那么下标在 i 到 j 之间的所有珠子都必须在这同一个背包中。
# 如果一个背包有下标从 i 到 j 的所有珠子，那么这个背包的价格是 weights[i] + weights[j] 。
# 一个珠子分配方案的 分数 是所有 k 个背包的价格之和。
# 请你返回所有分配方案中，最大分数 与 最小分数 的 差值 为多少。
# !1 <= k <= weights.length <= 1e5
# !1 <= weights[i] <= 1e9

# dp?
# !1e5是为了告诉你它很简单，不需要dp
# 1e5基本上都是On or Onlogn

# !分割成k-1个子数组 => 插入k-1个隔板
# !注意切片负数索引的陷阱 sum(nums[-1:]) 和 sum(nums[-0:])

from typing import List


class Solution:
    def putMarbles(self, weights: List[int], k: int) -> int:
        select = k - 1
        if select == 0:
            return 0
        sums = [a + b for a, b in zip(weights, weights[1:])]
        sums.sort()
        minSum = sum(sums[:select]) + weights[0] + weights[-1]
        maxSum = sum(sums[-select:]) + weights[0] + weights[-1]
        return maxSum - minSum


print(Solution().putMarbles(weights=[1, 3, 5, 1], k=2))
# [25,74,16,51,12,48,15,5]
# 1
print(Solution().putMarbles(weights=[25, 74, 16, 51, 12, 48, 15, 5], k=1))
