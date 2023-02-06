# 330. 按要求补齐数组
# https://leetcode.cn/problems/patching-array/
# 1798. 你能构造出连续值的最大数目
# https://leetcode.cn/problems/maximum-number-of-consecutive-values-you-can-make/

# 给你一个长度为 n 的整数数组 coins ，它代表你拥有的 n 个硬币。
# 如果你从这些硬币中选出一部分硬币，它们的和为 x ，那么称，你可以 构造 出 x 。
# 请返回从 0 开始（包括 0 ），你最多能 构造 出多少个连续整数。
# 你可能有多个相同值的硬币。
from typing import List

# !按照递增顺序使用硬币c与当前的[0,upper] 构造出下一个范围 [c,upper+c]
class Solution:
    def getMaximumConsecutive(self, coins: List[int]) -> int:
        upper = 0
        for num in sorted(coins):
            if num > upper + 1:
                break
            upper += num
        return upper + 1


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
