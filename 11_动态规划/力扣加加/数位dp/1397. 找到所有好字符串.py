# 好字符串 的定义为：它的长度为 n ，字典序大于等于 s1 ，字典序小于等于 s2 ，且不包含 evil 为子字符串。
# s1.length == n
# s2.length == n
# s1 <= s2
# 1 <= n <= 500

# https://leetcode.com/problems/find-all-good-strings/discuss/1133347/Python3-dp-and-kmp-...-finally
class Solution:
    def findGoodStrings(self, n: int, s1: str, s2: str, evil: str) -> int:
        ...


print(Solution().findGoodStrings(n=2, s1="aa", s2="da", evil="b"))
# 输出：51
# 解释：总共有 25 个以 'a' 开头的好字符串："aa"，"ac"，"ad"，...，"az"。
# 还有 25 个以 'c' 开头的好字符串："ca"，"cc"，"cd"，...，"cz"。
# 最后，还有一个以 'd' 开头的好字符串："da"。

# kmp+数位dp 直接放弃
