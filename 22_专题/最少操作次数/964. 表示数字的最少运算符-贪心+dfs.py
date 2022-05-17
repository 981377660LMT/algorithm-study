from functools import lru_cache
from math import log

# 2 <= x <= 100
# 1 <= target <= 2 * 10^8
# 我们希望编写一个能使表达式等于给定的目标值 target 且运算符最少的表达式。返回所用运算符的最少数量。
# 其中每个运算符 op1，op2，… 可以是加、减、乘、除（+，-，*，或是 /）之一

# https://leetcode-cn.com/problems/least-operators-to-express-number/comments/296268
# 感觉和LC 818的证明是一样的 https://leetcode.com/problems/race-car/
# 不停做“乘法”直到乘到正好小于target和正好大于target两个数。
# 小于target的那个数加一个“加号”，继续递归。大的那个数反过来走向target。
# 无论哪种方法，距离都在不停接近target。


class Solution:
    def leastOpsExpressTarget(self, x: int, target: int) -> int:
        @lru_cache(None)
        def dfs(need: int) -> int:
            """凑出need的最少需要的符号数"""
            if need < x:  # cur 个 x/x 或者 (x-cur)个 x/x
                return min(2 * need - 1, 2 * (x - need))

            # 乘法最快到达need
            cur, multiCount = x, 0
            while cur * x < need:
                cur *= x
                multiCount += 1

            diff1, diff2 = need - cur, cur * x - need  # 正着走，反着走

            res = dfs(diff1) + multiCount + 1
            if diff2 < diff1:
                res = min(res, dfs(diff2) + multiCount + 1 + 1)

            return res

        res = dfs(target)
        dfs.cache_clear()
        return res


print(Solution().leastOpsExpressTarget(x=3, target=19))
# 输出：5
# 解释：3 * 3 + 3 * 3 + 3 / 3 。表达式包含 5 个运算符。
