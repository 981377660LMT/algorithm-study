from typing import List

# 你可以通过按任意顺序连接其中某些数字来形成 3 的倍数，请你返回所能得到的最大的 3 的倍数。
# 由于答案可能不在整数数据类型范围内，请以字符串形式返回答案。
# 1 <= digits.length <= 10^4


# 总结：
# 按模分类  有点像状态机
class Solution:
    def largestMultipleOfThree(self, digits: List[int]) -> str:
        # 最大值每种模的情况
        dp = [-1, -1, -1]

        # 贪心
        for num in sorted(digits, reverse=True):
            # 前导，[0]表示没有前导
            for preVal in dp[:] + [0]:
                curVal = preVal * 10 + num
                dp[curVal % 3] = max(dp[curVal % 3], curVal)

        return str(dp[0]) if dp[0] != -1 else ''


print(Solution().largestMultipleOfThree(digits=[8, 1, 9]))
# 输出："981"
