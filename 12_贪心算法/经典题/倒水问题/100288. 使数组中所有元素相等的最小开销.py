# 100288. 使数组中所有元素相等的最小开销
# https://leetcode.cn/problems/minimum-cost-to-equalize-array/description/
# 给你一个整数数组 nums 和两个整数 cost1 和 cost2 。你可以执行以下 任一 操作 任意 次：
# 从 nums 中选择下标 i 并且将 nums[i] 增加 1 ，开销为 cost1。
# 选择 nums 中两个 不同 下标 i 和 j ，并且将 nums[i] 和 nums[j] 都 增加 1 ，开销为 cost2 。
# 你的目标是使数组中所有元素都 相等 ，请你返回需要的 最小开销 之和。
# 由于答案可能会很大，请你将它对 109 + 7 取余 后返回。


from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minCostToEqualizeArray(self, nums: List[int], cost1: int, cost2: int) -> int:
        nums.sort()
        max_ = nums[-1]
        if len(nums) == 1:
            return 0
        if len(nums) == 2:
            return (nums[1] - nums[0]) * cost1 % MOD

        # 倒水问题，水壶问题
        # https://leetcode.cn/problems/maximum-number-of-weeks-for-which-you-can-work/solutions/908442/ni-ke-yi-gong-zuo-de-zui-da-zhou-shu-by-rbidw/
        def solve(target: int) -> int:
            diff = [target - num for num in nums if num < target]
            diffSum = sum(diff)
            if cost1 * 2 <= cost2:
                return diffSum * cost1
            if not diff:
                return 0
            if len(diff) == 1:
                return diff[0] * cost1
            maxDiff = diff[0]
            otherDiff = diffSum - maxDiff
            if maxDiff <= otherDiff:
                if diffSum & 1 == 0:
                    return (diffSum // 2) * cost2
                return (diffSum // 2) * cost2 + cost1
            return otherDiff * cost2 + (maxDiff - otherDiff) * cost1

        res = solve(max_)
        res = min2(res, solve(max_ + 1))

        def getBestTarget() -> int:
            diff1 = max_ - nums[0]
            diff2 = sum(max_ - num for num in nums[1:])
            for target in range(max_, INF):
                if diff1 <= diff2:
                    return target
                diff1 += 1
                diff2 += len(nums) - 1

        best = getBestTarget()
        res = min2(res, solve(best))
        res = min2(res, solve(best + 1))
        if best - 1 >= max_:
            res = min2(res, solve(best - 1))
        return res % MOD


# [1,14,14,15] -> 21 最佳，答案20
# 2
# 1

# print(Solution().minCostToEqualizeArray([1, 14, 14, 15], 2, 1))
# [4,3,1,8] -》 答案8
# 5
# 1
print(Solution().minCostToEqualizeArray([2, 12, 24], 49, 4))
# [2,12,24] -》 答128
# 49
# 4
# [55,52,29,11]
# 18
# 2
# print(Solution().minCostToEqualizeArray([55, 52, 29, 11], 18, 2)) -118
