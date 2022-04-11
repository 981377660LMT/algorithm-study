# Palindrome Count
# 用给定的字符 能组成多少个长k的回文

# 长k的回文包含(k+1)>>1 个对，每个对有set(list(s))种取法
class Solution:
    def solve(self, s, k):
        return pow(len(set(list(s))), (k + 1) >> 1)


print(Solution().solve(s="ab", k=4))
# We can make these palindromes ["aaaa", "bbbb", "abba", "baab"].
