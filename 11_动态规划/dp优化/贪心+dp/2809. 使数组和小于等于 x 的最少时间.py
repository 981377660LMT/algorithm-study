# 给你两个长度相等下标从 0 开始的整数数组 nums1 和 nums2 。
# !每一秒，对于所有下标 0 <= i < nums1.length ，nums1[i] 的值都增加 nums2[i] 。
# 操作 完成后 ，你可以进行如下操作：
# !选择任一满足 0 <= i < nums1.length 的下标 i ，并使 nums1[i] = 0 。
# 同时给你一个整数 x 。
# 请你返回使 nums1 中所有元素之和 小于等于 x 所需要的 最少 时间，如果无法实现，那么返回 -1 。


# n<=1000
# nums[i]<=1000
# 0<=x<=1e6

# 枚举
# 数据范围暗示O(n^2),可能是dp
# !二分不对
# !贪心+dp
# !种菜模型
# https://leetcode.cn/problems/minimum-time-to-make-array-sum-at-most-x/solutions/2374920/jiao-ni-yi-bu-bu-si-kao-ben-ti-by-endles-2eho/


# !假设已经选好了要操作的元素，那么 num2[i] 越大，操作的时间就应该越靠后。
# 在第t秒,sum1+sum2*t的减少量的最大值相当于求解:
# !将nums2[i]从小到大排序后，从nums1中选一个长为t的子序列，子序列第j个数nums2[j]的系数为j，计算减少量的最大值。
# !=> dfsIndexRemain 01背包dp


from typing import List

INF = int(1e18)


def max(a, b):
    return a if a > b else b


class Solution:
    def minimumTime(self, nums1: List[int], nums2: List[int], x: int) -> int:
        n = len(nums1)
        goods = [(a, b) for a, b in zip(nums1, nums2)]
        goods.sort(key=lambda x: x[1])
        dp = [-INF] * (n + 1)  # 前i个物品选j个物品时的最大减少量
        dp[0] = 0
        for i, (base, delta) in enumerate(goods):
            ndp = dp[:]  # 不选
            for j in range(1, i + 1 + 1):  # 选
                ndp[j] = max(ndp[j], dp[j - 1] + base + j * delta)
            dp = ndp

        sum1 = sum(nums1)
        sum2 = sum(nums2)
        for t, v in enumerate(dp):
            if sum1 + sum2 * t - v <= x:
                return t
        return -1


# nums1 = [1,2,3], nums2 = [1,2,3], x = 4
print(Solution().minimumTime(nums1=[1, 2, 3], nums2=[1, 2, 3], x=4))
print(Solution().minimumTime(nums1=[1], nums2=[4], x=0))
