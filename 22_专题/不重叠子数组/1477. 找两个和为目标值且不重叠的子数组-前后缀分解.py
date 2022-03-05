from typing import List


# 请你在 arr 中找 两个互不重叠的子数组 且它们的和都等于 target 。
# 可能会有多种方案，请你返回满足要求的两个子数组长度和的 最小值 。
# 请返回满足要求的最小长度和，如果无法找到这样的两个子数组，请返回 -1 。

# 1 <= arr.length <= 10^5
# 1 <= arr[i] <= 1000  正数 前缀和递增

# 思路：哈希表+前缀和
# 利用前缀和+哈希表可以求出以某一位置为结尾的，和为目标值的子数组的最短长度。
# 对原数组正序和倒序分别求解一次，就可以枚举找到最终的答案。


class Solution:
    def minSumOfLengths(self, arr: List[int], target: int) -> int:
        def getMins(nums: List[int]) -> List[int]:
            """求出以每个位置为结尾的，和为目标值的子数组的最短长度"""
            n = len(nums)
            res = [int(1e20)] * n
            preSum, curSum = {0: -1}, 0

            for i, num in enumerate(nums):
                curSum += num
                if curSum - target in preSum:
                    res[i] = min(res[i], i - preSum[curSum - target])
                if i:
                    res[i] = min(res[i - 1], res[i])
                preSum[curSum] = i

            return res

        res = int(1e20)
        leftMins = getMins(arr)
        rightMins = getMins(arr[::-1])
        print(leftMins, rightMins)
        for i in range(len(arr) - 1):
            left, right = leftMins[i], rightMins[~(i + 1)]
            res = min(res, left + right)
        return res if res < int(1e19) else -1


print(Solution().minSumOfLengths(arr=[3, 1, 1, 1, 5, 1, 2, 1], target=3))  # 3

