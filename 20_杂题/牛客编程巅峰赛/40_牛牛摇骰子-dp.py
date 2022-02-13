# 一维数轴，初始位置在原点，每次可以选择向左或者向右移动0或3或7或11个单位
# N次询问，每次给出一个坐标arr[i]，求从0点走到arr[i]需要的最少次数

# n<=10^5 arr[i]<=10^9
class Solution:
    def MinimumTimes(self, arr: list[int]) -> list[int]:
        # write code here
        dp_init = [0, 3, 4, 1, 2, 3, 2, 1, 2, 3, 2, 1, 4, 3, 2, 3, 4, 3, 2, 3, 4, 3, 2, 5, 4, 3]
        max_num = max(arr)
        dp = [0 for i in range(max_num + 1)]
        # 初始化
        for i in range(1, 12):
            dp[i] = dp_init[i]

        # 动态转移方程
        for i in range(12, max_num + 1):
            dp[i] = min(dp[i - 11], dp[i - 7], dp[i - 3]) + 1

        res = []
        for i in arr:
            res.append(dp[i])
        return res

