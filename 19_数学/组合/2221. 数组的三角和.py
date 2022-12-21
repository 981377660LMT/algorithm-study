# 杨辉三角
# 给你一个下标从 0 开始的整数数组 nums ，其中 nums[i] 是 0 到 9 之间（两者都包含）的一个数字。
# nums 的 三角和 是执行以下操作以后最后剩下元素的值：

# !1.nums 初始包含 n 个元素。如果 n == 1 ，终止 操作。
# 否则，创建 一个新的下标从 0 开始的长度为 n - 1 的整数数组 newNums 。
# !2.对于满足 0 <= i < n - 1 的下标 i ，
# newNums[i] 赋值 为 (nums[i] + nums[i+1]) % 10 ，% 表示取余运算。
# !3.将 newNums 替换 数组 nums 。
# !4.从步骤 1 开始 重复 整个过程。
# 请你返回 nums 的三角和。


from typing import List


N = 1000
C = [[0] * (N + 10) for _ in range((N + 10))]
for i in range(N + 5):
    C[i][0] = 1
    for j in range(1, i + 1):
        C[i][j] = (C[i - 1][j - 1] + C[i - 1][j]) % 10  # 题目里是模10


class Solution:
    def triangularSum(self, nums: List[int]) -> int:
        """计算每个元素的贡献次数 -> 杨辉三角(二项式系数)

        找规律发现底层每个元素贡献次数为 C(n-1,i)
        O(n)
        """
        n = len(nums)
        res = 0
        for i, num in enumerate(nums):
            res += num * C[n - 1][i]
            res %= 10
        return res
