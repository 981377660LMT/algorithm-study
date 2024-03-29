# 给你一个长度为 n 、下标从 0 开始的整数数组 nums ，表示收集不同巧克力的成本。每个巧克力都对应一个不同的类型，最初，位于下标 i 的巧克力就对应第 i 个类型。

# 在一步操作中，你可以用成本 x 执行下述行为：


# 同时对于所有下标 0 <= i < n - 1 进行以下操作， 将下标 i 处的巧克力的类型更改为下标 (i + 1) 处的巧克力对应的类型。如果 i == n - 1 ，则该巧克力的类型将会变更为下标 0 处巧克力对应的类型。
# 假设你可以执行任意次操作，请返回收集所有类型巧克力所需的最小成本。


# 解释：最开始，巧克力的类型分别是 [0,1,2] 。我们可以用成本 1 购买第 1 个类型的巧克力。
# 接着，我们用成本 5 执行一次操作，巧克力的类型变更为 [2,0,1] 。我们可以用成本 1 购买第 0 个类型的巧克力。
# 然后，我们用成本 5 执行一次操作，巧克力的类型变更为 [1,2,0] 。我们可以用成本 1 购买第 2 个类型的巧克力。
# 因此，收集所有类型的巧克力需要的总成本是 (1 + 5 + 1 + 5 + 1) = 13 。可以证明这是一种最优方案。


# !遇到影响全局的操作 一般会想到枚举操作次数


from typing import List


class Solution:
    def minCost(self, nums: List[int], x: int) -> int:
        n = len(nums)
        res = sum(nums)
        mins = nums[:]
        for i in range(1, n):  # 旋转次数
            for j in range(n):
                mins[j] = min(mins[j], nums[(j + i) % n])
            res = min(res, sum(mins) + x * i)
        return res
