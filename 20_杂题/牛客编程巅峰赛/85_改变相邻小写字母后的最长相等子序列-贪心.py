# 给出一个仅包含小写字母的字符串s，你最多可以操作k次，
# 使得任意一个小写字母变为与其相邻的小写字母（ASCII码差值的绝对值为1），
# 请你求出可能的最长相等子序列（即求这个字符串修改至多k次后的的一个最长子序列，
# 且需要保证这个子序列中每个字母相等）。

#
#
# @param k int整型 表示最多的操作次数
# @param s string字符串 表示一个仅包含小写字母的字符串
# @return int整型
#
class Solution:
    def string2(self, k: int, s: str) -> int:
        n = len(s)
        res = 0
        ords = [ord(s[i]) - 97 for i in range(n)]
        # 计算修改成各个字符的成本
        for i in range(26):
            prices = [abs(i - ords[j]) for j in range(n)]
            # 贪心选择最小的成本
            prices.sort()
            cur = n
            money = k
            for j in range(n):
                if prices[j] > money:
                    cur = j
                    break
                money -= prices[j]
            res = max(cur, res)

        return res