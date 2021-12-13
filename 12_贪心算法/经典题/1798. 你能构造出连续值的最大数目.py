from typing import List

# 如果你从这些硬币中选出一部分硬币，它们的和为 x ，那么称，你可以 构造 出 x 。
# 请返回从 0 开始（包括 0 ），你最多能 构造 出多少个连续整数。
# 你可能有多个相同值的硬币。
# https://leetcode-cn.com/problems/maximum-number-of-consecutive-values-you-can-make/solution/python-5xing-dong-tai-gui-hua-si-lu-by-y-uxlg/


class Solution:
    def getMaximumConsecutive(self, coins: List[int]) -> int:
        endNum = 0
        for num in sorted(coins):
            if num > endNum + 1:
                break
            endNum += num
        return endNum + 1


print(Solution().getMaximumConsecutive(coins=[1, 1, 1, 4]))
# 输出：8
# 解释：你可以得到以下这些值：
# - 0：什么都不取 []
# - 1：取 [1]
# - 2：取 [1,1]
# - 3：取 [1,1,1]
# - 4：取 [4]
# - 5：取 [4,1]
# - 6：取 [4,1,1]
# - 7：取 [4,1,1,1]
# 从 0 开始，你可以构造出 8 个连续整数。
