# 01串翻转/01串反转
# 给你一个下标从 0 开始、长度为 n 的二进制字符串 s ，你可以对其执行两种操作：
# 选中一个下标 i 并且反转从下标 0 到下标 i（包括下标 0 和下标 i ）的所有字符，成本为 i + 1 。
# 选中一个下标 i 并且反转从下标 i 到下标 n - 1（包括下标 i 和下标 n - 1 ）的所有字符，成本为 n - i 。
# !返回使字符串内所有字符 相等 需要的 最小成本 。
# 反转 字符意味着：如果原来的值是 '0' ，则反转后值变为 '1' ，反之亦然。
# https://leetcode.cn/problems/minimum-cost-to-make-all-characters-equal/


# 方法1:前缀+后缀dp(麻烦)
# https://leetcode.cn/submissions/detail/435719613/
# 方法2:贪心(考察相邻元素)
# !如果s[i-1]!=s[i],那么必须反转其中一个.反转前缀=>i,反转后缀=>n-i
# 取最小值


class Solution:
    def minimumCost(self, s: str) -> int:
        n = len(s)
        res = 0
        for i, (a, b) in enumerate(zip(s, s[1:])):
            if a != b:
                res += min(i + 1, n - i - 1)
        return res


# s = "0011"
# s = "010101"
print(Solution().minimumCost("0011"))
print(Solution().minimumCost("010101"))
print(Solution().minimumCost("10101"))
print(Solution().minimumCost("1"))
