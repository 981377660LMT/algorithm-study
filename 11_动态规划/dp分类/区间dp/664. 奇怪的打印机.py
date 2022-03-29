from functools import lru_cache

# 有台奇怪的打印机有以下两个特殊要求：

# 打印机每次只能打印由 同一个字符 组成的序列。
# 每次可以在任意起始和结束位置打印新字符，并且会覆盖掉原来已有的字符。
# 给你一个字符串 s ，你的任务是计算这个打印机打印它需要的最少打印次数。
# 1 <= s.length <= 100


class Solution:
    def strangePrinter(self, s: str) -> int:
        @lru_cache(None)
        def dfs(s: str) -> int:
            if not s:
                return 0

            cost = dfs(s[:-1]) + 1  # 1.最坏情况：打一个新的要1块钱
            last_char = s[-1]
            for i, letter in enumerate(s[:-1]):
                if last_char == letter:
                    cost = min(cost, dfs(s[: i + 1]) + dfs(s[i + 1 : -1]))
            return cost

        return dfs(s)


print(Solution().strangePrinter("aaabbb"))
# 输出：2
# 解释：首先打印 "aaa" 然后打印 "bbb"。


# CABBA
#  CA | BB | A
# It was simply inserted with the cost of 1
# It was free from some previous step to the left that printed this character already (we can print extra character all the way till the end)
# 打一个新的要1快钱  cost = dfs(s[:-1]) + 1
# 打一个前面相同的不要钱   dfs(s[: i + 1]) + dfs(s[i + 1 : -1])
